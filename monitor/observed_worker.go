package monitor

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"sync"
)

const (
	Stopped = iota
	Running
	Failed
	Completed
)

func genUid() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

type MonitoredWorker struct {
	Itw   IterationWork
	wgrun sync.WaitGroup
	guid  string
	state int
	chsig chan int
	stwg  sync.WaitGroup
}

type IterationWork interface {
	DoWork() (bool, error)
	GetProgress() interface{}
}

func (mw *MonitoredWorker) wgoroute() {
	mw.wgrun.Add(1)
	mw.state = Running
	for {
		select {
		case newState := <-mw.chsig:
			if newState == Stopped {
				mw.state = newState
				log.Println("info: work stopped")
				mw.wgrun.Done()
				return
			}
		default:
			{
				isdone, err := mw.Itw.DoWork()
				if err != nil {
					log.Println("error: work failed")
					mw.state = Failed
					mw.wgrun.Done()
					return
				}
				if isdone {
					mw.state = Completed
					mw.wgrun.Done()
					log.Println("info: work done")
					return
				}
			}

		}
	}
}
func (mw MonitoredWorker) GetState() int {
	return mw.state
}
func (mw MonitoredWorker) GetId() string {
	return mw.guid

}
func (mw *MonitoredWorker) Start() (string, error) {
	mw.wgrun.Wait()
	if mw.state == Running {
		return "", errors.New("error: try start runing job")
	}
	mw.guid = genUid()
	mw.chsig = make(chan int, 1)
	go mw.wgoroute()
	return mw.guid, nil
}

func (mw *MonitoredWorker) Stop() error {
	if mw.state == Stopped {
		return errors.New("error: can't stop stoped work")
	}
	mw.chsig <- Stopped
	mw.wgrun.Wait()
	return nil

}
func (mw MonitoredWorker) GetProgress() interface{} {
	return mw.Itw.GetProgress()

}

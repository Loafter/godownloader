package monitor

import (
	"crypto/rand"
	"errors"
	"log"
	"sync"
	"fmt"
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
	log.Println("info: work start", mw.GetId())
	mw.wgrun.Add(1)
	defer func() {
		log.Print("info: realease work guid", mw.GetId())
		mw.wgrun.Done()
	}()

	mw.state = Running
	for {
		select {
		case newState := <-mw.chsig:
			if newState == Stopped {
				mw.state = newState
				log.Println("info: work stopped")
				return
			}
		default:
			{
				isdone, err := mw.Itw.DoWork()
				if err != nil {
					log.Println("error: guid", mw.guid, " work failed", err)
					mw.state = Failed
					return
				}
				if isdone {
					mw.state = Completed
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
func (mw *MonitoredWorker) GetId() string {
	if len(mw.guid)==0 {
		mw.guid = genUid()
	}
	return mw.guid

}
func (mw *MonitoredWorker) Start() error {
	if mw.state == Running {
		errors.New("error: try run runing job")
	}
	mw.chsig = make(chan int, 1)
	go mw.wgoroute()
	return nil
}

func (mw *MonitoredWorker) Stop() error {
	if mw.state == Stopped {
		panic("imposible start runing job")
	}
	mw.chsig <- Stopped
	mw.wgrun.Wait()
	close(mw.chsig)
	return nil

}
func (mw MonitoredWorker) GetProgress() interface{} {
	return mw.Itw.GetProgress()

}

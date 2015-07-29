package dtest

import (
	"testing"
	"godownloader/monitor"
	"time"
	"errors"
	"fmt"
	"math/rand"
	"log"
)

type TestWorkPool struct {
	From, id, To int32
}

func (tw TestWorkPool) GetProgress() interface{} {
	return tw.From

}


func (tw *TestWorkPool)BeforeRun()error{
	log.Println("info: exec before run")
	return nil
}
func (tw *TestWorkPool)AfterStop()error{
	log.Println("info: after stop")
	return nil
}


func (tw *TestWorkPool) DoWork() (bool, error) {
	time.Sleep(time.Millisecond*300)
	tw.From += 1
	log.Print(tw.From)
	if tw.From == tw.To {
		fmt.Println("done")
		return true, nil
	}
	if tw.From > tw.To {
		return false, errors.New("tw.From > tw.To")
	}
	return false, nil
}
func TestWorkerPool(t *testing.T) {
	wp := monitor.WorkerPool{}
	for i := 0; i < 20; i++ {
		mw := &monitor.MonitoredWorker{Itw:&TestWorkPool{From:0, To:20, id:rand.Int31()}}
		wp.AppendWork(mw)
	}
	wp.StartAll()
	time.Sleep(time.Second)
	log.Println("------------------Work Started------------------")
	log.Println(wp.GetAllProgress())
	log.Println("------------------Get All Progress--------------")
	time.Sleep(time.Second)
	wp.StopAll()
	log.Println("------------------Work Stop-------------------")

	time.Sleep(time.Second)
	wp.StartAll()
	time.Sleep(time.Second*5)
	wp.StopAll()
	wp.StartAll()
	wp.StopAll()
}

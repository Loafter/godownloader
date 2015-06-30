package dtest

import (
	"testing"
	"godownloader/monitor"
	"time"
	"math/rand"
	"errors"
	"log"
)

type TestWorkPool struct {
	From, sleep, To int32
}

func (tw TestWorkPool) GetProgress() interface{} {
	return tw.From

}
func (tw *TestWorkPool) DoWork() (bool, error) {
	time.Sleep(time.Millisecond*300)
	tw.From += 1
	if tw.From > tw.To {
		return false, errors.New("failed")
	}
	if tw.From == tw.To {
		log.Println("done")
		return true, nil
	}
	return false, nil
}


func TestWorkerPool(t *testing.T) {
	wp := monitor.WorkerPool{}
	wp.Init()
	for i := 0; i < 20; i++ {
		w :=TestWorkPool{From: rand.Int31n(5), To: 5 + rand.Int31n(26), sleep: 300 + rand.Int31n(1000)}
		mw:=monitor.MonitoredWorker{Itw:&w}
		wp.AppendWork(mw)
	}
	//test normal work
	wp.StartAll()
	time.Sleep(time.Second)
	wp.GetAllProgress()
	time.Sleep(time.Second)
	wp.StopAll()

	wp.StartAll()
	//start running work
	wp.StartAll()

	wp.StopAll()
	//stop stopped work
	wp.StopAll()

}

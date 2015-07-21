package dtest

import "godownloader/monitor"
import "testing"
import (
	"errors"
	"time"
	"log"
)

type TestWork struct {
	From, sleep, To int
}

func (tw TestWork) GetProgress() interface{} {
	return tw.From
}
func (tw *TestWork) DoWork() (bool, error) {
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

func TestWorker(t *testing.T) {
	tes := new(monitor.MonitoredWorker)
	itw := &TestWork{From: 1, To: 8, sleep: 300}
	tes.Itw = itw
	tes.Start()
	time.Sleep(time.Second * 1)
	if tes.GetState() != 1 {
		t.Error("Expected Running(1)")
		return
	}
	tes.Stop()
	if tes.GetState() != 0 {
		t.Error("Expected Stoped(0)")
		return
	}
	tes.Start()
	time.Sleep(time.Second * 9)
	if tes.GetState() != 3 {
		t.Error("Expected Comlete(3)")
		return
	}

	tes.Start()
	time.Sleep(time.Second * 1)
	if tes.GetState() != 2 {
		t.Error("Expected Failed(3)")
		return
	}
}

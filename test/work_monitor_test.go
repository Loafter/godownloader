package dtest

import "godownloader/monitor"
import "testing"
import (
	"log"
	"time"
	"errors"
)



type TestWork struct {
	 From, To int
}

func (tw TestWork)GetProgress() (interface{}){
	return  tw.From;

}
func (tw *TestWork)DoWork() (bool,error){
	time.Sleep(time.Millisecond * 300)
	tw.From += 1
	if tw.From > tw.To{
		return false,errors.New("failed")
	}
	if tw.From == tw.To {
		return true, nil
	}
	log.Println(tw)
	return false, nil
}

func TestWorker(*testing.T) {
	tes := new(monitor.MonitoredWorker)
	itw := &TestWork{From:1, To:8}
	tes.Itw = itw
	tes.Start()
	time.Sleep(time.Second * 1)
	log.Println("State:",tes.GetState())
	tes.Stop()
	log.Println("State:",tes.GetState())

	tes.Start()
	time.Sleep(time.Second * 9)
	log.Println("State:",tes.GetState())
	log.Println("Result:",tes.GetResult())

	tes.Start()
	time.Sleep(time.Second * 1)
	log.Println("State:",tes.GetState())
}


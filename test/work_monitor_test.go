package dtest

import "godownloader/monitor"
import "testing"
import (
	"errors"
	"log"
	"time"
)

func TFunc(from interface{}, to interface{}) (interface{}, interface{}, bool, error) {
	time.Sleep(time.Millisecond * 300)

	var f, t int
	var ok bool
	f, ok = from.(int)
	if !ok {
		log.Println("error:  missmatch from type")
		return nil, nil, false, errors.New("error:  missmatch from type")
	}
	t, ok = to.(int)
	if !ok {
		log.Println("error:  missmatch from type")
		return nil, nil, false, errors.New("error:  missmatch to type")
	}
	f += 1
	if f == t {
		return f, t, true, nil
	}
	log.Println(f, t)
	return f, t, false, nil

}
func TestWorker(*testing.T) {
	tes := new(monitor.MonitoredWorker)
	tes.Func = TFunc
	tes.Start(10, 20)
	time.Sleep(time.Second * 4)
	tes.Stop()

	tes.Start(5, 10)
	time.Sleep(time.Second * 8)
}

package dtest

import "testing"

import "godownloader/monitor"
import "godownloader/test"
import (
	"math/rand"
	"time"
)

func TestWorkerPool(t *testing.T) {
	wp := monitor.WorkerPool{}
	for i := 0; i < 20; i++ {
		w := dtest.TestWork{From: rand.Int31n(5), To: 5 + rand.Int31n(26), sleep: 300 + rand.Int31n(1000)}
		wp.AppendWork(w)
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

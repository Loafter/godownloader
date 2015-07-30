package monitor
import (
	"log"
)

type WorkerPool struct {
	workers map[string]*MonitoredWorker
}
func (wp *WorkerPool)AppendWork(iv *MonitoredWorker) {
	if wp.workers==nil{
		wp.workers=make(map[string]*MonitoredWorker)
	}
	wp.workers[iv.GetId()]=iv
}
func (wp *WorkerPool)StartAll() {

	for _, value := range wp.workers {
		value.Start()
	}
}

func (wp *WorkerPool)StopAll() {
	for _, value := range wp.workers {
		log.Println(value.Stop())
	}
}

func (wp *WorkerPool)GetAllProgress()interface{} {
	rs := make([]interface{}, 1)
	for _, value := range wp.workers {
		rs=append(rs, value.GetProgress())
	}
	return rs
}
package monitor
import "log"

type WorkerPool struct {
	dwks map[string]MonitoredWorker
}
func (wp *WorkerPool) Init(){
	wp.dwks=make(map[string]MonitoredWorker)
}
func (wp *WorkerPool)AppendWork(iv MonitoredWorker) {
	wp.dwks[iv.GetId()]=iv
}
func (wp *WorkerPool)StartAll() {
	for key, value := range wp.dwks {
		status:=value.Start()
		log.Println("info: job guid ",key," status ",status)
	}
}

func (wp *WorkerPool)StopAll() {
	for _, value := range wp.dwks {
		log.Println(value.Stop())
	}
}

func (wp *WorkerPool)GetAllProgress()[]interface{} {
	rs:=make([]interface{},1)
	for _, value := range wp.dwks {
		rs=append(rs,value.GetProgress())
	}
	return rs
}
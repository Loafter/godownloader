package monitor

import "log"
import "errors"

type WorkerMonitor struct {
	jobs      map[string]ObservedWorker
	wrkchn    chan string
	OnFailed  func(error)
	OnDone    func(wm *WorkerMonitor) interface{}
	OnStart   func(wm *WorkerMonitor) interface{}
	OnStopped func(wm *WorkerMonitor) interface{}
}

func (wm *WorkerMonitor) waitdone() {
	for {
		isclean := true
		for _, val := range wm.jobs {
			if val.GetState() != Stopped {
				<-wm.wrkchn
				isclean = false
			}
		}
		if isclean {
			return
		}
	}
}

func (wm *WorkerMonitor) monitor() {
	for {
		chw := <-wm.wrkchn
		if chw == "" {
			wm.waitdone()
			return
		}
		if worker, ok := wm.jobs[chw]; !ok {
			log.Println("error: can't find jobs with guid=%s", chw)
		} else {
			switch worker.GetState() {
			case Stopped:
				{
					return
				}

			}

		}
	}
}

func (wm *WorkerMonitor) StartJob(ow ObservedWorker, prm ...interface{}) error {
	if guid, err := ow.Start(wm.wrkchn, prm); err != nil {
		log.Println("error: can't start job guid=%s", guid)
		return errors.New("error: can't start job")
	} else {
		wm.jobs[guid] = ow
		log.Println("info:success start job guid=%s")
		return nil
	}
}

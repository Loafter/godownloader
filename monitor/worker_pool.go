package monitor

type WorkerPool struct {
	workers map[string]*MonitoredWorker
}

func (wp *WorkerPool) AppendWork(iv *MonitoredWorker) {
	if wp.workers == nil {
		wp.workers = make(map[string]*MonitoredWorker)
	}
	wp.workers[iv.GetId()] = iv
}
func (wp *WorkerPool) StartAll() []error {
	var errs []error
	for _, value := range wp.workers {
		if err := value.Start(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (wp *WorkerPool) StopAll() []error {
	var errs []error
	for _, value := range wp.workers {
		if err := value.Stop(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (wp *WorkerPool) GetAllProgress() interface{} {
	var pr []interface{}
	for _, value := range wp.workers {
		pr = append(pr, value.GetProgress())
	}
	return pr
}

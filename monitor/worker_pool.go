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
func (wp *WorkerPool) StartAll() error {

	for _, value := range wp.workers {
		if err := value.Start(); err != nil {
			return err
		}

	}
	return nil
}

func (wp *WorkerPool) StopAll() []error {
	var errs []error
	for _, value := range wp.workers {
		errs = append(errs, value.Stop())
	}
	return errs
}

func (wp *WorkerPool) GetAllProgress() interface{} {
	var errs []error
	for _, value := range wp.workers {
		errs = append(errs, value.Start())
	}
	return errs
}

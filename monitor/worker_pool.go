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

func (wp *WorkerPool) StopAll() error {
	for _, value := range wp.workers {
		if err := value.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (wp *WorkerPool) GetAllProgress() interface{} {
	rs := make([]interface{}, 1)
	for _, value := range wp.workers {
		rs = append(rs, value.GetProgress())
	}
	return rs
}

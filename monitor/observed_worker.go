package monitor

const (
	Stopped = iota
	Running
	Failed
	Completed
)

type ObservedWorker interface {
	GetState() int
	GetId() string
	Start(cn chan string, prm ...interface{}) (string, error)
	Stop() error
	GetDone() interface{}
	GetDescription() interface{}
	GetResult() (interface{}, error)
}

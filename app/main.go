package main

import "fmt"
import "godownloader/monitor"

func main() {
	dt := monitor.WorkerMonitor{}
	dt.StartJob(nil, nil)
	fmt.Printf("Hello, world.\n")
}

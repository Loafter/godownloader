package dtest
import "godownloader/monitor"
import "testing"
import "time"
func TestWorker(*testing.T) {
	tes:=new(monitor.MonitoredWorker)
	tes.Start(10,30)
	time.Sleep(time.Second*1)
	tes.Stop()
	time.Sleep(time.Second*2)
}

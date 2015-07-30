package dtest

import (
	"godownloader/http"
	"godownloader/monitor"
	"log"
	"os"
	"testing"
	"time"
)

func TestPartDownloadWorker(t *testing.T) {
	c, _ := httpclient.GetSize("http://ports.ubuntu.com/dists/precise/main/installer-powerpc/current/images/powerpc/netboot/mini.iso")
	f, _ := os.Create("part_download.data")
	defer f.Close()
	f.Truncate(c)
	dow := new(httpclient.PartialDownloader)
	dow.Init("http://ports.ubuntu.com/dists/precise/main/installer-powerpc/current/images/powerpc/netboot/mini.iso", f, 0, c)
	var i monitor.IterationWork
	i = dow
	/*i.BeforeRun()
	for {
		sta, _ := i.DoWork()
		if sta {
			return
		}
		log.Println(i.GetProgress())
	}*/
	mv := monitor.MonitoredWorker{Itw: i}
	log.Println(mv.Start())
	log.Println(mv.Start())
	time.Sleep(time.Second * 1)
	/*log.Println(mv.Stop())
	time.Sleep(time.Second*2)
	log.Println(mv.Start())
	log.Println(mv.Start())
	time.Sleep(time.Second*5)
	log.Println(mv.Stop())*/
}

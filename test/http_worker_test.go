package dtest

import (
	"testing"
	"godownloader/http"
	"os"

	"godownloader/monitor"
	"time"
)

func TestPartDownloadWorker(t *testing.T) {
	c, e := httpclient.GetSize("http://ports.ubuntu.com/dists/precise/main/installer-powerpc/current/images/powerpc/netboot/mini.iso")
	if (e!=nil) {
		t.Error("failed: Get size must be without error")
	}
	f,_:=os.Create("part_download.data")
	defer f.Close()
	f.Truncate(c);
	dow:=new(httpclient.PartialDownloader)
	dow.Init("http://ports.ubuntu.com/dists/precise/main/installer-powerpc/current/images/powerpc/netboot/mini.iso",f,0,0,c)
	mv:=monitor.MonitoredWorker{Itw:dow}
	mv.Start()
	time.Sleep(time.Second*10)
	mv.Stop()

}

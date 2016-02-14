package dtest

import (
	"godownloader/http"
	"godownloader/iotools"
	"godownloader/monitor"
	"testing"
	"time"
)

func TestMultiPartDownloadPool(t *testing.T) {
	partcount := 10
	pc := int64(partcount)
	url := "http://pinegrow.s3.amazonaws.com/PinegrowLinux64.2.2.zip"
	c, _ := httpclient.GetSize(url)
	f, _ := iotools.CreateSafeFile("ubuntu-15.04-snappy-amd64+generic.img.xz ")
	defer f.Close()
	f.Truncate(c)
	ps := c / pc
	wp := monitor.WorkerPool{}
	for i := int64(0); i < pc-1; i++ {
		//log.Println(ps*i, ps*i+ps)
		d := httpclient.CreatePartialDownloader(url, f, ps*i, ps*i, ps*i+ps)
		mv := monitor.MonitoredWorker{Itw: d}
		wp.AppendWork(&mv)
	}
	lastseg := c - (ps * (pc - 1))
	dow := httpclient.CreatePartialDownloader(url, f, lastseg,lastseg, c)
	mv := monitor.MonitoredWorker{Itw: dow}
	wp.AppendWork(&mv)
	wp.StartAll()
	time.Sleep(time.Second * 30)

}

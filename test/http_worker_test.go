package dtest

import (
	"godownloader/http"
	"godownloader/iotools"
	"godownloader/monitor"
	"log"
	"testing"
	"time"
)

func TestPartDownloadWorker(t *testing.T) {
	return
	url := "http://releases.ubuntu.com/14.04.2/ubuntu-14.04.2-server-amd64.list"
	c, _ := httpclient.GetSize(url)
	c = c / 2
	f, _ := iotools.CreateSafeFile("g_ubuntu-14.04.2-server-amd64.list")
	defer f.Close()
	log.Println(f.Truncate(c))
	dow := httpclient.CreatePartialDownloader(url, f, 0, c)
	mv := monitor.MonitoredWorker{Itw: dow}
	log.Println(mv.Start())
	log.Println(mv.Start())
	time.Sleep(time.Second * 1)
	log.Println(mv.Stop())
	time.Sleep(time.Second * 5)
	log.Println(mv.Start())
	log.Println(mv.Start())
	time.Sleep(time.Second * 5)
	log.Println(mv.Stop())
}

func TestMultiPartDownloadWorker(t *testing.T) {
	partcount := 10
	pc := int64(partcount)
	url := "http://releases.ubuntu.com/14.04.2/ubuntu-14.04.2-server-amd64.list"
	c, _ := httpclient.GetSize(url)
	f, _ := iotools.CreateSafeFile("gm_ubuntu-14.04.2-server-amd64.list")
	defer f.Close()
	f.Truncate(c)
	ps := c / pc
	for i := int64(0); i < pc-1; i++ {
		//log.Println(ps*i, ps*i+ps)
		d := httpclient.CreatePartialDownloader(url, f, ps*i, ps*i+ps)
		mv := monitor.MonitoredWorker{Itw: d}
		mv.Start()
	}
	lastseg := c - (ps * (pc - 1))
	dow := httpclient.CreatePartialDownloader(url, f, lastseg, c)
	mv := monitor.MonitoredWorker{Itw: dow}
	mv.Start()

	time.Sleep(time.Second * 15)

}

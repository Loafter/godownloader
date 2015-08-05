package dtest_test

import (
	"godownloader/http"
	"log"
	"testing"
	"time"
)

func TestDownload(t *testing.T) {
	dl, err := httpclient.CreateDownloader("http://ftp.nz.debian.org/debian/dists/Debian8.1/main/Contents-arm64.gz", "g_Contents-source64.gz", 15)
	if err != nil {
		t.Error("failed: can't create downloader")
	}
	err = dl.StartAll()
	if err != nil {
		t.Error("failed: can't start downloader")
	}
	/*time.Sleep(time.Second * 1)
	err = dl.StopAll()
	time.Sleep(time.Second * 1)
	err = dl.StartAll()
	if err != nil {
		t.Error("failed: can't stops downloader")
	}*/
	for {
		time.Sleep(time.Millisecond * 500)
		pr := dl.GetProgress()
		done := true
		for _, r := range pr {
			log.Print((r.To - r.Pos) / 1024)
		}
		log.Println("_________________________________________________________________________")
		for _, i := range pr {
			if i.Pos != i.To {
				done = false
				break
			}
		}

		if done {
			break
		}
	}
}

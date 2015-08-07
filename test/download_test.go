package dtest_test

import (
	"godownloader/http"
	"log"
	"testing"
	"time"
)

func TestDownload(t *testing.T) {
	dl, err := httpclient.CreateDownloader("http://pinegrow.s3.amazonaws.com/PinegrowLinux64.2.2.zip", "xubuntu-12.04.4-desktop-amd64.iso.zsync", 15)
	if err != nil {
		t.Error("failed: can't create downloader")
	}
	err = dl.StartAll()
	if err != nil {
		t.Error("failed: can't start downloader")
	}
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

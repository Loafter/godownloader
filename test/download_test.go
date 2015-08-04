package dtest_test

import (
	"godownloader/http"
	"testing"
	"time"
)

func TestDownload(t *testing.T) {
	dl, err := httpclient.CreateDownloader("http://mirror.yandex.ru/ubuntu-releases/15.04/ubuntu-15.04-snappy-amd64%2bgeneric.img.xz", "Test.dat", 7)
	if err != nil {
		t.Error("failed: can't create downloader")
	}
	err = dl.StartAll()
	if err != nil {
		t.Error("failed: can't start downloader")
	}
	time.Sleep(time.Second * 5)
	err = dl.StopAll()
	time.Sleep(time.Second * 5)
	err = dl.StartAll()
	if err != nil {
		t.Error("failed: can't stops downloader")
	}
	time.Sleep(time.Second * 7)
}

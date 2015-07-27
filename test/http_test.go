package dtest

import (
	"testing"
	"godownloader/http"
	"os"
	"log"
)

func TestMultiThreadSuppurt(t *testing.T) {
	return
	if _, e := httpclient.CheckMultipart("http://s0.cyberciti.org/images/misc/static/2012/11/ifdata-welcome-0.png"); (e!=nil) {
		t.Error("failed: CheckMultipart must be without error",e)
	}

	if _, e := httpclient.CheckMultipart("http://s0.cybedrciti.org/images/misc/static/2012/11/ifdata-welcome-0.png"); (e==nil) {
		t.Error("failed: CheckMultipart must with  beerror")
	}

	if _, e := httpclient.CheckMultipart("http://s0.cyberciti.org/images/misc/static/2012/11/ifdadta-welcome-0.png"); (e==nil) {
		t.Error("failed: CheckMultipart size must be with error")
	}
}
func TestGetSize(t *testing.T) {
	return
	if _, e := httpclient.GetSize("http://static.oper.ru/data/gallery/l1048752856.jpg"); (e!=nil) {
		t.Error("failed: Get size must be without error")
	}
	if _, e := httpclient.GetSize("http://static.oper.ru/data/gallery/l104d8752856.jpg"); (e==nil) {
		t.Error("failed: Get size must be with error")
	}
	if _, e := httpclient.GetSize("http://sdtatic.oper.ru/data/gallery/l1048752856.jpg"); (e==nil) {
		t.Error("failed: Get size must be with error")
	}

}

func TestPartDownload(t *testing.T) {
	c, e := httpclient.GetSize("http://s0.cyberciti.org/images/misc/static/2012/11/ifdata-welcome-0.png")
	if (e!=nil) {
		t.Error("failed: Get size must be without error")
	}
	f,_:=os.Create("part_download.data")
	defer f.Close()
	f.Truncate(c);

	dow:=httpclient.PartialDownloader{}
	dow.Init("http://s0.cyberciti.org/images/misc/static/2012/11/ifdata-welcome-0.png",f,0,0,c)
	for  {
		sta,_:=dow.DoWork()
		if !sta{
			return
		}
		log.Print(dow.GetProgress())
	}

}

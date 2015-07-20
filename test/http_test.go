package dtest

import (
	"testing"
	"godownloader/http"
)

func TestMultiThreadSuppurt(t *testing.T) {
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

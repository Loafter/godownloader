package dtest

import (
	"testing"
	"godownloader/http"
)
func TestGetSize(t *testing.T) {
 if _,e:=httpclient.GetSize("http://img1.goodfon.ru:80/original/1920x1080/3/4a/mlechnyy-put-kosmos-zvezdy-3734.jpg");(e!=nil){
	 t.Error("failed: Get size must be without error")
	}

	if _,e:=httpclient.GetSize("http://img1.goodffon.ru:80/original/1920x1080/3/4a/mlechnyy-put-kosmos-zvezdy-3734.jpg");(e==nil){
		t.Error("failed: Get size must with error")
	}

	if _,e:=httpclient.GetSize("http://img1.goodfon.ru:80/original/1920x1080/3/4a/mlechnyy-put-kdsfosmos-zvezdy-3734.jpg");(e==nil){
		t.Error("failed: Get size must be with error")
	}

}

func TestResume(t *testing.T) {

}
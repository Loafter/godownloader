package httpclient
import (
	"net/http"
	"log"
)
func TestMultithread(url string) (bool, error) {
	return false, nil
}

func GetSize(urls string) (int64, error) {
	cl:=http.Client{}
	resp,err:=cl.Head(urls)
	if err!=nil{
		return 0,err
	}
	log.Printf("info: file size is %d bytes \n",resp.ContentLength)
	return resp.ContentLength,nil
}
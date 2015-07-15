package httpclient
import (
	"net/http"
	"log"
	"errors"
)
func TestMultithread(url string) (bool, error) {
	return false, nil
}

func GetSize(urls string) (int64, error) {
	cl:=http.Client{}
	resp,err:=cl.Head(urls)
	if err!=nil{
		log.Printf("error: when try get file size %v \n",err)
		return 0,err
	}
	if resp.StatusCode!=200{
		log.Printf("error: file not found or moved status:",resp.StatusCode)
		return 0,errors.New("error: file not found or moved")
	}
	log.Printf("info: file size is %d bytes \n",resp.ContentLength)
	return resp.ContentLength,nil
}
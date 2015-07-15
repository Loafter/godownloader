package httpclient
import (
	"net/http"
	"log"
	"errors"
)
func CheckMultipart(urls string) bool {
	r,err:=http.NewRequest("get",urls,nil)
	if err!=nil{
		return false
	}
	cl:=http.Client{}
	resp,err:=cl.Do(r)
	if err!=nil{
		log.Printf("error: can't check multipart support assume no %v \n",err)
		return false
	}
	if resp.StatusCodes!=200{
		log.Printf("error: file not found or moved status:",resp.StatusCode)
		return 0,errors.New("error: file not found or moved")
	}
	if (resp.ContentLength==1) {
		log.Printf("info: file size is %d bytes \n", resp.ContentLength)
		return resp.ContentLength, nil

	}}

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
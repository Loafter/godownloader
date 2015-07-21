package httpclient
import (
	"net/http"
	"log"
	"errors"
	"os"
)
func CheckMultipart(urls string) (bool, error) {
	r, err := http.NewRequest("GET", urls, nil)
	if err!=nil {
		return false, err
	}
	r.Header.Add("Range", "bytes=0-0")
	cl := http.Client{}
	f, _ := os.Create("/home/andrew/Desktop/dum.txt")
	r.Write(f)
	defer f.Close()
	resp, err := cl.Do(r)
	if err!=nil {
		log.Printf("error: can't check multipart support assume no %v \n", err)
		return false, err
	}
	f1, _ := os.Create("/home/andrew/Desktop/res.txt")
	resp.Write(f1)
	if resp.StatusCode!=206 {
		return false, errors.New("error: file not found or moved status: "+ resp.Status)
	}
	if (resp.ContentLength==1) {
		log.Printf("info: file size is %d bytes \n", resp.ContentLength)
		return true, nil
	}
	return false, nil
}

func GetSize(urls string) (int64, error) {
	cl := http.Client{}
	resp, err := cl.Head(urls)
	if err!=nil {
		log.Printf("error: when try get file size %v \n", err)
		return 0, err
	}
	if resp.StatusCode!=200 {
		log.Printf("error: file not found or moved status:", resp.StatusCode)
		return 0, errors.New("error: file not found or moved")
	}
	log.Printf("info: file size is %d bytes \n", resp.ContentLength)
	return resp.ContentLength, nil
}

type PartialDownloader struct {
	from int64
	to   int64
	pos  int64
	rch  bool
	url  string
}
func (pd *PartialDownloader) Init(url string, from int64, pos int64, to int64) {
	pd.from=from
	pd.to=to
	pd.pos=pos
}
func (pd PartialDownloader) GetProgress() interface{} {
	return pd.from//, pd.to, pd.to
}
func (pd *PartialDownloader) DoWork() (bool, error) {
	//in last time we check resume support
	if !pd.rch {
		if nos, err := CheckMultipart(pd.url); nos {
			return false, err
		}
	}
	//assume resume support
	pd.rch=true
	//do download

	return true, nil
}

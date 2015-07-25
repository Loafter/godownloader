package httpclient
import (
	"net/http"
	"log"
	"errors"
	"os"
)
const MaxDownloadPortion=4096
func CheckMultipart(urls string) (bool, error) {
	r, err := http.NewRequest("GET", urls, nil)
	if err!=nil {
		return false, err
	}
	r.Header.Add("Range", "bytes=0-0")
	cl := http.Client{}
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

type DownloadProgress struct{
	from int64
	to   int64
	pos  int64
}
type PartialDownloader struct {
	dp DownloadProgress
	cli *http.Client
	rch  bool
	url  string
}
func (pd *PartialDownloader) Init(url string, from int64, pos int64, to int64) {
	pd.dp.from=from
	pd.dp.to=to
	pd.dp.pos=pos
}

func (pd PartialDownloader) GetProgress()DownloadProgress{
	return pd.dp
}
func constuctReqH(current int64,to int64)string{
	if to<current+MaxDownloadPortion{
		return "bytes="+current+"-"+to
	}

	return "bytes="+current+"-"+MaxDownloadPortion

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

	//check if our client is not created
	if pd.cli==nil{
		pd.cli=new(http.Client)
	}
	//create new req
	r, err := http.NewRequest("GET", pd.url, nil)
	//ok we construct query
	r.Header.Add("Range", constuctReqH(pd.dp.pos,pd.dp.to))
	if err!=nil {
		return false, err
	}
	//try send
	resp, err := pd.cli.Do(r)
	if err!=nil {
		log.Printf("error: error download part  file%v \n", err)
		return false, err
	}
	//check response
	if resp.StatusCode!=200 {
		log.Printf("error: file not found or moved status:", resp.StatusCode)
		return false, errors.New("error: file not found or moved")
	}
	pd.dp.pos+pd.dp.pos+MaxDownloadPortion
	if(pd.dp.pos==pd.dp.to){
		//ok download part complete normal
		return false,nil
	}
	//not full download try next segment
	return true, nil
}

package httpclient
import (
	"net/http"
	"log"
	"errors"
	"strconv"
	"os"
	"io/ioutil"
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
	client *http.Client
	isresume bool
	url  string
	file *os.File
}
func (pd *PartialDownloader) Init(url string,file  *os.File ,from int64, pos int64, to int64) {
	pd.file=file
	pd.url=url
	pd.dp.from=from
	pd.dp.to=to
	pd.dp.pos=pos
}

func (pd PartialDownloader) GetProgress()DownloadProgress{
	return pd.dp
}
func constuctReqH(current int64,to int64)string{
	if to<current+MaxDownloadPortion{
		return "bytes="+strconv.FormatInt(current,10)+"-"+strconv.FormatInt(to,10)
	}

	return "bytes="+ strconv.FormatInt(current,10)+"-"+strconv.FormatInt(current+MaxDownloadPortion,10)

}
func(pd *PartialDownloader)DownloadSergment()(bool,error){
	//in last time we check resume support
	if !pd.isresume {
		if nos, err := CheckMultipart(pd.url); !nos {
			return false, err
		}
	}
	//assume resume support
	pd.isresume=true
	//do download

	//check if our client is not created
	if pd.client==nil{
		pd.client=new(http.Client)
	}
	//create new req
	r, err := http.NewRequest("GET", pd.url, nil)
	//ok we construct query
	r.Header.Add("Range", constuctReqH(pd.dp.pos,pd.dp.to))
	a,_:=os.Create("req.txt")
	r.Write(a)
	a.Close()
	if err!=nil {
		return false, err
	}
	//try send request
	resp, err := pd.client.Do(r)
	if err!=nil {
		log.Printf("error: error download part file%v \n", err)
		return false, err
	}
	//check response
	if resp.StatusCode!=206 {
		log.Printf("error: file not found or moved status:", resp.StatusCode)
		return false, errors.New("error: file not found or moved")
	}
	//write flush data to disk
	dat,err:=ioutil.ReadAll(resp.Body)
	if (err!=nil){
		return false,err
	}
	c,err:=pd.file.WriteAt(dat,pd.dp.pos)
	if err!=nil{
		return false,err
	}
	pd.file.Sync()
	pd.dp.pos=pd.dp.pos+int64(c)
	if(pd.dp.pos==pd.dp.to){
		//ok download part complete normal
		return false,nil
	}
	//not full download next segment
	return true, nil
}
func (pd *PartialDownloader) DoWork() (bool, error) {
 return  pd.DownloadSergment()
}

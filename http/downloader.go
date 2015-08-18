package httpclient

import (
	"godownloader/iotools"
	"godownloader/monitor"
	"os"
	"strconv"
	"os/user"
)

type FileInfo struct {
	Size     int64  `json:"Size"`
	FileName string `json:"FileName"`
	Url      string `json:"Url"`
}
type Downloader struct {
	sf *iotools.SafeFile
	wp *monitor.WorkerPool
	Fi FileInfo
}

func (dl *Downloader) StopAll() []error {
	return dl.wp.StopAll()
}
func (dl *Downloader) StartAll() []error {
	return dl.wp.StartAll()
}
func (dl *Downloader) GetProgress() []DownloadProgress {
	pr := dl.wp.GetAllProgress().([]interface{})
	re := make([]DownloadProgress, len(pr))
	for i, val := range pr {
		re[i] = val.(DownloadProgress)
	}
	return re
}
func getDown()(string){
	usr, _ := user.Current()
	st:=strconv.QuoteRune(os.PathSeparator)
	st=st[1:len(st)-1]
	return usr.HomeDir+st+"Downloads"+st
}
func CreateDownloader(url string, fp string, seg int64) (dl *Downloader, err error) {
	c, err := GetSize(url)
	if err != nil {
		//can't get file size
		return nil, err
	}

	dfs:=getDown()+fp
	sf, err := iotools.CreateSafeFile(dfs)
	if err != nil {
		//can't create file on path
		return nil, err
	}

	if err := sf.Truncate(c); err != nil {
		//can't truncate file
		return nil, err
	}
	//create part-downloader foreach segment
	ps := c / seg
	wp := new(monitor.WorkerPool)
	for i := int64(0); i < seg-int64(1); i++ {
		d := CreatePartialDownloader(url, sf, ps*i, ps*i, ps*i+ps)
		mv := monitor.MonitoredWorker{Itw: d}
		wp.AppendWork(&mv)
	}
	lastseg := int64(ps * (seg - 1))
	dow := CreatePartialDownloader(url, sf, lastseg, lastseg, c)
	mv := monitor.MonitoredWorker{Itw: dow}

	//add to worker pool
	wp.AppendWork(&mv)
	d := Downloader{
		sf: sf,
		wp: wp,
		Fi: FileInfo{FileName: fp, Size: c, Url: url},
	}
	return &d, nil
}



func RestoreDownloader(url string, fp string, dp []DownloadProgress) (dl *Downloader, err error) {
	c, err := GetSize(url)
	if err != nil {
		//can't get file size
		return nil, err
	}
	dfs:=getDown()+fp
	sf, err := iotools.OpenSafeFile(dfs)
	if err != nil {
		//can't create file on path
		return nil, err
	}
	wp := new(monitor.WorkerPool)
	for _, r := range dp {
		dow := CreatePartialDownloader(url, sf, r.From, r.Pos, r.To)
		mv := monitor.MonitoredWorker{Itw: dow}

		//add to worker pool
		wp.AppendWork(&mv)

	}
	d := Downloader{
		sf: sf,
		wp: wp,
		Fi: FileInfo{FileName: fp, Size: c, Url: url},
	}
	return &d, nil
}

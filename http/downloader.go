package httpclient

import (
	"godownloader/iotools"
	"godownloader/monitor"
)

type Downloader struct {
	sf   *iotools.SafeFile
	size int64
	wp   *monitor.WorkerPool
}

func (dl *Downloader) StopAll() error {
	return dl.wp.StopAll()
}
func (dl *Downloader) StartAll() error {
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
func CreateDownloader(url string, fp string, seg int64) (dl *Downloader, err error) {
	c, err := GetSize(url)
	if err != nil {
		//can't get file size
		return nil, err
	}

	sf, err := iotools.CreateSafeFile(fp)
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
		d := CreatePartialDownloader(url, sf, ps*i, ps*i+ps)
		mv := monitor.MonitoredWorker{Itw: d}
		wp.AppendWork(&mv)
	}
	lastseg := int64(ps * (seg - 1))
	dow := CreatePartialDownloader(url, sf, lastseg, c)
	mv := monitor.MonitoredWorker{Itw: dow}

	//add to worker pool
	wp.AppendWork(&mv)
	d := Downloader{
		sf:   sf,
		size: c,
		wp:   wp,
	}
	return &d, nil
}

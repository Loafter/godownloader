package httpclient

import (
	"godownloader/iotools"
	//"godownloader/http"
)

type Downloader struct {
	sf   *iotools.SafeFile
	size int64
}

/*func CreaOteDownloader(url string,fp string,seg int) (dl *Downloader, err error) {
	c,err:=httpclient.GetSize(url)
	if err!=nil{
		return err
	}
	sf,err:= iotools.CreateSafeFile(fp)
	if err!=nil{
		return err
	}

	d:=Downloader{
		sf:sf,
		size:c,
	}
	return &d
}*/

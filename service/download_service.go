package DownloadService

import (
	"encoding/base64"
	"encoding/json"
	"godownloader/http"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type DJob struct {
	Id         int
	FileName   string
	Size       int64
	Downloaded int64
	Progress   int64
	Speed      int64
}

type NewJob struct {
	Url       string
	PartCount int64
	FilePath  string
}

type DServ struct {
	dls    []*httpclient.Downloader
	oplock sync.Mutex
}

func (srv *DServ) Start(listenPort int) error {
	http.HandleFunc("/progress.json", srv.progressJson)
	http.HandleFunc("/add_task", srv.addTask)
	http.HandleFunc("/remove_task", srv.removeTask)
	http.HandleFunc("/start_task", srv.startTask)
	http.HandleFunc("/stop_task", srv.stopTask)
	http.HandleFunc("/start_all_task", srv.startAllTask)
	http.HandleFunc("/stop_all_task", srv.stopAllTask)
	http.HandleFunc("/index.html", srv.index)
	if err := http.ListenAndServe(":"+strconv.Itoa(listenPort), nil); err != nil {
		return err
	}
	return nil
}

func (srv *DServ) SaveSettings(sf string) error {
	var ss ServiceSettings
	for _, i := range srv.dls {

		ss.Ds = append(ss.Ds, DownloadSettings{
			FI: i.Fi,
			Dp: i.GetProgress(),
		})
	}

	return ss.SaveToFile(sf)
}

func (srv *DServ) LoadSettings(sf string) error {
	ss, err := LoadFromFile(sf)
	if err != nil {
		log.Println("error: when try load settings", err)
		return err
	}
	log.Println(ss)
	for _, r := range ss.Ds {
		dl, err := httpclient.RestoreDownloader(r.FI.Url, r.FI.FileName, r.Dp)
		if err != nil {
			return err
		}
		srv.dls = append(srv.dls, dl)
	}
	return nil
}

const htmlData = "PCFkb2N0eXBlIGh0bWw+DQoNCjxodG1sPg0KDQo8aGVhZD4NCgk8dGl0bGU+R08gRE9XTkxPQUQ8L3RpdGxlPg0KCTxtZXRhIG5hbWU9InZpZXdwb3J0IiBjb250ZW50PSJ3aWR0aD1kZXZpY2Utd2lkdGgiPg0KCTxsaW5rIHJlbD0ic3R5bGVzaGVldCIgaHJlZj0iaHR0cHM6Ly9uZXRkbmEuYm9vdHN0cmFwY2RuLmNvbS9ib290c3dhdGNoLzMuMC4wL2pvdXJuYWwvYm9vdHN0cmFwLm1pbi5jc3MiPg0KCTxsaW5rIHJlbD0ic3R5bGVzaGVldCIgdHlwZT0idGV4dC9jc3MiIG1lZGlhPSJzY3JlZW4iDQoJCSAgaHJlZj0iaHR0cDovL3d3dy5ndXJpZGRvLm5ldC9kZW1vL2Nzcy90cmlyYW5kL3VpLmpxZ3JpZC1ib290c3RyYXAuY3NzIj4NCgk8c2NyaXB0IHR5cGU9InRleHQvamF2YXNjcmlwdCIgc3JjPSJodHRwczovL2FqYXguZ29vZ2xlYXBpcy5jb20vYWpheC9saWJzL2pxdWVyeS8yLjAuMy9qcXVlcnkubWluLmpzIj48L3NjcmlwdD4NCgk8c2NyaXB0IHR5cGU9InRleHQvamF2YXNjcmlwdCIgc3JjPSJodHRwczovL25ldGRuYS5ib290c3RyYXBjZG4uY29tL2Jvb3RzdHJhcC8zLjMuNC9qcy9ib290c3RyYXAubWluLmpzIj48L3NjcmlwdD4NCgk8c2NyaXB0IHR5cGU9InRleHQvamF2YXNjcmlwdCIgc3JjPSJodHRwOi8vd3d3Lmd1cmlkZG8ubmV0L2RlbW8vanMvdHJpcmFuZC9qcXVlcnkuanFHcmlkLm1pbi5qcyI+PC9zY3JpcHQ+DQoJPHNjcmlwdCB0eXBlPSJ0ZXh0L2phdmFzY3JpcHQiIHNyYz0iaHR0cDovL3d3dy5ndXJpZGRvLm5ldC9kZW1vL2pzL3RyaXJhbmQvaTE4bi9ncmlkLmxvY2FsZS1lbi5qcyI+PC9zY3JpcHQ+DQoJPGxpbmsgcmVsPSJzdHlsZXNoZWV0IiBocmVmPSIvL2NvZGUuanF1ZXJ5LmNvbS91aS8xLjExLjQvdGhlbWVzL3Ntb290aG5lc3MvanF1ZXJ5LXVpLmNzcyI+DQoJPHNjcmlwdCBzcmM9Imh0dHA6Ly9jb2RlLmpxdWVyeS5jb20vdWkvMS4xMS40L2pxdWVyeS11aS5qcyI+PC9zY3JpcHQ+DQoJPHN0eWxlIHR5cGU9InRleHQvY3NzIj4NCgkJYm9keSB7DQoJCXBhZGRpbmctdG9wOiA1MHB4Ow0KCQlwYWRkaW5nLWJvdHRvbTogMjBweDsNCgkJfQ0KDQoJCS50YWJsZSAucHJvZ3Jlc3Mgew0KCQltYXJnaW4tYm90dG9tOiAwcHg7DQoJCX0NCgk8L3N0eWxlPg0KCTxzY3JpcHQgdHlwZT0idGV4dC9qYXZhc2NyaXB0Ij4NCgkJZnVuY3Rpb24gVXBkYXRlVGFibGUoKSB7DQoJCQkkKCIjanFHcmlkIikNCgkJCQkuanFHcmlkKHsNCgkJCQkJdXJsOiAnaHR0cDovL2xvY2FsaG9zdDo5OTgxL3Byb2dyZXNzLmpzb24nLA0KCQkJCQltdHlwZTogIkdFVCIsDQoJCQkJCWFqYXhTdWJncmlkT3B0aW9uczogew0KCQkJCQkJYXN5bmM6IGZhbHNlDQoJCQkJCX0sDQoJCQkJCXN0eWxlVUk6ICdCb290c3RyYXAnLA0KCQkJCQlkYXRhdHlwZTogImpzb24iLA0KCQkJCQljb2xNb2RlbDogW3sNCgkJCQkJCWxhYmVsOiAnIycsDQoJCQkJCQluYW1lOiAnSWQnLA0KCQkJCQkJa2V5OiB0cnVlLA0KCQkJCQkJd2lkdGg6IDUNCgkJCQkJfSwgew0KCQkJCQkJbGFiZWw6ICdGaWxlIE5hbWUnLA0KCQkJCQkJbmFtZTogJ0ZpbGVOYW1lJywNCgkJCQkJCXdpZHRoOiAxNQ0KCQkJCQl9LCB7DQoJCQkJCQlsYWJlbDogJ1NpemUnLA0KCQkJCQkJbmFtZTogJ1NpemUnLA0KCQkJCQkJd2lkdGg6IDIwLA0KCQkJCQkJZm9ybWF0dGVyOiBGb3JtYXRCeXRlDQoJCQkJCX0sIHsNCgkJCQkJCWxhYmVsOiAnRG93bmxvYWRlZCcsDQoJCQkJCQluYW1lOiAnRG93bmxvYWRlZCcsDQoJCQkJCQl3aWR0aDogMjAsDQoJCQkJCQlmb3JtYXR0ZXI6IEZvcm1hdEJ5dGUNCgkJCQkJfSwgew0KCQkJCQkJbGFiZWw6ICclJywNCgkJCQkJCW5hbWU6ICdQcm9ncmVzcycsDQoJCQkJCQl3aWR0aDogNQ0KCQkJCQl9LCB7DQoJCQkJCQlsYWJlbDogJ1NwZWVkJywNCgkJCQkJCW5hbWU6ICdTcGVlZCcsDQoJCQkJCQl3aWR0aDogMTUsDQoJCQkJCQlmb3JtYXR0ZXI6IEZvcm1hdEJ5dGUNCgkJCQkJfSwgew0KCQkJCQkJbGFiZWw6ICdQcm9ncmVzcycsDQoJCQkJCQluYW1lOiAnUHJvZ3Jlc3MnLA0KCQkJCQkJZm9ybWF0dGVyOiBGb3JtYXRQcm9ncmVzc0Jhcg0KCQkJCQl9XSwNCgkJCQkJdmlld3JlY29yZHM6IHRydWUsDQoJCQkJCXJvd051bTogMjAsDQoJCQkJCXBhZ2VyOiAiI2pxR3JpZFBhZ2VyIg0KCQkJCX0pOw0KCQl9DQoNCgkJZnVuY3Rpb24gRml4VGFibGUoKSB7DQoJCQkkLmV4dGVuZCgkLmpncmlkLmFqYXhPcHRpb25zLCB7DQoJCQkJYXN5bmM6IGZhbHNlDQoJCQl9KQ0KCQkJJCgiI2pxR3JpZCIpDQoJCQkJLnNldEdyaWRXaWR0aCgkKHdpbmRvdykNCgkJCQkJLndpZHRoKCkgLSA1KQ0KCQkJJCgiI2pxR3JpZCIpDQoJCQkJLnNldEdyaWRIZWlnaHQoJCh3aW5kb3cpDQoJCQkJCS5oZWlnaHQoKSkNCgkJCSQod2luZG93KQ0KCQkJCS5iaW5kKCdyZXNpemUnLCBmdW5jdGlvbigpIHsNCgkJCQkJJCgiI2pxR3JpZCIpDQoJCQkJCQkuc2V0R3JpZFdpZHRoKCQod2luZG93KQ0KCQkJCQkJCS53aWR0aCgpIC0gNSk7DQoJCQkJCSQoIiNqcUdyaWQiKQ0KCQkJCQkJLnNldEdyaWRIZWlnaHQoJCh3aW5kb3cpDQoJCQkJCQkJLmhlaWdodCgpKQ0KCQkJCX0pDQoJCX0NCg0KCQlmdW5jdGlvbiBVcGRhdGVEYXRhKCkgew0KCQkJdmFyIGdyaWQgPSAkKCIjanFHcmlkIik7DQoJCQl2YXIgcm93S2V5ID0gZ3JpZC5qcUdyaWQoJ2dldEdyaWRQYXJhbScsICJzZWxyb3ciKTsNCgkJCSQoIiNqcUdyaWQiKS50cmlnZ2VyKCJyZWxvYWRHcmlkIik7DQoJCQlpZihyb3dLZXkpIHsNCgkJCQkkKCcjanFHcmlkJykuanFHcmlkKCJyZXNldFNlbGVjdGlvbiIpDQoJCQkJJCgnI2pxR3JpZCcpLmpxR3JpZCgnc2V0U2VsZWN0aW9uJywgcm93S2V5KTsNCgkJCX0NCgkJfQ0KDQoJCWZ1bmN0aW9uIEZvcm1hdFByb2dyZXNzQmFyKGNlbGxWYWx1ZSwgb3B0aW9ucywgcm93T2JqZWN0KSB7DQoJCQl2YXIgaW50VmFsID0gcGFyc2VJbnQoY2VsbFZhbHVlKTsNCg0KCQkJdmFyIGNlbGxIdG1sID0gJzxkaXYgY2xhc3M9InByb2dyZXNzIj48ZGl2IGNsYXNzPSJwcm9ncmVzcy1iYXIiIHN0eWxlPSJ3aWR0aDogJyArIGludFZhbCArICclOyI+PC9kaXY+PC9kaXY+Jw0KDQoJCQlyZXR1cm4gY2VsbEh0bWw7DQoJCX0NCg0KCQlmdW5jdGlvbiBGb3JtYXRCeXRlKGNlbGxWYWx1ZSwgb3B0aW9ucywgcm93T2JqZWN0KSB7DQoJCQl2YXIgaW50VmFsID0gcGFyc2VJbnQoY2VsbFZhbHVlKTsNCgkJCXZhciByYXMgPSAiIGIuIg0KCQkJaWYoaW50VmFsID4gMTAyNCkgew0KCQkJCWludFZhbCAvPSAxMDI0DQoJCQkJcmFzID0gIiBrYi4iDQoJCQl9DQoJCQlpZihpbnRWYWwgPiAxMDI0KSB7DQoJCQkJaW50VmFsIC89IDEwMjQNCgkJCQlyYXMgPSAiIE1iLiINCgkJCX0NCgkJCWlmKGludFZhbCA+IDEwMjQpIHsNCgkJCQlpbnRWYWwgLz0gMTAyNA0KCQkJCXJhcyA9ICIgR2IuIg0KCQkJfQ0KDQoJCQlpZihpbnRWYWwgPiAxMDI0KSB7DQoJCQkJaW50VmFsIC89IDEwMjQNCgkJCQlyYXMgPSAiIFRiLiINCgkJCX0NCgkJCXZhciBjZWxsSHRtbCA9IChpbnRWYWwpLnRvRml4ZWQoMSkgKyByYXM7DQoJCQlyZXR1cm4gY2VsbEh0bWw7DQoJCX0NCg0KCQlmdW5jdGlvbiBPbkxvYWQoKSB7DQoNCgkJCVVwZGF0ZVRhYmxlKCkNCgkJCUZpeFRhYmxlKCkNCgkJCXNldEludGVydmFsKFVwZGF0ZURhdGEsIDUwMCk7DQoJCX0NCg0KCQlmdW5jdGlvbiBBZGREb3dubG9hZCgpIHsNCgkJCXZhciByZXEgPSB7DQoJCQkJUGFydENvdW50OiBwYXJzZUludCgkKCIjcGFydF9jb3VudF9pZCIpLnZhbCgpKSwNCgkJCQlGaWxlUGF0aDogJCgiI3NhdmVfcGF0aF9pZCIpLnZhbCgpLA0KCQkJCVVybDogJCgiI3VybF9pZCIpLnZhbCgpDQoJCQl9Ow0KCQkJJC5hamF4KHsNCgkJCQkJdXJsOiAiL2FkZF90YXNrIiwNCgkJCQkJdHlwZTogIlBPU1QiLA0KCQkJCQlkYXRhOiBKU09OLnN0cmluZ2lmeShyZXEpLA0KCQkJCQlkYXRhVHlwZTogInRleHQiDQoJCQkJfSkNCgkJCQkuZXJyb3IoZnVuY3Rpb24oanNvbkRhdGEpIHsNCgkJCQkJY29uc29sZS5sb2coanNvbkRhdGEpDQoJCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIFJlbW92ZURvd25sb2FkKCkgew0KCQkJdmFyIGdyaWQgPSAkKCIjanFHcmlkIik7DQoJCQl2YXIgcm93S2V5ID0gcGFyc2VJbnQoZ3JpZC5qcUdyaWQoJ2dldEdyaWRQYXJhbScsICJzZWxyb3ciKSk7DQoJCQl2YXIgcmVxID0gcm93S2V5Ow0KCQkJJC5hamF4KHsNCgkJCQkJdXJsOiAiL3JlbW92ZV90YXNrIiwNCgkJCQkJdHlwZTogIlBPU1QiLA0KCQkJCQlkYXRhOiBKU09OLnN0cmluZ2lmeShyZXEpLA0KCQkJCQlkYXRhVHlwZTogInRleHQiDQoJCQkJfSkNCgkJCQkuZXJyb3IoZnVuY3Rpb24oanNvbkRhdGEpIHsNCgkJCQkJY29uc29sZS5sb2coanNvbkRhdGEpDQoJCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIFN0YXJ0RG93bmxvYWQoKSB7DQoJCQl2YXIgZ3JpZCA9ICQoIiNqcUdyaWQiKTsNCgkJCXZhciByb3dLZXkgPSBwYXJzZUludChncmlkLmpxR3JpZCgnZ2V0R3JpZFBhcmFtJywgInNlbHJvdyIpKTsNCgkJCXZhciByZXEgPSByb3dLZXk7DQoJCQkkLmFqYXgoew0KCQkJCQl1cmw6ICIvc3RhcnRfdGFzayIsDQoJCQkJCXR5cGU6ICJQT1NUIiwNCgkJCQkJZGF0YTogSlNPTi5zdHJpbmdpZnkocmVxKSwNCgkJCQkJZGF0YVR5cGU6ICJ0ZXh0Ig0KCQkJCX0pDQoJCQkJLmVycm9yKGZ1bmN0aW9uKGpzb25EYXRhKSB7DQoJCQkJCWNvbnNvbGUubG9nKGpzb25EYXRhKQ0KCQkJCX0pDQoJCX0NCg0KCQlmdW5jdGlvbiBTdGFydEFsbERvd25sb2FkKCkgew0KCQkJJC5hamF4KHsNCgkJCQkJdXJsOiAiL3N0YXJ0X2FsbF90YXNrIiwNCgkJCQkJdHlwZTogIlBPU1QiLA0KCQkJCQlkYXRhVHlwZTogInRleHQiDQoJCQkJfSkNCgkJCQkuZXJyb3IoZnVuY3Rpb24oanNvbkRhdGEpIHsNCgkJCQkJY29uc29sZS5sb2coanNvbkRhdGEpDQoJCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIFN0b3BBbGxEb3dubG9hZCgpIHsNCgkJCSQuYWpheCh7DQoJCQkJCXVybDogIi9zdG9wX2FsbF90YXNrIiwNCgkJCQkJdHlwZTogIlBPU1QiLA0KCQkJCQlkYXRhVHlwZTogInRleHQiDQoJCQkJfSkNCgkJCQkuZXJyb3IoZnVuY3Rpb24oanNvbkRhdGEpIHsNCgkJCQkJY29uc29sZS5sb2coanNvbkRhdGEpDQoJCQkJfSkNCgkJfQ0KDQoNCgkJZnVuY3Rpb24gT25DaGFuZ2VVcmwoKSB7DQoJCQl2YXIgZmlsZW5hbWUgPSAkKCIjdXJsX2lkIikudmFsKCkuc3BsaXQoJy8nKS5wb3AoKQ0KCQkJJCgiI3NhdmVfcGF0aF9pZCIpLnZhbChmaWxlbmFtZSkNCgkJfQ0KCTwvc2NyaXB0Pg0KPC9oZWFkPg0KDQo8Ym9keSBvbmxvYWQ9Ik9uTG9hZCgpIj4NCjxkaXYgY2xhc3M9Im5hdmJhciBuYXZiYXItaW52ZXJzZSBuYXZiYXItZml4ZWQtdG9wIj4NCgk8ZGl2IGNsYXNzPSJjb250YWluZXIiPg0KCQk8ZGl2IGNsYXNzPSJuYXZiYXItaGVhZGVyIj4NCgkJCTxidXR0b24gdHlwZT0iYnV0dG9uIiBjbGFzcz0ibmF2YmFyLXRvZ2dsZSIgZGF0YS10b2dnbGU9ImNvbGxhcHNlIiBkYXRhLXRhcmdldD0iLm5hdmJhci1jb2xsYXBzZSI+DQoJCQkJPHNwYW4gY2xhc3M9Imljb24tYmFyIj48L3NwYW4+PHNwYW4gY2xhc3M9Imljb24tYmFyIj48L3NwYW4+PHNwYW4gY2xhc3M9Imljb24tYmFyIj48L3NwYW4+DQoJCQk8L2J1dHRvbj4NCgkJCTxhIGNsYXNzPSJuYXZiYXItYnJhbmQiIGhyZWY9IiMiPkdPIERvd25sb2FkZXI8L2E+DQoJCTwvZGl2Pg0KCQk8ZGl2IGNsYXNzPSJuYXZiYXItY29sbGFwc2UgY29sbGFwc2UiPg0KCQkJPHVsIGNsYXNzPSJuYXYgbmF2YmFyLW5hdiI+DQoJCQkJPGxpIGNsYXNzPSJkcm9wZG93biI+DQoJCQkJCTxhIGhyZWY9IiMiIGNsYXNzPSJkcm9wZG93bi10b2dnbGUiIGRhdGEtdG9nZ2xlPSJkcm9wZG93biI+RmlsZSA8YiBjbGFzcz0iY2FyZXQiPjwvYj48L2E+DQoJCQkJCTx1bCBjbGFzcz0iZHJvcGRvd24tbWVudSI+DQoJCQkJCQk8bGk+DQoJCQkJCQkJPGEgZGF0YS10b2dnbGU9Im1vZGFsIiBkYXRhLXRhcmdldD0iI215TW9kYWwiPkFkZCBkb3dubG9hZDwvYT4NCgkJCQkJCTwvbGk+DQoJCQkJCQk8bGkgb25jbGljaz0iUmVtb3ZlRG93bmxvYWQoKSI+DQoJCQkJCQkJPGEgaHJlZj0iIyI+RGVsZXRlIGRvd25sb2FkPC9hPg0KCQkJCQkJPC9saT4NCgkJCQkJPC91bD4NCgkJCQk8L2xpPg0KCQkJCTxsaSBjbGFzcz0iZHJvcGRvd24iPg0KCQkJCQk8YSBocmVmPSIjIiBjbGFzcz0iZHJvcGRvd24tdG9nZ2xlIiBkYXRhLXRvZ2dsZT0iZHJvcGRvd24iPkFjdGlvbiA8YiBjbGFzcz0iY2FyZXQiPjwvYj48L2E+DQoJCQkJCTx1bCBjbGFzcz0iZHJvcGRvd24tbWVudSI+DQoJCQkJCQk8bGkgb25jbGljaz0iU3RhcnREb3dubG9hZCgpIj4NCgkJCQkJCQk8YSBocmVmPSIjIj5TdGFydDwvYT4NCgkJCQkJCTwvbGk+DQoJCQkJCQk8bGkgb25jbGljaz0iU3RvcERvd25sb2FkKCkiPg0KCQkJCQkJCTxhIGhyZWY9IiMiPlN0b3A8L2E+DQoJCQkJCQk8L2xpPg0KCQkJCQkJPGxpIGNsYXNzPSJkaXZpZGVyIj48L2xpPg0KCQkJCQkJPGxpIG9uY2xpY2s9IlN0YXJ0QWxsRG93bmxvYWQoKSI+DQoJCQkJCQkJPGEgaHJlZj0iIyI+U3RhcnQgYWxsPC9hPg0KCQkJCQkJPC9saT4NCgkJCQkJCTxsaSBvbmNsaWNrPSJTdG9wQWxsRG93bmxvYWQoKSI+DQoJCQkJCQkJPGEgaHJlZj0iIyI+U3RvcCBhbGw8L2E+DQoJCQkJCQk8L2xpPg0KCQkJCQk8L3VsPg0KCQkJCTwvbGk+DQoJCQkJPGxpPg0KCQkJCQk8YSBocmVmPSIjYWJvdXQiPkFib3V0PC9hPg0KCQkJCTwvbGk+DQoJCQk8L3VsPg0KCQk8L2Rpdj4NCgkJPCEtLS8ubmF2YmFyLWNvbGxhcHNlIC0tPg0KCTwvZGl2Pg0KPC9kaXY+DQo8L1A+DQo8dGFibGUgaWQ9ImpxR3JpZCI+PC90YWJsZT4NCg0KPCEtLSBNb2RhbCAtLT4NCjxkaXYgY2xhc3M9Im1vZGFsIGZhZGUiIGlkPSJteU1vZGFsIiByb2xlPSJkaWFsb2ciPg0KCTxkaXYgY2xhc3M9Im1vZGFsLWRpYWxvZyI+DQoNCgkJPCEtLSBNb2RhbCBjb250ZW50LS0+DQoJCTxkaXYgY2xhc3M9Im1vZGFsLWNvbnRlbnQiPg0KCQkJPGRpdiBjbGFzcz0ibW9kYWwtaGVhZGVyIj4NCgkJCQk8YnV0dG9uIHR5cGU9ImJ1dHRvbiIgY2xhc3M9ImNsb3NlIiBkYXRhLWRpc21pc3M9Im1vZGFsIj4mdGltZXM7PC9idXR0b24+DQoJCQkJPGg0IGNsYXNzPSJtb2RhbC10aXRsZSI+RW50ZXIgVXJsPC9oND4NCgkJCTwvZGl2Pg0KCQkJPGRpdiBjbGFzcz0ibW9kYWwtYm9keSI+DQoJCQkJPGRpdiBjbGFzcz0iZm9ybS1ncm91cCI+DQoJCQkJCTxsYWJlbCBjbGFzcz0iY29udHJvbC1sYWJlbCI+VXJsPC9sYWJlbD4NCg0KCQkJCQk8ZGl2IGNsYXNzPSJjb250cm9scyI+DQoJCQkJCQk8aW5wdXQgdHlwZT0idGV4dCIgb25jaGFuZ2U9Ik9uQ2hhbmdlVXJsKCkiIGlkPSJ1cmxfaWQiIGNsYXNzPSJmb3JtLWNvbnRyb2wiDQoJCQkJCQkJICAgdmFsdWU9Imh0dHA6Ly9taXJyb3IueWFuZGV4LnJ1L3VidW50dS1jZGltYWdlL3JlbGVhc2VzLzE1LjEwL2FscGhhLTIvc291cmNlL3dpbHktc3JjLTEuaXNvIj4NCgkJCQkJPC9kaXY+DQoJCQkJCTxsYWJlbCBjbGFzcz0iY29udHJvbC1sYWJlbCI+U2F2ZSBwYXRoPC9sYWJlbD4NCg0KCQkJCQk8ZGl2IGNsYXNzPSJjb250cm9scyI+DQoJCQkJCQk8aW5wdXQgdHlwZT0idGV4dCIgaWQ9InNhdmVfcGF0aF9pZCIgY2xhc3M9ImZvcm0tY29udHJvbCINCgkJCQkJCQkgICB2YWx1ZT0id2lseS1zcmMtMS5pc28iPg0KCQkJCQk8L2Rpdj4NCgkJCQkJPGxhYmVsIGNsYXNzPSJjb250cm9sLWxhYmVsIj5QYXJ0cyBjb3VudDwvbGFiZWw+DQoJCQkJCTxzZWxlY3QgY2xhc3M9ImZvcm0tY29udHJvbCIgaWQ9InBhcnRfY291bnRfaWQiPg0KCQkJCQkJPG9wdGlvbj4xPC9vcHRpb24+DQoJCQkJCQk8b3B0aW9uPjI8L29wdGlvbj4NCgkJCQkJCTxvcHRpb24+NDwvb3B0aW9uPg0KCQkJCQkJPG9wdGlvbj44PC9vcHRpb24+DQoJCQkJCQk8b3B0aW9uPjE2PC9vcHRpb24+DQoJCQkJCTwvc2VsZWN0Pg0KDQoJCQkJCTxkaXYgY2xhc3M9Im1vZGFsLWZvb3RlciI+DQoJCQkJCQk8YSBjbGFzcz0iYnRuIGJ0bi1wcmltYXJ5IiBvbmNsaWNrPSJBZGREb3dubG9hZCgpIiBkYXRhLWRpc21pc3M9Im1vZGFsIj5TdGFydCBkb3dubG9hZDwvYT4NCgkJCQkJPC9kaXY+DQoJCQkJPC9kaXY+DQoJCQk8L2Rpdj4NCgkJPC9kaXY+DQoNCgk8L2Rpdj4NCjwvZGl2Pg0KPC9ib2R5Pg0KDQo8L2h0bWw+"

func (srv *DServ) index(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type: text/html", "*")
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Println("warning: start page not found, return included page")
		val, _ := base64.StdEncoding.DecodeString(htmlData)
		rwr.Write(val)
		return
	}
	rwr.Write(content)
}

func (srv *DServ) addTask(rwr http.ResponseWriter, req *http.Request) {
	srv.oplock.Lock()
	defer func() {
		srv.oplock.Unlock()
		req.Body.Close()
	}()
	bodyData, err := ioutil.ReadAll(req.Body)
	rwr.Header().Set("Access-Control-Allow-Origin", "*")
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	var nj NewJob
	if err := json.Unmarshal(bodyData, &nj); err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	dl, err := httpclient.CreateDownloader(nj.Url, nj.FilePath, nj.PartCount)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	srv.dls = append(srv.dls, dl)
	js, _ := json.Marshal("ok")
	rwr.Write(js)
}

func (srv *DServ) startTask(rwr http.ResponseWriter, req *http.Request) {
	srv.oplock.Lock()
	defer func() {
		srv.oplock.Unlock()
		req.Body.Close()
	}()
	bodyData, err := ioutil.ReadAll(req.Body)
	rwr.Header().Set("Access-Control-Allow-Origin", "*")
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	var ind int
	if err := json.Unmarshal(bodyData, &ind); err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if !(len(srv.dls) > ind) {
		http.Error(rwr, "error: id is out of jobs list", http.StatusInternalServerError)
		return
	}

	if errs := srv.dls[ind].StartAll(); len(errs) > 0 {
		http.Error(rwr, "error: can't start all part", http.StatusInternalServerError)
		return
	}
	js, _ := json.Marshal("ok")
	rwr.Write(js)
}

func (srv *DServ) stopTask(rwr http.ResponseWriter, req *http.Request) {
	srv.oplock.Lock()
	defer func() {
		srv.oplock.Unlock()
		req.Body.Close()
	}()
	bodyData, err := ioutil.ReadAll(req.Body)
	rwr.Header().Set("Access-Control-Allow-Origin", "*")
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	var ind int
	if err := json.Unmarshal(bodyData, &ind); err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if !(len(srv.dls) > ind) {
		http.Error(rwr, "error: id is out of jobs list", http.StatusInternalServerError)
		return
	}

	srv.dls[ind].StopAll()
	js, _ := json.Marshal("ok")
	rwr.Write(js)
}

func (srv *DServ) startAllTask(rwr http.ResponseWriter, req *http.Request) {
	srv.oplock.Lock()
	defer func() {
		srv.oplock.Unlock()
		req.Body.Close()
	}()
	_, err := ioutil.ReadAll(req.Body)
	rwr.Header().Set("Access-Control-Allow-Origin", "*")
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	for _, e := range srv.dls {
		e.StartAll()
	}
	js, _ := json.Marshal("ok")
	rwr.Write(js)
}

func (srv *DServ) StopAllTask() {
	srv.oplock.Lock()
	defer func() {
		srv.oplock.Unlock()
	}()
	for _, e := range srv.dls {
		e.StopAll()
	}
}

func (srv *DServ) stopAllTask(rwr http.ResponseWriter, req *http.Request) {
	srv.oplock.Lock()
	defer func() {
		req.Body.Close()
	}()
	_, err := ioutil.ReadAll(req.Body)
	rwr.Header().Set("Access-Control-Allow-Origin", "*")
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	srv.StopAllTask()
	js, _ := json.Marshal("ok")
	rwr.Write(js)
}

func (srv *DServ) removeTask(rwr http.ResponseWriter, req *http.Request) {
	srv.oplock.Lock()
	defer func() {
		srv.oplock.Unlock()
		req.Body.Close()
	}()
	bodyData, err := ioutil.ReadAll(req.Body)
	rwr.Header().Set("Access-Control-Allow-Origin", "*")
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	var ind int
	if err := json.Unmarshal(bodyData, &ind); err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if !(len(srv.dls) > ind) {
		http.Error(rwr, "error: id is out of jobs list", http.StatusInternalServerError)
		return
	}

	log.Printf("try stop segment download %v", srv.dls[ind].StopAll())
	srv.dls = append(srv.dls[:ind], srv.dls[ind+1:]...)
	js, _ := json.Marshal("ok")
	rwr.Write(js)
}

func (srv *DServ) progressJson(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Access-Control-Allow-Origin", "*")
	jbs := make([]DJob, 0, len(srv.dls))
	for ind, i := range srv.dls {
		prs := i.GetProgress()
		var d int64
		var s int64
		for _, p := range prs {
			d = d + (p.Pos - p.From)
			s += p.Speed
		}
		j := DJob{
			Id:         ind,
			FileName:   i.Fi.FileName,
			Size:       i.Fi.Size,
			Progress:   (d * 100 / i.Fi.Size),
			Downloaded: d,
			Speed:      s,
		}
		jbs = append(jbs, j)
	}
	js, err := json.Marshal(jbs)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}
	rwr.Write(js)

}

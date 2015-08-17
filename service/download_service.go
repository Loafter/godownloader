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

const htmlData = "PCFkb2N0eXBlIGh0bWw+DQoNCjxodG1sPg0KDQo8aGVhZD4NCgk8dGl0bGU+R08gRE9XTkxPQUQ8L3RpdGxlPg0KCTxtZXRhIG5hbWU9InZpZXdwb3J0IiBjb250ZW50PSJ3aWR0aD1kZXZpY2Utd2lkdGgiPg0KCTxsaW5rIHJlbD0ic3R5bGVzaGVldCIgaHJlZj0iaHR0cHM6Ly9uZXRkbmEuYm9vdHN0cmFwY2RuLmNvbS9ib290c3dhdGNoLzMuMC4wL2pvdXJuYWwvYm9vdHN0cmFwLm1pbi5jc3MiPg0KCTxsaW5rIHJlbD0ic3R5bGVzaGVldCIgdHlwZT0idGV4dC9jc3MiIG1lZGlhPSJzY3JlZW4iIGhyZWY9Imh0dHA6Ly93d3cuZ3VyaWRkby5uZXQvZGVtby9jc3MvdHJpcmFuZC91aS5qcWdyaWQtYm9vdHN0cmFwLmNzcyI+DQoJPHNjcmlwdCB0eXBlPSJ0ZXh0L2phdmFzY3JpcHQiIHNyYz0iaHR0cHM6Ly9hamF4Lmdvb2dsZWFwaXMuY29tL2FqYXgvbGlicy9qcXVlcnkvMi4wLjMvanF1ZXJ5Lm1pbi5qcyI+PC9zY3JpcHQ+DQoJPHNjcmlwdCB0eXBlPSJ0ZXh0L2phdmFzY3JpcHQiIHNyYz0iaHR0cHM6Ly9uZXRkbmEuYm9vdHN0cmFwY2RuLmNvbS9ib290c3RyYXAvMy4zLjQvanMvYm9vdHN0cmFwLm1pbi5qcyI+PC9zY3JpcHQ+DQoJPHNjcmlwdCB0eXBlPSJ0ZXh0L2phdmFzY3JpcHQiIHNyYz0iaHR0cDovL3d3dy5ndXJpZGRvLm5ldC9kZW1vL2pzL3RyaXJhbmQvanF1ZXJ5LmpxR3JpZC5taW4uanMiPjwvc2NyaXB0Pg0KCTxzY3JpcHQgdHlwZT0idGV4dC9qYXZhc2NyaXB0IiBzcmM9Imh0dHA6Ly93d3cuZ3VyaWRkby5uZXQvZGVtby9qcy90cmlyYW5kL2kxOG4vZ3JpZC5sb2NhbGUtZW4uanMiPjwvc2NyaXB0Pg0KCTxsaW5rIHJlbD0ic3R5bGVzaGVldCIgaHJlZj0iLy9jb2RlLmpxdWVyeS5jb20vdWkvMS4xMS40L3RoZW1lcy9zbW9vdGhuZXNzL2pxdWVyeS11aS5jc3MiPg0KCTxzY3JpcHQgc3JjPSJodHRwOi8vY29kZS5qcXVlcnkuY29tL3VpLzEuMTEuNC9qcXVlcnktdWkuanMiPjwvc2NyaXB0Pg0KCTxzdHlsZSB0eXBlPSJ0ZXh0L2NzcyI+DQoJCWJvZHkgew0KCQkJcGFkZGluZy10b3A6IDUwcHg7DQoJCQlwYWRkaW5nLWJvdHRvbTogMjBweDsNCgkJfQ0KCQkNCgkJLnRhYmxlIC5wcm9ncmVzcyB7DQoJCQltYXJnaW4tYm90dG9tOiAwcHg7DQoJCX0NCgk8L3N0eWxlPg0KCTxzY3JpcHQgdHlwZT0idGV4dC9qYXZhc2NyaXB0Ij4NCgkJZnVuY3Rpb24gVXBkYXRlVGFibGUoKSB7DQoJCQkkKCIjanFHcmlkIikNCgkJCQkuanFHcmlkKHsNCgkJCQkJdXJsOiAnaHR0cDovL2xvY2FsaG9zdDo5OTgxL3Byb2dyZXNzLmpzb24nLA0KCQkJCQltdHlwZTogIkdFVCIsDQoJCQkJCWFqYXhTdWJncmlkT3B0aW9uczogew0KCQkJCQkJYXN5bmM6IGZhbHNlDQoJCQkJCX0sDQoJCQkJCXN0eWxlVUk6ICdCb290c3RyYXAnLA0KCQkJCQlkYXRhdHlwZTogImpzb24iLA0KCQkJCQljb2xNb2RlbDogW3sNCgkJCQkJCWxhYmVsOiAnIycsDQoJCQkJCQluYW1lOiAnSWQnLA0KCQkJCQkJa2V5OiB0cnVlLA0KCQkJCQkJd2lkdGg6IDUNCgkJCQkJfSwgew0KCQkJCQkJbGFiZWw6ICdGaWxlIE5hbWUnLA0KCQkJCQkJbmFtZTogJ0ZpbGVOYW1lJywNCgkJCQkJCXdpZHRoOiAxNQ0KCQkJCQl9LCB7DQoJCQkJCQlsYWJlbDogJ1NpemUnLA0KCQkJCQkJbmFtZTogJ1NpemUnLA0KCQkJCQkJd2lkdGg6IDIwLA0KCQkJCQkJZm9ybWF0dGVyOiBGb3JtYXRCeXRlDQoJCQkJCX0sIHsNCgkJCQkJCWxhYmVsOiAnRG93bmxvYWRlZCcsDQoJCQkJCQluYW1lOiAnRG93bmxvYWRlZCcsDQoJCQkJCQl3aWR0aDogMjAsDQoJCQkJCQlmb3JtYXR0ZXI6IEZvcm1hdEJ5dGUNCgkJCQkJfSwgew0KCQkJCQkJbGFiZWw6ICclJywNCgkJCQkJCW5hbWU6ICdQcm9ncmVzcycsDQoJCQkJCQl3aWR0aDogNQ0KCQkJCQl9LCB7DQoJCQkJCQlsYWJlbDogJ1NwZWVkJywNCgkJCQkJCW5hbWU6ICdTcGVlZCcsDQoJCQkJCQl3aWR0aDogMTUsDQoJCQkJCQlmb3JtYXR0ZXI6IEZvcm1hdEJ5dGUNCgkJCQkJfSwgew0KCQkJCQkJbGFiZWw6ICdQcm9ncmVzcycsDQoJCQkJCQluYW1lOiAnUHJvZ3Jlc3MnLA0KCQkJCQkJZm9ybWF0dGVyOiBGb3JtYXRQcm9ncmVzc0Jhcg0KCQkJCQl9XSwNCgkJCQkJdmlld3JlY29yZHM6IHRydWUsDQoJCQkJCXJvd051bTogMjAsDQoJCQkJCXBhZ2VyOiAiI2pxR3JpZFBhZ2VyIg0KCQkJCX0pOw0KCQl9DQoNCgkJZnVuY3Rpb24gRml4VGFibGUoKSB7DQoJCQkkLmV4dGVuZCgkLmpncmlkLmFqYXhPcHRpb25zLCB7DQoJCQkJYXN5bmM6IGZhbHNlDQoJCQl9KQ0KCQkJJCgiI2pxR3JpZCIpDQoJCQkJLnNldEdyaWRXaWR0aCgkKHdpbmRvdykNCgkJCQkJLndpZHRoKCkgLSA1KQ0KCQkJJCgiI2pxR3JpZCIpDQoJCQkJLnNldEdyaWRIZWlnaHQoJCh3aW5kb3cpDQoJCQkJCS5oZWlnaHQoKSkNCgkJCSQod2luZG93KQ0KCQkJCS5iaW5kKCdyZXNpemUnLCBmdW5jdGlvbigpIHsNCgkJCQkJJCgiI2pxR3JpZCIpDQoJCQkJCQkuc2V0R3JpZFdpZHRoKCQod2luZG93KQ0KCQkJCQkJCS53aWR0aCgpIC0gNSk7DQoJCQkJCSQoIiNqcUdyaWQiKQ0KCQkJCQkJLnNldEdyaWRIZWlnaHQoJCh3aW5kb3cpDQoJCQkJCQkJLmhlaWdodCgpKQ0KCQkJCX0pDQoJCX0NCg0KCQlmdW5jdGlvbiBVcGRhdGVEYXRhKCkgew0KCQkJdmFyIGdyaWQgPSAkKCIjanFHcmlkIik7DQoJCQl2YXIgcm93S2V5ID0gZ3JpZC5qcUdyaWQoJ2dldEdyaWRQYXJhbScsICJzZWxyb3ciKTsNCgkJCSQoIiNqcUdyaWQiKS50cmlnZ2VyKCJyZWxvYWRHcmlkIik7DQoJCQlpZihyb3dLZXkpIHsNCgkJCQkkKCcjanFHcmlkJykuanFHcmlkKCJyZXNldFNlbGVjdGlvbiIpDQoJCQkJJCgnI2pxR3JpZCcpLmpxR3JpZCgnc2V0U2VsZWN0aW9uJywgcm93S2V5KTsNCgkJCX0NCgkJfQ0KDQoJCWZ1bmN0aW9uIEZvcm1hdFByb2dyZXNzQmFyKGNlbGxWYWx1ZSwgb3B0aW9ucywgcm93T2JqZWN0KSB7DQoJCQl2YXIgaW50VmFsID0gcGFyc2VJbnQoY2VsbFZhbHVlKTsNCg0KCQkJdmFyIGNlbGxIdG1sID0gJzxkaXYgY2xhc3M9InByb2dyZXNzIj48ZGl2IGNsYXNzPSJwcm9ncmVzcy1iYXIiIHN0eWxlPSJ3aWR0aDogJyArIGludFZhbCArICclOyI+PC9kaXY+PC9kaXY+Jw0KDQoJCQlyZXR1cm4gY2VsbEh0bWw7DQoJCX0NCg0KCQlmdW5jdGlvbiBGb3JtYXRCeXRlKGNlbGxWYWx1ZSwgb3B0aW9ucywgcm93T2JqZWN0KSB7DQoJCQl2YXIgaW50VmFsID0gcGFyc2VJbnQoY2VsbFZhbHVlKTsNCgkJCXZhciByYXMgPSAiIGIuIg0KCQkJaWYoaW50VmFsID4gMTAyNCkgew0KCQkJCWludFZhbCAvPSAxMDI0DQoJCQkJcmFzID0gIiBrYi4iDQoJCQl9DQoJCQlpZihpbnRWYWwgPiAxMDI0KSB7DQoJCQkJaW50VmFsIC89IDEwMjQNCgkJCQlyYXMgPSAiIE1iLiINCgkJCX0NCgkJCWlmKGludFZhbCA+IDEwMjQpIHsNCgkJCQlpbnRWYWwgLz0gMTAyNA0KCQkJCXJhcyA9ICIgR2IuIg0KCQkJfQ0KDQoJCQlpZihpbnRWYWwgPiAxMDI0KSB7DQoJCQkJaW50VmFsIC89IDEwMjQNCgkJCQlyYXMgPSAiIFRiLiINCgkJCX0NCgkJCXZhciBjZWxsSHRtbCA9IChpbnRWYWwpLnRvRml4ZWQoMSkgKyByYXM7DQoJCQlyZXR1cm4gY2VsbEh0bWw7DQoJCX0NCg0KCQlmdW5jdGlvbiBPbkxvYWQoKSB7DQoNCgkJCVVwZGF0ZVRhYmxlKCkNCgkJCUZpeFRhYmxlKCkNCgkJCXNldEludGVydmFsKFVwZGF0ZURhdGEsIDUwMCk7DQoJCX0NCg0KCQlmdW5jdGlvbiBBZGREb3dubG9hZCgpIHsNCgkJCXZhciByZXEgPSB7DQoJCQkJUGFydENvdW50OiBwYXJzZUludCgkKCIjcGFydF9jb3VudF9pZCIpLnZhbCgpKSwNCgkJCQlGaWxlUGF0aDogJCgiI3NhdmVfcGF0aF9pZCIpLnZhbCgpLA0KCQkJCVVybDogJCgiI3VybF9pZCIpLnZhbCgpDQoJCQl9Ow0KCQkJJC5hamF4KHsNCgkJCQkJdXJsOiAiL2FkZF90YXNrIiwNCgkJCQkJdHlwZTogIlBPU1QiLA0KCQkJCQlkYXRhOiBKU09OLnN0cmluZ2lmeShyZXEpLA0KCQkJCQlkYXRhVHlwZTogInRleHQiDQoJCQkJfSkNCgkJCQkuZXJyb3IoZnVuY3Rpb24oanNvbkRhdGEpIHsNCgkJCQkJY29uc29sZS5sb2coanNvbkRhdGEpDQoJCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIFJlbW92ZURvd25sb2FkKCkgew0KCQkJdmFyIGdyaWQgPSAkKCIjanFHcmlkIik7DQoJCQl2YXIgcm93S2V5ID0gcGFyc2VJbnQoZ3JpZC5qcUdyaWQoJ2dldEdyaWRQYXJhbScsICJzZWxyb3ciKSk7DQoJCQl2YXIgcmVxID0gcm93S2V5Ow0KCQkJJC5hamF4KHsNCgkJCQkJdXJsOiAiL3JlbW92ZV90YXNrIiwNCgkJCQkJdHlwZTogIlBPU1QiLA0KCQkJCQlkYXRhOiBKU09OLnN0cmluZ2lmeShyZXEpLA0KCQkJCQlkYXRhVHlwZTogInRleHQiDQoJCQkJfSkNCgkJCQkuZXJyb3IoZnVuY3Rpb24oanNvbkRhdGEpIHsNCgkJCQkJY29uc29sZS5sb2coanNvbkRhdGEpDQoJCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIFN0YXJ0RG93bmxvYWQoKSB7DQoJCQl2YXIgZ3JpZCA9ICQoIiNqcUdyaWQiKTsNCgkJCXZhciByb3dLZXkgPSBwYXJzZUludChncmlkLmpxR3JpZCgnZ2V0R3JpZFBhcmFtJywgInNlbHJvdyIpKTsNCgkJCXZhciByZXEgPSByb3dLZXk7DQoJCQkkLmFqYXgoew0KCQkJCQl1cmw6ICIvc3RhcnRfdGFzayIsDQoJCQkJCXR5cGU6ICJQT1NUIiwNCgkJCQkJZGF0YTogSlNPTi5zdHJpbmdpZnkocmVxKSwNCgkJCQkJZGF0YVR5cGU6ICJ0ZXh0Ig0KCQkJCX0pDQoJCQkJLmVycm9yKGZ1bmN0aW9uKGpzb25EYXRhKSB7DQoJCQkJCWNvbnNvbGUubG9nKGpzb25EYXRhKQ0KCQkJCX0pDQoJCX0NCg0KCQlmdW5jdGlvbiBTdGFydEFsbERvd25sb2FkKCkgew0KCQkJJC5hamF4KHsNCgkJCQkJdXJsOiAiL3N0YXJ0X2FsbF90YXNrIiwNCgkJCQkJdHlwZTogIlBPU1QiLA0KCQkJCQlkYXRhVHlwZTogInRleHQiDQoJCQkJfSkNCgkJCQkuZXJyb3IoZnVuY3Rpb24oanNvbkRhdGEpIHsNCgkJCQkJY29uc29sZS5sb2coanNvbkRhdGEpDQoJCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIFN0b3BBbGxEb3dubG9hZCgpIHsNCgkJCSQuYWpheCh7DQoJCQkJCXVybDogIi9zdG9wX2FsbF90YXNrIiwNCgkJCQkJdHlwZTogIlBPU1QiLA0KCQkJCQlkYXRhVHlwZTogInRleHQiDQoJCQkJfSkNCgkJCQkuZXJyb3IoZnVuY3Rpb24oanNvbkRhdGEpIHsNCgkJCQkJY29uc29sZS5sb2coanNvbkRhdGEpDQoJCQkJfSkNCgkJfQ0KDQoNCgkJZnVuY3Rpb24gU3RvcERvd25sb2FkKCkgew0KCQkJdmFyIGdyaWQgPSAkKCIjanFHcmlkIik7DQoJCQl2YXIgcm93S2V5ID0gcGFyc2VJbnQoZ3JpZC5qcUdyaWQoJ2dldEdyaWRQYXJhbScsICJzZWxyb3ciKSk7DQoJCQl2YXIgcmVxID0gcm93S2V5Ow0KCQkJJC5hamF4KHsNCgkJCQkJdXJsOiAiL3N0b3BfdGFzayIsDQoJCQkJCXR5cGU6ICJQT1NUIiwNCgkJCQkJZGF0YTogSlNPTi5zdHJpbmdpZnkocmVxKSwNCgkJCQkJZGF0YVR5cGU6ICJ0ZXh0Ig0KCQkJCX0pDQoJCQkJLmVycm9yKGZ1bmN0aW9uKGpzb25EYXRhKSB7DQoJCQkJCWNvbnNvbGUubG9nKGpzb25EYXRhKQ0KCQkJCX0pDQoJCX0NCgk8L3NjcmlwdD4NCjwvaGVhZD4NCg0KPGJvZHkgb25sb2FkPSJPbkxvYWQoKSI+DQoJPGRpdiBjbGFzcz0ibmF2YmFyIG5hdmJhci1pbnZlcnNlIG5hdmJhci1maXhlZC10b3AiPg0KCQk8ZGl2IGNsYXNzPSJjb250YWluZXIiPg0KCQkJPGRpdiBjbGFzcz0ibmF2YmFyLWhlYWRlciI+DQoJCQkJPGJ1dHRvbiB0eXBlPSJidXR0b24iIGNsYXNzPSJuYXZiYXItdG9nZ2xlIiBkYXRhLXRvZ2dsZT0iY29sbGFwc2UiIGRhdGEtdGFyZ2V0PSIubmF2YmFyLWNvbGxhcHNlIj4NCgkJCQkJPHNwYW4gY2xhc3M9Imljb24tYmFyIj48L3NwYW4+PHNwYW4gY2xhc3M9Imljb24tYmFyIj48L3NwYW4+PHNwYW4gY2xhc3M9Imljb24tYmFyIj48L3NwYW4+DQoJCQkJPC9idXR0b24+DQoJCQkJPGEgY2xhc3M9Im5hdmJhci1icmFuZCIgaHJlZj0iIyI+R08gRG93bmxvYWRlcjwvYT4NCgkJCTwvZGl2Pg0KCQkJPGRpdiBjbGFzcz0ibmF2YmFyLWNvbGxhcHNlIGNvbGxhcHNlIj4NCgkJCQk8dWwgY2xhc3M9Im5hdiBuYXZiYXItbmF2Ij4NCgkJCQkJPGxpIGNsYXNzPSJkcm9wZG93biI+DQoJCQkJCQk8YSBocmVmPSIjIiBjbGFzcz0iZHJvcGRvd24tdG9nZ2xlIiBkYXRhLXRvZ2dsZT0iZHJvcGRvd24iPkZpbGUgPGIgY2xhc3M9ImNhcmV0Ij48L2I+PC9hPg0KCQkJCQkJPHVsIGNsYXNzPSJkcm9wZG93bi1tZW51Ij4NCgkJCQkJCQk8bGk+DQoJCQkJCQkJCTxhIGRhdGEtdG9nZ2xlPSJtb2RhbCIgZGF0YS10YXJnZXQ9IiNteU1vZGFsIj5BZGQgZG93bmxvYWQ8L2E+DQoJCQkJCQkJPC9saT4NCgkJCQkJCQk8bGkgb25jbGljaz0iUmVtb3ZlRG93bmxvYWQoKSI+DQoJCQkJCQkJCTxhIGhyZWY9IiMiPkRlbGV0ZSBkb3dubG9hZDwvYT4NCgkJCQkJCQk8L2xpPg0KCQkJCQkJPC91bD4NCgkJCQkJPC9saT4NCgkJCQkJPGxpIGNsYXNzPSJkcm9wZG93biI+DQoJCQkJCQk8YSBocmVmPSIjIiBjbGFzcz0iZHJvcGRvd24tdG9nZ2xlIiBkYXRhLXRvZ2dsZT0iZHJvcGRvd24iPkFjdGlvbiA8YiBjbGFzcz0iY2FyZXQiPjwvYj48L2E+DQoJCQkJCQk8dWwgY2xhc3M9ImRyb3Bkb3duLW1lbnUiPg0KCQkJCQkJCTxsaSBvbmNsaWNrPSJTdGFydERvd25sb2FkKCkiPg0KCQkJCQkJCQk8YSBocmVmPSIjIj5TdGFydDwvYT4NCgkJCQkJCQk8L2xpPg0KCQkJCQkJCTxsaSBvbmNsaWNrPSJTdG9wRG93bmxvYWQoKSI+DQoJCQkJCQkJCTxhIGhyZWY9IiMiPlN0b3A8L2E+DQoJCQkJCQkJPC9saT4NCgkJCQkJCQk8bGkgY2xhc3M9ImRpdmlkZXIiPjwvbGk+DQoJCQkJCQkJPGxpIG9uY2xpY2s9IlN0YXJ0QWxsRG93bmxvYWQoKSI+DQoJCQkJCQkJCTxhIGhyZWY9IiMiPlN0YXJ0IGFsbDwvYT4NCgkJCQkJCQk8L2xpPg0KCQkJCQkJCTxsaSBvbmNsaWNrPSJTdG9wQWxsRG93bmxvYWQoKSI+DQoJCQkJCQkJCTxhIGhyZWY9IiMiPlN0b3AgYWxsPC9hPg0KCQkJCQkJCTwvbGk+DQoJCQkJCQk8L3VsPg0KCQkJCQk8L2xpPg0KCQkJCQk8bGk+DQoJCQkJCQk8YSBocmVmPSIjYWJvdXQiPkFib3V0PC9hPg0KCQkJCQk8L2xpPg0KCQkJCTwvdWw+DQoJCQk8L2Rpdj4NCgkJCTwhLS0vLm5hdmJhci1jb2xsYXBzZSAtLT4NCgkJPC9kaXY+DQoJPC9kaXY+DQoJPC9QPg0KCTx0YWJsZSBpZD0ianFHcmlkIj48L3RhYmxlPg0KDQoJPCEtLSBNb2RhbCAtLT4NCgk8ZGl2IGNsYXNzPSJtb2RhbCBmYWRlIiBpZD0ibXlNb2RhbCIgcm9sZT0iZGlhbG9nIj4NCgkJPGRpdiBjbGFzcz0ibW9kYWwtZGlhbG9nIj4NCg0KCQkJPCEtLSBNb2RhbCBjb250ZW50LS0+DQoJCQk8ZGl2IGNsYXNzPSJtb2RhbC1jb250ZW50Ij4NCgkJCQk8ZGl2IGNsYXNzPSJtb2RhbC1oZWFkZXIiPg0KCQkJCQk8YnV0dG9uIHR5cGU9ImJ1dHRvbiIgY2xhc3M9ImNsb3NlIiBkYXRhLWRpc21pc3M9Im1vZGFsIj4mdGltZXM7PC9idXR0b24+DQoJCQkJCTxoNCBjbGFzcz0ibW9kYWwtdGl0bGUiPkVudGVyIFVybDwvaDQ+DQoJCQkJPC9kaXY+DQoJCQkJPGRpdiBjbGFzcz0ibW9kYWwtYm9keSI+DQoJCQkJCTxkaXYgY2xhc3M9ImZvcm0tZ3JvdXAiPg0KCQkJCQkJPGxhYmVsIGNsYXNzPSJjb250cm9sLWxhYmVsIj5Vcmw8L2xhYmVsPg0KDQoJCQkJCQk8ZGl2IGNsYXNzPSJjb250cm9scyI+DQoJCQkJCQkJPGlucHV0IHR5cGU9InRleHQiIGlkPSJ1cmxfaWQiIGNsYXNzPSJmb3JtLWNvbnRyb2wiIHZhbHVlPSJodHRwOi8vdWNtaXJyb3IuY2FudGVyYnVyeS5hYy5uei9saW51eC91YnVudHUvdHJ1c3R5L3d1YmkuZXhlIj4NCgkJCQkJCTwvZGl2Pg0KCQkJCQkJPGxhYmVsIGNsYXNzPSJjb250cm9sLWxhYmVsIj5TYXZlIHBhdGg8L2xhYmVsPg0KDQoJCQkJCQk8ZGl2IGNsYXNzPSJjb250cm9scyI+DQoJCQkJCQkJPGlucHV0IHR5cGU9InRleHQiIGlkPSJzYXZlX3BhdGhfaWQiIGNsYXNzPSJmb3JtLWNvbnRyb2wiIHZhbHVlPSIuL3d1YmkuZXhlIj4NCgkJCQkJCTwvZGl2Pg0KCQkJCQkJPGxhYmVsIGNsYXNzPSJjb250cm9sLWxhYmVsIj5QYXJ0cyBjb3VudDwvbGFiZWw+DQoJCQkJCQk8c2VsZWN0IGNsYXNzPSJmb3JtLWNvbnRyb2wiIGlkPSJwYXJ0X2NvdW50X2lkIj4NCgkJCQkJCQk8b3B0aW9uPjE8L29wdGlvbj4NCgkJCQkJCQk8b3B0aW9uPjI8L29wdGlvbj4NCgkJCQkJCQk8b3B0aW9uPjQ8L29wdGlvbj4NCgkJCQkJCQk8b3B0aW9uPjg8L29wdGlvbj4NCgkJCQkJCQk8b3B0aW9uPjE2PC9vcHRpb24+DQoJCQkJCQk8L3NlbGVjdD4NCg0KCQkJCQkJPGRpdiBjbGFzcz0ibW9kYWwtZm9vdGVyIj4NCgkJCQkJCQk8YSBjbGFzcz0iYnRuIGJ0bi1wcmltYXJ5IiBvbmNsaWNrPSJBZGREb3dubG9hZCgpIiBkYXRhLWRpc21pc3M9Im1vZGFsIj5TdGFydCBkb3dubG9hZDwvYT4NCgkJCQkJCTwvZGl2Pg0KCQkJCQk8L2Rpdj4NCgkJCQk8L2Rpdj4NCgkJCTwvZGl2Pg0KDQoJCTwvZGl2Pg0KCTwvZGl2Pg0KPC9ib2R5Pg0KDQo8L2h0bWw+"

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

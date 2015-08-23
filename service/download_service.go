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

const htmlData = "PCFkb2N0eXBlIGh0bWw+DQoNCjxodG1sPg0KDQo8aGVhZD4NCgk8dGl0bGU+R08gRE9XTkxPQUQ8L3RpdGxlPg0KCTxtZXRhIG5hbWU9InZpZXdwb3J0IiBjb250ZW50PSJ3aWR0aD1kZXZpY2Utd2lkdGgiPg0KCTxsaW5rIHJlbD0ic3R5bGVzaGVldCIgaHJlZj0iaHR0cHM6Ly9uZXRkbmEuYm9vdHN0cmFwY2RuLmNvbS9ib290c3dhdGNoLzMuMC4wL2pvdXJuYWwvYm9vdHN0cmFwLm1pbi5jc3MiPg0KCTxsaW5rIHJlbD0ic3R5bGVzaGVldCIgdHlwZT0idGV4dC9jc3MiIG1lZGlhPSJzY3JlZW4iIGhyZWY9Imh0dHA6Ly93d3cuZ3VyaWRkby5uZXQvZGVtby9jc3MvdHJpcmFuZC91aS5qcWdyaWQtYm9vdHN0cmFwLmNzcyI+DQoJPHNjcmlwdCB0eXBlPSJ0ZXh0L2phdmFzY3JpcHQiIHNyYz0iaHR0cHM6Ly9hamF4Lmdvb2dsZWFwaXMuY29tL2FqYXgvbGlicy9qcXVlcnkvMi4wLjMvanF1ZXJ5Lm1pbi5qcyI+PC9zY3JpcHQ+DQoJPHNjcmlwdCB0eXBlPSJ0ZXh0L2phdmFzY3JpcHQiIHNyYz0iaHR0cHM6Ly9uZXRkbmEuYm9vdHN0cmFwY2RuLmNvbS9ib290c3RyYXAvMy4zLjQvanMvYm9vdHN0cmFwLm1pbi5qcyI+PC9zY3JpcHQ+DQoJPHNjcmlwdCB0eXBlPSJ0ZXh0L2phdmFzY3JpcHQiIHNyYz0iaHR0cDovL3d3dy5ndXJpZGRvLm5ldC9kZW1vL2pzL3RyaXJhbmQvanF1ZXJ5LmpxR3JpZC5taW4uanMiPjwvc2NyaXB0Pg0KCTxzY3JpcHQgdHlwZT0idGV4dC9qYXZhc2NyaXB0IiBzcmM9Imh0dHA6Ly93d3cuZ3VyaWRkby5uZXQvZGVtby9qcy90cmlyYW5kL2kxOG4vZ3JpZC5sb2NhbGUtZW4uanMiPjwvc2NyaXB0Pg0KCTxsaW5rIHJlbD0ic3R5bGVzaGVldCIgaHJlZj0iLy9jb2RlLmpxdWVyeS5jb20vdWkvMS4xMS40L3RoZW1lcy9zbW9vdGhuZXNzL2pxdWVyeS11aS5jc3MiPg0KCTxzY3JpcHQgc3JjPSJodHRwOi8vY29kZS5qcXVlcnkuY29tL3VpLzEuMTEuNC9qcXVlcnktdWkuanMiPjwvc2NyaXB0Pg0KCTxzdHlsZSB0eXBlPSJ0ZXh0L2NzcyI+DQoJCWJvZHkgew0KCQlwYWRkaW5nLXRvcDogNTBweDsNCgkJcGFkZGluZy1ib3R0b206IDIwcHg7DQoJCX0NCg0KCQkudGFibGUgLnByb2dyZXNzIHsNCgkJbWFyZ2luLWJvdHRvbTogMHB4Ow0KCQl9DQoJPC9zdHlsZT4NCgk8c2NyaXB0IHR5cGU9InRleHQvamF2YXNjcmlwdCI+DQoJCWZ1bmN0aW9uIFVwZGF0ZVRhYmxlKCkgew0KCQkJJCgiI2pxR3JpZCIpDQoJCQkJLmpxR3JpZCh7DQoJCQkJCXVybDogJ2h0dHA6Ly9sb2NhbGhvc3Q6OTk4MS9wcm9ncmVzcy5qc29uJywNCgkJCQkJbXR5cGU6ICJHRVQiLA0KCQkJCQlhamF4U3ViZ3JpZE9wdGlvbnM6IHsNCgkJCQkJCWFzeW5jOiBmYWxzZQ0KCQkJCQl9LA0KCQkJCQlzdHlsZVVJOiAnQm9vdHN0cmFwJywNCgkJCQkJZGF0YXR5cGU6ICJqc29uIiwNCgkJCQkJY29sTW9kZWw6IFt7DQoJCQkJCQlsYWJlbDogJyMnLA0KCQkJCQkJbmFtZTogJ0lkJywNCgkJCQkJCWtleTogdHJ1ZSwNCgkJCQkJCXdpZHRoOiA1DQoJCQkJCX0sIHsNCgkJCQkJCWxhYmVsOiAnRmlsZSBOYW1lJywNCgkJCQkJCW5hbWU6ICdGaWxlTmFtZScsDQoJCQkJCQl3aWR0aDogMTUNCgkJCQkJfSwgew0KCQkJCQkJbGFiZWw6ICdTaXplJywNCgkJCQkJCW5hbWU6ICdTaXplJywNCgkJCQkJCXdpZHRoOiAyMCwNCgkJCQkJCWZvcm1hdHRlcjogRm9ybWF0Qnl0ZQ0KCQkJCQl9LCB7DQoJCQkJCQlsYWJlbDogJ0Rvd25sb2FkZWQnLA0KCQkJCQkJbmFtZTogJ0Rvd25sb2FkZWQnLA0KCQkJCQkJd2lkdGg6IDIwLA0KCQkJCQkJZm9ybWF0dGVyOiBGb3JtYXRCeXRlDQoJCQkJCX0sIHsNCgkJCQkJCWxhYmVsOiAnJScsDQoJCQkJCQluYW1lOiAnUHJvZ3Jlc3MnLA0KCQkJCQkJd2lkdGg6IDUNCgkJCQkJfSwgew0KCQkJCQkJbGFiZWw6ICdTcGVlZCcsDQoJCQkJCQluYW1lOiAnU3BlZWQnLA0KCQkJCQkJd2lkdGg6IDE1LA0KCQkJCQkJZm9ybWF0dGVyOiBGb3JtYXRCeXRlDQoJCQkJCX0sIHsNCgkJCQkJCWxhYmVsOiAnUHJvZ3Jlc3MnLA0KCQkJCQkJbmFtZTogJ1Byb2dyZXNzJywNCgkJCQkJCWZvcm1hdHRlcjogRm9ybWF0UHJvZ3Jlc3NCYXINCgkJCQkJfV0sDQoJCQkJCXZpZXdyZWNvcmRzOiB0cnVlLA0KCQkJCQlyb3dOdW06IDIwLA0KCQkJCQlwYWdlcjogIiNqcUdyaWRQYWdlciINCgkJCQl9KTsNCgkJfQ0KDQoJCWZ1bmN0aW9uIEZpeFRhYmxlKCkgew0KCQkJJC5leHRlbmQoJC5qZ3JpZC5hamF4T3B0aW9ucywgew0KCQkJCWFzeW5jOiBmYWxzZQ0KCQkJfSkNCgkJCSQoIiNqcUdyaWQiKQ0KCQkJCS5zZXRHcmlkV2lkdGgoJCh3aW5kb3cpDQoJCQkJCS53aWR0aCgpIC0gNSkNCgkJCSQoIiNqcUdyaWQiKQ0KCQkJCS5zZXRHcmlkSGVpZ2h0KCQod2luZG93KQ0KCQkJCQkuaGVpZ2h0KCkpDQoJCQkkKHdpbmRvdykNCgkJCQkuYmluZCgncmVzaXplJywgZnVuY3Rpb24oKSB7DQoJCQkJCSQoIiNqcUdyaWQiKQ0KCQkJCQkJLnNldEdyaWRXaWR0aCgkKHdpbmRvdykNCgkJCQkJCQkud2lkdGgoKSAtIDUpOw0KCQkJCQkkKCIjanFHcmlkIikNCgkJCQkJCS5zZXRHcmlkSGVpZ2h0KCQod2luZG93KQ0KCQkJCQkJCS5oZWlnaHQoKSkNCgkJCQl9KQ0KCQl9DQoNCgkJZnVuY3Rpb24gVXBkYXRlRGF0YSgpIHsNCgkJCXZhciBncmlkID0gJCgiI2pxR3JpZCIpOw0KCQkJdmFyIHJvd0tleSA9IGdyaWQuanFHcmlkKCdnZXRHcmlkUGFyYW0nLCAic2Vscm93Iik7DQoJCQkkKCIjanFHcmlkIikudHJpZ2dlcigicmVsb2FkR3JpZCIpOw0KCQkJaWYocm93S2V5KSB7DQoJCQkJJCgnI2pxR3JpZCcpLmpxR3JpZCgicmVzZXRTZWxlY3Rpb24iKQ0KCQkJCSQoJyNqcUdyaWQnKS5qcUdyaWQoJ3NldFNlbGVjdGlvbicsIHJvd0tleSk7DQoJCQl9DQoJCX0NCg0KCQlmdW5jdGlvbiBGb3JtYXRQcm9ncmVzc0JhcihjZWxsVmFsdWUsIG9wdGlvbnMsIHJvd09iamVjdCkgew0KCQkJdmFyIGludFZhbCA9IHBhcnNlSW50KGNlbGxWYWx1ZSk7DQoNCgkJCXZhciBjZWxsSHRtbCA9ICc8ZGl2IGNsYXNzPSJwcm9ncmVzcyI+PGRpdiBjbGFzcz0icHJvZ3Jlc3MtYmFyIiBzdHlsZT0id2lkdGg6ICcgKyBpbnRWYWwgKyAnJTsiPjwvZGl2PjwvZGl2PicNCg0KCQkJcmV0dXJuIGNlbGxIdG1sOw0KCQl9DQoNCgkJZnVuY3Rpb24gRm9ybWF0Qnl0ZShjZWxsVmFsdWUsIG9wdGlvbnMsIHJvd09iamVjdCkgew0KCQkJdmFyIGludFZhbCA9IHBhcnNlSW50KGNlbGxWYWx1ZSk7DQoJCQl2YXIgcmFzID0gIiBiLiINCgkJCWlmKGludFZhbCA+IDEwMjQpIHsNCgkJCQlpbnRWYWwgLz0gMTAyNA0KCQkJCXJhcyA9ICIga2IuIg0KCQkJfQ0KCQkJaWYoaW50VmFsID4gMTAyNCkgew0KCQkJCWludFZhbCAvPSAxMDI0DQoJCQkJcmFzID0gIiBNYi4iDQoJCQl9DQoJCQlpZihpbnRWYWwgPiAxMDI0KSB7DQoJCQkJaW50VmFsIC89IDEwMjQNCgkJCQlyYXMgPSAiIEdiLiINCgkJCX0NCg0KCQkJaWYoaW50VmFsID4gMTAyNCkgew0KCQkJCWludFZhbCAvPSAxMDI0DQoJCQkJcmFzID0gIiBUYi4iDQoJCQl9DQoJCQl2YXIgY2VsbEh0bWwgPSAoaW50VmFsKS50b0ZpeGVkKDEpICsgcmFzOw0KCQkJcmV0dXJuIGNlbGxIdG1sOw0KCQl9DQoNCgkJZnVuY3Rpb24gT25Mb2FkKCkgew0KDQoJCQlVcGRhdGVUYWJsZSgpDQoJCQlGaXhUYWJsZSgpDQoJCQlzZXRJbnRlcnZhbChVcGRhdGVEYXRhLCA1MDApOw0KCQl9DQoNCgkJZnVuY3Rpb24gQWRkRG93bmxvYWQoKSB7DQoJCQl2YXIgcmVxID0gew0KCQkJCVBhcnRDb3VudDogcGFyc2VJbnQoJCgiI3BhcnRfY291bnRfaWQiKS52YWwoKSksDQoJCQkJRmlsZVBhdGg6ICQoIiNzYXZlX3BhdGhfaWQiKS52YWwoKSwNCgkJCQlVcmw6ICQoIiN1cmxfaWQiKS52YWwoKQ0KCQkJfTsNCgkJCSQuYWpheCh7DQoJCQkJCXVybDogIi9hZGRfdGFzayIsDQoJCQkJCXR5cGU6ICJQT1NUIiwNCgkJCQkJZGF0YTogSlNPTi5zdHJpbmdpZnkocmVxKSwNCgkJCQkJZGF0YVR5cGU6ICJ0ZXh0Ig0KCQkJCX0pDQoJCQkJLmVycm9yKGZ1bmN0aW9uKGpzb25EYXRhKSB7DQoJCQkJCWNvbnNvbGUubG9nKGpzb25EYXRhKQ0KCQkJCX0pDQoJCX0NCg0KCQlmdW5jdGlvbiBSZW1vdmVEb3dubG9hZCgpIHsNCgkJCXZhciBncmlkID0gJCgiI2pxR3JpZCIpOw0KCQkJdmFyIHJvd0tleSA9IHBhcnNlSW50KGdyaWQuanFHcmlkKCdnZXRHcmlkUGFyYW0nLCAic2Vscm93IikpOw0KCQkJdmFyIHJlcSA9IHJvd0tleTsNCgkJCSQuYWpheCh7DQoJCQkJCXVybDogIi9yZW1vdmVfdGFzayIsDQoJCQkJCXR5cGU6ICJQT1NUIiwNCgkJCQkJZGF0YTogSlNPTi5zdHJpbmdpZnkocmVxKSwNCgkJCQkJZGF0YVR5cGU6ICJ0ZXh0Ig0KCQkJCX0pDQoJCQkJLmVycm9yKGZ1bmN0aW9uKGpzb25EYXRhKSB7DQoJCQkJCWNvbnNvbGUubG9nKGpzb25EYXRhKQ0KCQkJCX0pDQoJCX0NCg0KCQlmdW5jdGlvbiBTdGFydERvd25sb2FkKCkgew0KCQkJdmFyIGdyaWQgPSAkKCIjanFHcmlkIik7DQoJCQl2YXIgcm93S2V5ID0gcGFyc2VJbnQoZ3JpZC5qcUdyaWQoJ2dldEdyaWRQYXJhbScsICJzZWxyb3ciKSk7DQoJCQl2YXIgcmVxID0gcm93S2V5Ow0KCQkJJC5hamF4KHsNCgkJCQkJdXJsOiAiL3N0YXJ0X3Rhc2siLA0KCQkJCQl0eXBlOiAiUE9TVCIsDQoJCQkJCWRhdGE6IEpTT04uc3RyaW5naWZ5KHJlcSksDQoJCQkJCWRhdGFUeXBlOiAidGV4dCINCgkJCQl9KQ0KCQkJCS5lcnJvcihmdW5jdGlvbihqc29uRGF0YSkgew0KCQkJCQljb25zb2xlLmxvZyhqc29uRGF0YSkNCgkJCQl9KQ0KCQl9DQoNCgkJZnVuY3Rpb24gU3RvcERvd25sb2FkKCkgew0KCQkJdmFyIGdyaWQgPSAkKCIjanFHcmlkIik7DQoJCQl2YXIgcm93S2V5ID0gcGFyc2VJbnQoZ3JpZC5qcUdyaWQoJ2dldEdyaWRQYXJhbScsICJzZWxyb3ciKSk7DQoJCQl2YXIgcmVxID0gcm93S2V5Ow0KCQkJJC5hamF4KHsNCgkJCQkJdXJsOiAiL3N0b3BfdGFzayIsDQoJCQkJCXR5cGU6ICJQT1NUIiwNCgkJCQkJZGF0YTogSlNPTi5zdHJpbmdpZnkocmVxKSwNCgkJCQkJZGF0YVR5cGU6ICJ0ZXh0Ig0KCQkJCX0pDQoJCQkJLmVycm9yKGZ1bmN0aW9uKGpzb25EYXRhKSB7DQoJCQkJCWNvbnNvbGUubG9nKGpzb25EYXRhKQ0KCQkJCX0pDQoJCX0NCg0KCQlmdW5jdGlvbiBTdGFydEFsbERvd25sb2FkKCkgew0KCQkJJC5hamF4KHsNCgkJCQkJdXJsOiAiL3N0YXJ0X2FsbF90YXNrIiwNCgkJCQkJdHlwZTogIlBPU1QiLA0KCQkJCQlkYXRhVHlwZTogInRleHQiDQoJCQkJfSkNCgkJCQkuZXJyb3IoZnVuY3Rpb24oanNvbkRhdGEpIHsNCgkJCQkJY29uc29sZS5sb2coanNvbkRhdGEpDQoJCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIFN0b3BBbGxEb3dubG9hZCgpIHsNCgkJCSQuYWpheCh7DQoJCQkJCXVybDogIi9zdG9wX2FsbF90YXNrIiwNCgkJCQkJdHlwZTogIlBPU1QiLA0KCQkJCQlkYXRhVHlwZTogInRleHQiDQoJCQkJfSkNCgkJCQkuZXJyb3IoZnVuY3Rpb24oanNvbkRhdGEpIHsNCgkJCQkJY29uc29sZS5sb2coanNvbkRhdGEpDQoJCQkJfSkNCgkJfQ0KDQoNCgkJZnVuY3Rpb24gT25DaGFuZ2VVcmwoKSB7DQoJCQl2YXIgZmlsZW5hbWUgPSAkKCIjdXJsX2lkIikudmFsKCkuc3BsaXQoJy8nKS5wb3AoKQ0KCQkJJCgiI3NhdmVfcGF0aF9pZCIpLnZhbChmaWxlbmFtZSkNCgkJfQ0KCTwvc2NyaXB0Pg0KPC9oZWFkPg0KDQo8Ym9keSBvbmxvYWQ9Ik9uTG9hZCgpIj4NCjxkaXYgY2xhc3M9Im5hdmJhciBuYXZiYXItaW52ZXJzZSBuYXZiYXItZml4ZWQtdG9wIj4NCgk8ZGl2IGNsYXNzPSJjb250YWluZXIiPg0KCQk8ZGl2IGNsYXNzPSJuYXZiYXItaGVhZGVyIj4NCgkJCTxidXR0b24gdHlwZT0iYnV0dG9uIiBjbGFzcz0ibmF2YmFyLXRvZ2dsZSIgZGF0YS10b2dnbGU9ImNvbGxhcHNlIiBkYXRhLXRhcmdldD0iLm5hdmJhci1jb2xsYXBzZSI+DQoJCQkJPHNwYW4gY2xhc3M9Imljb24tYmFyIj48L3NwYW4+PHNwYW4gY2xhc3M9Imljb24tYmFyIj48L3NwYW4+PHNwYW4gY2xhc3M9Imljb24tYmFyIj48L3NwYW4+DQoJCQk8L2J1dHRvbj4NCgkJCTxhIGNsYXNzPSJuYXZiYXItYnJhbmQiIGhyZWY9IiMiPkdPIERvd25sb2FkZXI8L2E+DQoJCTwvZGl2Pg0KCQk8ZGl2IGNsYXNzPSJuYXZiYXItY29sbGFwc2UgY29sbGFwc2UiPg0KCQkJPHVsIGNsYXNzPSJuYXYgbmF2YmFyLW5hdiI+DQoJCQkJPGxpIGNsYXNzPSJkcm9wZG93biI+DQoJCQkJCTxhIGhyZWY9IiMiIGNsYXNzPSJkcm9wZG93bi10b2dnbGUiIGRhdGEtdG9nZ2xlPSJkcm9wZG93biI+RmlsZSA8YiBjbGFzcz0iY2FyZXQiPjwvYj48L2E+DQoJCQkJCTx1bCBjbGFzcz0iZHJvcGRvd24tbWVudSI+DQoJCQkJCQk8bGk+DQoJCQkJCQkJPGEgZGF0YS10b2dnbGU9Im1vZGFsIiBkYXRhLXRhcmdldD0iI215TW9kYWwiPkFkZCBkb3dubG9hZDwvYT4NCgkJCQkJCTwvbGk+DQoJCQkJCQk8bGkgb25jbGljaz0iUmVtb3ZlRG93bmxvYWQoKSI+DQoJCQkJCQkJPGEgaHJlZj0iIyI+RGVsZXRlIGRvd25sb2FkPC9hPg0KCQkJCQkJPC9saT4NCgkJCQkJPC91bD4NCgkJCQk8L2xpPg0KCQkJCTxsaSBjbGFzcz0iZHJvcGRvd24iPg0KCQkJCQk8YSBocmVmPSIjIiBjbGFzcz0iZHJvcGRvd24tdG9nZ2xlIiBkYXRhLXRvZ2dsZT0iZHJvcGRvd24iPkFjdGlvbiA8YiBjbGFzcz0iY2FyZXQiPjwvYj48L2E+DQoJCQkJCTx1bCBjbGFzcz0iZHJvcGRvd24tbWVudSI+DQoJCQkJCQk8bGkgb25jbGljaz0iU3RhcnREb3dubG9hZCgpIj4NCgkJCQkJCQk8YSBocmVmPSIjIj5TdGFydDwvYT4NCgkJCQkJCTwvbGk+DQoJCQkJCQk8bGkgb25jbGljaz0iU3RvcERvd25sb2FkKCkiPg0KCQkJCQkJCTxhIGhyZWY9IiMiPlN0b3A8L2E+DQoJCQkJCQk8L2xpPg0KCQkJCQkJPGxpIGNsYXNzPSJkaXZpZGVyIj48L2xpPg0KCQkJCQkJPGxpIG9uY2xpY2s9IlN0YXJ0QWxsRG93bmxvYWQoKSI+DQoJCQkJCQkJPGEgaHJlZj0iIyI+U3RhcnQgYWxsPC9hPg0KCQkJCQkJPC9saT4NCgkJCQkJCTxsaSBvbmNsaWNrPSJTdG9wQWxsRG93bmxvYWQoKSI+DQoJCQkJCQkJPGEgaHJlZj0iIyI+U3RvcCBhbGw8L2E+DQoJCQkJCQk8L2xpPg0KCQkJCQk8L3VsPg0KCQkJCTwvbGk+DQoJCQkJPGxpPg0KCQkJCQk8YSBocmVmPSIjYWJvdXQiPkFib3V0PC9hPg0KCQkJCTwvbGk+DQoJCQk8L3VsPg0KCQk8L2Rpdj4NCgkJPCEtLS8ubmF2YmFyLWNvbGxhcHNlIC0tPg0KCTwvZGl2Pg0KPC9kaXY+DQo8L1A+DQo8dGFibGUgaWQ9ImpxR3JpZCI+PC90YWJsZT4NCg0KPCEtLSBNb2RhbCAtLT4NCjxkaXYgY2xhc3M9Im1vZGFsIGZhZGUiIGlkPSJteU1vZGFsIiByb2xlPSJkaWFsb2ciPg0KCTxkaXYgY2xhc3M9Im1vZGFsLWRpYWxvZyI+DQoNCgkJPCEtLSBNb2RhbCBjb250ZW50LS0+DQoJCTxkaXYgY2xhc3M9Im1vZGFsLWNvbnRlbnQiPg0KCQkJPGRpdiBjbGFzcz0ibW9kYWwtaGVhZGVyIj4NCgkJCQk8YnV0dG9uIHR5cGU9ImJ1dHRvbiIgY2xhc3M9ImNsb3NlIiBkYXRhLWRpc21pc3M9Im1vZGFsIj4mdGltZXM7PC9idXR0b24+DQoJCQkJPGg0IGNsYXNzPSJtb2RhbC10aXRsZSI+RW50ZXIgVXJsPC9oND4NCgkJCTwvZGl2Pg0KCQkJPGRpdiBjbGFzcz0ibW9kYWwtYm9keSI+DQoJCQkJPGRpdiBjbGFzcz0iZm9ybS1ncm91cCI+DQoJCQkJCTxsYWJlbCBjbGFzcz0iY29udHJvbC1sYWJlbCI+VXJsPC9sYWJlbD4NCg0KCQkJCQk8ZGl2IGNsYXNzPSJjb250cm9scyI+DQoJCQkJCQk8aW5wdXQgdHlwZT0idGV4dCIgb25jaGFuZ2U9Ik9uQ2hhbmdlVXJsKCkiIGlkPSJ1cmxfaWQiIGNsYXNzPSJmb3JtLWNvbnRyb2wiIHZhbHVlPSJodHRwOi8vbWlycm9yLnlhbmRleC5ydS91YnVudHUtY2RpbWFnZS9yZWxlYXNlcy8xNS4xMC9hbHBoYS0yL3NvdXJjZS93aWx5LXNyYy0xLmlzbyI+DQoJCQkJCTwvZGl2Pg0KCQkJCQk8bGFiZWwgY2xhc3M9ImNvbnRyb2wtbGFiZWwiPlNhdmUgcGF0aDwvbGFiZWw+DQoNCgkJCQkJPGRpdiBjbGFzcz0iY29udHJvbHMiPg0KCQkJCQkJPGlucHV0IHR5cGU9InRleHQiIGlkPSJzYXZlX3BhdGhfaWQiIGNsYXNzPSJmb3JtLWNvbnRyb2wiIHZhbHVlPSJ3aWx5LXNyYy0xLmlzbyI+DQoJCQkJCTwvZGl2Pg0KCQkJCQk8bGFiZWwgY2xhc3M9ImNvbnRyb2wtbGFiZWwiPlBhcnRzIGNvdW50PC9sYWJlbD4NCgkJCQkJPHNlbGVjdCBjbGFzcz0iZm9ybS1jb250cm9sIiBpZD0icGFydF9jb3VudF9pZCI+DQoJCQkJCQk8b3B0aW9uPjE8L29wdGlvbj4NCgkJCQkJCTxvcHRpb24+Mjwvb3B0aW9uPg0KCQkJCQkJPG9wdGlvbj40PC9vcHRpb24+DQoJCQkJCQk8b3B0aW9uPjg8L29wdGlvbj4NCgkJCQkJCTxvcHRpb24+MTY8L29wdGlvbj4NCgkJCQkJPC9zZWxlY3Q+DQoNCgkJCQkJPGRpdiBjbGFzcz0ibW9kYWwtZm9vdGVyIj4NCgkJCQkJCTxhIGNsYXNzPSJidG4gYnRuLXByaW1hcnkiIG9uY2xpY2s9IkFkZERvd25sb2FkKCkiIGRhdGEtZGlzbWlzcz0ibW9kYWwiPlN0YXJ0IGRvd25sb2FkPC9hPg0KCQkJCQk8L2Rpdj4NCgkJCQk8L2Rpdj4NCgkJCTwvZGl2Pg0KCQk8L2Rpdj4NCg0KCTwvZGl2Pg0KPC9kaXY+DQo8L2JvZHk+DQoNCjwvaHRtbD4="

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
	srv.StartAllTask()
	js, _ := json.Marshal("ok")
	rwr.Write(js)
}

func (srv *DServ) StopAllTask() {
	srv.oplock.Lock()
	defer func() {
		srv.oplock.Unlock()
	}()
	for _, e := range srv.dls {
		log.Println("info stopall result:", e.StopAll())
	}
}

func (srv *DServ) StartAllTask() {
	srv.oplock.Lock()
	defer func() {
		srv.oplock.Unlock()
	}()
	for _, e := range srv.dls {
		log.Println("info start all result:", e.StartAll())
	}
}
func (srv *DServ) stopAllTask(rwr http.ResponseWriter, req *http.Request) {
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

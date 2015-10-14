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
	http.HandleFunc("/", srv.Redirect)
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

const htmlData = "PCFkb2N0eXBlIGh0bWw+DQoNCjxodG1sPg0KDQo8aGVhZD4NCgk8dGl0bGU+R08gRE9XTkxPQUQ8L3RpdGxlPg0KCTxtZXRhIG5hbWU9InZpZXdwb3J0IiBjb250ZW50PSJ3aWR0aD1kZXZpY2Utd2lkdGgiPg0KCTxsaW5rIHJlbD0ic3R5bGVzaGVldCIgaHJlZj0iaHR0cHM6Ly9uZXRkbmEuYm9vdHN0cmFwY2RuLmNvbS9ib290c3dhdGNoLzMuMC4wL2pvdXJuYWwvYm9vdHN0cmFwLm1pbi5jc3MiPg0KCTxsaW5rIHJlbD0ic3R5bGVzaGVldCIgdHlwZT0idGV4dC9jc3MiIG1lZGlhPSJzY3JlZW4iDQoJCSAgaHJlZj0iaHR0cDovL3d3dy5ndXJpZGRvLm5ldC9kZW1vL2Nzcy90cmlyYW5kL3VpLmpxZ3JpZC1ib290c3RyYXAuY3NzIj4NCgk8c2NyaXB0IHR5cGU9InRleHQvamF2YXNjcmlwdCIgc3JjPSJodHRwczovL2FqYXguZ29vZ2xlYXBpcy5jb20vYWpheC9saWJzL2pxdWVyeS8yLjAuMy9qcXVlcnkubWluLmpzIj48L3NjcmlwdD4NCgk8c2NyaXB0IHR5cGU9InRleHQvamF2YXNjcmlwdCIgc3JjPSJodHRwczovL25ldGRuYS5ib290c3RyYXBjZG4uY29tL2Jvb3RzdHJhcC8zLjMuNC9qcy9ib290c3RyYXAubWluLmpzIj48L3NjcmlwdD4NCgk8c2NyaXB0IHR5cGU9InRleHQvamF2YXNjcmlwdCIgc3JjPSJodHRwOi8vd3d3Lmd1cmlkZG8ubmV0L2RlbW8vanMvdHJpcmFuZC9qcXVlcnkuanFHcmlkLm1pbi5qcyI+PC9zY3JpcHQ+DQoJPHNjcmlwdCB0eXBlPSJ0ZXh0L2phdmFzY3JpcHQiIHNyYz0iaHR0cDovL3d3dy5ndXJpZGRvLm5ldC9kZW1vL2pzL3RyaXJhbmQvaTE4bi9ncmlkLmxvY2FsZS1lbi5qcyI+PC9zY3JpcHQ+DQoJPGxpbmsgcmVsPSJzdHlsZXNoZWV0IiBocmVmPSIvL2NvZGUuanF1ZXJ5LmNvbS91aS8xLjExLjQvdGhlbWVzL3Ntb290aG5lc3MvanF1ZXJ5LXVpLmNzcyI+DQoJPHNjcmlwdCBzcmM9Imh0dHA6Ly9jb2RlLmpxdWVyeS5jb20vdWkvMS4xMS40L2pxdWVyeS11aS5qcyI+PC9zY3JpcHQ+DQoJPHN0eWxlIHR5cGU9InRleHQvY3NzIj4NCgkJYm9keSB7DQoJCXBhZGRpbmctdG9wOiA1MHB4Ow0KCQlwYWRkaW5nLWJvdHRvbTogMjBweDsNCgkJfQ0KDQoJCS50YWJsZSAucHJvZ3Jlc3Mgew0KCQltYXJnaW4tYm90dG9tOiAwcHg7DQoJCX0NCgk8L3N0eWxlPg0KCTxzY3JpcHQgdHlwZT0idGV4dC9qYXZhc2NyaXB0Ij4NCgkJZnVuY3Rpb24gVXBkYXRlVGFibGUoKSB7DQoJCQkkKCIjanFHcmlkIikNCgkJCQkuanFHcmlkKHsNCgkJCQkJdXJsOiAnaHR0cDovLycrbG9jYXRpb24uaG9zdG5hbWUrJzo5OTgxL3Byb2dyZXNzLmpzb24nLA0KCQkJCQltdHlwZTogIkdFVCIsDQoJCQkJCWFqYXhTdWJncmlkT3B0aW9uczogew0KCQkJCQkJYXN5bmM6IGZhbHNlDQoJCQkJCX0sDQoJCQkJCXN0eWxlVUk6ICdCb290c3RyYXAnLA0KCQkJCQlkYXRhdHlwZTogImpzb24iLA0KCQkJCQljb2xNb2RlbDogW3sNCgkJCQkJCWxhYmVsOiAnIycsDQoJCQkJCQluYW1lOiAnSWQnLA0KCQkJCQkJa2V5OiB0cnVlLA0KCQkJCQkJd2lkdGg6IDUNCgkJCQkJfSwgew0KCQkJCQkJbGFiZWw6ICdGaWxlIE5hbWUnLA0KCQkJCQkJbmFtZTogJ0ZpbGVOYW1lJywNCgkJCQkJCXdpZHRoOiAxNQ0KCQkJCQl9LCB7DQoJCQkJCQlsYWJlbDogJ1NpemUnLA0KCQkJCQkJbmFtZTogJ1NpemUnLA0KCQkJCQkJd2lkdGg6IDIwLA0KCQkJCQkJZm9ybWF0dGVyOiBGb3JtYXRCeXRlDQoJCQkJCX0sIHsNCgkJCQkJCWxhYmVsOiAnRG93bmxvYWRlZCcsDQoJCQkJCQluYW1lOiAnRG93bmxvYWRlZCcsDQoJCQkJCQl3aWR0aDogMjAsDQoJCQkJCQlmb3JtYXR0ZXI6IEZvcm1hdEJ5dGUNCgkJCQkJfSwgew0KCQkJCQkJbGFiZWw6ICclJywNCgkJCQkJCW5hbWU6ICdQcm9ncmVzcycsDQoJCQkJCQl3aWR0aDogNQ0KCQkJCQl9LCB7DQoJCQkJCQlsYWJlbDogJ1NwZWVkJywNCgkJCQkJCW5hbWU6ICdTcGVlZCcsDQoJCQkJCQl3aWR0aDogMTUsDQoJCQkJCQlmb3JtYXR0ZXI6IEZvcm1hdFNwZWVkQnl0ZQ0KCQkJCQl9LCB7DQoJCQkJCQlsYWJlbDogJ1Byb2dyZXNzJywNCgkJCQkJCW5hbWU6ICdQcm9ncmVzcycsDQoJCQkJCQlmb3JtYXR0ZXI6IEZvcm1hdFByb2dyZXNzQmFyDQoJCQkJCX1dLA0KCQkJCQl2aWV3cmVjb3JkczogdHJ1ZSwNCgkJCQkJcm93TnVtOiAyMCwNCgkJCQkJcGFnZXI6ICIjanFHcmlkUGFnZXIiDQoJCQkJfSk7DQoJCX0NCg0KCQlmdW5jdGlvbiBGaXhUYWJsZSgpIHsNCgkJCSQuZXh0ZW5kKCQuamdyaWQuYWpheE9wdGlvbnMsIHsNCgkJCQlhc3luYzogZmFsc2UNCgkJCX0pDQoJCQkkKCIjanFHcmlkIikNCgkJCQkuc2V0R3JpZFdpZHRoKCQod2luZG93KQ0KCQkJCQkud2lkdGgoKSAtIDUpDQoJCQkkKCIjanFHcmlkIikNCgkJCQkuc2V0R3JpZEhlaWdodCgkKHdpbmRvdykNCgkJCQkJLmhlaWdodCgpKQ0KCQkJJCh3aW5kb3cpDQoJCQkJLmJpbmQoJ3Jlc2l6ZScsIGZ1bmN0aW9uKCkgew0KCQkJCQkkKCIjanFHcmlkIikNCgkJCQkJCS5zZXRHcmlkV2lkdGgoJCh3aW5kb3cpDQoJCQkJCQkJLndpZHRoKCkgLSA1KTsNCgkJCQkJJCgiI2pxR3JpZCIpDQoJCQkJCQkuc2V0R3JpZEhlaWdodCgkKHdpbmRvdykNCgkJCQkJCQkuaGVpZ2h0KCkpDQoJCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIFVwZGF0ZURhdGEoKSB7DQoJCQl2YXIgZ3JpZCA9ICQoIiNqcUdyaWQiKTsNCgkJCXZhciByb3dLZXkgPSBncmlkLmpxR3JpZCgnZ2V0R3JpZFBhcmFtJywgInNlbHJvdyIpOw0KCQkJJCgiI2pxR3JpZCIpLnRyaWdnZXIoInJlbG9hZEdyaWQiKTsNCgkJCWlmKHJvd0tleSkgew0KCQkJCSQoJyNqcUdyaWQnKS5qcUdyaWQoInJlc2V0U2VsZWN0aW9uIikNCgkJCQkkKCcjanFHcmlkJykuanFHcmlkKCdzZXRTZWxlY3Rpb24nLCByb3dLZXkpOw0KCQkJfQ0KCQl9DQoNCgkJZnVuY3Rpb24gRm9ybWF0UHJvZ3Jlc3NCYXIoY2VsbFZhbHVlLCBvcHRpb25zLCByb3dPYmplY3QpIHsNCgkJCXZhciBpbnRWYWwgPSBwYXJzZUludChjZWxsVmFsdWUpOw0KDQoJCQl2YXIgY2VsbEh0bWwgPSAnPGRpdiBjbGFzcz0icHJvZ3Jlc3MiPjxkaXYgY2xhc3M9InByb2dyZXNzLWJhciIgc3R5bGU9IndpZHRoOiAnICsgaW50VmFsICsgJyU7Ij48L2Rpdj48L2Rpdj4nDQoNCgkJCXJldHVybiBjZWxsSHRtbDsNCgkJfQ0KDQoJCWZ1bmN0aW9uIEZvcm1hdEJ5dGUoY2VsbFZhbHVlLCBvcHRpb25zLCByb3dPYmplY3QpIHsNCgkJCXZhciBpbnRWYWwgPSBwYXJzZUludChjZWxsVmFsdWUpOw0KCQkJdmFyIHJhcyA9ICIgQi4iDQoJCQlpZihpbnRWYWwgPiAxMDI0KSB7DQoJCQkJaW50VmFsIC89IDEwMjQNCgkJCQlyYXMgPSAiIEtCLiINCgkJCX0NCgkJCWlmKGludFZhbCA+IDEwMjQpIHsNCgkJCQlpbnRWYWwgLz0gMTAyNA0KCQkJCXJhcyA9ICIgTUIuIg0KCQkJfQ0KCQkJaWYoaW50VmFsID4gMTAyNCkgew0KCQkJCWludFZhbCAvPSAxMDI0DQoJCQkJcmFzID0gIiBHQi4iDQoJCQl9DQoNCgkJCWlmKGludFZhbCA+IDEwMjQpIHsNCgkJCQlpbnRWYWwgLz0gMTAyNA0KCQkJCXJhcyA9ICIgVEIuIg0KCQkJfQ0KCQkJdmFyIGNlbGxIdG1sID0gKGludFZhbCkudG9GaXhlZCgxKSArIHJhczsNCgkJCXJldHVybiBjZWxsSHRtbDsNCgkJfQ0KCQlmdW5jdGlvbiBGb3JtYXRTcGVlZEJ5dGUoY2VsbFZhbHVlLCBvcHRpb25zLCByb3dPYmplY3QpIHsNCgkJCXZhciBpbnRWYWwgPSBwYXJzZUludChjZWxsVmFsdWUpOw0KCQkJdmFyIHJhcyA9ICIgQi9zZWMuIg0KCQkJaWYoaW50VmFsID4gMTAyNCkgew0KCQkJCWludFZhbCAvPSAxMDI0DQoJCQkJcmFzID0gIiBLQi9zZWMuIg0KCQkJfQ0KCQkJaWYoaW50VmFsID4gMTAyNCkgew0KCQkJCWludFZhbCAvPSAxMDI0DQoJCQkJcmFzID0gIiBNQi9zZWMuIg0KCQkJfQ0KCQkJaWYoaW50VmFsID4gMTAyNCkgew0KCQkJCWludFZhbCAvPSAxMDI0DQoJCQkJcmFzID0gIiBHQi9zZWMiDQoJCQl9DQoNCgkJCWlmKGludFZhbCA+IDEwMjQpIHsNCgkJCQlpbnRWYWwgLz0gMTAyNA0KCQkJCXJhcyA9ICIgVEIuIg0KCQkJfQ0KCQkJdmFyIGNlbGxIdG1sID0gKGludFZhbCkudG9GaXhlZCgxKSArIHJhczsNCgkJCXJldHVybiBjZWxsSHRtbDsNCgkJfQ0KDQoJCWZ1bmN0aW9uIE9uTG9hZCgpIHsNCg0KCQkJVXBkYXRlVGFibGUoKQ0KCQkJRml4VGFibGUoKQ0KCQkJc2V0SW50ZXJ2YWwoVXBkYXRlRGF0YSwgNTAwKTsNCgkJfQ0KDQoJCWZ1bmN0aW9uIEFkZERvd25sb2FkKCkgew0KCQkJdmFyIHJlcSA9IHsNCgkJCQlQYXJ0Q291bnQ6IHBhcnNlSW50KCQoIiNwYXJ0X2NvdW50X2lkIikudmFsKCkpLA0KCQkJCUZpbGVQYXRoOiAkKCIjc2F2ZV9wYXRoX2lkIikudmFsKCksDQoJCQkJVXJsOiAkKCIjdXJsX2lkIikudmFsKCkNCgkJCX07DQoJCQkkLmFqYXgoew0KCQkJCQl1cmw6ICIvYWRkX3Rhc2siLA0KCQkJCQl0eXBlOiAiUE9TVCIsDQoJCQkJCWRhdGE6IEpTT04uc3RyaW5naWZ5KHJlcSksDQoJCQkJCWRhdGFUeXBlOiAidGV4dCINCgkJCQl9KQ0KCQkJCS5lcnJvcihmdW5jdGlvbihqc29uRGF0YSkgew0KCQkJCQljb25zb2xlLmxvZyhqc29uRGF0YSkNCgkJCQl9KQ0KCQl9DQoNCgkJZnVuY3Rpb24gUmVtb3ZlRG93bmxvYWQoKSB7DQoJCQl2YXIgZ3JpZCA9ICQoIiNqcUdyaWQiKTsNCgkJCXZhciByb3dLZXkgPSBwYXJzZUludChncmlkLmpxR3JpZCgnZ2V0R3JpZFBhcmFtJywgInNlbHJvdyIpKTsNCgkJCXZhciByZXEgPSByb3dLZXk7DQoJCQkkLmFqYXgoew0KCQkJCQl1cmw6ICIvcmVtb3ZlX3Rhc2siLA0KCQkJCQl0eXBlOiAiUE9TVCIsDQoJCQkJCWRhdGE6IEpTT04uc3RyaW5naWZ5KHJlcSksDQoJCQkJCWRhdGFUeXBlOiAidGV4dCINCgkJCQl9KQ0KCQkJCS5lcnJvcihmdW5jdGlvbihqc29uRGF0YSkgew0KCQkJCQljb25zb2xlLmxvZyhqc29uRGF0YSkNCgkJCQl9KQ0KCQl9DQoNCgkJZnVuY3Rpb24gU3RhcnREb3dubG9hZCgpIHsNCgkJCXZhciBncmlkID0gJCgiI2pxR3JpZCIpOw0KCQkJdmFyIHJvd0tleSA9IHBhcnNlSW50KGdyaWQuanFHcmlkKCdnZXRHcmlkUGFyYW0nLCAic2Vscm93IikpOw0KCQkJdmFyIHJlcSA9IHJvd0tleTsNCgkJCSQuYWpheCh7DQoJCQkJCXVybDogIi9zdGFydF90YXNrIiwNCgkJCQkJdHlwZTogIlBPU1QiLA0KCQkJCQlkYXRhOiBKU09OLnN0cmluZ2lmeShyZXEpLA0KCQkJCQlkYXRhVHlwZTogInRleHQiDQoJCQkJfSkNCgkJCQkuZXJyb3IoZnVuY3Rpb24oanNvbkRhdGEpIHsNCgkJCQkJY29uc29sZS5sb2coanNvbkRhdGEpDQoJCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIFN0b3BEb3dubG9hZCgpIHsNCgkJCXZhciBncmlkID0gJCgiI2pxR3JpZCIpOw0KCQkJdmFyIHJvd0tleSA9IHBhcnNlSW50KGdyaWQuanFHcmlkKCdnZXRHcmlkUGFyYW0nLCAic2Vscm93IikpOw0KCQkJdmFyIHJlcSA9IHJvd0tleTsNCgkJCSQuYWpheCh7DQoJCQkJCXVybDogIi9zdG9wX3Rhc2siLA0KCQkJCQl0eXBlOiAiUE9TVCIsDQoJCQkJCWRhdGE6IEpTT04uc3RyaW5naWZ5KHJlcSksDQoJCQkJCWRhdGFUeXBlOiAidGV4dCINCgkJCQl9KQ0KCQkJCS5lcnJvcihmdW5jdGlvbihqc29uRGF0YSkgew0KCQkJCQljb25zb2xlLmxvZyhqc29uRGF0YSkNCgkJCQl9KQ0KCQl9DQoNCgkJZnVuY3Rpb24gU3RhcnRBbGxEb3dubG9hZCgpIHsNCgkJCSQuYWpheCh7DQoJCQkJCXVybDogIi9zdGFydF9hbGxfdGFzayIsDQoJCQkJCXR5cGU6ICJQT1NUIiwNCgkJCQkJZGF0YVR5cGU6ICJ0ZXh0Ig0KCQkJCX0pDQoJCQkJLmVycm9yKGZ1bmN0aW9uKGpzb25EYXRhKSB7DQoJCQkJCWNvbnNvbGUubG9nKGpzb25EYXRhKQ0KCQkJCX0pDQoJCX0NCg0KCQlmdW5jdGlvbiBTdG9wQWxsRG93bmxvYWQoKSB7DQoJCQkkLmFqYXgoew0KCQkJCQl1cmw6ICIvc3RvcF9hbGxfdGFzayIsDQoJCQkJCXR5cGU6ICJQT1NUIiwNCgkJCQkJZGF0YVR5cGU6ICJ0ZXh0Ig0KCQkJCX0pDQoJCQkJLmVycm9yKGZ1bmN0aW9uKGpzb25EYXRhKSB7DQoJCQkJCWNvbnNvbGUubG9nKGpzb25EYXRhKQ0KCQkJCX0pDQoJCX0NCg0KDQoJCWZ1bmN0aW9uIE9uQ2hhbmdlVXJsKCkgew0KCQkJdmFyIGZpbGVuYW1lID0gJCgiI3VybF9pZCIpLnZhbCgpLnNwbGl0KCcvJykucG9wKCkNCgkJCSQoIiNzYXZlX3BhdGhfaWQiKS52YWwoZmlsZW5hbWUpDQoJCX0NCgk8L3NjcmlwdD4NCjwvaGVhZD4NCg0KPGJvZHkgb25sb2FkPSJPbkxvYWQoKSI+DQo8ZGl2IGNsYXNzPSJuYXZiYXIgbmF2YmFyLWludmVyc2UgbmF2YmFyLWZpeGVkLXRvcCI+DQoJPGRpdiBjbGFzcz0iY29udGFpbmVyIj4NCgkJPGRpdiBjbGFzcz0ibmF2YmFyLWhlYWRlciI+DQoJCQk8YnV0dG9uIHR5cGU9ImJ1dHRvbiIgY2xhc3M9Im5hdmJhci10b2dnbGUiIGRhdGEtdG9nZ2xlPSJjb2xsYXBzZSIgZGF0YS10YXJnZXQ9Ii5uYXZiYXItY29sbGFwc2UiPg0KCQkJCTxzcGFuIGNsYXNzPSJpY29uLWJhciI+PC9zcGFuPjxzcGFuIGNsYXNzPSJpY29uLWJhciI+PC9zcGFuPjxzcGFuIGNsYXNzPSJpY29uLWJhciI+PC9zcGFuPg0KCQkJPC9idXR0b24+DQoJCQk8YSBjbGFzcz0ibmF2YmFyLWJyYW5kIiBocmVmPSIjIj5HTyBEb3dubG9hZGVyPC9hPg0KCQk8L2Rpdj4NCgkJPGRpdiBjbGFzcz0ibmF2YmFyLWNvbGxhcHNlIGNvbGxhcHNlIj4NCgkJCTx1bCBjbGFzcz0ibmF2IG5hdmJhci1uYXYiPg0KCQkJCTxsaSBjbGFzcz0iZHJvcGRvd24iPg0KCQkJCQk8YSBocmVmPSIjIiBjbGFzcz0iZHJvcGRvd24tdG9nZ2xlIiBkYXRhLXRvZ2dsZT0iZHJvcGRvd24iPkZpbGUgPGIgY2xhc3M9ImNhcmV0Ij48L2I+PC9hPg0KCQkJCQk8dWwgY2xhc3M9ImRyb3Bkb3duLW1lbnUiPg0KCQkJCQkJPGxpPg0KCQkJCQkJCTxhIGRhdGEtdG9nZ2xlPSJtb2RhbCIgZGF0YS10YXJnZXQ9IiNteU1vZGFsIj5BZGQgZG93bmxvYWQ8L2E+DQoJCQkJCQk8L2xpPg0KCQkJCQkJPGxpIG9uY2xpY2s9IlJlbW92ZURvd25sb2FkKCkiPg0KCQkJCQkJCTxhIGhyZWY9IiMiPkRlbGV0ZSBkb3dubG9hZDwvYT4NCgkJCQkJCTwvbGk+DQoJCQkJCTwvdWw+DQoJCQkJPC9saT4NCgkJCQk8bGkgY2xhc3M9ImRyb3Bkb3duIj4NCgkJCQkJPGEgaHJlZj0iIyIgY2xhc3M9ImRyb3Bkb3duLXRvZ2dsZSIgZGF0YS10b2dnbGU9ImRyb3Bkb3duIj5BY3Rpb24gPGIgY2xhc3M9ImNhcmV0Ij48L2I+PC9hPg0KCQkJCQk8dWwgY2xhc3M9ImRyb3Bkb3duLW1lbnUiPg0KCQkJCQkJPGxpIG9uY2xpY2s9IlN0YXJ0RG93bmxvYWQoKSI+DQoJCQkJCQkJPGEgaHJlZj0iIyI+U3RhcnQ8L2E+DQoJCQkJCQk8L2xpPg0KCQkJCQkJPGxpIG9uY2xpY2s9IlN0b3BEb3dubG9hZCgpIj4NCgkJCQkJCQk8YSBocmVmPSIjIj5TdG9wPC9hPg0KCQkJCQkJPC9saT4NCgkJCQkJCTxsaSBjbGFzcz0iZGl2aWRlciI+PC9saT4NCgkJCQkJCTxsaSBvbmNsaWNrPSJTdGFydEFsbERvd25sb2FkKCkiPg0KCQkJCQkJCTxhIGhyZWY9IiMiPlN0YXJ0IGFsbDwvYT4NCgkJCQkJCTwvbGk+DQoJCQkJCQk8bGkgb25jbGljaz0iU3RvcEFsbERvd25sb2FkKCkiPg0KCQkJCQkJCTxhIGhyZWY9IiMiPlN0b3AgYWxsPC9hPg0KCQkJCQkJPC9saT4NCgkJCQkJPC91bD4NCgkJCQk8L2xpPg0KCQkJCTxsaT4NCgkJCQkJPGEgaHJlZj0iI2Fib3V0Ij5BYm91dDwvYT4NCgkJCQk8L2xpPg0KCQkJPC91bD4NCgkJPC9kaXY+DQoJCTwhLS0vLm5hdmJhci1jb2xsYXBzZSAtLT4NCgk8L2Rpdj4NCjwvZGl2Pg0KPC9QPg0KPHRhYmxlIGlkPSJqcUdyaWQiPjwvdGFibGU+DQoNCjwhLS0gTW9kYWwgLS0+DQo8ZGl2IGNsYXNzPSJtb2RhbCBmYWRlIiBpZD0ibXlNb2RhbCIgcm9sZT0iZGlhbG9nIj4NCgk8ZGl2IGNsYXNzPSJtb2RhbC1kaWFsb2ciPg0KDQoJCTwhLS0gTW9kYWwgY29udGVudC0tPg0KCQk8ZGl2IGNsYXNzPSJtb2RhbC1jb250ZW50Ij4NCgkJCTxkaXYgY2xhc3M9Im1vZGFsLWhlYWRlciI+DQoJCQkJPGJ1dHRvbiB0eXBlPSJidXR0b24iIGNsYXNzPSJjbG9zZSIgZGF0YS1kaXNtaXNzPSJtb2RhbCI+JnRpbWVzOzwvYnV0dG9uPg0KCQkJCTxoNCBjbGFzcz0ibW9kYWwtdGl0bGUiPkVudGVyIFVybDwvaDQ+DQoJCQk8L2Rpdj4NCgkJCTxkaXYgY2xhc3M9Im1vZGFsLWJvZHkiPg0KCQkJCTxkaXYgY2xhc3M9ImZvcm0tZ3JvdXAiPg0KCQkJCQk8bGFiZWwgY2xhc3M9ImNvbnRyb2wtbGFiZWwiPlVybDwvbGFiZWw+DQoNCgkJCQkJPGRpdiBjbGFzcz0iY29udHJvbHMiPg0KCQkJCQkJPGlucHV0IHR5cGU9InRleHQiIG9uY2hhbmdlPSJPbkNoYW5nZVVybCgpIiBpZD0idXJsX2lkIiBjbGFzcz0iZm9ybS1jb250cm9sIg0KCQkJCQkJCSAgIHZhbHVlPSJodHRwOi8vbWlycm9yLnlhbmRleC5ydS91YnVudHUtY2RpbWFnZS9yZWxlYXNlcy8xNS4xMC9hbHBoYS0yL3NvdXJjZS93aWx5LXNyYy0xLmlzbyI+DQoJCQkJCTwvZGl2Pg0KCQkJCQk8bGFiZWwgY2xhc3M9ImNvbnRyb2wtbGFiZWwiPlNhdmUgcGF0aDwvbGFiZWw+DQoNCgkJCQkJPGRpdiBjbGFzcz0iY29udHJvbHMiPg0KCQkJCQkJPGlucHV0IHR5cGU9InRleHQiIGlkPSJzYXZlX3BhdGhfaWQiIGNsYXNzPSJmb3JtLWNvbnRyb2wiIHZhbHVlPSJ3aWx5LXNyYy0xLmlzbyI+DQoJCQkJCTwvZGl2Pg0KCQkJCQk8bGFiZWwgY2xhc3M9ImNvbnRyb2wtbGFiZWwiPlBhcnRzIGNvdW50PC9sYWJlbD4NCgkJCQkJPHNlbGVjdCBjbGFzcz0iZm9ybS1jb250cm9sIiBpZD0icGFydF9jb3VudF9pZCI+DQoJCQkJCQk8b3B0aW9uPjE8L29wdGlvbj4NCgkJCQkJCTxvcHRpb24+Mjwvb3B0aW9uPg0KCQkJCQkJPG9wdGlvbj40PC9vcHRpb24+DQoJCQkJCQk8b3B0aW9uPjg8L29wdGlvbj4NCgkJCQkJCTxvcHRpb24+MTY8L29wdGlvbj4NCgkJCQkJPC9zZWxlY3Q+DQoNCgkJCQkJPGRpdiBjbGFzcz0ibW9kYWwtZm9vdGVyIj4NCgkJCQkJCTxhIGNsYXNzPSJidG4gYnRuLXByaW1hcnkiIG9uY2xpY2s9IkFkZERvd25sb2FkKCkiIGRhdGEtZGlzbWlzcz0ibW9kYWwiPlN0YXJ0IGRvd25sb2FkPC9hPg0KCQkJCQk8L2Rpdj4NCgkJCQk8L2Rpdj4NCgkJCTwvZGl2Pg0KCQk8L2Rpdj4NCg0KCTwvZGl2Pg0KPC9kaXY+DQo8L2JvZHk+DQoNCjwvaHRtbD4="

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
func (service *DServ) Redirect(responseWriter http.ResponseWriter, request *http.Request) {
	http.Redirect(responseWriter, request, "/index.html", 301)
}

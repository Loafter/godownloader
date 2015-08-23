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

const htmlData = "PCFkb2N0eXBlIGh0bWw+DQoNCjxodG1sPg0KDQo8aGVhZD4NCgk8dGl0bGU+R08gRE9XTkxPQUQ8L3RpdGxlPg0KCTxtZXRhIG5hbWU9InZpZXdwb3J0IiBjb250ZW50PSJ3aWR0aD1kZXZpY2Utd2lkdGgiPg0KCTxsaW5rIHJlbD0ic3R5bGVzaGVldCIgaHJlZj0iaHR0cHM6Ly9uZXRkbmEuYm9vdHN0cmFwY2RuLmNvbS9ib290c3dhdGNoLzMuMC4wL2pvdXJuYWwvYm9vdHN0cmFwLm1pbi5jc3MiPg0KCTxsaW5rIHJlbD0ic3R5bGVzaGVldCIgdHlwZT0idGV4dC9jc3MiIG1lZGlhPSJzY3JlZW4iDQoJCSAgaHJlZj0iaHR0cDovL3d3dy5ndXJpZGRvLm5ldC9kZW1vL2Nzcy90cmlyYW5kL3VpLmpxZ3JpZC1ib290c3RyYXAuY3NzIj4NCgk8c2NyaXB0IHR5cGU9InRleHQvamF2YXNjcmlwdCIgc3JjPSJodHRwczovL2FqYXguZ29vZ2xlYXBpcy5jb20vYWpheC9saWJzL2pxdWVyeS8yLjAuMy9qcXVlcnkubWluLmpzIj48L3NjcmlwdD4NCgk8c2NyaXB0IHR5cGU9InRleHQvamF2YXNjcmlwdCIgc3JjPSJodHRwczovL25ldGRuYS5ib290c3RyYXBjZG4uY29tL2Jvb3RzdHJhcC8zLjMuNC9qcy9ib290c3RyYXAubWluLmpzIj48L3NjcmlwdD4NCgk8c2NyaXB0IHR5cGU9InRleHQvamF2YXNjcmlwdCIgc3JjPSJodHRwOi8vd3d3Lmd1cmlkZG8ubmV0L2RlbW8vanMvdHJpcmFuZC9qcXVlcnkuanFHcmlkLm1pbi5qcyI+PC9zY3JpcHQ+DQoJPHNjcmlwdCB0eXBlPSJ0ZXh0L2phdmFzY3JpcHQiIHNyYz0iaHR0cDovL3d3dy5ndXJpZGRvLm5ldC9kZW1vL2pzL3RyaXJhbmQvaTE4bi9ncmlkLmxvY2FsZS1lbi5qcyI+PC9zY3JpcHQ+DQoJPGxpbmsgcmVsPSJzdHlsZXNoZWV0IiBocmVmPSIvL2NvZGUuanF1ZXJ5LmNvbS91aS8xLjExLjQvdGhlbWVzL3Ntb290aG5lc3MvanF1ZXJ5LXVpLmNzcyI+DQoJPHNjcmlwdCBzcmM9Imh0dHA6Ly9jb2RlLmpxdWVyeS5jb20vdWkvMS4xMS40L2pxdWVyeS11aS5qcyI+PC9zY3JpcHQ+DQoJPHN0eWxlIHR5cGU9InRleHQvY3NzIj4NCgkJYm9keSB7DQoJCXBhZGRpbmctdG9wOiA1MHB4Ow0KCQlwYWRkaW5nLWJvdHRvbTogMjBweDsNCgkJfQ0KDQoJCS50YWJsZSAucHJvZ3Jlc3Mgew0KCQltYXJnaW4tYm90dG9tOiAwcHg7DQoJCX0NCgk8L3N0eWxlPg0KCTxzY3JpcHQgdHlwZT0idGV4dC9qYXZhc2NyaXB0Ij4NCgkJZnVuY3Rpb24gVXBkYXRlVGFibGUoKSB7DQoJCQkkKCIjanFHcmlkIikNCgkJCQkuanFHcmlkKHsNCgkJCQkJdXJsOiAnaHR0cDovL2xvY2FsaG9zdDo5OTgxL3Byb2dyZXNzLmpzb24nLA0KCQkJCQltdHlwZTogIkdFVCIsDQoJCQkJCWFqYXhTdWJncmlkT3B0aW9uczogew0KCQkJCQkJYXN5bmM6IGZhbHNlDQoJCQkJCX0sDQoJCQkJCXN0eWxlVUk6ICdCb290c3RyYXAnLA0KCQkJCQlkYXRhdHlwZTogImpzb24iLA0KCQkJCQljb2xNb2RlbDogW3sNCgkJCQkJCWxhYmVsOiAnIycsDQoJCQkJCQluYW1lOiAnSWQnLA0KCQkJCQkJa2V5OiB0cnVlLA0KCQkJCQkJd2lkdGg6IDUNCgkJCQkJfSwgew0KCQkJCQkJbGFiZWw6ICdGaWxlIE5hbWUnLA0KCQkJCQkJbmFtZTogJ0ZpbGVOYW1lJywNCgkJCQkJCXdpZHRoOiAxNQ0KCQkJCQl9LCB7DQoJCQkJCQlsYWJlbDogJ1NpemUnLA0KCQkJCQkJbmFtZTogJ1NpemUnLA0KCQkJCQkJd2lkdGg6IDIwLA0KCQkJCQkJZm9ybWF0dGVyOiBGb3JtYXRCeXRlDQoJCQkJCX0sIHsNCgkJCQkJCWxhYmVsOiAnRG93bmxvYWRlZCcsDQoJCQkJCQluYW1lOiAnRG93bmxvYWRlZCcsDQoJCQkJCQl3aWR0aDogMjAsDQoJCQkJCQlmb3JtYXR0ZXI6IEZvcm1hdEJ5dGUNCgkJCQkJfSwgew0KCQkJCQkJbGFiZWw6ICclJywNCgkJCQkJCW5hbWU6ICdQcm9ncmVzcycsDQoJCQkJCQl3aWR0aDogNQ0KCQkJCQl9LCB7DQoJCQkJCQlsYWJlbDogJ1NwZWVkJywNCgkJCQkJCW5hbWU6ICdTcGVlZCcsDQoJCQkJCQl3aWR0aDogMTUsDQoJCQkJCQlmb3JtYXR0ZXI6IEZvcm1hdEJ5dGUNCgkJCQkJfSwgew0KCQkJCQkJbGFiZWw6ICdQcm9ncmVzcycsDQoJCQkJCQluYW1lOiAnUHJvZ3Jlc3MnLA0KCQkJCQkJZm9ybWF0dGVyOiBGb3JtYXRQcm9ncmVzc0Jhcg0KCQkJCQl9XSwNCgkJCQkJdmlld3JlY29yZHM6IHRydWUsDQoJCQkJCXJvd051bTogMjAsDQoJCQkJCXBhZ2VyOiAiI2pxR3JpZFBhZ2VyIg0KCQkJCX0pOw0KCQl9DQoNCgkJZnVuY3Rpb24gRml4VGFibGUoKSB7DQoJCQkkLmV4dGVuZCgkLmpncmlkLmFqYXhPcHRpb25zLCB7DQoJCQkJYXN5bmM6IGZhbHNlDQoJCQl9KQ0KCQkJJCgiI2pxR3JpZCIpDQoJCQkJLnNldEdyaWRXaWR0aCgkKHdpbmRvdykNCgkJCQkJLndpZHRoKCkgLSA1KQ0KCQkJJCgiI2pxR3JpZCIpDQoJCQkJLnNldEdyaWRIZWlnaHQoJCh3aW5kb3cpDQoJCQkJCS5oZWlnaHQoKSkNCgkJCSQod2luZG93KQ0KCQkJCS5iaW5kKCdyZXNpemUnLCBmdW5jdGlvbigpIHsNCgkJCQkJJCgiI2pxR3JpZCIpDQoJCQkJCQkuc2V0R3JpZFdpZHRoKCQod2luZG93KQ0KCQkJCQkJCS53aWR0aCgpIC0gNSk7DQoJCQkJCSQoIiNqcUdyaWQiKQ0KCQkJCQkJLnNldEdyaWRIZWlnaHQoJCh3aW5kb3cpDQoJCQkJCQkJLmhlaWdodCgpKQ0KCQkJCX0pDQoJCX0NCg0KCQlmdW5jdGlvbiBVcGRhdGVEYXRhKCkgew0KCQkJdmFyIGdyaWQgPSAkKCIjanFHcmlkIik7DQoJCQl2YXIgcm93S2V5ID0gZ3JpZC5qcUdyaWQoJ2dldEdyaWRQYXJhbScsICJzZWxyb3ciKTsNCgkJCSQoIiNqcUdyaWQiKS50cmlnZ2VyKCJyZWxvYWRHcmlkIik7DQoJCQlpZihyb3dLZXkpIHsNCgkJCQkkKCcjanFHcmlkJykuanFHcmlkKCJyZXNldFNlbGVjdGlvbiIpDQoJCQkJJCgnI2pxR3JpZCcpLmpxR3JpZCgnc2V0U2VsZWN0aW9uJywgcm93S2V5KTsNCgkJCX0NCgkJfQ0KDQoJCWZ1bmN0aW9uIEZvcm1hdFByb2dyZXNzQmFyKGNlbGxWYWx1ZSwgb3B0aW9ucywgcm93T2JqZWN0KSB7DQoJCQl2YXIgaW50VmFsID0gcGFyc2VJbnQoY2VsbFZhbHVlKTsNCg0KCQkJdmFyIGNlbGxIdG1sID0gJzxkaXYgY2xhc3M9InByb2dyZXNzIj48ZGl2IGNsYXNzPSJwcm9ncmVzcy1iYXIiIHN0eWxlPSJ3aWR0aDogJyArIGludFZhbCArICclOyI+PC9kaXY+PC9kaXY+Jw0KDQoJCQlyZXR1cm4gY2VsbEh0bWw7DQoJCX0NCg0KCQlmdW5jdGlvbiBGb3JtYXRCeXRlKGNlbGxWYWx1ZSwgb3B0aW9ucywgcm93T2JqZWN0KSB7DQoJCQl2YXIgaW50VmFsID0gcGFyc2VJbnQoY2VsbFZhbHVlKTsNCgkJCXZhciByYXM9IiBiLiINCgkJCWlmIChpbnRWYWw+MTAyNCl7DQoJCQkJaW50VmFsLz0xMDI0DQoJCQkJcmFzPSIga2IuIg0KCQkJfQ0KCQkJaWYgKGludFZhbD4xMDI0KXsNCgkJCQlpbnRWYWwvPTEwMjQNCgkJCQlyYXM9IiBNYi4iDQoJCQl9DQoJCQlpZiAoaW50VmFsPjEwMjQpew0KCQkJCWludFZhbC89MTAyNA0KCQkJCXJhcz0iIEdiLiINCgkJCX0NCg0KCQkJaWYgKGludFZhbD4xMDI0KXsNCgkJCQlpbnRWYWwvPTEwMjQNCgkJCQlyYXM9IiBUYi4iDQoJCQl9DQoJCQl2YXIgY2VsbEh0bWwgPSAoaW50VmFsKS50b0ZpeGVkKDEpK3JhczsNCgkJCXJldHVybiBjZWxsSHRtbDsNCgkJfQ0KDQoJCWZ1bmN0aW9uIE9uTG9hZCgpIHsNCg0KCQkJVXBkYXRlVGFibGUoKQ0KCQkJRml4VGFibGUoKQ0KCQkJc2V0SW50ZXJ2YWwoVXBkYXRlRGF0YSwgNTAwKTsNCgkJfQ0KDQoJCWZ1bmN0aW9uIEFkZERvd25sb2FkKCkgew0KCQkJdmFyIHJlcSA9IHsNCgkJCQlQYXJ0Q291bnQ6IHBhcnNlSW50KCQoIiNwYXJ0X2NvdW50X2lkIikudmFsKCkpLA0KCQkJCUZpbGVQYXRoOiAkKCIjc2F2ZV9wYXRoX2lkIikudmFsKCksDQoJCQkJVXJsOiAkKCIjdXJsX2lkIikudmFsKCkNCgkJCX07DQoJCQkkLmFqYXgoew0KCQkJCQl1cmw6ICIvYWRkX3Rhc2siLA0KCQkJCQl0eXBlOiAiUE9TVCIsDQoJCQkJCWRhdGE6IEpTT04uc3RyaW5naWZ5KHJlcSksDQoJCQkJCWRhdGFUeXBlOiAidGV4dCINCgkJCQl9KQ0KCQkJCS5lcnJvcihmdW5jdGlvbihqc29uRGF0YSkgew0KCQkJCQljb25zb2xlLmxvZyhqc29uRGF0YSkNCgkJCQl9KQ0KCQl9DQoNCgkJZnVuY3Rpb24gUmVtb3ZlRG93bmxvYWQoKSB7DQoJCQl2YXIgZ3JpZCA9ICQoIiNqcUdyaWQiKTsNCgkJCXZhciByb3dLZXkgPSBwYXJzZUludChncmlkLmpxR3JpZCgnZ2V0R3JpZFBhcmFtJywgInNlbHJvdyIpKTsNCgkJCXZhciByZXEgPSByb3dLZXk7DQoJCQkkLmFqYXgoew0KCQkJCQl1cmw6ICIvcmVtb3ZlX3Rhc2siLA0KCQkJCQl0eXBlOiAiUE9TVCIsDQoJCQkJCWRhdGE6IEpTT04uc3RyaW5naWZ5KHJlcSksDQoJCQkJCWRhdGFUeXBlOiAidGV4dCINCgkJCQl9KQ0KCQkJCS5lcnJvcihmdW5jdGlvbihqc29uRGF0YSkgew0KCQkJCQljb25zb2xlLmxvZyhqc29uRGF0YSkNCgkJCQl9KQ0KCQl9DQoNCgkJZnVuY3Rpb24gU3RhcnREb3dubG9hZCgpIHsNCgkJCXZhciBncmlkID0gJCgiI2pxR3JpZCIpOw0KCQkJdmFyIHJvd0tleSA9IHBhcnNlSW50KGdyaWQuanFHcmlkKCdnZXRHcmlkUGFyYW0nLCAic2Vscm93IikpOw0KCQkJdmFyIHJlcSA9IHJvd0tleTsNCgkJCSQuYWpheCh7DQoJCQkJCXVybDogIi9zdGFydF90YXNrIiwNCgkJCQkJdHlwZTogIlBPU1QiLA0KCQkJCQlkYXRhOiBKU09OLnN0cmluZ2lmeShyZXEpLA0KCQkJCQlkYXRhVHlwZTogInRleHQiDQoJCQkJfSkNCgkJCQkuZXJyb3IoZnVuY3Rpb24oanNvbkRhdGEpIHsNCgkJCQkJY29uc29sZS5sb2coanNvbkRhdGEpDQoJCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIFN0YXJ0QWxsRG93bmxvYWQoKSB7DQoJCQkkLmFqYXgoew0KCQkJCQl1cmw6ICIvc3RhcnRfYWxsX3Rhc2siLA0KCQkJCQl0eXBlOiAiUE9TVCIsDQoJCQkJCWRhdGFUeXBlOiAidGV4dCINCgkJCQl9KQ0KCQkJCS5lcnJvcihmdW5jdGlvbihqc29uRGF0YSkgew0KCQkJCQljb25zb2xlLmxvZyhqc29uRGF0YSkNCgkJCQl9KQ0KCQl9DQoNCgkJZnVuY3Rpb24gU3RvcEFsbERvd25sb2FkKCkgew0KCQkJJC5hamF4KHsNCgkJCQkJdXJsOiAiL3N0b3BfYWxsX3Rhc2siLA0KCQkJCQl0eXBlOiAiUE9TVCIsDQoJCQkJCWRhdGFUeXBlOiAidGV4dCINCgkJCQl9KQ0KCQkJCS5lcnJvcihmdW5jdGlvbihqc29uRGF0YSkgew0KCQkJCQljb25zb2xlLmxvZyhqc29uRGF0YSkNCgkJCQl9KQ0KCQl9DQoNCg0KCQlmdW5jdGlvbiBPbkNoYW5nZVVybCgpIHsNCgkJCXZhciBmaWxlbmFtZSA9ICQoIiN1cmxfaWQiKS52YWwoKS5zcGxpdCgnLycpLnBvcCgpDQoJCQkkKCIjc2F2ZV9wYXRoX2lkIikudmFsKGZpbGVuYW1lKQ0KCQl9DQoJPC9zY3JpcHQ+DQo8L2hlYWQ+DQoNCjxib2R5IG9ubG9hZD0iT25Mb2FkKCkiPg0KPGRpdiBjbGFzcz0ibmF2YmFyIG5hdmJhci1pbnZlcnNlIG5hdmJhci1maXhlZC10b3AiPg0KCTxkaXYgY2xhc3M9ImNvbnRhaW5lciI+DQoJCTxkaXYgY2xhc3M9Im5hdmJhci1oZWFkZXIiPg0KCQkJPGJ1dHRvbiB0eXBlPSJidXR0b24iIGNsYXNzPSJuYXZiYXItdG9nZ2xlIiBkYXRhLXRvZ2dsZT0iY29sbGFwc2UiIGRhdGEtdGFyZ2V0PSIubmF2YmFyLWNvbGxhcHNlIj4NCgkJCQk8c3BhbiBjbGFzcz0iaWNvbi1iYXIiPjwvc3Bhbj48c3BhbiBjbGFzcz0iaWNvbi1iYXIiPjwvc3Bhbj48c3BhbiBjbGFzcz0iaWNvbi1iYXIiPjwvc3Bhbj4NCgkJCTwvYnV0dG9uPg0KCQkJPGEgY2xhc3M9Im5hdmJhci1icmFuZCIgaHJlZj0iIyI+R08gRG93bmxvYWRlcjwvYT4NCgkJPC9kaXY+DQoJCTxkaXYgY2xhc3M9Im5hdmJhci1jb2xsYXBzZSBjb2xsYXBzZSI+DQoJCQk8dWwgY2xhc3M9Im5hdiBuYXZiYXItbmF2Ij4NCgkJCQk8bGkgY2xhc3M9ImRyb3Bkb3duIj4NCgkJCQkJPGEgaHJlZj0iIyIgY2xhc3M9ImRyb3Bkb3duLXRvZ2dsZSIgZGF0YS10b2dnbGU9ImRyb3Bkb3duIj5GaWxlIDxiIGNsYXNzPSJjYXJldCI+PC9iPjwvYT4NCgkJCQkJPHVsIGNsYXNzPSJkcm9wZG93bi1tZW51Ij4NCgkJCQkJCTxsaT4NCgkJCQkJCQk8YSBkYXRhLXRvZ2dsZT0ibW9kYWwiIGRhdGEtdGFyZ2V0PSIjbXlNb2RhbCI+QWRkIGRvd25sb2FkPC9hPg0KCQkJCQkJPC9saT4NCgkJCQkJCTxsaSBvbmNsaWNrPSJSZW1vdmVEb3dubG9hZCgpIj4NCgkJCQkJCQk8YSBocmVmPSIjIj5EZWxldGUgZG93bmxvYWQ8L2E+DQoJCQkJCQk8L2xpPg0KCQkJCQk8L3VsPg0KCQkJCTwvbGk+DQoJCQkJPGxpIGNsYXNzPSJkcm9wZG93biI+DQoJCQkJCTxhIGhyZWY9IiMiIGNsYXNzPSJkcm9wZG93bi10b2dnbGUiIGRhdGEtdG9nZ2xlPSJkcm9wZG93biI+QWN0aW9uIDxiIGNsYXNzPSJjYXJldCI+PC9iPjwvYT4NCgkJCQkJPHVsIGNsYXNzPSJkcm9wZG93bi1tZW51Ij4NCgkJCQkJCTxsaSBvbmNsaWNrPSJTdGFydERvd25sb2FkKCkiPg0KCQkJCQkJCTxhIGhyZWY9IiMiPlN0YXJ0PC9hPg0KCQkJCQkJPC9saT4NCgkJCQkJCTxsaSBvbmNsaWNrPSJTdG9wRG93bmxvYWQoKSI+DQoJCQkJCQkJPGEgaHJlZj0iIyI+U3RvcDwvYT4NCgkJCQkJCTwvbGk+DQoJCQkJCQk8bGkgY2xhc3M9ImRpdmlkZXIiPjwvbGk+DQoJCQkJCQk8bGkgb25jbGljaz0iU3RhcnRBbGxEb3dubG9hZCgpIj4NCgkJCQkJCQk8YSBocmVmPSIjIj5TdGFydCBhbGw8L2E+DQoJCQkJCQk8L2xpPg0KCQkJCQkJPGxpIG9uY2xpY2s9IlN0b3BBbGxEb3dubG9hZCgpIj4NCgkJCQkJCQk8YSBocmVmPSIjIj5TdG9wIGFsbDwvYT4NCgkJCQkJCTwvbGk+DQoJCQkJCTwvdWw+DQoJCQkJPC9saT4NCgkJCQk8bGk+DQoJCQkJCTxhIGhyZWY9IiNhYm91dCI+QWJvdXQ8L2E+DQoJCQkJPC9saT4NCgkJCTwvdWw+DQoJCTwvZGl2Pg0KCQk8IS0tLy5uYXZiYXItY29sbGFwc2UgLS0+DQoJPC9kaXY+DQo8L2Rpdj4NCjwvUD4NCjx0YWJsZSBpZD0ianFHcmlkIj48L3RhYmxlPg0KDQo8IS0tIE1vZGFsIC0tPg0KPGRpdiBjbGFzcz0ibW9kYWwgZmFkZSIgaWQ9Im15TW9kYWwiIHJvbGU9ImRpYWxvZyI+DQoJPGRpdiBjbGFzcz0ibW9kYWwtZGlhbG9nIj4NCg0KCQk8IS0tIE1vZGFsIGNvbnRlbnQtLT4NCgkJPGRpdiBjbGFzcz0ibW9kYWwtY29udGVudCI+DQoJCQk8ZGl2IGNsYXNzPSJtb2RhbC1oZWFkZXIiPg0KCQkJCTxidXR0b24gdHlwZT0iYnV0dG9uIiBjbGFzcz0iY2xvc2UiIGRhdGEtZGlzbWlzcz0ibW9kYWwiPiZ0aW1lczs8L2J1dHRvbj4NCgkJCQk8aDQgY2xhc3M9Im1vZGFsLXRpdGxlIj5FbnRlciBVcmw8L2g0Pg0KCQkJPC9kaXY+DQoJCQk8ZGl2IGNsYXNzPSJtb2RhbC1ib2R5Ij4NCgkJCQk8ZGl2IGNsYXNzPSJmb3JtLWdyb3VwIj4NCgkJCQkJPGxhYmVsIGNsYXNzPSJjb250cm9sLWxhYmVsIj5Vcmw8L2xhYmVsPg0KDQoJCQkJCTxkaXYgY2xhc3M9ImNvbnRyb2xzIj4NCgkJCQkJCTxpbnB1dCB0eXBlPSJ0ZXh0IiBvbmNoYW5nZT0iT25DaGFuZ2VVcmwoKSIgaWQ9InVybF9pZCIgY2xhc3M9ImZvcm0tY29udHJvbCINCgkJCQkJCQkgICB2YWx1ZT0iaHR0cDovL2VkaXN0cmlidXRpb24uaGVhbHRoLmdlLmNvbS9wdWIvZGlzdHJvLzIwMTQxMTI0MTQwMTA2MC5HRUhDXzIwNTYyNzUtMDI1LklTTyI+DQoJCQkJCTwvZGl2Pg0KCQkJCQk8bGFiZWwgY2xhc3M9ImNvbnRyb2wtbGFiZWwiPlNhdmUgcGF0aDwvbGFiZWw+DQoNCgkJCQkJPGRpdiBjbGFzcz0iY29udHJvbHMiPg0KCQkJCQkJPGlucHV0IHR5cGU9InRleHQiIGlkPSJzYXZlX3BhdGhfaWQiIGNsYXNzPSJmb3JtLWNvbnRyb2wiDQoJCQkJCQkJICAgdmFsdWU9IjIwMTQxMTI0MTQwMTA2MC5HRUhDXzIwNTYyNzUtMDI1LklTTyI+DQoJCQkJCTwvZGl2Pg0KCQkJCQk8bGFiZWwgY2xhc3M9ImNvbnRyb2wtbGFiZWwiPlBhcnRzIGNvdW50PC9sYWJlbD4NCgkJCQkJPHNlbGVjdCBjbGFzcz0iZm9ybS1jb250cm9sIiBpZD0icGFydF9jb3VudF9pZCI+DQoJCQkJCQk8b3B0aW9uPjE8L29wdGlvbj4NCgkJCQkJCTxvcHRpb24+Mjwvb3B0aW9uPg0KCQkJCQkJPG9wdGlvbj40PC9vcHRpb24+DQoJCQkJCQk8b3B0aW9uPjg8L29wdGlvbj4NCgkJCQkJCTxvcHRpb24+MTY8L29wdGlvbj4NCgkJCQkJPC9zZWxlY3Q+DQoNCgkJCQkJPGRpdiBjbGFzcz0ibW9kYWwtZm9vdGVyIj4NCgkJCQkJCTxhIGNsYXNzPSJidG4gYnRuLXByaW1hcnkiIG9uY2xpY2s9IkFkZERvd25sb2FkKCkiIGRhdGEtZGlzbWlzcz0ibW9kYWwiPlN0YXJ0IGRvd25sb2FkPC9hPg0KCQkJCQk8L2Rpdj4NCgkJCQk8L2Rpdj4NCgkJCTwvZGl2Pg0KCQk8L2Rpdj4NCg0KCTwvZGl2Pg0KPC9kaXY+DQo8L2JvZHk+DQoNCjwvaHRtbD4="

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

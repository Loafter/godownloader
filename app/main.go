package main

import (
	"godownloader/service"
	"log"
	"os"
	"os/signal"
	"os/user"
	"runtime"
	"strconv"
	"syscall"
)

func getSetPath() string {
	usr, _ := user.Current()
	st := strconv.QuoteRune(os.PathSeparator)
	st = st[1 : len(st)-1]
	return usr.HomeDir + st + ".godownload"
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	gdownsrv := new(DownloadService.DServ)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		func() {
			gdownsrv.StopAllTask()
			log.Println("info: save setting ", gdownsrv.SaveSettings(getSetPath()))
		}()
		os.Exit(1)
	}()
	gdownsrv.LoadSettings(getSetPath())
	log.Println("GUI located add http://localhost:9981/index.html")
	log.Println(gdownsrv.Start(9981))
}

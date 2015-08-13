package main

import (
	"godownloader/service"
	"log"
)

func main() {
	EaSrvCmp := new(DownloadService.DServ)
	log.Println("GUI located add http://localhost:9981/index.html")
	log.Println(EaSrvCmp.Start(9981))
}


package main
import (

	"sync"
	"os"
	"math/rand"
	"runtime"
	"log"
)
func WriteTo(f *os.File,d []byte,syn *sync.WaitGroup){
	syn.Add(1)
	defer syn.Done()
	_,err:=f.WriteAt(d,rand.Int63n(1024))
	if(err!=nil){
		log.Println(err)
		return
	}
}
func main() {
	runtime.GOMAXPROCS(8)

	syn:=new(sync.WaitGroup)
	testf,_:=os.Create("/home/andrew/Desktop/tes.txt")
	defer testf.Close()
	testf.Truncate(1024*4)
	d:=[]byte("Test Test")
	for e:=0;e<1060;e++{
		go WriteTo(testf,d,syn)

	}
	syn.Wait()

	testf.Sync()
log.Println("exit")

}

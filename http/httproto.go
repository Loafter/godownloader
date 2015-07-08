package http
import "net"
func TestMultithread(url string) (bool, error) {
	return false, nil
}

func GetSize(adprot string, url string) (uint, error) {
	conn, err := net.Dial("tcp", adprot)
	defer func() {
		conn.Close()
	}()
	if err!=nil {
		return 0, err
	}
	conn.Write()
	return 0, nil
}
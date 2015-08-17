package iotools

import (
	"os"
	"sync"
)

type SafeFile struct {
	*os.File
	lock sync.Mutex
}

func (sf *SafeFile) WriteAt(b []byte, off int64) (n int, err error) {
	sf.lock.Lock()
	defer sf.lock.Unlock()
	return sf.File.WriteAt(b, off)
}

func (sf *SafeFile) Sync() error {
	sf.lock.Lock()
	defer sf.lock.Unlock()
	return sf.File.Sync()
}
func OpenSafeFile(name string) (file *SafeFile, err error) {
	f, err := os.OpenFile(name, os.O_RDWR, 0666)
	return &SafeFile{File: f}, err
}

func CreateSafeFile(name string) (file *SafeFile, err error) {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	return &SafeFile{File: f}, err
}

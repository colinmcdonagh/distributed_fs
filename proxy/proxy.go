package proxy

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type proxy struct {
	controllerAddr string
}

func (p *proxy) fileAddr(file string) string {
	return fmt.Sprintf("%s/%s", p.controllerAddr, file)
}

func (p *proxy) download(file string) *os.File {
	resp, err := http.Get(fmt.Sprintf(p.fileAddr(file)))
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	fileBytes, _ := ioutil.ReadAll(resp.Body)
	filePtr, err := os.Create(file)
	if err != nil {
		// handle error
	}
	filePtr.Write(fileBytes)
	return filePtr
}

func (p *proxy) upload(localfile, remotefile string) {

	var r io.Reader
	r, _ = os.Open(localfile)
	resp, err := http.Post(p.fileAddr(remotefile), "text/plain", r)

}

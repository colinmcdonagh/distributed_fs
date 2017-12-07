package proxy

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// TODO: get download integrated into this file.

type proxy struct {
	controllerAddr string
}

func (p *proxy) fileAddr(file string) string {
	return fmt.Sprintf("%s/%s", p.controllerAddr, file)
}

func (p *proxy) download(path string) (string, error) {
	// check index
	if serverInt, exists := filePath_serverInt[path]; exists {
		fileResp, err := http.Get(fmt.Sprintf("%s%s",
			fileServerAddrs[serverInt],
			path))
		if err != nil {
			return "", err
		}
		defer fileResp.Body.Close()
		file, _ := ioutil.ReadAll(fileResp.Body)
		return string(file), nil
	}
	return "", fmt.Errorf("file %s doesn`t exist\n", path)
}

func (p *proxy) upload(localfile, remotefile string) {

	var r io.Reader
	r, _ = os.Open(localfile)
	resp, err := http.Post(p.fileAddr(remotefile), "text/plain", r)

}

package proxy

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Specification
/*

 */

// Assumptions
/*

 the organiser (directory service) address is specified as an address, and not
 an http request, i.e. as www.example.com and not http://www.example.com

 */

type Proxy struct {
	organiserAddr string
}

func New(orgAddr string) Proxy {
	return Proxy{
		organiserAddr: orgAddr,
	}
}

func (p *Proxy) OrganiserAddr() string {
	return p.organiserAddr
}

func (p *Proxy) Download(path string) (string, error) {

	// get download server address from organiser.
	resp, err := http.Get(fmt.Sprintf(`http://%s/%s`, p.organiserAddr, path))
	if err != nil {
		return "", fmt.Errorf("error encountered trying to retrieve download server addr: %s",
			err)
	}
	defer resp.Body.Close()

	dLServerAddrAndLf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error encountered reading response from organiser: %s", err)
	}

	splitMsg := strings.Split(string(dLServerAddrAndLf), ",")
	dLServerAddr := splitMsg[0]
	lF := splitMsg[1] 
	// local filename

	// make request to file server itself.
	resp, err = http.Get(fmt.Sprintf(`http://%s/%s`, dLServerAddr, lF))
	if err != nil {
		return "", fmt.Errorf("error encountered requesting file %s from "+
			"server %s: %s", lF, dLServerAddr, err)
	}
	defer resp.Body.Close()

	fileContents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error encountered reading response from " +
			"file server %s: %s", dLServerAddr, err)
	}

	return string(fileContents), nil
}

func (p *Proxy) Upload(localFilePath, officialFilePath, serverLocalFileName string) error {

	requestSuffix := ""
	if serverLocalFileName != "" {
		requestSuffix = fmt.Sprintf("&local_filepath=%s", serverLocalFileName)
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/%s?upload%s", 
		p.organiserAddr, officialFilePath, requestSuffix))
	if err != nil {
		return fmt.Errorf("error encountered retrieving upload server addr:%s",
			err)
	}
	defer resp.Body.Close()

	dlServerAddrAndLf, _ := ioutil.ReadAll(resp.Body)

	splitMsg := strings.Split(string(dlServerAddrAndLf), ",")
	dLServerAddr := splitMsg[0]
	lF := splitMsg[1]
	// local filename

	// make request to the upload server itself
	var r io.Reader
	r, _ = os.Open(localFilePath)

	resp, err = http.Post(fmt.Sprintf(`http://%s/%s?upload&local_filepath=%s`,
	 	p.organiserAddr, dLServerAddr, lF), "text/plain", r)

	uploadResp, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("response from upload server was: %s", uploadResp)

	return nil
}

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

 can determine file server responded with an error if there's no ',' comma
 in the response, since it's used in order to separate returned values.

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

	serverMsgBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error encountered reading response from organiser: %s", err)
	}
	serverMsg := string(serverMsgBytes)

	// could use regex matching of msgs here instead, but the basic
	// determination is if there's a comma in the returned msg.
	if !strings.Contains(serverMsg, ",") {
		return "", fmt.Errorf("file server returned error msg: %s", serverMsg)
	}

	dLServerAddrAndLf := strings.Split(serverMsg, ",")
	dLServerAddr := dLServerAddrAndLf[0]
	lF := dLServerAddrAndLf[1]
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
		return "", fmt.Errorf("error encountered reading response from "+
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
		return fmt.Errorf("error encountered retrieving upload server addr: %s",
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
	r, err = os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("error opening local file: %s", err)
	}

	resp, err = http.Post(fmt.Sprintf("http://%s/%s",
		dLServerAddr, lF), "text/plain", r)
	if err != nil {
		return fmt.Errorf("error encountered trying to upload: %s", err)
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error encountered trying to read server response "+
			"when uploading: %s", err)
	}

	return nil
}

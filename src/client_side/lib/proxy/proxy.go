package proxy

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Assumptions
/*
 can determine file server responded with an error if there's no ',' comma
 in the response, since it's used in order to separate returned values.
*/

// Proxy interacts with the distributed file system services on behalf of client
// applications and command line interfaces.
type Proxy struct {
	directorySrvAddr string
	lockSrvAddr      string
	cacheDir         string
}

// New returns a proxy with specified parameters.
func New(directorySrvAddr, lockSrvAddr, cacheDir string) Proxy {
	// create cache directory if it doesn't already exist.
	if cacheDir != "" {
		if _, err := os.Stat(cacheDir); !os.IsNotExist(err) {
			os.Mkdir(cacheDir, os.ModePerm)
		}
	}

	return Proxy{
		directorySrvAddr: directorySrvAddr,
		lockSrvAddr:      lockSrvAddr,
		cacheDir:         cacheDir,
	}
}

// Lock locks a specific global file.
func (p *Proxy) Lock(filePath string) error {
	if p.lockSrvAddr == "" {
		return fmt.Errorf("no lock server defined in proxy")
	}

	resp, err := http.Get(fmt.Sprintf(`http://%s/%s`, p.lockSrvAddr, filePath))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// was having an issue detecting the server's http error response when already
	// locked, so instead just write back "0" as the response.
	success, _ := ioutil.ReadAll(resp.Body)
	if string(success) != "1" {
		return fmt.Errorf("lock could not be taken")
	}
	return nil
}

// Unlock unlocks a specific global file.
func (p *Proxy) Unlock(filePath string) error {
	if p.lockSrvAddr == "" {
		return fmt.Errorf("no lock server defined in proxy")
	}

	resp, err := http.Get(fmt.Sprintf(`http://%s/%s?unlock`, p.lockSrvAddr, filePath))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// will always be able to unlock.
	_, _ = ioutil.ReadAll(resp.Body)
	return nil
}

func createFile(filename string, fileContents []byte) error {
	dir := filepath.Dir(filename)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	lF, err := os.Create(filename)
	if err != nil {
		return err
	}

	_, err = lF.Write(fileContents)
	if err != nil {
		return err
	}
	return nil
}

// Download downloads a global file from a file server, but only saves it in
// memory as opposed to in a file.
// Therefore, there is less requirement for cleaning up downloaded files where
// persitence isn't required. (for example using the `cat` command)
func (p *Proxy) Download(filePath string) (string, error) {
	if p.directorySrvAddr == "" {
		return "", fmt.Errorf("direcotry server address not specified")
	}

	// get download file server address and local file name from directory service.
	resp, err := http.Get(fmt.Sprintf(`http://%s/%s`, p.directorySrvAddr, filePath))
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

	// could use regexes here instead, but the basic criteria for
	// determining the server address and local file name is the presence of a comma.
	if !strings.Contains(serverMsg, ",") {
		return "", fmt.Errorf("file server returned error msg: %s", serverMsg)
	}

	dLServerAddrAndLf := strings.Split(serverMsg, ",")
	dLServerAddr := dLServerAddrAndLf[0]
	lF := dLServerAddrAndLf[1]
	// local filename on file server

	var fileContents []byte
	cachedFilepath := filepath.Join(p.cacheDir, lF)
	if _, err = os.Stat(cachedFilepath); !os.IsNotExist(err) {
		// if file exists locally.
		fileContents, err = ioutil.ReadFile(cachedFilepath)
		if err != nil {
			return "", fmt.Errorf("error reading cached file %s: %s",
				lF, err)
		}
	} else {
		// make request to file server itself.
		resp, err = http.Get(fmt.Sprintf(`http://%s/%s`, dLServerAddr, lF))
		if err != nil {
			return "", fmt.Errorf("error encountered requesting file %s from "+
				"server %s: %s", lF, dLServerAddr, err)
		}
		defer resp.Body.Close()

		fileContents, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error encountered reading response from "+
				"file server %s: %s", dLServerAddr, err)
		}
		if p.cacheDir != "" {
			// save file if we're using caching.
			err = createFile(cachedFilepath, fileContents)
			if err != nil {
				return "", fmt.Errorf("couldn't create cached file: %s", err)
			}
		}

	}

	return string(fileContents), nil
}

// Upload uploads a file to a file server.
func (p *Proxy) Upload(uploadFile, gloablFilePath string) error {
	if p.directorySrvAddr == "" {
		return fmt.Errorf("directory server address not specified")
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/%s?upload",
		p.directorySrvAddr, gloablFilePath))
	if err != nil {
		return fmt.Errorf("error encountered retrieving upload server addr: %s", err)
	}
	defer resp.Body.Close()

	dlServerAddrAndLf, _ := ioutil.ReadAll(resp.Body)
	splitMsg := strings.Split(string(dlServerAddrAndLf), ",")
	dLServerAddr := splitMsg[0]
	lF := splitMsg[1]
	// local filename on file server

	// make request to the upload server itself
	var u io.Reader
	u, err = os.Open(uploadFile)
	if err != nil {
		return fmt.Errorf("error opening local file: %s", err)
	}

	resp, err = http.Post(fmt.Sprintf("http://%s/%s",
		dLServerAddr, lF), "text/plain", u)
	if err != nil {
		return fmt.Errorf("error encountered trying to upload: %s", err)
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error encountered trying to read server response "+
			"when uploading: %s", err)
	}

	if p.cacheDir != "" {
		// create cached file for skpping downloading in the future.
		c, err := os.Create(filepath.Join(p.cacheDir, lF))
		if err != nil {
			return fmt.Errorf("error encountered creating cache file: %s", err)
		}
		_, err = io.Copy(c, u)
		if err != nil {
			return fmt.Errorf("error encountered copying over upload file to cache file: %s", err)
		}
	}

	return nil
}

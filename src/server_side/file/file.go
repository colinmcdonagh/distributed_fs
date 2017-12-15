// file server
package main

import (
	"flag"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	portPtr := flag.String("port", "", "port to run file server on")
	flag.Parse()

	http.HandleFunc("/", handleQuery)
	http.ListenAndServe(":"+*portPtr, nil)
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	var response string
	filePath := strings.TrimLeft(r.URL.Path, "/")

	if r.Method == "GET" {
		// 0. proxy looking to download file.

		fileContents, err := getLocalFile(filePath)
		if err != nil {
			// 0.1 but there was an error trying to access such a file.
			// presume that the error is because the file doesn't exist.
			fmt.Printf("error encountered when trying to access local file %s: %s",
				filePath, err)
			http.Error(w, fmt.Sprintf("file %s doesn't exist to be downloaded",
				filePath), 007)
			return
		}
		// 0.2 file does eixst.
		response = fmt.Sprintf(fileContents)

	} else if r.Method == "POST" {
		// 1. proxy looking to upload a file.

		// first need to read the file being sent over http.
		fileContents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("error encountered reading body of posted file %s: %s",
				filePath, err)
			return
		}

		// create directory.
		dir := filepath.Dir(filePath)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			http.Error(w, fmt.Sprintf("dir %s could not be created on the server: %s",
				dir, err), 990)
			return
		}

		err = write(filePath, string(fileContents))
		if err != nil {
			// 1.1.1 can't create the local file and write to it.
			http.Error(w, fmt.Sprintf("file %s could not be created on the server: %s",
				filePath, err), 991)
			return
		}
		response = "1"

	} else {
		// 2. don't support this request type
		response = fmt.Sprintf("cannot handle %s type requests\n", r.Method)
	}

	fmt.Fprintf(w, html.EscapeString(response))
}

func getLocalFile(path string) (string, error) {

	lF, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("encountered error when trying to open file %s "+
			"for reading: %s\n", path, err)
	}

	fileContents := make([]byte, 1000000)
	numBytesRead, err := lF.Read(fileContents)
	if err != nil {
		if err == io.EOF {
			return "", nil
		}
		return "", fmt.Errorf("encountered error when trying to read file %s: %s",
			path, err)
	}

	return string(fileContents[:numBytesRead]), nil
}

func write(path, contents string) error {
	lF, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("encountered error when trying to create file: %s", err)
	}

	_, err = lF.Write([]byte(contents))
	if err != nil {
		return fmt.Errorf("encountered errror when trying to write to a file: %s", err)
	}

	return nil
}

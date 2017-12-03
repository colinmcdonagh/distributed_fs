// Allow for transparent access

package main

import (
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	fileServerAddrs  []string
	currUploadServer int

	filePath_serverInt map[string]int
)

func main() {

	fileServerAddrsPtr := flag.String("fileServers", "", "file server addresses")
	flag.Parse()

	fileServerAddrs = strings.Split(*fileServerAddrsPtr, ",")

	// sockets
	// listen on specific port for file_servers to establish themselves.

	// whenever new file_server, add it to the list.

	// keep an index of which files have been found on which file servers.
	// but if no copy, then check each.

	http.HandleFunc("/", handleQuery)
	http.ListenAndServe(":8080", nil)
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	response := ""

	v := r.URL.Query()

	if len(v) == 0 {
		file, err := download(r.URL.Path)
		if err != nil {
			response = fmt.Sprintf("query returned error: %s", err)
		} else {
			response = file
		}
	} else if localFile, ok := v["local_file"]; ok {
		// upload func
		// add to indices
		response = fmt.Sprintf("PUT %s into %s\n", localFile[0], r.URL.Path)
	} else {
		// error
		response = "specified params != local_file\n"
	}
	fmt.Fprintln(w, html.EscapeString(response))
}

func download(path string) (string, error) {
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

// organiser

package main

import (
	"flag"
	"fmt"
	"html"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

var (
	fileServerAddrs       []string
	currUploadServerIndex int

	filePath_serverAddrs map[string][]string
)

func init() {
	filePath_serverAddrs = make(map[string][]string, 0)
}

func main() {

	fileServerAddrsPtr := flag.String("fileSrvAddrs", "", "file server addresses")
	flag.Parse()

	fileServerAddrs = strings.Split(*fileServerAddrsPtr, ",")

	// let's start off with a random server with position n in the fileServerAddrs
	// array which will be responsible for downloads, and from this the upload
	// server will be selected as n plus a certain offset.

	currUploadServerIndex = rand.Intn(len(fileServerAddrs)) - 1
	// before first use of currUploadServer, it will be incremented.
	// therefore, can't == -1.

	http.HandleFunc("/", handleQuery)
	http.ListenAndServe(":8080", nil)
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	var response string
	v := r.URL.Query()

	filePath := strings.TrimLeft(r.URL.Path, "/")

	if len(v) == 0 {
		// 0. it's requesting what server to download from.
		if serverAddr, ok := filePath_serverAddrs[filePath]; ok {
			latestVersion := len(serverAddr) - 1
			// 0.1 file being downloaded exists.
			response = fmt.Sprintf("%s,%s", serverAddr[latestVersion],
				fmt.Sprintf("%s%s", filePath, strconv.Itoa(latestVersion)),
			)
		} else {
			// 0.2 file being downloaded doesn't exist.
			http.Error(w, fmt.Sprintf("file %s doesn't exist", filePath), 71)
		}

	} else if _, ok := v["upload"]; ok {
		// 1. it's requesting what server to upload to.
		currUploadServerIndex = (currUploadServerIndex + 1) %
			len(fileServerAddrs)

		fPServers := filePath_serverAddrs[filePath] // make copy for convenience.
		uploadServer := fileServerAddrs[currUploadServerIndex]
		fPServers = append(fPServers, uploadServer)

		response = fmt.Sprintf("%s,%s",
			uploadServer, fmt.Sprintf("%s%s", filePath, strconv.Itoa(len(fPServers)-1)),
		)
		filePath_serverAddrs[filePath] = fPServers // copy back over.

	} else {
		fmt.Printf("params specified but no upload param specified\n")
		http.Error(w, "url parameters specified but no `upload` param", 72)
	}

	fmt.Fprintf(w, html.EscapeString(response))
}

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

// Specification
/*
 organiser needs to keep track of file server addresses,
 and what server is due to be uploaded to next.

 Also needs to keep track of which file path corresponds to
 which server.
*/

// Assumptions
/*
 an assumption that will be made is that file servers won't be
 able to be added after the original addresses are specified when
 starting up the organiser.

 there is only one organiser.

 a proxy that requests to upload a file will do so, with the
 given local filename specified in the upload server address request.

 can only set local filename when creating the file. After that, would
 have to introduce some sort of remove file command implemented by
 the file servers.
*/

var (
	fileServerAddrs       []string
	currUploadServerIndex int

	filePath_serverAddrs map[string][]string
)

func init() {
	filePath_serverAddrs = make(map[string][]string, 0)
}

func main() {

	fileServerAddrsPtr := flag.String("fileServers", "", "file server addresses")
	flag.Parse()

	fileServerAddrs = strings.Split(*fileServerAddrsPtr, ",")

	// let's start off with a random server with position n in the fileServerAddrs
	// array which will be responsible for downloads, and from this the upload
	// server will be selected as n plus a certain offset.

	currUploadServerIndex = rand.Intn(len(fileServerAddrs)) - 1
	// before first use of currUploadServer, it will be incremented.
	// therefore, can't == -1.

	http.HandleFunc("/", handleQuery)
	http.ListenAndServe(":8081", nil)
}

// two types of base queries that can be made to the organiser are
// that of downloading and uploading.
// both make a query to the path that they are interested in.
// uploads specify an upload parameter.
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
			fmt.Printf("file %s doesn't exist.\n", filePath)
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

// organiser

package main

import (
	"flag"
	"fmt"
	"html"
	"math/rand"
	"net/http"
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

	filePath_serverAddr    map[string]string
	filePath_localFileName map[string]string
)

func init() {
	filePath_serverAddr = make(map[string]string, 0)
	filePath_localFileName = make(map[string]string, 0)
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
	http.ListenAndServe(":8080", nil)
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
		if serverAddr, ok := filePath_serverAddr[filePath]; ok {
			// 0.1 file being downloaded exists.
			response = fmt.Sprintf("%s,%s", serverAddr,
				filePath_localFileName[filePath])
		} else {
			// 0.2 file being downloaded doesn't exist.
			fmt.Printf("file %s doesn't exist.\n", filePath)
			http.Error(w, fmt.Sprintf("file %s doesn't exist", filePath), 71)
		}

	} else if _, ok := v["upload"]; ok {
		// 1. it's requesting what server to upload to.

		lF := filePath

		if lFA, ok := v["local_filepath"]; ok {
			lF = lFA[0]
		} // local file which official file will be mapped to.

		// first need to check if the file exists on a server already.
		if serverAddr, ok := filePath_serverAddr[filePath]; ok {
			// 1.1 file to be upload already exists.
			response = fmt.Sprintf("%s,%s",
				serverAddr,
				filePath_localFileName[filePath])
		} else {
			// 1.2 file is being created for the first time.
			currUploadServerIndex = (currUploadServerIndex + 1) %
				len(fileServerAddrs)
			response = fmt.Sprintf("%s,%s",
				fileServerAddrs[currUploadServerIndex], lF)
			filePath_serverAddr[filePath] = fileServerAddrs[currUploadServerIndex]
			filePath_localFileName[filePath] = lF
		}
	} else {
		fmt.Printf("params specified but no upload param specified\n")
		http.Error(w, "url parameters specified but no `upload` param", 72)
	}

	fmt.Fprintf(w, html.EscapeString(response))
}

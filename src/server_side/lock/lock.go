package main

import (
	"flag"
	"fmt"
	"html"
	"net/http"
	"strings"
)

var (
	lockedFiles map[string]bool
)

func init() {
	lockedFiles = make(map[string]bool)
}

func main() {
	lockPort := flag.String("port", "", "port to run lock server on.")
	flag.Parse()

	http.HandleFunc("/", handleQuery)
	http.ListenAndServe(":"+*lockPort, nil)
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	filePath := strings.TrimLeft(r.URL.Path, "/")

	if len(v) == 0 {
		if locked, ok := lockedFiles[filePath]; ok {
			if locked {
				// was having some difficulty on proxy side detecting an http error
				// when using http.Error(...), so instead just write "0" as response
				// when failing.
				fmt.Fprintf(w, html.EscapeString("0"))
				return
			}
		}
		fmt.Fprintf(w, html.EscapeString("1"))
		lockedFiles[filePath] = true

	} else if _, ok := v["unlock"]; ok {
		lockedFiles[filePath] = false
		fmt.Fprintf(w, html.EscapeString("1"))
	} else {
		http.Error(w, "received unrecognised query param", 91)
	}
}

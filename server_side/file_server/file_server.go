package main

import (
  "flag"
	"fmt"
	"html"
	"net/http"
)

func main() {
  portPtr := flag.String("port", "", "port to run file server on")
  flag.Parse()

  http.HandleFunc("/", handleQuery)
  http.ListenAndServe(":"+*portPtr, nil)
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
  response := ""

  if r.Method == "GET" {
		file, err := fetch(r.URL.Path)
		if err != nil {
			response = fmt.Sprintf("query returned error: %s\n", err)
		} else {
			response = file
		}
	} else if r.Method == "POST" {
		// upload func
		// add to indices
    err := create(r.URL.Path)
    if err != nil {
      response = fmt.Sprintf("couldn't create file: %s\n", err)
    }
    response = "created file\n"
	} else {
		// error
		response = fmt.Sprintf("cannot handle %s type requests\n", r.Method)
	}

  fmt.Fprintln(w, html.EscapeString(response))
}

func fetch(path string) (string, error) {
  return path, nil
}

func create(path string) error {
  return nil
}

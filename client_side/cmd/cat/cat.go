package main

import (
	"fmt"
	"strings"

	"../../proxy"

	"os"
)

const ORG_ADDR = "127.0.0.1:8080"

func main() {
	proxee := proxy.New(ORG_ADDR, "")

	if len(os.Args) != 2 {
		fmt.Println("Please enter one argument, the file to displayed.")
		return
	}
	file := strings.TrimSpace(os.Args[1])

	fileContents, err := proxee.Download(file)
	if err != nil {
		fmt.Printf("error trying to download %s: %s\n", file, err)
		return
	}
	fmt.Println(fileContents)
}

package main

import (
	"fmt"
	"os"
	"strings"

	"../../lib/proxy"
)

// Assumptions
/*
 */

const ORG_ADDR = "127.0.0.1:8081"

func main() {

	proxee := proxy.New(ORG_ADDR, "", "cache")
	if len(os.Args) != 3 {
		fmt.Printf("Please specify the file to copy and the global file identifier.\n"+
			"Received the following arguments: %v\n", os.Args)
	}

	localFile := strings.TrimSpace(os.Args[1])
	globalFile := strings.TrimSpace(os.Args[2])

	err := proxee.Upload(localFile, globalFile)
	if err != nil {
		fmt.Printf("error trying to upload local file %s to global file %s: %s.\n",
			localFile, globalFile, err)
	}
}

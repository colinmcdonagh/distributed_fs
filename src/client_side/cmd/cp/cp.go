package main

import (
	"fmt"
	"os"
	"strings"

	"../../config"
	"../../lib/proxy"
)

func main() {
	proxee := proxy.New(config.DirSrvAddr, config.LockSrvAddr, config.CacheDir)
	if len(os.Args) != 3 {
		fmt.Printf("Please specify the file to copy and the global file identifier.\n"+
			"Received the following arguments: %v\n", os.Args)
	}

	localFile := strings.TrimSpace(os.Args[1])
	globalFile := strings.TrimSpace(os.Args[2])

	err := proxee.Upload(localFile, globalFile)
	if err != nil {
		os.Exit(1)
	}
}

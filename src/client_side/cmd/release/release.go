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
	if len(os.Args) != 2 {
		fmt.Printf("Please specify the file to release.\n")
	}

	releaseFile := strings.TrimSpace(os.Args[1])

	_ = proxee.Unlock(releaseFile)
}

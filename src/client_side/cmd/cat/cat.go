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
		fmt.Println("Please enter one argument, the file to displayed.")
		return
	}
	file := strings.TrimSpace(os.Args[1])

	fileContents, err := proxee.Download(file)
	if err != nil {
		fmt.Printf("error trying to download %s: %s\n", file, err)
		os.Exit(1)
	}
	fmt.Printf(fileContents)
}

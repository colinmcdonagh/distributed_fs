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

const LOCK_ADDR = "127.0.0.1:9104"

func main() {

	proxee := proxy.New("", LOCK_ADDR, "")
	if len(os.Args) != 2 {
		fmt.Printf("Please specify the file to release.\n")
	}

	releaseFile := strings.TrimSpace(os.Args[1])

	_ = proxee.Unlock(releaseFile)
	fmt.Println("unlocked.")
}

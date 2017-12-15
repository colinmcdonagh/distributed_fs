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

const LOCK_ADDR = "127.0.0.1:8104"

func main() {

	proxee := proxy.New("", LOCK_ADDR, "")
	if len(os.Args) != 2 {
		fmt.Printf("Please specify the file to take.\n")
	}

	takeFile := strings.TrimSpace(os.Args[1])

	err := proxee.Lock(takeFile)
	if err != nil {
		fmt.Printf("didn't achieve lock on file: %s\n", err)
		return
	}
	fmt.Println("success.")
}

package main

import (
	"./proxy"

	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

// Assumptions
/*
 exposing what a file server will save a file as does away with
 file transparency, but states that files may be named differently
 locally. Since I'm not using a database to store file - server address
 couples, I simulate the ability to upload files to a particular local
 path on a server, even though it's official path is still the path 
 that it's identified as.
*/

var (
	proxee proxy.Proxy

	downloadRegex *regexp.Regexp
	uploadRegex *regexp.Regexp
)

func init() {
	downloadRegex = regexp.MustCompile(`^cat ([^\S]*)$`)
	uploadRegex = regexp.MustCompile(`^cp ([^\S]*) ([^\S]*):([^\S]*|$)`)
	// local_filepath official_filepath:server_local_filepath
}

func main() {

	organiserAddrPtr := flag.String("orgAddr", "", "address of organiser")
	flag.Parse()

	proxee = proxy.New(*organiserAddrPtr)

	for {
		// read input.
		// call corresponding function.
		reader := bufio.NewReader(os.Stdin)
		cmd, _ := reader.ReadString('\n')
		if matches := downloadRegex.FindStringSubmatch(cmd); len(matches) > 0 {
			fileContents, err := proxee.Download(matches[1])
			if err != nil {
				fmt.Printf("error trying to download %s: %s", matches[1], err)
				continue
			}
			fmt.Println(fileContents)

		} else if matches := uploadRegex.FindStringSubmatch(cmd); len(matches) > 0 {
			
			var lF string
			if len(matches) == 4 {
				lF = matches[3]
			}

			err := proxee.Upload(matches[1], matches[2], lF)
			if err != nil {
				fmt.Printf("error trying to upload local file %s to %s: %s", 
					matches[1], matches[2])
				continue
			}

		}
	}

}

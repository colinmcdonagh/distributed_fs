#!/bin/bash

# call this script from where files will be copied over to distrib fs.

mkdir target

go build cmdline/cp.go -o target/cp
go build cmdline/cat.go -o target/cat
go build cmdline/release.go -o target/release
go build cmdline/take.go -o target/take

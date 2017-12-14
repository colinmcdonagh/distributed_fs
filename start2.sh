#!/bin/bash

# run this script from where file servers should be keeping their files.
# can split this up into two different scripts if need to have file servers
# write to different directories.

go run server_side/directory/directory.go -fileServers "127.0.0.1:${1},127.0.0.1:${2}" &

go run server_side/file/file.go -port $1 &
go run server_side/file/file.go -port $2 &

go run server_side/lock/lock.go - port $3

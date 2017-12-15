#!/bin/bash

# client side.
go build -o bin/client_side/cp src/client_side/cmd/cp/cp.go
go build -o bin/client_side/cat src/client_side/cmd/cat/cat.go
go build -o bin/client_side/release src/client_side/cmd/release/release.go
go build -o bin/client_side/take src/client_side/cmd/take/take.go

# server side.
go build -o bin/server_side/directory src/server_side/directory/directory.go
go build -o bin/server_side/file src/server_side/file/file.go
go build -o bin/server_side/lock src/server_side/lock/lock.go

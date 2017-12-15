# Distributed File System written in Go.

Colin McDonagh - 13322897

Features implemented:
* Distributed Transparent File Access
* Directory Service
* Caching
* Lock Service

## Prerequisites
* Go
* Vim

## Build Steps
1. Determine what ports to run the lock server and directory server on.
For e.g. 8081 and 8084 respectively.

2. Set LockSrvAddr and DirSrvAddr in src/client_side/config/config.go

3. sh build.sh

4. mkdir fs_home && cd fs_home

5. open up a terminal window in this directory for each of the directory server,
lock server, and file servers.

6. ../bin/server_side/lock -port %lockPort (e.g. -port 8081)

7. ../bin/server_side/file -port %filePort (for each file server, e.g.,
  run this twice with 8082 and 8083 as %filePort)

8. ../bin/server_side/directory -fileSrvAddrs %fsAddr1,%fsAddr2,...,%fsAddrn
(e.g. -fileSrvAddrs 127.0.0.1:8082,127.0.0.1:8083)

9. open another terminal window and cd example_app

10. sh edit_file.sh %filename (e.g., sh edit_file.sh test.file)

11. open another terminal window in this directory and access the same file using sh edit_file.sh
in order to illustrate use of lock server.

## Architecture

* Uses a download / upload model.

* Directory server keeps track of versions of files, and on what server they're
stored on.

* When a file is edited, a new version of the file is created. Although this is
not a very scalable feature at all, it allows for more easily transitioning in the
future to using diffs instead of whole files, like git etc.

* Proxy implements client side file caching. The directory to cache under is relative
to where the application using the proxy is called from. The name of the cache directory
itself is specified in src/client_side/config/config.go and is `.cache` by default.

* Includes a sample script for editing files, and common command line tools to run such
  as `cp`, `cat` as well as `take` and `release` for locking files.

* Would be relatively easy to expand current project to display older versions of files
 using the file system's cat command.

* If all the services are up and running, adding another file service requires restarting
  the directory service with the updated list of file servers.

## Notes

Changes to lock server address of directory server address require rebuilding.

Please contact me if any environment-centric build issues arise for you.

Cheers

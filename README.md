Basic Distributed File System written in Go.

Features a command line program (example/cmdline.go), which when run, supports the following
commands for transparently accessing distributed file system:
	* ls
	* cat %file
	* cp %local_file %remote_file

TODO:

1. get cmdline working.
2. lock service.
3. security service.
4. possibly transaction service.

Basic Distributed File System written in Go.

Features a command line program (example/cmdline.go), which when run, supports the following
commands for transparently accessing distributed file system:
	* ls
	* cat %file
	* mv %local_file %remote_file

TODO:

1. modify controller to only return server addresses
2. modify proxy to first make request to controller, and then directly to file servers
3. modify file server to actually serve up shtuff.

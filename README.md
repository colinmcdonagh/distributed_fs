# Distributed File System written in Go.

## Prerequisites
* Go
* vim

## Build Steps

Features a command line clients (client_side/cmd), which when run, supports the following
commands for transparently accessing distributed file system:
* lock %file
* unlock %file
* cat %file
* cp %local_file %remote_file

Also includes a script to launch vim and edit files transparently. (client_side/text_editor)

TODO:
1. ability to version (but not using this in edit)
2. a good README, comments, and code quality.

allow usage of local filenames on proxy. This corresponds to using versions in practice.

some times get an EOF error... can't replicate the error however.
Could be from copying over empty files.
So, when doing so, please write something to the file.

can't add servers, static...

if build script hangs, make sure to check correct config.

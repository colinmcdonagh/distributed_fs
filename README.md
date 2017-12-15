Basic Distributed File System written in Go.

Interesting feature: saves versions of files, although not diffs, but whole files.
(however, could be easily expanded).

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

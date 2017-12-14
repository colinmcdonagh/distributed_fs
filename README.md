Basic Distributed File System written in Go.

Features a command line clients (client_side/cmd), which when run, supports the following
commands for transparently accessing distributed file system:
* lock %file
* unlock %file
* cat %file
* cp %local_file %remote_file

TODO:
1. caching.
2. possibly some sort of replication.
3. ability to version.
4. a good README, comments, and code quality.

Basic Distributed File System written in Go.

Features a command line program (example/cmdline.go), which when run, supports the following
commands for transparently accessing distributed file system:
* ls
* cat %file
* cp %local_file %remote_file

TODO:
1.
. Get rid of option for specifying localfilename, and instead
  extend this functionality into the organiser, appending a version
  number to the file, which would allow for using diffs in the future.
  
2.
  1. lock service for sure.
  2. security service.
  3. (optionally) caching in part, maybe half-done.

3. etc.
Exporting path.

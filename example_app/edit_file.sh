#!/bin/bash
# script for editing files in vim.

# lock file we're interested in.
../bin/client_side/take $1
if [ $? -ne 0 ]; then
  # if already taken.
  echo "file is currently being edited by another party."
  exit 1
fi

# create a tmp directory for editing file.
mkdir -p tmp

# get current file.
FILE="$(../bin/client_side/cat ${1})"

if [ $? -eq 0 ]; then
  # if file exists.
  echo "${FILE}" | vim - -c ":file tmp/tempfile"
else
    # if file doesn't exist yet.
    vim tmp/tempfile
fi
# if file doesn't exist, i.e., vim didn't create it because it was empty.
[ -s tmp/tempfile ] || touch tmp/tempfile

# upload new version of file.
../bin/client_side/cp tmp/tempfile $1
rm -rf tmp

# release lock on file.
../bin/client_side/release $1

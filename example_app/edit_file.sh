#!/bin/bash

# run this script to open distributed file in vim.

mkdir -p tmp
# see if cat returns an error.
# if it does, well start afresh.
FILE="$(../bin/client_side/cat ${1})"
if [ $? -eq 0 ]; then
  echo "${FILE}" | vim - -c ":file tmp/tempfile"
else
    vim tmp/tempfile
fi

# if file is empty, create it in order to upload empty file.
[ -s tmp/tempfile ] || touch tmp/tempfile

../bin/client_side/cp tmp/tempfile ${1}
if [ $? -eq 0 ]; then
  echo OK
else
  echo FAIL
fi
rm -rf tmp

#!/bin/bash

# run this script to open distributed file in vim.

mkdir -p tmp
../bin/client_side/cat $1 | vim - -c "tmp/tempfile"
../bin/client_side/cp tmp/tempfile $1
rm -rf tmp

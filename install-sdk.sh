#!/bin/bash

ARCH="386"
VERSION="1.9.2"

if [[ `uname -a` == *x86_64* ]]
then
    ARCH="amd64"
fi

file=go_appengine_sdk_linux_$ARCH-$VERSION.zip
echo "downloading '$file'"

wget https://googleappengine.googlecode.com/files/$file -nv
unzip -q $file -d .
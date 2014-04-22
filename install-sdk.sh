#!/bin/bash

ARCH="386"
VERSION="1.9.2"
SVERSION=$VERSION | sed 's//\./g'

if [[ `uname -a` == *x86_64* ]]
then
    ARCH="amd64"
fi

FILE=go_appengine_sdk_linux_$ARCH-$VERSION.zip
echo "downloading '$FILE'"

wget https://commondatastorage.googleapis.com/appengine-sdks/featured/$FILE -nv
wget https://console.developers.google.com/m/cloudstorage/b/appengine-sdks/o/deprecated/$SVERSION/$FILE -nv
unzip -q $FILE -d .
#!/bin/sh

MONGODB_FILE_NAME="mongodb-src-r4.2.0.tar.gz"
MONGODB_URL="https://fastdl.mongodb.org/src/$MONGODB_FILE_NAME"
ESCAPE_SEQ="\033[0"
DEFAULT="$ESCAPE_SEQ;39m"
CYAN="$ESCAPE_SEQ;36m"
RED="$ESCAPE_SEQ;31m"
YELLOW="$ESCAPE_SEQ;33m"
MAGENTA="$ESCAPE_SEQ;35m"

echo $CYAN"\nDownloading Mongodb Tar File: $DEFAULT$MONGODB_URL"
echo $MAGENTA
curl -O $MONGODB_URL

echo $CYAN"Decompressing Mongodb Tar File: $DEFAULT$MONGODB_FILE_NAME"
tar -zxf $MONGODB_FILE_NAME

LOCAL_MONGODB_PATH="$(pwd)/mongodb-macos-x86_64-4.2.0"
LOCAL_MONGODB_BIN_PATH="$LOCAL_MONGODB_PATH/bin"
echo $CYAN"Unpacked Mongodb Tar File at: $DEFAULT$LOCAL_MONGODB_PATH"
echo $RED"ADDING TO PATH: $DEFAULT$LOCAL_MONGODB_BIN_PATH"
export PATH=$PATH:"$LOCAL_MONGODB_BIN_PATH"

MONGODB_DATA_PATH="/data/db"
echo $RED"REQUIRES SUDO PERMISSIONS"$CYAN" Creating Mongodb Data Directory: $DEFAULT$MONGODB_DATA_PATH"
sudo mkdir -p /data/db

echo $RED"REQUIRES SUDO PERMISSIONS"$CYAN" Setting Mongodb Data Directory Read+Write Permissions for '$USER'"
sudo chmod -R a+rw /data/db

echo $RED"CLEANING UP"$CYAN" Removing Mongodb Tar File: $DEFAULT$MONGODB_FILE_NAME"
rm $MONGODB_FILE_NAME

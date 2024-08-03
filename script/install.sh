#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin:/opt/homebrew/bin
export PATH

startTime=`date +%s`
echo "安装中..."

DST_DIR=/opt/dztasks

sysName=`uname`
sysArch=`arch`

echo $sysName, $sysArch

VERSION='1.0'

ARCH="amd64"
if [ "$sysArch" == "x86_64" ];then
	ARCH="amd64"
elif [ "$sysArch" == "aarch64" ];then
	ARCH="arm64"
elif [ "$sysArch" == "arm64" ];then
	ARCH="arm64"
else
	ARCH="amd64"
fi

SYSTEM=linux
if [ "$sysName" == "Darwin" ];then
	SYSTEM=darwin
fi


# https://github.com/midoks/dztasks/releases/download/1.0/dztasks_v1.0_linux_amd64.tar.gz
DZTASKS_URL=https://github.com/midoks/dztasks/releases/download/${VERSION}
FILE_NAME=dztasks_v${VERSION}_${SYSTEM}_${ARCH}.tar.gz
TMP_DIR=/tmp

if [ ! -f $TMP_DIR/${FILE_NAME} ];then
	wget --no-check-certificate -O $TMP_DIR/${FILE_NAME} ${DZTASKS_URL}/${FILE_NAME}
fi

mkdir -p $DST_DIR

cd $DST_DIR && tar zxvf $TMP_DIR/$FILE_NAME

if [ -f $TMP_DIR/install.sh ];then
	rm -rf $TMP_DIR/install.sh
fi

if [ -f $TMP_DIR/${FILE_NAME} ];then
	rm -rf $TMP_DIR/${FILE_NAME}
fi


endTime=`date +%s`
((outTime=(${endTime}-${startTime})/60))
echo -e "安装耗时:\033[32m $outTime \033[0mMinute!"



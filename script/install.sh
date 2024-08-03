#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin:/opt/homebrew/bin
export PATH

startTime=`date +%s`
echo "安装中..."


sysName=`uname`
sysArch=`arch`

echo $sysName, $sysArch




endTime=`date +%s`
((outTime=(${endTime}-${startTime})/60))
echo -e "安装耗时:\033[32m $outTime \033[0mMinute!"



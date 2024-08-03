#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin:/opt/homebrew/bin
export PATH

startTime=`date +%s`
echo "安装中..."

DST_DIR=/opt/dztasks

sysName=`uname`
sysArch=`arch`

echo $sysName, $sysArch

systemd_dir=/lib/systemd/system
if [ ! -d /usr/lib/systemd/system ];then
	systemd_dir=/usr/lib/systemd/system
fi

systemctl stop dztasks
systemctl disable dztasks

if [ -f ${systemd_dir}/dztasks.service ];then
	rm -rf ${systemd_dir}/dztasks.service
fi

systemctl daemon-reload



endTime=`date +%s`
((outTime=(${endTime}-${startTime})/60))
echo -e "卸载耗时:\033[32m $outTime \033[0mMinute!"



#!/bin/sh

if [ $# -ne 2 ]; then
	echo "Usage: sh $0 <ip> <app>"
	return
fi

adb connect $1
adb -s $1 uninstall $2

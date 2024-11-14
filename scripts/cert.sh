#!/bin/sh

LOG_FILE='/app/logging/log/operation.privapp.log'
#adb connect $1

#msg=$(adb root)
#sh /app/logging/agent/helper/log.sh D "Changing to root: $msg" $0 $LOG_FILE

#msg=$(adb connect $1)
#sh /app/logging/agent/helper/log.sh D "Connecting to device: $msg" $0 $LOG_FILE

adb -s $1 push cert/c8750f0d.0 /sdcard
adb -s $1 push cert/install.sh /sdcard
adb -s $1 shell su -c sh /sdcard/install.sh

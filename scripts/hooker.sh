#!/bin/sh
if [ $# -ne 1 ]
then
    echo 'Usage: sh hooker.sh <target_ip>'
    return
fi

LOG_FILE='/app/logging/log/operation.privapp.log'

#msg=$(adb connect $1)
#sh /app/logging/agent/helper/log.sh D "Connecting to device: $msg" $0 $LOG_FILE
adb -s $1 push intercept/pinning/frida-server-12.7.3-android-x86_64 /sdcard/frida-server
adb -s $1 push intercept/pinning/install.sh /sdcard
adb -s $1 shell su -c "sh /sdcard/install.sh"
adb -s $1 push ~/.mitmproxy/mitmproxy-ca-cert.cer /data/local/tmp/cert-der.cr
adb -s $1 shell su -c '/data/local/frida-server &'

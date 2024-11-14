#!/bin/sh

LOG_FILE='/app/logging/log/operation.privapp.log'

pid=$(ps aux | grep -v grep | grep mitmdump | awk '{print $2}')
pid_frida=$(ps aux | grep -v grep | grep fridactl | awk '{print $2}')

kill -9 $pid
if [ $? -eq 0 ]; then
        sh /app/logging/agent/helper/log.sh D "kill mitmproxy success " $0 $LOG_FILE
else
        sh /app/logging/agent/helper/log.sh E "kill mitmproxy fail " $0 $LOG_FILE
fi

kill -9 $pid_frida
if [ $? -eq 0 ]; then
        sh /app/logging/agent/helper/log.sh D "kill fridactl success" $0 $LOG_FILE
else
        sh /app/logging/agent/helper/log.sh E "kill fridactl fail" $0 $LOG_FILE
fi

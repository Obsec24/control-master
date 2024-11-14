#!/usr/bin/env python3
import imp
import tools 
import os 

LOG_FILE='/app/logging/log/operation.privapp.log'
HELPER_JSON_LOGGER = '/app/logging/agent/helper/log.py'

#configure json logger
assert os.path.isfile(HELPER_JSON_LOGGER), '%s  is not a valid file or path to file' % HELPER_JSON_LOG
log = imp.load_source('log', HELPER_JSON_LOGGER)
logger =  log.init_logger(LOG_FILE)

(success, result) = tools.call_sh_output("ps aux | grep -v grep | grep mitmdump | awk '{print $2}'")
print(success, result)

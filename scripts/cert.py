import tools
import imp
import configparser 
import os

#Logging init
LOG_FILE = '/app/logging/log/operation.privapp.log'
HELPER_JSON_LOGGER = '/app/logging-master/agent/helper/log.py'
assert os.path.isfile(HELPER_JSON_LOGGER), '%s  is not a valid file or path to file' % HELPER_JSON_LOGGER
log = imp.load_source('log', HELPER_JSON_LOGGER)
logger =  log.init_logger(LOG_FILE)

#reading device config
CONFIG_FILE = '/app/scripts/testing.config'
assert os.path.isfile(CONFIG_FILE), '%s  is not a valid file or path to file' % CONFIG_FILE
config = configparser.ConfigParser()
config.read(CONFIG_FILE)
assert 'device' in config.sections(), 'Config file %s does not contain an [device] section' % CONFIG_FILE
assert 'device_serial' in config['device'], 'Config file %s does not have an device_serial value in the device section' % CONFIG_FILE
device_serial = config['device']['device_serial']
print(device_serial)
#logger.error("Google Play Authentication failure")
#loger.debug("Reading device properties")
#logger.warning("APK download failed.", extra={'apk': 'nombre de la APK'})
#logger.info("Successful APK download", extra={'apk': 'nombre de la APK'})


#tools.adb

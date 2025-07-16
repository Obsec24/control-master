#!/usr/bin/env python3
import shutil
from datetime import datetime, timedelta
import time

import tools
import imp
import os
import argparse
import sys
import tarfile

# Tools config file
TOOLS_FILE = '/app/scripts/testing.config'
# Logging init
FILE_LOGS = '/app/logging/log/operation.privapp.log'
HELPER_JSON_LOGGER = '/app/logging-master/agent/helper/log.py'
TERMINAL_PASSWORD = "5131"
DIR_OUTPUT = "/app/logging/log"

DEVICE_NOT_CONNECTED_ERROR = 10
APP_INSTALL_FAIL_ERROR = 20
MITM_PROXY_START_ERROR = 30
APK_NOT_FOUND = 50

# configure json logger
assert os.path.isfile(HELPER_JSON_LOGGER), '%s  is not a valid file or path to file' % HELPER_JSON_LOGGER
log = imp.load_source('log', HELPER_JSON_LOGGER)
logger = log.init_logger(FILE_LOGS)


# logger.error("Google Play Authentication failure")
# loger.debug("Reading device properties")
# logger.warning("APK download failed.", extra={'apk': 'nombre de la APK'})
# logger.info("Successful APK download", extra={'apk': 'nombre de la APK'})

def parse_args():
    parser = argparse.ArgumentParser(description='Automatic traffic analysis')
    parser.add_argument('--app', '-a', help='App package name', required=True)
    parser.add_argument('--device', '-d', help='Android device ID or IP', required=True)
    parser.add_argument('--label', '-l', help='Testing label', required=True)
    parser.add_argument('--timeout', '-t', help='Testing time', required=True, type=int)
    parser.add_argument('--permissions', '-p', action='store_true', help='Gran all permissions')
    parser.add_argument('--reboot', '-r', action='store_true', help='Reboot before testing')
    parser.add_argument('--screenshots', '-s', action='store_true', help='Capture screenshots')
    return parser.parse_args()


def compress_pngs(png_dir, outfile, delete_pngs=True):
    # Find all PNGs in the directory and compress them as <png_dir>/<outfile>.tar.bz2
    png_files = [os.path.join(png_dir, x) for x in os.listdir(png_dir) if x.endswith('.png')]

    if len(png_files) > 0:
        tar = tarfile.open(os.path.join(png_dir, outfile), 'w:gz')
        for png in png_files:
            tar.add(png, arcname=os.path.basename(png))
        tar.close()

    if delete_pngs:
        for png in png_files:
            os.remove(png)


if __name__ == '__main__':
    # parsing parameters
    args = parse_args()
    app = args.app
    device = args.device
    timeout = args.timeout
    perm = args.permissions
    reboot = args.reboot
    screen = args.screenshots
    testing_label = args.label

    # initializing connection to android device
    assert os.path.isfile(TOOLS_FILE), '%s  is not a valid file or path to file' % TOOLS_FILE
    tools.init(TOOLS_FILE, device)

    # checking availability of APK
    if not os.path.isfile("base.apk"):
        logger.error('APK to be installed not found',
                  extra={'testing_label': testing_label, 'apk': app, 'container': 'traffic', 'device': device})
        sys.exit(APK_NOT_FOUND)
    # creating output folder
    version = tools.aapt_version_code("base.apk")
    data_dir = os.path.join(DIR_OUTPUT, app, version)
    if not os.path.isdir(data_dir):
        os.makedirs(data_dir)
        os.chmod(data_dir, 0o666)
    data_dir = os.path.join(data_dir, testing_label)
    if not os.path.isdir(data_dir):
        os.makedirs(data_dir)
        os.chmod(data_dir, 0o666)
    # Redirect print statements to file
    orig_sysout = sys.stdout
    log_file = os.path.join(data_dir, '%s-%s-first-phase-%s.log' % (app, version, testing_label))
    f = None
    f = open(log_file, 'w')
    sys.stdout = f

    # checking if device is connected
    (success, result) = tools.adb_shell(['getprop ro.serialno'])
    if not success:
        logger.error("No device connected", extra={'testing_label': testing_label, 'version': version,
                                                                     'apk': app, 'container': 'traffic', 'device': device})
        sys.exit(DEVICE_NOT_CONNECTED_ERROR)
    #configuring mobile device
    (success, result) = tools.adb_shell(['svc wifi enable'])
    (success, result) = tools.adb_shell(['svc data disable'])
    (success, result) = tools.adb_shell(['settings put secure location_providers_allowed +gps'])
    time.sleep(3)

    # installing app
    failure = False
    for aux in range(1):  # I disabled the second attempt
        if failure:
            # Trying to install the app disabling "verify apps over USB"
            tools.adb_shell(['settings put global verifier_verify_adb_installs 0'], retry_limit=0)
            logger.info('Disabled the verification of apps over USB', extra={'testing_label': testing_label, 'version': version,
                                                                     'apk': app, 'container': 'traffic', 'device': device})
            time.sleep(2)
        if perm:
            #(success, result, permissions, perm_nogranted) = tools.adb_install_auto('base.apk', grant_all_perms=perm)
            (success, result, permissions, perm_nogranted) = tools.adb_install('base.apk', grant_all_perms=perm)
        else:
            #(success, result) = tools.adb_install_auto('base.apk', perm)
            (success, result) = tools.adb_install('base.apk', grant_all_perms=perm)
        if not success:
            failure = True
        else:
            break

    if not success:
        logger.error('App install failed', extra={'testing_label': testing_label, 'version': version,
                                                                     'apk': app, 'container': 'traffic',
                                                  'exception_message': result, 'device': device})
        sys.exit(APP_INSTALL_FAIL_ERROR)
    else:
        if perm:
            logger.debug('App install successfully',
                        extra={'testing_label': testing_label, 'apk': app, 'version': version,
                               'permissions': permissions,'container': 'traffic',
                               'non-granted': perm_nogranted, 'device': device})
        else:
            logger.debug('App install successfully', extra={'testing_label': testing_label, 'apk': app,
                                                           'version': version, 'container': 'traffic', 'device': device})
    # if set reboot, then reboot
    if reboot:
        tools.adb_reboot(wait=True)  # True waits phone is rebooted

    # starting proxy
    (success, result) = tools.call_sh('nohup mitmdump -s intercept/inspect_requests.py --set app={} &'.format(app),
                                      timeout_secs=20)
    if not success:
        logger.error('Mitmproxy start failed',
                     extra={'testing_label': testing_label, 'apk': app, 'version': version, 'container': 'traffic',
                            'exception_message': result, 'device': device})
        sys.exit(MITM_PROXY_START_ERROR)
    else:
        logger.debug('Mitmproxy start succesfully', extra={'testing_label': testing_label, 'apk': app,
                                                          'version': version, 'container': 'traffic', 'device': device})
    # starting frida server
    tools.adb_shell(['su', '-c', '/data/local/frida-server'], retry_limit=0)
    # starting frida client
    (success, result) = tools.call_sh(
        '/app/intercept-master/pinning/fridactl.py {} {} {} {} &'.format(device, app, testing_label, version), timeout_secs=10)
    if not success:
        logger.error('Fridactl start failed, pinning-protected traffic will not be captured',
                     extra={'testing_label': testing_label, 'apk': app, 'version': version, 'container': 'traffic',
                            'exception_message': result, 'device': device})
    else:
        logger.debug('Fridactl start successfully', extra={'testing_label': testing_label, 'apk': app,
                                                        'version': version, 'container': 'traffic', 'device': device})

    # time.sleep(timeout // 3)
    # screen_file = os.path.join(data_dir,
    #                            '%s-%s-test-%s-fp-start.png' % (app, version, testing_label))
    # tools.adb_screenshot(screen_file)
    #
    # time.sleep(timeout // 3)
    # screen_file = os.path.join(data_dir,
    #                            '%s-%s-test-%s-fp-middle.png' % (app, version, testing_label))
    # tools.adb_screenshot(screen_file)

    time.sleep(timeout)
    # screen_file = os.path.join(data_dir,
    #                            '%s-%s-test-%s-fp-end.png' % (app, version, testing_label))
    #tools.adb_screenshot(screen_file)

    # logger.info('Capture of first phase traffic finished', extra={'testing_label': testing_label, 'apk': app,
    #                                                               'version': version, 'container': 'traffic', 'device': device})
    if f is not None:
        sys.stdout = orig_sysout
        f.close()
    sys.exit(0)

#!/usr/bin/env python3

import argparse
import datetime
import imp
import os
import sys
import tarfile
import time

import tools

TOOLS_FILE = '/app/scripts/testing.config'
# Logging init
FILE_LOGS = '/app/logging/log/operation.privapp.log'
HELPER_JSON_LOGGER = '/app/logging/agent/helper/log.py'
DIR_OUTPUT = "/app/logging/log"
SCREENSHOT_MANUAL_TIMEOUT = 10
# configure json logger
assert os.path.isfile(HELPER_JSON_LOGGER), '%s  is not a valid file or path to file' % HELPER_JSON_LOGGER
log = imp.load_source('log', HELPER_JSON_LOGGER)
logger = log.init_logger(FILE_LOGS)


def parse_args():
    parser = argparse.ArgumentParser(description='Automatic traffic analysis')
    parser.add_argument('--app', '-a', help='App package name', required=True)
    parser.add_argument('--device', '-d', help='Android device ID or IP', required=True)
    parser.add_argument('--monkey', '-m', action='store_true', help='Use monkey')
    parser.add_argument('--label', '-l', help='Testing label', required=True)
    parser.add_argument('--timeout', '-t', help='Testing time', required=True, type=int)
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
    monkey = args.monkey
    timeout = args.timeout
    testing_label = args.label
    # initializing connection to android device
    assert os.path.isfile(TOOLS_FILE), '%s  is not a valid file or path to file' % TOOLS_FILE
    tools.init(TOOLS_FILE, device)
    # creating output folder
    version = tools.aapt_version_code("base.apk")
    data_dir = os.path.join(DIR_OUTPUT, app, version)
    if not os.path.isdir(data_dir):
        os.makedirs(data_dir)
        os.chmod(data_dir, 0o776)
    data_dir = os.path.join(data_dir, 'testing-%s' % testing_label)
    if not os.path.isdir(data_dir):
        os.makedirs(data_dir)
        os.chmod(data_dir, 0o776)

    # Redirect print statements to file
    orig_sysout = sys.stdout
    log_file = os.path.join(data_dir, '%s-%s-testing-second-phase-%s.log' % (app, version, testing_label))
    f = None
    f = open(log_file, 'w')
    sys.stdout = f

    # checking if device is connected
    (success, result) = tools.adb_shell(['getprop ro.serialno'])
    if not success:
        logger.error("No device connected, stopping testing in second phase",
                     extra={'testing_label': testing_label, 'apk': app,
                            'version': version, 'container': 'traffic', 'device': device})
        sys.exit(10)

    if monkey:
        (permissions, permissions_nogranted) = tools.adb_grant_permission("base.apk")
        logger.info('Permissions granted',
                    extra={'testing_label': testing_label, 'apk': app, 'version': version, 'container': 'traffic',
                           'device': device, 'permissions': permissions, 'non-granted': permissions_nogranted})
    else:
        permissions = tools.aapt_permissions('base.apk')
        logger.info('Permissions requested',
                    extra={'testing_label': testing_label, 'apk': app, 'version': version, 'container': 'traffic',
                           'device': device, 'permissions': permissions})
    # Explore the app
    end_time = datetime.datetime.now() + datetime.timedelta(seconds=timeout)
    screen_count = 0
    while datetime.datetime.now() < end_time:
        screen_count = screen_count + 1
        if monkey:
            tools.adb_monkey(app, delay_ms=500)
        else:
            time.sleep(SCREENSHOT_MANUAL_TIMEOUT)
        # screen_file = os.path.join(data_dir,
        #                            '%s-%s-test-%s-sp-%d.png' % (app, version, testing_label, screen_count))
        # tools.adb_screenshot(screen_file)
    # logger.info('Capture in second phase finished', extra={'testing_label': testing_label, 'apk': app,
    #                                                               'version': version, 'container': 'traffic', 'device': device})
    # Close the log file and restore stdout
    if f is not None:
        sys.stdout = orig_sysout
        f.close()
    sys.exit(0)

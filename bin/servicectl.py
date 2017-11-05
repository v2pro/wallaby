#!/usr/bin/env python2.7
import os
import subprocess
import sys
import urllib2
import json
import datetime
import signal
import time

# ===== biz-module requirements =====
# the whole process start with: (get called by deploy system)
#       biz-module/control.sh start

# biz-module/control.sh start:
#   export WALLABY_SERVICES="tag1:8005 tag2:8006 tag3:8007"
#   wallaby/servicectl.py start

# biz-module/wcontrol.sh is required,
#   start|stop must be provided so that the biz process can be managed by wallaby

# ===== instructions =====
# servicectl provides:
# 1. public command interface: exposed to biz-module
# 2. wallaby proxy commands
# 3. wallaby app commands

# Please see help() for more information.

# ===== ENVIRONMENTS =====
# passed from biz modules, e.g., biz-module/control.sh

# WALLABY_SERVICES:
#       format: {tag-name}:{port} {tag-name2}:{port2}...
#       e.g., "tag1:8005 tag2:8006 tag3:8007"

# WALLABY_BIZ_MODULE_DIR: the directory of biz module


# ===== GLOBALS =====
# all the biz modules will be deployed under WALLABY_DEPLOY_ROOT
WALLABY_DEPLOY_ROOT="/tmp/wallaby_services/versions"

# all the biz module process pids are saved in WALLABY_SERVICE_PROC_DIR
WALLABY_SERVICE_PROC_DIR="/tmp/wallaby_services/proc"

WALLABY_PIDFILE="/tmp/wallaby-proxy.pid"

# server address for getting/setting running app versions
WALLABY_SERVER="http://127.0.0.1:8869"

def help():
    print "\n===== Manual =====\n"
    print sys.argv[0], "{start|stop|status|restart|appstart|appstop|appstatus|proxystart|proxystop|proxystatus}"

def getBizModuleDir():
    print os.environ["WALLABY_BIZ_MODULE_DIR"]
    return os.environ["WALLABY_BIZ_MODULE_DIR"]

def getRunningVersions():
    global WALLABY_SERVER
    resp = urllib2.urlopen(WALLABY_SERVER+"/list").read()
    versions = json.loads(resp)
    if versions is None:
        return []
    return versions

def notifyNewVersion(data):
    global WALLABY_SERVER
    req = urllib2.Request(WALLABY_SERVER+'/set')
    req.add_header('Content-Type', 'application/json')
    response = urllib2.urlopen(req, json.dumps(data))

def getNextVersion(runningMap):
    serviceList = os.environ["WALLABY_SERVICES"].split()
    newVersion = {}
    for service in serviceList:
        tag, port = service.split(":")

        address = "127.0.0.1:"+ port
        if address in runningMap:
            print tag, address, "is running, try next"
        else:
            newVersion = {
                "address": "127.0.0.1:"+ port,
                "name": "test",
                "version": "1.0.3",
                "status": "Running",
                "tag": tag
            }
            break
    return newVersion

def appStart(version, appPath):
    if not os.path.exists(WALLABY_SERVICE_PROC_DIR):
        os.makedirs(WALLABY_SERVICE_PROC_DIR)
    try:
        port = version["address"].split(":")[1]
        tag = version["tag"]
        print "\n\n===== START APP: %s:%s" % (tag, port)
        pidfile = WALLABY_SERVICE_PROC_DIR + "/" + tag + ".pid"
        startCmd = ["%s/wcontrol.sh" % (appPath), "start", tag]
        print " ".join(startCmd)
        print "\n\n===== APP OUTPUT"
        pid = subprocess.Popen(startCmd).pid
        f = open(pidfile, 'w')
        f.write(str(pid))
        print "\n\n===== APP is RUNNING =====\npid: %d\npidfile: %s" % (pid, pidfile)
        return True
    except subprocess.CalledProcessError as err:
        print err
        return False

def appStop(tag):
    pidfile = WALLABY_SERVICE_PROC_DIR + "/" + tag + ".pid"
    f = open(pidfile, 'r')
    pid = int(f.read())
    print "\n\n===== STOP APP"
    try:
        os.kill(pid, 0)
        os.kill(pid, signal.SIGKILL)
        print "APP is STOPPED, pid: %d" % pid
    except OSError:
        print "APP is DOWN, pid: %d" % pid

def appStatus(tag):
    pidfile = WALLABY_SERVICE_PROC_DIR + "/" + tag + ".pid"
    f = open(pidfile, 'r')
    pid = int(f.read())
    try:
        os.kill(pid, 0)
        print "APP is RUNNING, pid: %d" % pid
    except OSError:
        print "APP is DOWN, pid: %d" % pid

def appStatusCmd():
    print "\n===== APP STATUS:"
    versions = getRunningVersions()
    for v in versions:
        appStatus(v["tag"])

def appStartCmd():
    print "not implemented"
    pass

def appStopCmd():
    if len(sys.argv) < 2:
        print "usage: %s appstop {tag}" % (sys.argv[0])
        exit(1)
    tag=sys.argv[2]
    appStop(tag)

def cleanOldRelease():
    print "no implemented"

def deploy(version, dst):
    print "\n\n===== DEPLOY ", version
    if not os.path.exists(dst):
        os.makedirs(dst)
    deploycmd="rsync -avz %s/. %s" % (getBizModuleDir(), dst)
    print deploycmd
    try:
        subprocess.check_call(deploycmd, shell=True)
        print "DEPLOY SUCCESS"
        return True
    except subprocess.CalledProcessError as err:
        print "DEPLOY FAIL", err
        return False

def proxyStartCmd():
    print "\n\n===== START WALLABY PROXY at 8869"
    global WALLABY_PIDFILE
    try:
        f = open(WALLABY_PIDFILE, 'r')
        pid = int(f.read())
        os.kill(pid, 0)
        print "WALLABY PROXY is ALREADY RUNNING, pid: %d" % pid
    except OSError:
        pid = subprocess.Popen(["proxy"]).pid
        print "===== WALLABY PROXY is RUNNING =====\npid: %d\npidfile: %s" % (pid, WALLABY_PIDFILE)
        print "please wait 2 seconds..."
        time.sleep(2)
        f = open(WALLABY_PIDFILE, 'w')
        f.write(str(pid))

def proxyStopCmd():
    global WALLABY_PIDFILE
    f = open(WALLABY_PIDFILE, 'r')
    pid = int(f.read())
    print "\n\n===== STOP WALLABY PROXY at 8869"
    try:
        os.kill(pid, 0)
        os.kill(pid, signal.SIGKILL)
        print "WALLABY PROXY is STOPPED, pid: %d" % pid
    except OSError:
        print "WALLABY PROXY is DOWN, pid: %d" % pid

def proxyStatusCmd():
    print "\n===== WALLABY PROXY STATUS:"
    global WALLABY_PIDFILE
    f = open(WALLABY_PIDFILE, 'r')
    pid = int(f.read())
    try:
        os.kill(pid, 0)
        print "WALLABY PROXY is RUNNING, pid: %d" % pid
    except OSError:
        print "WALLABY PROXY is DOWN, pid: %d" % pid

# ===== command interface =====
def start():
    proxyStartCmd()
    print "\n\n===== WALLABY START"
    print "get running versions..."
    runningVersions = getRunningVersions()
    runningMap = {}
    for v in runningVersions:
        runningMap[v["address"]] = v
        print v["name"], ":", v["address"], v["version"], v["status"]

    if len(os.environ["WALLABY_SERVICES"]) == 0:
        print "WALLABY_SERVICES not found, exit"
        exit(1)
    print "WALLABY_SERVICES:", os.environ["WALLABY_SERVICES"]
    newVersion = getNextVersion(runningMap)
    if newVersion:
        release = datetime.datetime.now().strftime("%Y%m%d%H%M%S")
        newVersion["release"] = release
        print "found new version", newVersion
        global WALLABY_DEPLOY_ROOT
        dst = WALLABY_DEPLOY_ROOT + "/" + newVersion["release"]
        if deploy(newVersion, dst):
            if appStart(newVersion, dst):
                notifyNewVersion(newVersion)
    else:
        print "all versions are running, exit"

def stop():
    proxyStopCmd()
    print "===== WALLABY STOP"

def restart():
    proxyStopCmd()
    proxyStartCmd()

def status():
    proxyStatusCmd()

def main():
    if len(sys.argv) <= 1:
        help()
        exit(1)

    cmd=sys.argv[1]
    if cmd == 'start':
        start()
    elif cmd == 'stop':
        stop()
    elif cmd == 'restart':
        restart()
    elif cmd == 'status':
        status()
    elif cmd == 'appstatus':
        appStatusCmd()
    elif cmd == 'appstart':
        appStartCmd()
    elif cmd == 'appstop':
        appStopCmd()
    elif cmd == 'proxystatus':
        proxyStatusCmd()
    elif cmd == 'proxystart':
        proxyStartCmd()
    elif cmd == 'proxystop':
        proxyStopCmd()
    else:
        help()
        exit(1)

if __name__ == "__main__":
    main()
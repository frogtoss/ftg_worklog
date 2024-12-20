#!/usr/bin/env python3

import os
import sys
import json
import shutil
import subprocess

from os.path import join as path_join

def is_host_windows():
    return os.name == 'nt'

def exe(name):
    if is_host_windows():
        return name + '.exe'
    else:
        return name

def fatal(msg):
    print(msg, file=sys.stderr)
    sys.exit(1)

def message(msg):
    if not find_arg_0param('--quiet'):
        print(msg)

def find_arg_0param(expected_arg):
    for arg in sys.argv:
        if arg == expected_arg:
            return True

    return False

def shell(cmd):
    if not find_arg_0param('--quiet'):
        print(' '.join(cmd))
    cp = subprocess.run(cmd)

    if cp.returncode != 0:
        fatal("%s failed" % ' '.join(cmd))

def shell_backtick(cmd, shell):
    return subprocess.run(cmd, capture_output=True, shell=shell).stdout

def get_installed_executable_path(exe):
    gopath = path_join(os.environ.get('GOPATH'), 'bin')
    if not os.environ.get('GOBIN') is None:
        gopath = os.environ.get('GOBIN')

    return path_join(gopath, exe)



def get_host_os_tools_bin():
    return os.getenv("FTG_TOOLS_BIN_DIR")

#
# main
#

EXE=exe("ftgworklog")

os.chdir(path_join('cmd', 'ftgworklog'))

shell(['go', 'install'])

if os.environ.get('FTG_PROJECT_ROOT') is None:
    sys.exit(0)

if not find_arg_0param('--skip-install'):
    src_path = get_installed_executable_path(EXE)
    dst_path = path_join(get_host_os_tools_bin(), EXE)
    shutil.copy2(src_path, dst_path)
    message("%s installed to '%s'" % (EXE, dst_path))

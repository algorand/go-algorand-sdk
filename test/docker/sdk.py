#!/usr/bin/env python3

import subprocess
import sys

default_dirs = { 
    'features_dir': '/opt/go/src/github.com/algorand/go-algorand-sdk/test/features',
    'source': '/opt/go/src/github.com/algorand/go-algorand-sdk',
    'docker': '/opt/go/src/github.com/algorand/go-algorand-sdk/test/docker',
}

def setup_sdk():
    """
    Setup go cucumber environment.
    """    
    subprocess.check_call(['go generate %s/...' % default_dirs['source']], shell=True)

def test_sdk():
    sys.stdout.flush()
    subprocess.check_call(['behave %s -f progress2' % default_dirs['test']], shell=True)
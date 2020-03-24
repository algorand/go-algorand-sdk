#!/usr/bin/env python3

import subprocess
import sys

default_dirs = { 
    'features_dir': '/opt/go/src/github.com/algorand/go-algorand-sdk/test/features',
    'source': '/opt/go/src/github.com/algorand/go-algorand-sdk',
    'docker': '/opt/go/src/github.com/algorand/go-algorand-sdk/test/docker',
    'test': '/opt/go/src/github.com/algorand/go-algorand-sdk/test'
}

def setup_sdk():
    """
    Setup go cucumber environment.
    """    
    subprocess.check_call(['go generate %s/...' % default_dirs['source']], shell=True)

def test_sdk():
    sys.stdout.flush()
    subprocess.check_call(['go test'], shell=True, cwd=default_dirs['test'])
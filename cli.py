#!/bin/python

import argparse
import subprocess
parser = argparse.ArgumentParser(
                    prog='kvmcli',
                    description='A Python script for managing virtual machines in a KVM-based cluster.',
                    epilog='Enjoy')

parser.add_argument('--info', action="store_true", help='Print information about your cluster')
parser.add_argument('--apply', metavar='YAML_PATH', help='apply configuration from file')

args = parser.parse_args()

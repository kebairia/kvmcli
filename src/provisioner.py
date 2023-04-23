#!/bin/python
import os
import subprocess
from pathlib import Path
import colorama
import yaml

import ansible

def load_yaml_data(file_path):

    """Loads data from a YAML file and returns it."""

    try:
        with open(file_path, "r") as server_yaml:
            data = yaml.safe_load(server_yaml)
    except yaml.YAMLError as exc:
        print(f"Failed to parse YAML file: {exc}")
        return None
    return data
data = load_yaml_data("./servers.yml")
for vm in data['vms']:
    print(vm['provisioner'])


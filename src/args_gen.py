#!/bin/python
import yaml
import logging
from pathlib import Path

from config import config


def load_yaml_data(file_path):
    """Loads data from a YAML file and returns it."""

    try:
        with open(file_path, "r") as server_yaml:
            data = yaml.safe_load(server_yaml)
    except yaml.YAMLError as exc:
        logging.error(f"Failed to parse YAML file: {exc}")
        return None
    return data


def create_virt_install_args(index, vm):
    """Creates a dictionary of arguments for the virt-install command."""

    name = vm['info'].get('name', f"{config.IMAGE_NAME}-{index+1}")
    operating_system = str(vm['info'].get('os', 'generic'))
    ram = str(vm['info'].get("ram", config.RAM))
    cpus = str(vm['info'].get("cpus", config.CPUS))
    bridge = vm['network']['interface'].get("bridge", config.BRIDGE_NAME)
    mac_address = vm['network']['interface'].get(
        "mac_address", config.MAC_ADDRESS)
    disk_size = vm['storage']['disk'].get("size", config.SIZE_GB)
    disk = Path(config.IMAGES_PATH) / \
        f"{name}.{vm['storage']['disk']['format']}"

    virt_install_args = {
        "name": f"{name}",
        "network": f"bridge={bridge},model=virtio,mac={mac_address}",
        "disk": f"path={disk},size={disk_size}",
        "ram": f"{ram}",
        "vcpus": f"{cpus}",
        "os-variant": f"{operating_system}",
        "sound": config.SOUND_TYPE,
        "rng": config.RNG_DEVICE,
        "virt-type": config.VIRTUALIZATION_TYPE,
        "import": "",
        "wait": "0",
        "quiet": "",
        "connect": "qemu:///system"
    }
    return virt_install_args

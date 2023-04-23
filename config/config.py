#!/bin/python
import tomlkit
from typing import Dict

# Define constants for the keys in the TOML file
YAML_PATH_KEY = "yaml_path"
TEMPLATE_NAME_KEY = "template_name"
ARTIFACTS_PATH_KEY = "artifacts_path"
IMAGES_PATH_KEY = "images_path"
IMAGE_NAME_KEY = "image_name"
RAM_KEY = "ram"
CPUS_KEY = "cpus"
SOUND_TYPE_KEY = "sound_type"
RNG_DEVICE_KEY = "rng_device"
MAC_ADDRESS_KEY = "mac_address"
BRIDGE_NAME_KEY = "bridge_name"
SIZE_GB_KEY = "size_gb"
VIRTUALIZATION_TYPE_KEY = "virtualization_type"

def load_config(file_path: str) -> Dict[str, dict]:
    """
    Load the TOML configuration file at the specified path and return a dictionary
    containing the configuration values.
    """
    with open(file_path, "r") as f:
        config = tomlkit.load(f)
    return config

# Load the configuration file
config = load_config("config.cfg")

# Read the configuration values
YAML_PATH: str = config[YAML_PATH_KEY]
TEMPLATE_NAME: str = config[TEMPLATE_NAME_KEY]
ARTIFACTS_PATH: str = config["image"][ARTIFACTS_PATH_KEY]
IMAGES_PATH: str = config["image"][IMAGES_PATH_KEY]
IMAGE_NAME: str = config["image"][IMAGE_NAME_KEY]
RAM: int = config["hardware"][RAM_KEY]
CPUS: int = config["hardware"][CPUS_KEY]
SOUND_TYPE: str = config["hardware"][SOUND_TYPE_KEY]
RNG_DEVICE: str = config["hardware"][RNG_DEVICE_KEY]
MAC_ADDRESS: str = config["network"][MAC_ADDRESS_KEY]
BRIDGE_NAME: str = config["network"][BRIDGE_NAME_KEY]
SIZE_GB: int = config["disk"][SIZE_GB_KEY]
VIRTUALIZATION_TYPE: str = config["platform"][VIRTUALIZATION_TYPE_KEY]

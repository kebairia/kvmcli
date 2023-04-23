#!/bin/python
import os
import time
import shutil
from pathlib import Path

from config import config
from src.misc import print_message


def copy_image(name: str, vm: dict) -> dict:
    """
    Copies a virtual machine image to a specified path.

    Args:
        name (str): Name of the image file to copy.
        vm (dict): Dictionary containing information about the virtual machine.

    Returns:
        dict: A dictionary containing the destination path and a boolean indicating
              whether the source image file exists.
    """
    images_path = config.IMAGES_PATH
    if images_path is None:
        print_message("ERROR", f"Could not find the default images path --> {default_images_path}")
        return {'dest_image': None, 'source_image_exists': False, 'dest_image_exists': False}

    dest_image = Path(images_path) / f"{name}.{vm['storage']['disk']['format']}"
    source_image = Path(config.ARTIFACTS_PATH) / f"{vm['info']['image']}.{vm['storage']['disk']['format']}"
    if not dest_image.exists():
        if source_image.exists():
            try:
                print_message("INFO", f"Copying new VM to {dest_image}")
                shutil.copy2(str(source_image), str(dest_image))
            except shutil.SameFileError:
                print_message("WARNING", f"File already exists at {dest_image}. Skipping copy.")
        else:
            print_message("ERROR", f"File not found - {source_image}")
            return {'dest_image': None, 'source_image_exists': False, 'dest_image_exists': False}
    else:
        print_message("ERROR", f"File already exists at {dest_image}. Skipping copy.")
        return {'dest_image': None, 'source_image_exists': False, 'dest_image_exists': False}

    return {'dest_image': str(dest_image), 'source_image_exists': source_image.exists(), 'dest_image_exists': dest_image.exists() }

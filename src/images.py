#!/bin/python
import shutil
import logging
from pathlib import Path

from config import config

# from src.misc import print_message

logging.basicConfig(level=logging.INFO, format="%(levelname)s: %(message)s")


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
        logging.error(f"Could not find the default images path --> {images_path}")
        return {
            "dest_image": None,
            "source_image_exists": False,
            "dest_image_exists": False,
        }

    dest_image = Path(images_path) / f"{name}.{vm['storage']['disk']['format']}"

    source_image = (
        Path(config.ARTIFACTS_PATH)
        / f"{vm['info']['image']}.{vm['storage']['disk']['format']}"
    )

    if not dest_image.exists():
        if source_image.exists():
            try:
                logging.info(f"Copying new VM to {dest_image}")
                shutil.copy(str(source_image), str(dest_image))
            except shutil.SameFileError:
                logging.warning(f"File already exists at {dest_image}. Skipping copy.")
        else:
            logging.error(f"File not found - {source_image}")
            return {
                "dest_image": None,
                "source_image_exists": False,
                "dest_image_exists": False,
            }
    else:
        logging.error(f"File already exists at {dest_image}. Skipping copy.")
        return {
            "dest_image": None,
            "source_image_exists": False,
            "dest_image_exists": False,
        }

    return {
        "dest_image": str(dest_image),
        "source_image_exists": source_image.exists(),
        "dest_image_exists": dest_image.exists(),
    }

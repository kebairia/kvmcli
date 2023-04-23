#!/bin/python

from rich.console import Console
from rich.table import Table
from args_gen import create_virt_install_args, load_yaml_data
import config

def report(yaml_file):

    data = load_yaml_data(yaml_file)
    if data is None:
        exit()

    table = Table(title=f"{yaml_file}".upper())

    columns = ["SERVERS", "SYSTEM", "RAM", "CPUS", "BRIDGE", "MAC ADDRESS", "DISK SIZE"]

    for column in columns:
        table.add_column(column)

    for index, vm in enumerate(data['vms']):
        vm_args = create_virt_install_args(index, vm)
        row = [
            f"{vm_args['name']}",
            f"{vm['info']['os']}",
            f"{vm_args['ram']} MB",
            f"{vm_args['vcpus']}",
            f"{vm['network']['interface']['bridge']}",
            f"{vm['network']['interface']['mac_address']}",
            f"{vm['storage']['disk']['size']} GB"
        ]
        table.add_row(*row, style='bright_green')

    console = Console()
    console.print(table)

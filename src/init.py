#!/bin/python
import yaml
from config import config
import colorama
colorama.init()  
def create_template():
    template = config.TEMPLATE_NAME
    dict_file = {
            'version': 1.0,
            'vms': 
            [{'info': 
              {'name': '<NODE NAME>',
               'image': 'rocky9.1',
               'ram': 1536,
               'cpus': 1,
               'os': 'rocky9'},
              'network': 
              {'interface': 
               {'bridge': 'virbr1',
                'mac_address': '02:A3:10:00:00:XX'}},
              'storage': 
              {'disk': 
               {'size': 30,
                'type': 'SSD',
                'format': 'qcow2'}
               }
              }
             ]
            }
    with open(rf'{template}', 'w') as file:
        print(f"template file with the name {colorama.Style.DIM}{template}{colorama.Style.RESET_ALL} is created !")
        yaml.dump(dict_file, file)
create_template()

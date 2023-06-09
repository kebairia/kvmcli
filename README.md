<a name="readme-top"></a>



<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![GPL-3.0 License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]



<!-- PROJECT LOGO -->
<br />
<div align="center">
  <!-- <a href="https://github.com/kebairia/kvmcli"> -->
  <!--   <img src="images/logo.png" alt="Logo" width="80" height="80"> -->
  <!-- </a> -->

<h3 align="center">KVMcli</h3>

  <p align="center">
    A Python script for managing virtual machines in a KVM-based environment.
    <br />
    <a href="https://github.com/kebairia/kvmcli"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/kebairia/kvmcli">View Demo</a>
    ·
    <a href="https://github.com/kebairia/kvmcli/issues">Report Bug</a>
    ·
    <a href="https://github.com/kebairia/kvmcli/issues">Request Feature</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <!-- <li><a href="#prerequisites">Prerequisites</a></li> -->
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>


<!-- ABOUT THE PROJECT -->
## About The Project

| Demo                                                                                                                       |
| :--:                                                                                                                       |
| [![Deploy Demo](https://asciinema.org/a/4BaCROzRjSZVLZMgKpMeptbvi.svg)](https://asciinema.org/a/4BaCROzRjSZVLZMgKpMeptbvi)|
| Deploy with a single command (after updating the configuration files)                                                      |

<!-- Here's a blank template to get started: To avoid retyping too much info. Do a search and replace with your text editor for the following: `kebairia`, `kvmcli`, `twitter_handle`, `linkedin_username`, `email_client`, `email`, `project_title`, `project_description` -->

<a href="#readme-top"><img src="https://img.shields.io/badge/UP-ED2B2A?style=for-the-badge&logo=acclaim&logoColor=white" align="right"/> </a>



### Built With

| Technologies                         | Used for                                                         |
|:-------------------------------------|:-----------------------------------------------------------------|
| [![Python][Python.icon]][Python.url] | Python is the primary programming language used in this project. |
| [![YAML][YAML.icon]][YAML.url]       | YAML is utilized to describe the virtual machines cluster.        |
| [![TOML][TOML.icon]][TOML.url]       | TOML is used for managing configuration files in the project.    |

<a href="#readme-top"><img src="https://img.shields.io/badge/UP-ED2B2A?style=for-the-badge&logo=acclaim&logoColor=white" align="right"/> </a>


<!-- GETTING STARTED -->
## Getting Started

<!-- ### Prerequisites -->
<!-- This is an example of how to list things you need to use the software and how to install them. -->
<!-- * npm -->
<!--   ```sh -->
<!--   npm install npm@latest -g -->
<!--   ``` -->

### Installation

1. Clone the repo


   ```sh
   git clone https://github.com/kebairia/kvmcli.git
   ```
2. Install the project


   ```sh
   make install
   ```

<a href="#readme-top"><img src="https://img.shields.io/badge/UP-ED2B2A?style=for-the-badge&logo=acclaim&logoColor=white" align="right"/> </a>


<!-- USAGE EXAMPLES -->
## Usage

<!-- Use this space to show useful examples of how a project can be used. Additional screenshots, code examples and demos work well in this space. You may also link to more resources. -->

<!-- _For more examples, please refer to the [Documentation](https://example.com)_ -->

**Note**: This project is still under development, but you can use it to provision VMs with different operating systems.

The project has the following structure:

- A YAML file named `servers.yml` [^1]  for describing your cluster.
- A configuration file named `config.cfg` for assigning default values, such as the name of the main YAML file.
- The `kvmcli` command line tool is the main tool for the project.

### Configuring default values

Modify `config.cfg` with the default values that you need. You can specify the path to the server YAML file, the path to the artifacts and images directories, the name of the image that you want to use, on so on.

Here's an example of how to configure the default values in `config.cfg`:

```toml

# KVMCLI provisioner script configuration file
 
[misc]
# Number of forks to use in a multi-process application
fork = 5  

# YAML configurations
[yaml]
# Default path for YAML file
path = "servers.yml"
# Name of template YAML file used when initializing script
template = "servers.yml"

# Image configurations
[image]
artifacts_path = "/home/zakaria/dox/homelab/artifacts"
images_path = "/home/zakaria/dox/homelab/images"
image_name = "homelab"

...
```

### Launching the provisioning process
The `kvmcli` command is used for launching the provisioning process. You can use it to create a template, print information about your cluster, apply configuration from a YAML file, or ignore a specific node.

#### Creating a template
To begin, utilize the `--init` option for generating a template that will serve as a reference for creating your virtual machines."


```sh
kvmcli --init
```


```
Template file with the name `template.yml` is created !
```

This will create a template file named template.yml. The content of the template will be like the following:

```yaml
version: 1.0
vms:
- info:
    name: node1
    image: rocky9.1
    ram: 1536
    cpus: 1
    os: rocky9
  network:
    interface:
      bridge: virbr1
      mac_address: 02:A3:10:00:00:10
  storage:
    disk:
      format: qcow2
      size: 30
      type: SSD
```

You can use the `--info` option to print the content of the template file in a pretty table:
It uses the default value of `template_name` from the `config.cfg` configuration file


```sh
kvmcli --info
```

If you want to use another file as a reference, use the `-f` or `--file` option:


```sh
kvmcli --info -f template.yml
```


```
                                  TEMPLATE.YML                                 
 ┏━━━━━━━━━┳━━━━━━━━┳━━━━━━━━━┳━━━━━━┳━━━━━━━━┳━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━┓
 ┃ SERVERS ┃ SYSTEM ┃ RAM     ┃ CPUS ┃ BRIDGE ┃ MAC ADDRESS       ┃ DISK SIZE ┃
 ┡━━━━━━━━━╇━━━━━━━━╇━━━━━━━━━╇━━━━━━╇━━━━━━━━╇━━━━━━━━━━━━━━━━━━━╇━━━━━━━━━━━┩
 │ node1   │ rocky9 │ 1536 MB │ 1    │ virbr1 │ 02:A3:10:00:00:10 │ 30 GB     │
 └─────────┴────────┴─────────┴──────┴────────┴───────────────────┴───────────┘
```
#### Applying configuration from a YAML file
When you're happy with the result, you can start provisioning using the `-a` or `--apply` option:


```sh
kvmcli --apply -f template.yml
```


```
 INFO: Copying new VM to /home/zakaria/dox/homelab/images/node1.qcow2
 INFO: Provisioning a new VM named node1
 
 INFO: All VMs provisioned successfully!
```

The `--ignore` flag is used to exclude specific nodes from the provisioning process when applying a configuration from a YAML file using the `kvmcli` tool. 

For example, running `kvmcli --apply -f template.yml --ignore node1` will apply the configuration defined in `template.yml`, but exclude the `node1` node from being provisioned.


```sh
kvmcli --apply -f template.yml --ignore node1
```

-h, --help Show the help message and exit.


```
usage: kvmcli [-h] [-I] [-i] [-a] [-f YAML_FILE] [--ignore NODE_NAME]

A Python script for managing virtual machines in a KVM-based environment.

options:
  -h, --help            show this help message and exit
  -I, --info            Print information about your cluster
  -i, --init            Create template file
  -a, --apply           apply configuration from YAML_FILE
  -f YAML_FILE, --file YAML_FILE
                        Specify a yaml file
  --ignore NODE_NAME    Ignore NODE NAME

Enjoy
```

<a href="#readme-top"><img src="https://img.shields.io/badge/UP-ED2B2A?style=for-the-badge&logo=acclaim&logoColor=white" align="right"/> </a>

<!-- ROADMAP -->
## Roadmap

- [x] Print report for the cluster
- [x] Provision multiple VMs with different Operating Systems
- [x] Enhancing command line tool
- [x] Multiprocessing support
- [ ] Provisioning system
- [ ] Logging system

See the [open issues](https://github.com/kebairia/kvmcli/issues) for a full list of proposed features (and known issues).

<a href="#readme-top"><img src="https://img.shields.io/badge/UP-ED2B2A?style=for-the-badge&logo=acclaim&logoColor=white" align="right"/> </a>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<a href="#readme-top"><img src="https://img.shields.io/badge/UP-ED2B2A?style=for-the-badge&logo=acclaim&logoColor=white" align="right"/> </a>

<!-- LICENSE -->
## License

Distributed under the GPL-3.0 License. See `LICENSE.md` for more information.

<a href="#readme-top"><img src="https://img.shields.io/badge/UP-ED2B2A?style=for-the-badge&logo=acclaim&logoColor=white" align="right"/> </a>



<!-- CONTACT -->
## Contact

<!-- Your Name - [@twitter_handle](https://twitter.com/twitter_handle) - email@email_client.com -->

Project Link: [https://github.com/kebairia/kvmcli](https://github.com/github_username/repo_name)


<!-- ACKNOWLEDGMENTS -->
## Acknowledgments 

* [README Template](https://github.com/othneildrew/Best-README-Template) 
<!--* []() -->
<!--* []() -->

<a href="#readme-top"><img src="https://img.shields.io/badge/UP-ED2B2A?style=for-the-badge&logo=acclaim&logoColor=white" align="right"/> </a>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/kebairia/kvmcli.svg?style=for-the-badge
[contributors-url]: https://github.com/kebairia/kvmcli/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/kebairia/kvmcli.svg?style=for-the-badge
[forks-url]: https://github.com/kebairia/kvmcli/network/members
[stars-shield]: https://img.shields.io/github/stars/kebairia/kvmcli.svg?style=for-the-badge
[stars-url]: https://github.com/kebairia/kvmcli/stargazers
[issues-shield]: https://img.shields.io/github/issues/kebairia/kvmcli.svg?style=for-the-badge
[issues-url]: https://github.com/kebairia/kvmcli/issues
[license-shield]: https://img.shields.io/github/license/kebairia/kvmcli.svg?style=for-the-badge
[license-url]: https://github.com/kebairia/kvmcli/blob/main/LICENSE.md
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/zakaria-kebairia

[blog-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[blog-url]: https://kebairia.github.io-kebairia

[Python.icon]: https://img.shields.io/badge/python-4584b6?style=for-the-badge&logo=python&logoColor=white
[Python.url]:  https://www.python.org/

[YAML.icon]: https://img.shields.io/badge/yaml-red?style=for-the-badge&logo=yaml&logoColor=whte
[YAML.url]: https://yaml.org/

[TOML.icon]: https://img.shields.io/badge/toml-9d4626?style=for-the-badge&logo=toml&logoColor=whte
[TOML.url]: https://toml.io/

<!-- [UP.icon]: <a href="#readme-top"><img src="https://img.shields.io/badge/UP-ED2B2A?style=for-the-badge&logo=acclaim&logoColor=white" align="right"/> </a> -->
<!-- [UP.url]:  https://github.com/kebairia/kvmcli#readme-top -->

### Footnotes

[^1]: Or any name you like, you can change that by changing the `path` and `template_name` keys under `[yaml]` section in `config.cfg`

    ```yaml
    [yaml]
    # Default path for YAML file
    path = "servers.yml"
    # Name of template YAML file used when initializing script
    template = "servers.yml"
    ```

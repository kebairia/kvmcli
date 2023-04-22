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
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]



<!-- PROJECT LOGO -->
<br />
<div align="center">
  <!-- <a href="https://github.com/kebairia/kvmcli"> -->
  <!--   <img src="images/logo.png" alt="Logo" width="80" height="80"> -->
  <!-- </a> -->

<h3 align="center">KVMcli</h3>

  <p align="center">
    A Python script for managing virtual machines in a KVM-based cluster.
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
    <!-- <li><a href="#roadmap">Roadmap</a></li> -->
    <!-- <li><a href="#contributing">Contributing</a></li> -->
    <li><a href="#license">License</a></li>
    <!-- <li><a href="#contact">Contact</a></li> -->
    <!-- <li><a href="#acknowledgments">Acknowledgments</a></li> -->
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

| Demo                                                                                                                       |
| :--:                                                                                                                       |
| [![Deploy Demo](https://asciinema.org/a/0yJKkTA0pFMSjjxdrOytsZnos.svg)](https://asciinema.org/a/0yJKkTA0pFMSjjxdrOytsZnos) |
| Deploy with a single command (after updating the configuration files)                                                      |

<!-- Here's a blank template to get started: To avoid retyping too much info. Do a search and replace with your text editor for the following: `kebairia`, `kvmcli`, `twitter_handle`, `linkedin_username`, `email_client`, `email`, `project_title`, `project_description` -->

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

[![Python][Python.icon]][Python.url]
[![YAML][YAML.icon]][YAML.url]
[![TOML][TOML.icon]][TOML.url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

This is an example of how you may give instructions on setting up your project locally.
To get a local copy up and running follow these simple example steps.

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
2. Install required packages
   ```sh
   pip install -r requirements.txt
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

<!-- Use this space to show useful examples of how a project can be used. Additional screenshots, code examples and demos work well in this space. You may also link to more resources. -->

<!-- _For more examples, please refer to the [Documentation](https://example.com)_ -->
Here's a sample usage section based on the bullet points you provided:

**Note** This project is still under development, but you can use it to provision VMs with different operating systems.


1. The project has the following structure:
   - A YAML file for describing your cluster, named `servers.yml` by default.
   - A config file named `config.cfg` for assigning default values, such as the name of the main YAML file.
   - `kvmcli` is the main command line tool for the project.

2. Modify the `servers.yml` file with the nodes (VMs) that you want. You can also choose another name for this file and update it in the `config.cfg` file.

**Example**:
```yaml
version: 1.0

vms:
  - info:
      name: admin1
      image: ubuntu22.04
      ram: 2048
      cpus: 2
      os: ubuntu22.04
    network:
      interface:
        bridge: virbr1
        mac_address: "02:A3:10:00:00:02"
    storage:
      disk:
        size: 30
        type: SSD
        format: qcow2
```

3. Modify `config.cfg` with the default values that you need. 
```toml
# TOML Configuration file for provisioner script

# Path to server YAML file
yaml_path = "servers.yml"

# Image configurations
[image]
artifacts_path = "/home/zakaria/dox/homelab/artifacts"
images_path = "/home/zakaria/dox/homelab/images"
image_name = "homelab"
```

4. `kvmcli` is the command for launching the provisioning process. Currently, it is just a command for provisioning VMs from `servers.yml`. 
    For other feature, you can test each function by itself by `python <function>.py`.


Here's an example of how to provision VMs using `kvmcli`:
``` sh
./kvmcli
```
This will create the VMs specified in the `servers.yml` file. You can then connect to the VMs using a remote desktop client or SSH.


To get a table for all the VMs listed in `servers.yml` execute:
``` sh
py ./info.py
```
```
                                    SERVERS.YML
┏━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━┳━━━━━━┳━━━━━━━━┳━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━┓
┃ SERVERS ┃ SYSTEM      ┃ RAM     ┃ CPUS ┃ BRIDGE ┃ MAC ADDRESS       ┃ DISK SIZE ┃
┡━━━━━━━━━╇━━━━━━━━━━━━━╇━━━━━━━━━╇━━━━━━╇━━━━━━━━╇━━━━━━━━━━━━━━━━━━━╇━━━━━━━━━━━┩
│ admin1  │ ubuntu22.04 │ 2048 MB │ 2    │ virbr1 │ 02:A3:10:00:00:02 │ 30 GB     │
└─────────┴─────────────┴─────────┴──────┴────────┴───────────────────┴───────────┘

```


<!-- For more information on the available commands, run `kvmcli --help`. -->





<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ROADMAP -->
## Roadmap
    
- [x] Print report for the cluster
- [x] Provision multiple VMs with different Operating Systems
- [ ] Enhancing command line tool
- [ ] Logging system

See the [open issues](https://github.com/kebairia/kvmcli/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>



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

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
<!-- ## License -->

<!-- Distributed under the MIT License. See `LICENSE.txt` for more information. -->

<!-- <p align="right">(<a href="#readme-top">back to top</a>)</p> -->



<!-- CONTACT -->
## Contact

Your Name - [@twitter_handle](https://twitter.com/twitter_handle) - email@email_client.com

Project Link: [https://github.com/kebairia/kvmcli](https://github.com/github_username/repo_name)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->
<!-- ## Acknowledgments -->

<!-- * []() -->
<!-- * []() -->
<!-- * []() -->

<!-- <p align="right">(<a href="#readme-top">back to top</a>)</p> -->



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
[license-url]: https://github.com/kebairia/kvmcli/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/zakaria-kebairia
[homelab-arch-screenshot]: ./arch.png

[Python.icon]: https://img.shields.io/badge/python-4584b6?style=for-the-badge&logo=python&logoColor=ffde57
[Python.url]:  https://www.python.org/

[YAML.icon]: https://img.shields.io/badge/yaml-red?style=for-the-badge&logo=yaml&logoColor=whte
[YAML.url]: https://yaml.org/

[TOML.icon]: https://img.shields.io/badge/toml-9d4626?style=for-the-badge&logo=toml&logoColor=whte
[TOML.url]: https://toml.io/

[UP.icon]: https://img.shields.io/badge/UP-ED2B2A?style=for-the-badge&logo=^&logoColor=ffde57
[UP.url]:  https://github.com/kebairia/kvmcli#readme-top


# kvmcli

**kvmcli** is a command-line interface (CLI) tool for managing KVM (Kernel-based Virtual Machine) virtual machines. Inspired by Kubernetes’ kubectl, kvmcli provides a declarative way to create, delete, and manage VMs using YAML configuration files.

## Features

- **Declarative VM Management:** Define virtual machines using a YAML manifest.
- **Create and Delete VMs:** Easily provision new VMs or remove existing ones.
- **List VM Information:** Retrieve details such as CPU, memory, disk, network status, and more.
- **Overlay Disk Creation:** Automatically create overlay disk images for VMs.
- **Modular Commands:** Built using Cobra for a structured and extendable CLI.

## Prerequisites

- [Go](https://golang.org) 1.16 or later.
- [libvirt](https://libvirt.org/) installed and properly configured.
- [qemu-img](https://www.qemu.org/docs/master/tools/qemu-img.html) for handling disk image overlays.
- A valid YAML configuration file to define your VM(s).

## Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/kebairia/kvmcli.git
cd kvmcli
go build -o kvmcli .
```

Alternatively, install using:

```bash
go install github.com/kebairia/kvmcli@latest
```

## Usage

kvmcli uses a command structure similar to kubectl. Below are some example commands:

### Create VM(s)

Provision VM(s) defined in your YAML configuration file:

```bash
kvmcli create -f /path/to/vm-config.yaml
```

### Delete VM(s)

Delete VM(s) as defined in your configuration file. Use the --all flag to delete all VMs:

```bash

kvmcli delete -f /path/to/vm-config.yaml
# Or, to delete all VMs:
kvmcli delete --all
```

### List VM Information

Display information about your VMs:

```bash
kvmcli get vm -f /path/to/vm-config.yaml
```

Additional subcommands include:

- Snapshots: kvmcli get snapshot

- Networks: kvmcli get network

### Initialize a YAML Template

Generate a template file with a sample VM definition:

```bash
kvmcli init
```

This command creates a starting YAML file to help you define your VM(s).

## Configuration File Format

The YAML configuration file defines the virtual machine properties. An example configuration:

```yaml
apiVersion: kvmcli/v1
kind: VirtualMachine
metadata:
  name: myvm01
  namespace: homelab
  labels:
    role: role01
    environment: production
spec:
  cpu: 2
  memory: 4096
  image: "rocky95"
  disk:
    size: "20G"
  network:
    name: homelab
    macAddress: "02:A3:10:00:01:01"
  autostart: true
```

This example is parsed by the configuration loader defined in `load_manifest.go`.

## Project Structure

```graphql

kvmcli/
├── cmd/                        # Cobra command definitions
│   ├── root.go                 # Root command and subcommand registration.
│   ├── create_cmd.go           # VM creation command.
│   ├── delete_cmd.go           # VM deletion command.
│   └── list_cmd.go             # VM info retrieval command.
├── internal/
│   ├── config/                 # YAML configuration parsing.
│   ├── logger/                 # Logging utility
│   ├── operations/             # libvirt connection and initialization.
│   │   └── init_cmd.go         # Command for initializing a YAML template.
│   └── operations/vms/         # VM operations
│       ├── create.go           # Creating VMs and overlay disk images.
│       ├── delete.go           # Deleting VMs and cleaning up disk images.
│       ├── list.go             # Listing and retrieving VM information.
│       └── prepare_vm.go       # Preparing VMs by creating overlay images.
└── README.md
```

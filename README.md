# kvmcli

**kvmcli** is a modern, declarative CLI tool for managing KVM (Kernel-based Virtual Machine) infrastructure. Inspired by Terraform and Kubernetes, it uses **HCL (HashiCorp Configuration Language)** to define infrastructure as code.

## Key Features

- **Infrastructure as Code**: Define VMs, Networks, and Storage Pools using simple, readable HCL files.
- **State Management**: Tracks resource state in a local SQLite database to prevent drift and manage lifecycle.
- **Declarative Networking**: Configure bridge networks and DHCP ranges easily.
- **Data Sources**: Reference existing resources (like pre-existing networks or storage pools) using `data` blocks.
- **Cluster Management**: Group VMs into clusters with defined start/stop orders.

## Prerequisites

- **Linux** with KVM/QEMU enabled.
- **libvirt** daemon running.
- **Go** 1.22+ (to build from source).

## Installation

```bash
git clone https://github.com/kebairia/kvmcli.git
cd kvmcli
go build -o kvmcli .
sudo mv kvmcli /usr/local/bin/
```

## Quick Start

### 1. Define your Infrastructure

Create a `main.hcl` file:

```hcl
# Define a Storage Pool
store "default" {
  namespace = "homelab"

  # Images managed by this store
  image "ubuntu-22.04" {
    url = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.disk1.img"
  }
}

# Define a Network with DHCP
network "services" {
  namespace = "homelab"
  mode      = "nat"
  cidr      = "192.168.100.0/24"

  dhcp {
    start = "192.168.100.10"
    end   = "192.168.100.200"
  }
}

# Define a Virtual Machine
vm "web-server-01" {
  namespace = "homelab"
  cpu       = 2
  memory    = 4096 # MB

  # Reference the store and network defined above
  image     = "ubuntu-22.04"
  # store     = store.default.name
  # network   = network.services.name
  store     = store.default
  network   = network.services
}
```

### 2. Apply Configuration

Provision your resources:

```bash
kvmcli create -f main.hcl
```

### 3. Manage Resources

List created resources:

```bash
kvmcli get vm
kvmcli get network
# or
kvmcli get net
```

Delete resources:

```bash
kvmcli delete -f main.hcl
# Or delete all resources tracked in the database
kvmcli delete --all
```

## Advanced Usage

### Data Sources

Reference resources that already exist in the database but are not defined in the current file. This is useful for sharing resources across multiple HCL files.

```hcl
# Look up an existing network named "default"
data "network" "default" {}

vm "worker-01" {
  # Use the looked-up network name
  network = data.network.default
  # ...
}
```

## Project Structure

- `cmd/`: Entry points and CLI command definitions (Cobra).
- `internal/config/`: HCL parser and configuration structs (`vms.Config`, `network.Config`, etc.).
- `internal/database/`: SQLite state management (`database.VirtualMachine`, `database.Network`).
- `internal/network/`: Libvirt network management logic.
- `internal/vms/`: VM lifecycle management.
- `internal/store/`: Storage pool and image management.

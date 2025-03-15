package utils

import (
	"encoding/xml"
)

// Define constants for reusable values
const (
	DomainTypeKVM   = "kvm"
	ArchX86_64      = "x86_64"
	MachineQ35      = "pc-q35-7.2"
	BootDeviceHD    = "hd"
	DiskTypeFile    = "file"
	DiskDeviceDisk  = "disk"
	DriverNameQEMU  = "qemu"
	DiskFormatQCOW2 = "qcow2"
	TargetDevVDA    = "vda"
	VirtIO          = "virtio"
	NetTypeNetwork  = "network"
	GraphicsTypeVNC = "vnc"
)

// Domain represents the root domain element
type Domain struct {
	XMLName xml.Name `xml:"domain"`
	Type    string   `xml:"type,attr"`
	Name    string   `xml:"name"`
	Memory  Memory   `xml:"memory"`
	VCPU    VCPU     `xml:"vcpu"`
	OS      OS       `xml:"os"`
	CPU     CPU      `xml:"cpu"`
	Devices Devices  `xml:"devices"`
}

// Memory defines the memory configuration
type Memory struct {
	Unit  string `xml:"unit,attr"`
	Value int    `xml:",chardata"`
}

// VCPU defines the vCPU configuration
type VCPU struct {
	Placement string `xml:"placement,attr"`
	Value     int    `xml:",chardata"`
}

// OS represents the OS configuration
type OS struct {
	Type OSType `xml:"type"`
	Boot Boot   `xml:"boot"`
}

// OSType represents the OS type attributes
type OSType struct {
	Arch    string `xml:"arch,attr"`
	Machine string `xml:"machine,attr"`
	Value   string `xml:",chardata"`
}

// Boot represents the boot configuration
type Boot struct {
	Dev string `xml:"dev,attr"`
}

// CPU represents the CPU configuration
type CPU struct {
	Mode       string `xml:"mode,attr"`
	Check      string `xml:"check,attr"`
	Migratable string `xml:"migratable,attr"`
}

// Devices holds all device configurations
type Devices struct {
	Emulator  string    `xml:"emulator"`
	Disk      Disk      `xml:"disk"`
	Interface Interface `xml:"interface"`
	Channel   Channel   `xml:"channel"`
	Serial    Serial    `xml:"serial"`
	Console   Console   `xml:"console"`
	Graphics  Graphics  `xml:"graphics"`
}

// Disk represents the disk configuration
type Disk struct {
	Type   string     `xml:"type,attr"`
	Device string     `xml:"device,attr"`
	Driver DiskDriver `xml:"driver"`
	Source DiskSource `xml:"source"`
	Target DiskTarget `xml:"target"`
}

// DiskDriver represents the disk driver configuration
type DiskDriver struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

// DiskSource represents the source file for the disk
type DiskSource struct {
	File string `xml:"file,attr"`
}

// DiskTarget represents the target device
type DiskTarget struct {
	Dev string `xml:"dev,attr"`
	Bus string `xml:"bus,attr"`
}

// Interface represents the network interface configuration
type Interface struct {
	Type   string     `xml:"type,attr"`
	MAC    MACAddress `xml:"mac"`
	Source NetSource  `xml:"source"`
	Model  NetModel   `xml:"model"`
}

// MACAddress represents the MAC address of the interface
type MACAddress struct {
	Address string `xml:"address,attr"`
}

// NetSource represents the network source
type NetSource struct {
	Network string `xml:"network,attr"`
}

// NetModel represents the network model type
type NetModel struct {
	Type string `xml:"type,attr"`
}

// Channel represents the SPICE channel
type Channel struct {
	Type    string        `xml:"type,attr"`
	Target  ChannelTarget `xml:"target"`
	Address VirtioAddress `xml:"address"`
}

// ChannelTarget represents the target for a SPICE channel
type ChannelTarget struct {
	Type string `xml:"type,attr"`
	Name string `xml:"name,attr"`
}

// VirtioAddress represents the Virtio serial address
type VirtioAddress struct {
	Type       string `xml:"type,attr"`
	Controller string `xml:"controller,attr"`
	Bus        string `xml:"bus,attr"`
	Port       string `xml:"port,attr"`
}

type Serial struct {
	Type   string       `xml:"type,attr"`
	Target SerialTarget `xml:"target"`
}
type SerialTarget struct {
	Port string `xml:"port,attr"`
}

type Console struct {
	Type   string        `xml:"type,attr"`
	Target ConsoleTarget `xml:"target"`
}
type ConsoleTarget struct {
	Type string `xml:"type,attr"`
	Port string `xml:"port,attr"`
}

// Graphics represents the graphics configuration
type Graphics struct {
	Type     string         `xml:"type,attr"`
	AutoPort string         `xml:"autoport,attr"`
	Listen   GraphicsListen `xml:"listen"`
	Image    ImageSettings  `xml:"image"`
}

// GraphicsListen represents the graphics listen type
type GraphicsListen struct {
	Type string `xml:"type,attr"`
}

// ImageSettings represents image compression settings
type ImageSettings struct {
	Compression string `xml:"compression,attr"`
}

// Constructor for new Domains
func NewDomain(name string, mem int, cpu int, source string, mac_address string) Domain {
	return Domain{
		Type: DomainTypeKVM,
		Name: name,
		Memory: Memory{
			Unit:  "MiB",
			Value: mem,
		},
		VCPU: VCPU{
			Placement: "static",
			Value:     cpu,
		},
		OS: OS{
			Type: OSType{
				Arch:    ArchX86_64,
				Machine: MachineQ35,
				Value:   "hvm",
			},
			Boot: Boot{
				Dev: BootDeviceHD,
			},
		},
		CPU: CPU{
			Mode:       "host-passthrough",
			Check:      "none",
			Migratable: "on",
		},
		Devices: Devices{
			Emulator: "/usr/bin/qemu-system-x86_64",
			Disk: Disk{
				Type:   DiskTypeFile,
				Device: DiskDeviceDisk,
				Driver: DiskDriver{
					Name: DriverNameQEMU,
					Type: DiskFormatQCOW2,
				},
				Source: DiskSource{
					File: source,
				},
				Target: DiskTarget{
					Dev: TargetDevVDA,
					Bus: VirtIO,
				},
			},
			Interface: Interface{
				Type: NetTypeNetwork,
				MAC: MACAddress{
					Address: mac_address,
				},
				Source: NetSource{
					Network: "homelab",
				},
				Model: NetModel{
					Type: VirtIO,
				},
			},
			Channel: Channel{
				Type: "spicevmc",
				Target: ChannelTarget{
					Type: "virtio",
					Name: "com.redhat.spice.0",
				},
				Address: VirtioAddress{
					Type:       "virtio-serial",
					Controller: "0",
					Bus:        "0",
					Port:       "2",
				},
			},
			Serial: Serial{
				Type: "pty",
				Target: SerialTarget{
					Port: "0",
				},
			},
			Console: Console{
				Type: "pty",
				Target: ConsoleTarget{
					Type: "serial",
					Port: "0",
				},
			},
			Graphics: Graphics{
				Type:     "spice",
				AutoPort: "yes",
				Listen: GraphicsListen{
					Type: "address",
				},
				Image: ImageSettings{
					Compression: "off",
				},
			},
		},
	}
}

// GenerateXML method for Domain struct
func (d *Domain) GenerateXML() ([]byte, error) {
	return xml.MarshalIndent(d, "", "  ")
}

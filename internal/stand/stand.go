package stand

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/digitalocean/go-libvirt"
	"gitlab.kuberfleet.io/standman/internal/config"
	"libvirt.org/go/libvirtxml"
)

func NodeUp(socket string, node *config.Node) error {
	c, err := net.DialTimeout("unix", socket, 2*time.Second)
	if err != nil {
		return fmt.Errorf("can't up domain: %w", err)
	}

	l := libvirt.New(c)
	if err := l.Connect(); err != nil {
		return fmt.Errorf("can't up domain: %w", err)
	}

	s := nodeToSpec(node)
	x, err := s.Marshal()
	if err != nil {
		return fmt.Errorf("can't up domain: %w", err)
	}

	d, err := l.DomainDefineXML(x)
	if err != nil {
		return fmt.Errorf("can't up domain: %w", err)
	}

	err = l.DomainCreate(d)
	if err != nil {
		return fmt.Errorf("can't up domain: %w", err)
	}

	if err := l.Disconnect(); err != nil {
		return fmt.Errorf("can't up domain: %w", err)
	}

	return nil
}

func NodeDown(socket string, node *config.Node) error {
	c, err := net.DialTimeout("unix", socket, 2*time.Second)
	if err != nil {
		return fmt.Errorf("can't down domain: %w", err)
	}

	l := libvirt.New(c)
	if err := l.Connect(); err != nil {
		return fmt.Errorf("can't down domain: %w", err)
	}

	s := nodeToSpec(node)
	x, err := s.Marshal()
	if err != nil {
		return fmt.Errorf("can't up domain: %w", err)
	}

	d, err := l.DomainDefineXML(x)
	if err != nil {
		return fmt.Errorf("can't up domain: %w", err)
	}

	err = l.DomainDestroy(d)
	if err != nil {
		return fmt.Errorf("can't up domain: %w", err)
	}

	if err := l.Disconnect(); err != nil {
		return fmt.Errorf("can't down domain: %w", err)
	}

	return nil
}

func nodeToSpec(node *config.Node) *libvirtxml.Domain {
	disks := make([]libvirtxml.DomainDisk, 0)
	for _, v := range node.Disks {
		disks = append(
			disks,
			libvirtxml.DomainDisk{
				Driver: &libvirtxml.DomainDiskDriver{
					Type: "qcow2",
				},
				Source: &libvirtxml.DomainDiskSource{
					File: &libvirtxml.DomainDiskSourceFile{
						File: fmt.Sprintf(
							"%s/%s",
							pwd(),
							v.Source,
						),
					},
				},
				Target: &libvirtxml.DomainDiskTarget{
					Dev: v.Target,
				},
				IOTune: &libvirtxml.DomainDiskIOTune{
					ReadIopsSec:  uint64(v.Read),
					WriteIopsSec: uint64(v.Write),
				},
			},
		)
	}

	interfaces := make([]libvirtxml.DomainInterface, 0)
	for _, v := range node.Networks {
		interfaces = append(
			interfaces,
			libvirtxml.DomainInterface{
				Source: &libvirtxml.DomainInterfaceSource{
					Bridge: &libvirtxml.DomainInterfaceSourceBridge{
						Bridge: v.Source,
					},
				},
				/*Target: &libvirtxml.DomainInterfaceTarget{
					Dev: v.Target,
				},*/
			},
		)
	}

	return &libvirtxml.Domain{
		Type: "qemu",
		UUID: node.UUID,
		Name: node.Name,
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{
				Type:    "hvm",
				Arch:    "",
				Machine: "",
			},
		},
		VCPU: &libvirtxml.DomainVCPU{
			Value: uint(node.Resources.CPU.Cores * node.Resources.CPU.Threads),
		},
		CPU: &libvirtxml.DomainCPU{
			Topology: &libvirtxml.DomainCPUTopology{
				Sockets: 1,
				Cores:   node.Resources.CPU.Cores,
				Threads: node.Resources.CPU.Threads,
			},
		},
		Memory: &libvirtxml.DomainMemory{
			Value: uint(node.Resources.Memory),
			Unit:  "B",
		},
		Devices: &libvirtxml.DomainDeviceList{
			Disks:      disks,
			Interfaces: interfaces,
			Graphics: []libvirtxml.DomainGraphic{
				{
					VNC: &libvirtxml.DomainGraphicVNC{
						Port: node.Display.VNC.Port,
					},
				},
			},
		},

		// BlockIOTune: &libvirtxml.DomainBlockIOTune{
		// 	Device: []libvirtxml.DomainBlockIOTuneDevice{
		// 		{
		// 			Path: "/dev/sda",
		// 			ReadIopsSec: 1000,
		// 			WriteIopsSec: 1000,
		// 		},
		// 	},
		// },
	}
}

func pwd() string {
	e, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return filepath.Dir(e)
}

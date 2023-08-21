package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

var ErrConfigurationNotExists = errors.New("сonfiguration not exists")

type CPU struct {
	Cores   int `json:"cores"`
	Threads int `json:"threads"`
}

type Memory int

type Resources struct {
	CPU    CPU    `json:"cpu"`
	Memory Memory `json:"memory"`
}

type Disk struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Read   int    `json:"read"`
	Write  int    `json:"write"`
}

type Disks []Disk

type Network struct {
	Source  string `json:"source"`
	Target  string `json:"target"`
	Address string `json:"address"`
}

type Networks []Network

type VNC struct {
	Port int `json:"port"`
}

type Display struct {
	VNC VNC `json:"vnc"`
}

type Node struct {
	UUID          string                 `json:"uuid"`
	Name          string                 `json:"name"`
	Image         string                 `json:"image"`
	Configuration map[string]interface{} `json:"configuration"`
	Resources     Resources              `json:"resources"`
	Disks         Disks                  `json:"disks"`
	Networks      Networks               `json:"networks"`
	Display       Display                `json:"display"`
}

type Nodes []Node

type Bridge struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Bridges []Bridge

type Config struct {
	Version int     `json:"version"`
	Nodes   Nodes   `json:"nodes"`
	Bridges Bridges `json:"bridges"`
}

func Read(path string) (*Config, error) {
	var c Config
	if !isExist(path) {
		return nil, ErrConfigurationNotExists
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

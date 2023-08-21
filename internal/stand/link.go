package stand

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

func BridgeUp(name string, address string) error {
	la := netlink.NewLinkAttrs()
	la.Name = name
	err := netlink.LinkAdd(
		&netlink.Bridge{LinkAttrs: la},
	)
	if err != nil {
		return fmt.Errorf("can't up bridge: %w", err)
	}

	br, err := netlink.LinkByName(name)
	if err != nil {
		return fmt.Errorf("can't up bridge: %w", err)
	}

	err = netlink.LinkSetUp(br)
	if err != nil {
		return fmt.Errorf("can't up bridge: %w", err)
	}

	addr, err := netlink.ParseAddr(address)
	if err != nil {
		return fmt.Errorf("can't up bridge: %w", err)
	}

	err = netlink.AddrAdd(br, addr)
	if err != nil {
		return fmt.Errorf("can't up bridge: %w", err)
	}

	return nil
}

func BridgeDown(name string) error {
	br, err := netlink.LinkByName(name)
	if err != nil {
		return fmt.Errorf("can't down bridge: %w", err)
	}

	err = netlink.LinkSetDown(br)
	if err != nil {
		return fmt.Errorf("can't down bridge: %w", err)
	}

	err = netlink.LinkDel(br)
	if err != nil {
		return fmt.Errorf("can't down bridge: %w", err)
	}

	return nil
}

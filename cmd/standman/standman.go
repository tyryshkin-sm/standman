package main

import (
	"fmt"

	"github.com/spf13/cobra"

	config "gitlab.kuberfleet.io/standman/internal/config"
	stand "gitlab.kuberfleet.io/standman/internal/stand"
)

var rootCmd = &cobra.Command{
	Use:   "standman",
	Short: "QEMU Stand Manager",
	Long:  "Manages QEMU nodes declared in configuration",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Up stand",
	Long:  "Starts QEMU nodes, creates bridge network and connects nodes to network.",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			fmt.Println(err)
			return
		}

		socketPath, err := cmd.Flags().GetString("socket")
		if err != nil {
			fmt.Println(err)
			return
		}

		c, err := config.Read(configPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, b := range c.Bridges {
			err = stand.BridgeUp(b.Name, b.Address)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		for _, n := range c.Nodes {
			err = stand.NodeUp(socketPath, &n)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Down stand",
	Long:  "Disconnects nodes from network, stop QEMU nodes, deletes bridge network.",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			fmt.Println(err)
			return
		}

		socketPath, err := cmd.Flags().GetString("socket")
		if err != nil {
			fmt.Println(err)
			return
		}

		c, err := config.Read(configPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, n := range c.Nodes {
			err = stand.NodeDown(socketPath, &n)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		for _, b := range c.Bridges {
			err = stand.BridgeDown(b.Name)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	},
}

func main() {
	rootCmd.PersistentFlags().String("config", "standman.yaml", "configuration file")
	rootCmd.PersistentFlags().String("socket", "/var/run/libvirt/libvirt-sock", "libvirt socket")

	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)

	rootCmd.Execute()
}

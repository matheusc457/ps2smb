package cmd

import (
	"fmt"
	"os"

	"github.com/matheusc457/ps2smb/internal/network"
	"github.com/spf13/cobra"
)

var interfacesCmd = &cobra.Command{
	Use:   "interfaces",
	Short: "List available network interfaces",
	Long:  `Displays all network interfaces with their IP addresses. Useful for selecting which interface to use with the --interface flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runInterfaces(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(interfacesCmd)
}

func runInterfaces() error {
	interfaces, err := network.ListInterfaces()
	if err != nil {
		return fmt.Errorf("failed to list interfaces: %v", err)
	}

	if len(interfaces) == 0 {
		fmt.Println("No active network interfaces found")
		return nil
	}

	fmt.Println("Available Network Interfaces:")
	fmt.Println("=============================")
	fmt.Println()

	for name, ip := range interfaces {
		fmt.Printf("  %s\n", name)
		fmt.Printf("    IP Address: %s\n", ip)
		fmt.Println()
	}

	fmt.Println("Usage:")
	fmt.Printf("  ps2smb info --interface <name>\n")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Printf("  ps2smb info --interface enp3s0\n")

	return nil
}

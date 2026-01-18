package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ps2smb",
	Short: "Configure SMB shares for PlayStation 2 network gaming",
	Long: `ps2smb is a command-line tool that automates the setup and management
of Samba servers optimized for PlayStation 2 network gaming via OPL.

It handles server configuration, network detection, and provides
step-by-step instructions for connecting your PS2.`,
	Example: `  # Initialize Samba server for PS2
  sudo ps2smb init

  # Show connection information
  sudo ps2smb info

  # Show connection info using specific network interface
  sudo ps2smb info --interface enp3s0

  # List available network interfaces
  ps2smb interfaces`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Remove default flags that aren't needed
	rootCmd.CompletionOptions.DisableDefaultCmd = false
}

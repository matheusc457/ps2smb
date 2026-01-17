package cmd

import (
	"fmt"
	"os"

	"github.com/matheusc457/ps2smb/internal/config"
	"github.com/matheusc457/ps2smb/internal/network"
	"github.com/matheusc457/ps2smb/internal/samba"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show connection information for PS2",
	Long:  `Displays IP address, share details, and instructions for configuring OPL on your PlayStation 2.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runInfo(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo() error {
	// Check if ps2smb is configured
	if !config.Exists() {
		return fmt.Errorf("ps2smb is not configured. Please run 'sudo ps2smb init' first")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}

	// Get local IP
	ip, err := network.GetLocalIP()
	if err != nil {
		return fmt.Errorf("failed to detect IP address: %v", err)
	}

	// Check Samba status
	sambaRunning := samba.IsSambaRunning()
	statusSymbol := "✓"
	statusText := "Running"
	if !sambaRunning {
		statusSymbol = "✗"
		statusText = "Not Running"
	}

	// Format SMB path
	smbPath := network.FormatSMBPath(ip, cfg.ShareName)

	// Display information
	fmt.Println("PS2SMB Connection Information")
	fmt.Println("=============================")
	fmt.Println()
	fmt.Printf("Server Status: %s %s\n", statusSymbol, statusText)
	fmt.Printf("IP Address: %s\n", ip)
	fmt.Printf("Share Name: %s\n", cfg.ShareName)
	fmt.Printf("Games Path: %s\n", cfg.GamesPath)
	fmt.Println()

	// Authentication info
	fmt.Println("Authentication:")
	if cfg.UseGuest {
		fmt.Println("  Type: Guest (no password required)")
	} else {
		fmt.Printf("  Type: User authentication\n")
		fmt.Printf("  User: %s\n", cfg.SambaUser)
		fmt.Println("  Password: (set during init)")
	}
	fmt.Println()

	fmt.Printf("SMB Path: %s\n", smbPath)
	fmt.Println()

	// PS2 Configuration Instructions
	fmt.Println("Configure on your PS2 (OPL):")
	fmt.Println("=============================")
	fmt.Println("1. Go to Network Settings in OPL")
	fmt.Println("2. Set these values:")
	fmt.Printf("   - IP Address Type: Static or DHCP\n")
	fmt.Printf("   - SMB Server: %s\n", ip)
	fmt.Printf("   - SMB Share: %s\n", cfg.ShareName)
	
	if cfg.UseGuest {
		fmt.Println("   - SMB User: (leave empty)")
		fmt.Println("   - SMB Password: (leave empty)")
	} else {
		fmt.Printf("   - SMB User: %s\n", cfg.SambaUser)
		fmt.Println("   - SMB Password: (password you set)")
	}
	
	fmt.Println("3. Save and reconnect")
	fmt.Println()

	fmt.Println("Place your games in:")
	fmt.Printf("  DVD: %s/DVD/\n", cfg.GamesPath)
	fmt.Printf("  CD: %s/CD/\n", cfg.GamesPath)
	fmt.Println()

	if !sambaRunning {
		fmt.Println("WARNING: Samba service is not running!")
		fmt.Println("Start it with: sudo systemctl start smb")
		fmt.Println()
	}

	// Show all IPs if multiple interfaces
	allIPs, err := network.GetAllLocalIPs()
	if err == nil && len(allIPs) > 1 {
		fmt.Println("Available network interfaces:")
		for _, ip := range allIPs {
			fmt.Printf("  - %s\n", ip)
		}
		fmt.Println()
	}

	return nil
}

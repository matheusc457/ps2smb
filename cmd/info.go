package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/matheusc457/ps2smb/internal/config"
	"github.com/matheusc457/ps2smb/internal/network"
	"github.com/matheusc457/ps2smb/internal/samba"
	"github.com/spf13/cobra"
)

var (
	useNetBIOS       bool
	interfaceName    string
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
	infoCmd.Flags().BoolVarP(&useNetBIOS, "netbios", "n", false, "Use NetBIOS name instead of IP address")
	infoCmd.Flags().StringVarP(&interfaceName, "interface", "i", "", "Specify network interface (e.g., eth0, enp3s0)")
}

func getHostname() (string, error) {
	// Method 1: Try /etc/hostname (works on most Linux)
	data, err := os.ReadFile("/etc/hostname")
	if err == nil && len(data) > 0 {
		return strings.ToUpper(strings.TrimSpace(string(data))), nil
	}
	
	// Method 2: Try hostnamectl command (systemd-based)
	cmd := exec.Command("hostnamectl", "hostname")
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		return strings.ToUpper(strings.TrimSpace(string(output))), nil
	}
	
	// Method 3: Try hostname command (traditional)
	cmd = exec.Command("hostname")
	output, err = cmd.Output()
	if err == nil && len(output) > 0 {
		return strings.ToUpper(strings.TrimSpace(string(output))), nil
	}
	
	// Method 4: Fallback to os.Hostname() (Go builtin)
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return strings.ToUpper(name), nil
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
	var ip string
	
	if interfaceName != "" {
		// Use specified interface
		var ipErr error
		ip, ipErr = network.GetIPFromInterface(interfaceName)
		if ipErr != nil {
			return fmt.Errorf("failed to get IP from interface %s: %v", interfaceName, ipErr)
		}
	} else {
		// Auto-detect IP
		var ipErr error
		ip, ipErr = network.GetLocalIP()
		if ipErr != nil {
			return fmt.Errorf("failed to detect IP address: %v", ipErr)
		}
	}

	// Get hostname for NetBIOS
	hostname, hostnameErr := getHostname()

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
	if hostname != "" {
		fmt.Printf("NetBIOS Name: %s\n", hostname)
	}
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
	fmt.Println()
	fmt.Println("1. Go to 'Network Config' in OPL main menu")
	fmt.Println()
	fmt.Println("2. PS2 Network Settings:")
	fmt.Println("   - IP address type: DHCP (recommended)")
	fmt.Println()
	fmt.Println("   OR if using Static IP:")
	fmt.Printf("     - IP address: 192.168.1.10 (must be in same network as %s)\n", ip)
	fmt.Println("       * Choose an available IP in the same range")
	fmt.Printf("       * Example: If PC is %s, PS2 could be 192.168.1.10, 192.168.1.20, etc\n", ip)
	fmt.Println("     - Mask: 255.255.255.0")
	fmt.Printf("     - Gateway: %s (your PC's IP for direct connection, or router IP)\n", ip)
	fmt.Println()
	
	fmt.Println("3. SMB Server Settings:")
	
	if useNetBIOS && hostname != "" && hostnameErr == nil {
		fmt.Println("   - Address type: NetBIOS")
		fmt.Printf("   - Address: %s (hostname in UPPERCASE)\n", hostname)
	} else {
		if useNetBIOS && (hostname == "" || hostnameErr != nil) {
			fmt.Println("   WARNING: Could not get hostname, using IP instead")
		}
		fmt.Println("   - Address type: IP")
		fmt.Printf("   - Address: %s\n", ip)
	}
	
	fmt.Printf("   - Share: %s\n", cfg.ShareName)
	fmt.Println("   - Port: 445 (default, don't change)")
	
	if cfg.UseGuest {
		fmt.Println("   - User: (leave empty)")
		fmt.Println("   - Password: (leave empty)")
	} else {
		fmt.Printf("   - User: %s\n", cfg.SambaUser)
		fmt.Println("   - Password: (password you set during init)")
	}
	fmt.Println()
	
	fmt.Println("4. Advanced Settings (if needed):")
	fmt.Println("   - For direct crossover cable connection:")
	fmt.Println("     - Ethernet operation mode: 100Mbit half-duplex")
	fmt.Println()
	
	fmt.Println("5. Save settings and select 'Reconnect'")
	fmt.Println()

	fmt.Println("Place your game ISOs in:")
	fmt.Printf("  DVD games: %s/DVD/\n", cfg.GamesPath)
	fmt.Printf("  CD games:  %s/CD/\n", cfg.GamesPath)
	fmt.Println()

	if !sambaRunning {
		fmt.Println("WARNING: Samba service is not running!")
		fmt.Println("Start it with: sudo systemctl start smb")
		fmt.Println()
	}

	// Show all IPs if multiple interfaces
	allIPs, ipsErr := network.GetAllLocalIPs()
	if ipsErr == nil && len(allIPs) > 1 {
		fmt.Println("Available network interfaces:")
		for _, ip := range allIPs {
			fmt.Printf("  - %s\n", ip)
		}
		fmt.Println()
		fmt.Println("Tip: Use the interface connected to your PS2")
		fmt.Println()
	}

	// Show usage tip
	if !useNetBIOS {
		fmt.Println("Tip: You can use NetBIOS instead of IP with: ps2smb info --netbios")
	} else {
		fmt.Println("Tip: You can use IP address instead with: ps2smb info")
	}

	return nil
}

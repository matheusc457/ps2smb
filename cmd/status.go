package cmd

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/matheusc457/ps2smb/internal/config"
	"github.com/matheusc457/ps2smb/internal/samba"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check PS2 SMB server status",
	Long:  `Performs health checks on the Samba server and configuration to ensure everything is ready for PS2 connection.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runStatus(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus() error {
	fmt.Println("PS2SMB Status Check")
	fmt.Println("===================")
	fmt.Println()

	allOK := true

	// Check 1: Is ps2smb configured?
	fmt.Print("Configuration exists... ")
	if !config.Exists() {
		printStatus(false)
		fmt.Println("  Run 'sudo ps2smb init' to configure")
		return nil
	}
	printStatus(true)

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}

	// Check 2: Is Samba installed?
	fmt.Print("Samba installed... ")
	if !samba.IsSambaInstalled() {
		printStatus(false)
		allOK = false
		fmt.Println("  Install Samba to continue")
	} else {
		printStatus(true)
	}

	// Check 3: Is Samba service running?
	fmt.Print("Samba service running... ")
	if !samba.IsSambaRunning() {
		printStatus(false)
		allOK = false
		fmt.Println("  Start with: sudo systemctl start smb")
	} else {
		printStatus(true)
	}

	// Check 4: Does games directory exist?
	fmt.Printf("Games directory (%s)... ", cfg.GamesPath)
	if _, err := os.Stat(cfg.GamesPath); os.IsNotExist(err) {
		printStatus(false)
		allOK = false
		fmt.Printf("  Directory does not exist\n")
	} else {
		printStatus(true)
	}

	// Check 5: Do DVD and CD subdirectories exist?
	fmt.Print("DVD directory... ")
	dvdPath := cfg.GamesPath + "/DVD"
	if _, err := os.Stat(dvdPath); os.IsNotExist(err) {
		printStatus(false)
		allOK = false
	} else {
		printStatus(true)
	}

	fmt.Print("CD directory... ")
	cdPath := cfg.GamesPath + "/CD"
	if _, err := os.Stat(cdPath); os.IsNotExist(err) {
		printStatus(false)
		allOK = false
	} else {
		printStatus(true)
	}

	// Check 6: Is port 445 reachable?
	fmt.Print("Port 445 (SMB) reachable... ")
	portOpen := checkPort("localhost", 445)
	if !portOpen {
		printStatus(false)
		allOK = false
		fmt.Println("  Port may be blocked by firewall")
		fmt.Println("  Open with: sudo ufw allow 445")
	} else {
		printStatus(true)
	}

	// Summary
	fmt.Println()
	fmt.Println("Summary:")
	fmt.Println("--------")
	if allOK {
		fmt.Println("All checks passed! Your PS2 SMB server is ready.")
		fmt.Println("\nRun 'ps2smb info' to see connection details.")
	} else {
		fmt.Println("Some checks failed. Please fix the issues above.")
	}

	return nil
}

func printStatus(ok bool) {
	if ok {
		fmt.Println("✓")
	} else {
		fmt.Println("✗")
	}
}

func checkPort(host string, port int) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

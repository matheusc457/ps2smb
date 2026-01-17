package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/matheusc457/ps2smb/internal/config"
	"github.com/matheusc457/ps2smb/internal/samba"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize and configure Samba for PS2",
	Long:  `Sets up Samba server with optimized settings for PlayStation 2 network gaming via OPL.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runInit(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit() error {
	fmt.Println("PS2SMB Initialization")
	fmt.Println("=====================\n")

	// Check if already configured
	if config.Exists() {
		fmt.Println("Warning: ps2smb is already configured.")
		if !askYesNo("Do you want to reconfigure?") {
			fmt.Println("Initialization cancelled.")
			return nil
		}
	}

	// Check root privileges
	if !samba.IsRoot() {
		return fmt.Errorf("this command requires root privileges. Please run with sudo")
	}

	// Detect distro
	fmt.Println("Detecting Linux distribution...")
	distro, err := samba.DetectDistro()
	if err != nil {
		return fmt.Errorf("failed to detect distribution: %v", err)
	}
	fmt.Printf("Detected: %s\n\n", distro.Name)

	// Check if Samba is installed
	if !samba.IsSambaInstalled() {
		fmt.Println("Samba is not installed on your system.")
		if distro.InstallCmd != "" {
			fmt.Printf("You can install it with:\n  %s\n\n", distro.InstallCmd)
			if askYesNo("Would you like to install Samba now?") {
				fmt.Println("Please run the install command above and then run 'ps2smb init' again.")
				return nil
			}
		}
		return fmt.Errorf("samba is required but not installed")
	}
	fmt.Println("Samba is installed.\n")

	// Ask for games path
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the path where PS2 games will be stored [/home/ps2games]: ")
	gamesPath, _ := reader.ReadString('\n')
	gamesPath = strings.TrimSpace(gamesPath)
	if gamesPath == "" {
		gamesPath = "/home/ps2games"
	}

	// Ask about authentication
	fmt.Println("\nAuthentication options:")
	fmt.Println("1. Guest access (no password required)")
	fmt.Println("2. User authentication (more secure)")
	fmt.Print("Choose option [1]: ")
	authChoice, _ := reader.ReadString('\n')
	authChoice = strings.TrimSpace(authChoice)
	if authChoice == "" {
		authChoice = "1"
	}

	useGuest := authChoice == "1"
	sambaUser := ""

	// Backup existing config
	fmt.Println("\nBacking up existing Samba configuration...")
	if err := samba.BackupConfig(); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	// Add PS2 share
	fmt.Println("Adding PS2 share to Samba configuration...")
	if err := samba.AddPS2Share(gamesPath, useGuest); err != nil {
		return fmt.Errorf("failed to add PS2 share: %v", err)
	}

	// Create Samba user if needed
	if !useGuest {
		fmt.Println("\nCreating Samba user 'ps2user'...")
		fmt.Println("You will be prompted to set a password.")
		if err := samba.CreateSambaUser("ps2user", ""); err != nil {
			return fmt.Errorf("failed to create Samba user: %v", err)
		}
		sambaUser = "ps2user"
	}

	// Enable and restart Samba
	fmt.Println("\nEnabling Samba service...")
	if err := samba.EnableSamba(); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	fmt.Println("Restarting Samba service...")
	if err := samba.RestartSamba(); err != nil {
		return fmt.Errorf("failed to restart Samba: %v", err)
	}

	// Save configuration
	cfg := &config.Config{
		GamesPath:     gamesPath,
		ShareName:     samba.ShareName,
		UseGuest:      useGuest,
		SambaUser:     sambaUser,
		ConfigVersion: "1.0",
	}

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %v", err)
	}

	// Success message
	fmt.Println("\n========================================")
	fmt.Println("Configuration completed successfully!")
	fmt.Println("========================================")
	fmt.Printf("\nGames directory: %s\n", gamesPath)
	fmt.Printf("Share name: %s\n", samba.ShareName)
	if useGuest {
		fmt.Println("Authentication: Guest (no password)")
	} else {
		fmt.Printf("Authentication: User (%s)\n", sambaUser)
	}
	fmt.Println("\nRun 'ps2smb info' to see connection details for your PS2.")

	return nil
}

func askYesNo(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/N): ", question)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

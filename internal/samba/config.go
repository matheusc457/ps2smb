package samba

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	SmbConfPath = "/etc/samba/smb.conf"
	ShareName   = "PS2"
)

// BackupConfig creates a backup of smb.conf
func BackupConfig() error {
	if !IsRoot() {
		return fmt.Errorf("root privileges required")
	}

	backupPath := fmt.Sprintf("%s.backup.%d", SmbConfPath, time.Now().Unix())
	
	input, err := os.ReadFile(SmbConfPath)
	if err != nil {
		// If file doesn't exist, that's okay
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read smb.conf: %v", err)
	}

	err = os.WriteFile(backupPath, input, 0644)
	if err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	fmt.Printf("Backup created: %s\n", backupPath)
	return nil
}

// RemovePS2Share removes existing PS2 share from smb.conf
func RemovePS2Share() error {
	if !IsRoot() {
		return fmt.Errorf("root privileges required")
	}

	data, err := os.ReadFile(SmbConfPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist, nothing to remove
		}
		return fmt.Errorf("failed to read smb.conf: %v", err)
	}

	content := string(data)
	
	// Find [PS2] section
	startIdx := strings.Index(content, "[PS2]")
	if startIdx == -1 {
		return nil // PS2 section doesn't exist
	}

	// Find the next section or end of file
	endIdx := len(content)
	nextSection := strings.Index(content[startIdx+5:], "[")
	if nextSection != -1 {
		endIdx = startIdx + 5 + nextSection
	}

	// Remove the PS2 section
	newContent := content[:startIdx] + content[endIdx:]

	// Write back to file
	err = os.WriteFile(SmbConfPath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write smb.conf: %v", err)
	}

	fmt.Println("Removed existing PS2 configuration from smb.conf")
	return nil
}

// AddPS2Share adds PS2 share configuration to smb.conf
func AddPS2Share(gamesPath string, useGuest bool) error {
	if !IsRoot() {
		return fmt.Errorf("root privileges required")
	}

	// Create games directory structure if it doesn't exist
	if err := os.MkdirAll(gamesPath, 0755); err != nil {
		return fmt.Errorf("failed to create games directory: %v", err)
	}

	// Create DVD and CD subdirectories for OPL
	dvdPath := filepath.Join(gamesPath, "DVD")
	cdPath := filepath.Join(gamesPath, "CD")

	if err := os.MkdirAll(dvdPath, 0755); err != nil {
		return fmt.Errorf("failed to create DVD directory: %v", err)
	}

	if err := os.MkdirAll(cdPath, 0755); err != nil {
		return fmt.Errorf("failed to create CD directory: %v", err)
	}

	fmt.Println("Created DVD and CD directories for OPL")

	// Remove existing PS2 share if present
	if err := RemovePS2Share(); err != nil {
		return fmt.Errorf("failed to remove existing PS2 share: %v", err)
	}

	// Build share configuration
	shareConfig := fmt.Sprintf(`

[%s]
   comment = PlayStation 2 Games
   path = %s
   browseable = yes
   read only = yes
   create mask = 0644
   directory mask = 0755
`, ShareName, gamesPath)

	if useGuest {
		shareConfig += "   guest ok = yes\n"
		shareConfig += "   public = yes\n"
	} else {
		shareConfig += "   guest ok = no\n"
		shareConfig += "   valid users = ps2user\n"
	}

	// Append to smb.conf
	f, err := os.OpenFile(SmbConfPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open smb.conf: %v", err)
	}
	defer f.Close()

	if _, err = f.WriteString(shareConfig); err != nil {
		return fmt.Errorf("failed to write to smb.conf: %v", err)
	}

	fmt.Println("PS2 share configuration added to smb.conf")
	return nil
}

// EnableSMBv1 enables SMB v1 protocol (required for PS2)
func EnableSMBv1() error {
	if !IsRoot() {
		return fmt.Errorf("root privileges required")
	}

	// Simple append to global section
	globalConfig := "\n   min protocol = NT1\n"
	
	f, err := os.OpenFile(SmbConfPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open smb.conf: %v", err)
	}
	defer f.Close()

	if _, err = f.WriteString(globalConfig); err != nil {
		return fmt.Errorf("failed to enable SMB v1: %v", err)
	}

	return nil
}

// CreateSambaUser creates a Samba user
func CreateSambaUser(username, password string) error {
	if !IsRoot() {
		return fmt.Errorf("root privileges required")
	}

	// Create system user if doesn't exist
	cmd := exec.Command("useradd", "-M", "-s", "/usr/sbin/nologin", username)
	_ = cmd.Run() // Ignore error if user already exists

	// Set Samba password
	cmd = exec.Command("smbpasswd", "-a", username)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create Samba user: %v", err)
	}

	// Enable the user
	cmd = exec.Command("smbpasswd", "-e", username)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to enable Samba user: %v", err)
	}

	return nil
}

// RestartSamba restarts the Samba service
func RestartSamba() error {
	if !IsRoot() {
		return fmt.Errorf("root privileges required")
	}

	serviceName := GetSambaServiceName()
	cmd := exec.Command("systemctl", "restart", serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart Samba: %v", err)
	}

	fmt.Println("Samba service restarted successfully")
	return nil
}

// EnableSamba enables Samba to start on boot
func EnableSamba() error {
	if !IsRoot() {
		return fmt.Errorf("root privileges required")
	}

	serviceName := GetSambaServiceName()
	cmd := exec.Command("systemctl", "enable", serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to enable Samba: %v", err)
	}

	return nil
}

package samba

import (
	"os"
	"os/exec"
	"strings"
)

type Distro struct {
	Name          string
	PackageManager string
	InstallCmd     string
	ServiceManager string
}

// DetectDistro detects the Linux distribution
func DetectDistro() (*Distro, error) {
	// Try to read /etc/os-release
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return nil, err
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	
	distroID := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			distroID = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
			break
		}
	}

	distro := &Distro{
		ServiceManager: "systemctl", // Most modern distros use systemd
	}

	switch distroID {
	case "ubuntu", "debian", "linuxmint", "pop":
		distro.Name = "Debian-based"
		distro.PackageManager = "apt"
		distro.InstallCmd = "sudo apt update && sudo apt install -y samba"
	case "arch", "manjaro", "endeavouros":
		distro.Name = "Arch-based"
		distro.PackageManager = "pacman"
		distro.InstallCmd = "sudo pacman -S --noconfirm samba"
	case "fedora", "rhel", "centos":
		distro.Name = "Fedora-based"
		distro.PackageManager = "dnf"
		distro.InstallCmd = "sudo dnf install -y samba"
	default:
		distro.Name = "Unknown"
		distro.PackageManager = "unknown"
		distro.InstallCmd = ""
	}

	return distro, nil
}

// IsSambaInstalled checks if Samba is installed
func IsSambaInstalled() bool {
	_, err := exec.LookPath("smbd")
	return err == nil
}

// IsRoot checks if running as root
func IsRoot() bool {
	return os.Geteuid() == 0
}

// IsSambaRunning checks if Samba service is running
func IsSambaRunning() bool {
	serviceName := GetSambaServiceName()
	cmd := exec.Command("systemctl", "is-active", serviceName)
	err := cmd.Run()
	return err == nil
}

// GetSambaServiceName returns the correct service name for the distro
func GetSambaServiceName() string {
	distro, err := DetectDistro()
	if err != nil {
		return "smbd" // default
	}

	// Arch uses 'smb', most others use 'smbd'
	if distro.PackageManager == "pacman" {
		return "smb"
	}

	return "smbd"
}

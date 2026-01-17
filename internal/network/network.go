package network

import (
	"fmt"
	"net"
)

// GetLocalIP returns the primary local IP address
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no local IP address found")
}

// GetAllLocalIPs returns all local IP addresses with interface names
func GetAllLocalIPs() ([]string, error) {
	var ips []string
	
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		// Skip loopback and down interfaces
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil {
					ips = append(ips, fmt.Sprintf("%s (%s)", ipnet.IP.String(), iface.Name))
				}
			}
		}
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("no network interfaces found")
	}

	return ips, nil
}

// FormatSMBPath formats the SMB path for Windows/PS2 style
func FormatSMBPath(ip, shareName string) string {
	return fmt.Sprintf("\\\\%s\\%s", ip, shareName)
}

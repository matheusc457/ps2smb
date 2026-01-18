# ps2smb

A command-line tool to configure and manage SMB shares for PlayStation 2 network gaming on Linux.

## Overview

ps2smb automates the setup of a Samba server optimized for PlayStation 2 network gaming via OPL (Open PS2 Loader). It handles server configuration, network detection, and provides step-by-step instructions for connecting your PS2.

## Features

- Automatic Samba detection and configuration
- Multi-distribution support (Debian/Ubuntu, Arch, Fedora)
- SMB v1 protocol configuration for PS2 compatibility
- Guest or password-based authentication
- Automatic backup of existing Samba configuration
- Network interface selection for multi-adapter systems
- NetBIOS and IP address modes
- OPL-compliant directory structure creation

## Installation

### Prerequisites

- Linux distribution (tested on Arch, Ubuntu, Debian, Fedora)
- Go 1.21 or higher (for building from source)
- Root/sudo access for configuration

### From Source

```bash
# Clone the repository
git clone https://github.com/matheusc457/ps2smb
cd ps2smb

# Build the binary
go build -o ps2smb ./cmd/ps2smb

# Optional: Install system-wide
sudo cp ps2smb /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/matheusc457/ps2smb/cmd/ps2smb@latest
```

## Usage

### Initial Configuration

Configure the Samba server for PS2 access:

```bash
sudo ps2smb init
```

This command will:
- Detect your Linux distribution
- Verify Samba installation
- Create an optimized SMB share configuration
- Set up the games directory with DVD and CD subdirectories
- Configure authentication (guest or password-based)
- Enable and start the Samba service

### View Connection Information

Display network details and OPL configuration instructions:

```bash
sudo ps2smb info
```

Options:
- `--netbios, -n`: Use NetBIOS hostname instead of IP address
- `--interface, -i <name>`: Specify network interface to use

Examples:
```bash
sudo ps2smb info --netbios
sudo ps2smb info --interface enp3s0
```

### List Network Interfaces

View all available network interfaces:

```bash
ps2smb interfaces
```

Useful for identifying which network adapter to use with the `--interface` flag.

## Directory Structure

After initialization, ps2smb creates the following structure:

```
/your/games/path/
├── DVD/    # Place DVD game ISOs here
└── CD/     # Place CD game ISOs here
```

## Configuration Files

- User configuration: `~/.config/ps2smb/config.json`
- Samba configuration: `/etc/samba/smb.conf`
- Configuration backups: `/etc/samba/smb.conf.backup.<timestamp>`

## Network Setup

### Direct Connection (Crossover Cable)
Connect your PC directly to the PS2 using an Ethernet crossover cable. The PS2 and PC must be on the same network subnet.

### Network Connection (Router/Switch)
Connect both the PC and PS2 to the same network via a router or switch. This is the recommended setup as it allows DHCP configuration.

## Troubleshooting

### Samba Service Not Running
```bash
sudo systemctl start smb
sudo systemctl status smb
```

### Check Firewall Settings
Ensure port 445 is open:
```bash
sudo ufw allow 445
```

### Verify Share Access
```bash
smbclient -L localhost -N
```

## Contributing

Contributions are welcome. Please submit issues or pull requests via GitHub.

## License

GNU General Public License v3.0 - see [LICENSE](LICENSE) for details.

## Disclaimer

This tool is intended for use with legally obtained PS2 game backups. Users are responsible for ensuring they have the right to use any game files.

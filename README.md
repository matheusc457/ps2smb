# ps2smb

A CLI tool to easily configure SMB shares for playing PS2 games over network on Linux.

## What is this?

ps2smb automates the setup of a Samba server optimized for PlayStation 2 network gaming via OPL (Open PS2 Loader). No more manual configuration headaches!

## Features

- Automatic Samba detection and configuration
- Multi-distro support (Debian/Ubuntu, Arch, Fedora)
- SMB v1 protocol for PS2 compatibility
- Guest or password authentication options
- Automatic backup of existing Samba configuration

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/matheusc457/ps2smb
cd ps2smb

# Build
go build -o ps2smb ./cmd/ps2smb

# Optional: Install to system
sudo cp ps2smb /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/matheusc457/ps2smb/cmd/ps2smb@latest
```

## Usage

### Initial Setup

Configure Samba server for PS2 (requires sudo):

```bash
sudo ps2smb init
```

This will:
- Detect your Linux distribution
- Check if Samba is installed
- Configure SMB share optimized for PS2
- Create games directory
- Set up authentication (guest or password)
- Start and enable Samba service

### View Connection Info

```bash
ps2smb info
```

Shows the IP address and settings needed to configure OPL on your PS2.

## Requirements

- Linux (tested on Arch, Ubuntu, Debian, Fedora)
- Go 1.21 or higher (for building from source)
- Samba (can be installed during setup)
- Root/sudo access for initial configuration

## Roadmap

### v0.1.0 (Current)
- [x] Project structure
- [x] CLI framework with Cobra
- [x] Init command with Samba configuration
- [x] Multi-distro support
- [x] User configuration management

### v0.2.0 (Planned)
- [ ] Info command (network detection)
- [ ] Status command (health checks)
- [ ] List command (scan games directory)
- [ ] Add command (game management)

### v1.0.0 (Future)
- [ ] Automatic game renaming
- [ ] Firewall configuration assistance
- [ ] Connection testing
- [ ] Enhanced error handling

## Contributing

Contributions are welcome! Feel free to:
- Report bugs by opening an issue
- Suggest new features
- Submit pull requests

## License

GNU General Public License v3.0 - see [LICENSE](LICENSE) for details.

## Disclaimer

This tool is intended for use with legally obtained PS2 game backups only. Users are responsible for ensuring they have the right to use any game files.

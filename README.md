# ps2smb

A CLI tool to easily configure SMB shares for playing PS2 games over network on Linux.

## What is this?

ps2smb automates the setup of a Samba server optimized for PlayStation 2 network gaming via OPL (Open PS2 Loader). No more manual configuration headaches!

## Features

- Automatic Samba installation and configuration
- Optimized SMB v1 settings for PS2 compatibility
- Network detection and connection info
- Game library management
- Status checks and troubleshooting

## Quick Start

```bash
# Install
git clone https://github.com/matheusc457/ps2smb && cd ps2smb/

# Compile
go build -o ps2smb ./cmd/ps2smb/

# Execute
./ps2smb

```

## Requirements

- Linux (tested on Ubuntu, Debian, Arch, Fedora)
- Samba (will be installed automatically if not present)
- Root/sudo access for initial setup

## Commands

```bash
ps2smb init       # Configure Samba server for PS2
ps2smb info       # Show connection information
ps2smb list       # List available games
ps2smb add        # Add game to library
ps2smb status     # Check server status
```

## Roadmap

- [x] Project setup
- [ ] Basic CLI structure
- [ ] Samba detection and configuration
- [ ] Network detection
- [ ] Game management
- [ ] Status checks
- [ ] Multi-distro support

## Contributing

Contributions are welcome! Feel free to open issues or submit PRs.

## License

GNU General Public License v3.0 - see LICENSE for details.

## Disclaimer

This tool is for use with legally obtained PS2 game backups only.

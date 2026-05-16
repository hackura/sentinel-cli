# Hackura Sentinel CLI

Production-grade CLI tool for Hackura Sentinel AI.

## Features
- Secure device-based authentication via web portal.
- Real-time and background security scans.
- Batch scanning from files.
- Formatted output with risk levels and threat signals.

## Installation
```bash
make build
sudo mv sentinel /usr/local/bin/
```

## Usage
### Login
```bash
sentinel login
```

### Scan
```bash
sentinel scan https://example.com
sentinel scan --file targets.txt
sentinel scan --async https://example.com
```

### View Results
```bash
sentinel results list
sentinel results show <id>
```

## Security
- Tokens are stored in `~/.sentinel/token` with 0600 permissions.
- Device IDs are masked in output.
- All communications are over HTTPS.

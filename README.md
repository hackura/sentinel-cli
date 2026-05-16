<p align="center">
  <img src="assets/logo.png" width="128" alt="Hackura Sentinel Logo">
</p>

# Hackura Sentinel CLI

Production-grade CLI tool for Hackura Sentinel AI.

## Features
- Secure device-based authentication via web portal.
- Real-time and background security scans.
- Batch scanning from files.
- Formatted output with risk levels and threat signals.

## Installation
### Quick Install (Linux & macOS)
```bash
curl -fsSL https://sentinel.hackura.app/install.sh | bash
```

### Manual Build
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

## Legal
### License
Distributed under the MIT License. See `LICENSE` for more information.

### Privacy Policy
Your privacy is critical to us. View our full policy at [sentinel.hackura.app/privacy](https://sentinel.hackura.app/privacy).

### Terms of Service
By using this CLI, you agree to our terms at [sentinel.hackura.app/terms](https://sentinel.hackura.app/terms).

---
© 2026 Hackura AI Systems. All rights reserved.

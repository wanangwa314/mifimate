# MiFiMate

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)](https://golang.org/dl/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)]()

A free and open-source cross-platform desktop application for managing ZTE MiFi devices. Built with Go and Fyne, this application provides a modern graphical interface to monitor and configure your mobile WiFi hotspot.

**Tested on:** ZTE MF927U

## Features

### ✅ Working Features

- **Device Monitoring**
  - Real-time network status (3G/4G/5G)
  - Signal strength indicator
  - Battery level monitoring
  - Network statistics (upload/download speeds)
  - Connected devices list
  - WAN IP address display

- **WiFi Management**
  - View current WiFi configuration (SSID, security mode)
  - Modify WiFi settings (SSID, password, security mode)
  - Hide SSID option
  - Channel selection
  - Maximum clients configuration

- **SMS Management**
  - Read SMS messages
  - Delete SMS messages
  - View message timestamps
  - Recent messages widget on dashboard

- **Device Control**
  - Remote device reboot
  - Remote device shutdown
  - Network connect/disconnect

- **Application Features**
  - Auto-login on startup
  - Auto-refresh status polling (every 3 seconds)
  - Modern card-based UI design
  - Configuration file management
  - Adjustable log levels

## Installation

### From Source

#### Prerequisites

- Go 1.21 or higher
- GCC compiler (for CGO)
- Platform-specific dependencies:
  - **Linux**: `libgl1-mesa-dev xorg-dev`
  - **macOS**: Xcode command line tools
  - **Windows**: MinGW-w64

#### Build Instructions

```bash
# Clone the repository
git clone https://github.com/yourusername/mifimate.git
cd mifimate

# Install dependencies
go mod download

# Build the application
go build -o mifimate .

# Run the application
./mifimate
```

### Cross-Platform Builds

```bash
# Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -o mifimate.exe

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o mifimate-mac

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o mifimate-mac-arm64

# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o mifimate-linux
```

## Quick Start

1. **Connect to your MiFi device's WiFi network**

2. **Launch the application**
   ```bash
   ./mifimate
   ```

3. **Configure on first run**

   The application will create a configuration file at:
   - Linux/macOS: `~/.config/mifimate/config.yaml`
   - Windows: `%APPDATA%\mifimate\config.yaml`

   Edit the file to set your device password:
   ```yaml
   device:
     default_ip: "192.168.1.1"
     username: "admin"
     password: "your_password_here"  # Set your device password
   ```

4. **Restart and connect**

   The application will automatically attempt to connect on startup.

## Configuration

### Device Settings

```yaml
device:
  default_ip: "192.168.1.1"        # MiFi device IP address
  username: "admin"                 # Device admin username
  password: ""                      # Device admin password
  connection_timeout: 10            # Connection timeout in seconds
  poll_interval: 3                  # Status refresh interval in seconds
  auto_reconnect: true              # Auto-reconnect on connection loss
```

### Application Settings

```yaml
app:
  theme: "system"                   # UI theme: light, dark, system
  language: "en"                    # Interface language
  auto_start: false                 # Start on system boot
  minimize_to_tray: true            # Minimize to system tray
  show_notifications: true          # Show desktop notifications
  log_level: "info"                 # Logging level: debug, info, warn, error
```

## Supported Devices

### Confirmed Working
- **ZTE MF927U** ✅ (Fully tested)

### Potentially Compatible
- Other ZTE MiFi devices using the same web interface
- Devices accessible via `192.168.1.1` with similar API endpoints

**Note:** This application is a port of the ZTE MiFi web interface to a native desktop application. It should work with most modern ZTE MiFi devices, but has only been thoroughly tested on the MF927U model.

## Troubleshooting

### Connection Issues

**Problem:** Cannot connect to device

**Solutions:**
- Verify you're connected to the device's WiFi network
- Check device IP address (default: 192.168.1.1)
- Ensure device is powered on
- Verify credentials in config file
- Check firewall settings

### Build Issues

**Problem:** Build fails with Fyne errors

**Solutions:**
```bash
# Install platform dependencies (Linux)
sudo apt-get install libgl1-mesa-dev xorg-dev

# Clean and rebuild
go clean -cache
go mod tidy
go build -o mifimate .
```

## Acknowledgments

- **ZTE** - For the original MiFi web interface that inspired this project

## Disclaimer

This project is not affiliated with, endorsed by, or connected to ZTE Corporation. All product names, logos, and brands are property of their respective owners.

Use this software at your own risk. The authors are not responsible for any damage to your device or data loss.

## Building the Application

### Prerequisites
- Go 1.21 or higher
- Fyne dependencies (see [Fyne documentation](https://docs.fyne.io/started/))

### Build Commands

```bash
# Development build and run
go run main.go

# Production build
go build -o mifimate .

# Run the compiled binary
./mifimate
```

### Cross-platform Compilation

```bash
# For Windows
GOOS=windows GOARCH=amd64 go build -o mifimate.exe

# For macOS
GOOS=darwin GOARCH=amd64 go build -o mifimate-mac

# For Linux
GOOS=linux GOARCH=amd64 go build -o mifimate-linux
```

## Configuration

On first run, the application creates a default configuration file at:
- Linux/macOS: `~/.config/mifimate/config.yaml`
- Windows: `%APPDATA%\mifimate\config.yaml`

### Default Configuration
```yaml
device:
  default_ip: "192.168.1.1"
  username: "admin"
  password: ""
  connection_timeout: 10
  poll_interval: 3
  auto_reconnect: true

app:
  theme: "system"
  language: "en"
  auto_start: false
  minimize_to_tray: true
  show_notifications: true
  log_level: "info"
```

## Usage

1. **Launch the application**
   ```bash
   ./mifimate
   ```

2. **Configure password**
   - Edit the config file and set your device password
   - Or set it through the UI (password dialog to be implemented in Phase 2)

3. **Connect to device**
   - Click the "Connect" button
   - Application will attempt to reach the device at the configured IP
   - On successful login, device status will be displayed

4. **View device information**
   - Network type (3G/4G/5G)
   - Signal strength with quality indicator
   - Battery level with status
   - WAN IP address
   - Number of connected devices
   - Real-time upload/download speeds

5. **Refresh status**
   - Click "Refresh" to update device information manually

6. **Disconnect**
   - Click "Disconnect" to log out from the device

## API Endpoints Used

The application interacts with the MiFi device using these endpoints:

- `POST /goform/goform_set_cmd_process` - For login, logout, and configuration changes
- `GET /goform/goform_get_cmd_process` - For retrieving device status and information

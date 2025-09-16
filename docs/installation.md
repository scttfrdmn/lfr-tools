# Installation Guide

This guide covers various methods to install lfr on your system.

## Prerequisites

### AWS Configuration

Before using lfr, ensure you have AWS credentials configured:

```bash
# Using AWS CLI
aws configure

# Or set environment variables
export AWS_ACCESS_KEY_ID=your-access-key
export AWS_SECRET_ACCESS_KEY=your-secret-key
export AWS_DEFAULT_REGION=us-east-1
```

Required AWS permissions:
- IAM: Create/delete users, groups, policies
- Lightsail for Research: Create/manage instances, key pairs, volumes
- EC2: Describe regions and VPCs (for Lightsail for Research integration)

## Installation Methods

### 1. Homebrew (Recommended for macOS/Linux)

```bash
# Add the tap
brew tap scttfrdmn/lfr-tools

# Install lfr
brew install lfr

# Verify installation
lfr version
```

### 2. GitHub Releases

Download the latest binary for your platform from the [releases page](https://github.com/scttfrdmn/lfr-tools/releases).

#### Linux/macOS

```bash
# Download and install (replace with latest version and your platform)
curl -L https://github.com/scttfrdmn/lfr-tools/releases/download/v1.0.0/lfr_Linux_x86_64.tar.gz | tar xz
sudo mv lfr /usr/local/bin/
chmod +x /usr/local/bin/lfr
```

#### Windows

Download the Windows executable and add it to your PATH.

### 3. Docker

```bash
# Pull the image
docker pull ghcr.io/scttfrdmn/lfr-tools:latest

# Run with AWS credentials
docker run --rm -it \
  -e AWS_ACCESS_KEY_ID \
  -e AWS_SECRET_ACCESS_KEY \
  -e AWS_DEFAULT_REGION \
  ghcr.io/scttfrdmn/lfr-tools:latest --help

# Create an alias for convenience
alias lfr='docker run --rm -it -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -e AWS_DEFAULT_REGION ghcr.io/scttfrdmn/lfr-tools:latest'
```

### 4. Go Install (Build from Source)

```bash
# Install directly from source
go install github.com/scttfrdmn/lfr-tools@latest

# Verify installation
lfr-tools version
```

### 5. Manual Build

```bash
# Clone the repository
git clone https://github.com/scttfrdmn/lfr-tools.git
cd lfr-tools

# Build the binary
make build

# Install to GOPATH/bin
make install
```

## Verification

After installation, verify lfr is working:

```bash
# Check version
lfr version

# View help
lfr --help

# Test AWS connectivity
lfr instances list
```

## Configuration

Create a configuration file for default settings:

```bash
# Create config directory
mkdir -p ~/.config

# Create configuration file
cat > ~/.lfr-tools.yaml << EOF
aws:
  profile: "default"
  region: "us-east-1"

defaults:
  blueprint: "ubuntu_22_04"
  bundle: "nano_2_0"
  idle_threshold: 120

ssh:
  key_path: "~/.ssh/lfr-tools"
  config_path: "~/.ssh/config.d/lfr-tools"
EOF
```

## Shell Completion

Enable shell completion for better CLI experience:

### Bash

```bash
# Add to ~/.bashrc
echo 'source <(lfr-tools completion bash)' >> ~/.bashrc
source ~/.bashrc
```

### Zsh

```bash
# Add to ~/.zshrc
echo 'source <(lfr-tools completion zsh)' >> ~/.zshrc
source ~/.zshrc
```

### Fish

```bash
# Add to fish config
lfr-tools completion fish | source
```

## Updating

### Homebrew

```bash
brew update
brew upgrade lfr-tools
```

### Docker

```bash
docker pull ghcr.io/scttfrdmn/lfr-tools:latest
```

### Go Install

```bash
go install github.com/scttfrdmn/lfr-tools@latest
```

## Uninstalling

### Homebrew

```bash
brew uninstall lfr-tools
brew untap scttfrdmn/lfr-tools
```

### Manual

```bash
# Remove binary
sudo rm /usr/local/bin/lfr-tools

# Remove configuration (optional)
rm -rf ~/.lfr-tools.yaml
```

## Troubleshooting

### Permission Issues

If you encounter permission errors:

```bash
# Linux/macOS
sudo chown $USER:$USER /usr/local/bin/lfr-tools
chmod +x /usr/local/bin/lfr-tools
```

### AWS Credential Issues

Verify AWS credentials are properly configured:

```bash
aws sts get-caller-identity
```

### Network Issues

If experiencing connectivity issues:

```bash
# Test basic connectivity
curl -I https://lightsail.us-east-1.amazonaws.com

# Check if behind corporate firewall
lfr --debug instances list
```

## Next Steps

- [Quick Start Guide](../README.md#quick-start)
- [Configuration Reference](configuration.md)
- [Command Reference](commands.md)
- [Troubleshooting Guide](troubleshooting.md)
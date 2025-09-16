# lfr-tools

[![CI](https://github.com/scttfrdmn/lfr-tools/actions/workflows/ci.yml/badge.svg)](https://github.com/scttfrdmn/lfr-tools/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scttfrdmn/lfr-tools)](https://goreportcard.com/report/github.com/scttfrdmn/lfr-tools)
[![Release](https://img.shields.io/github/release/scttfrdmn/lfr-tools.svg)](https://github.com/scttfrdmn/lfr-tools/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive CLI tool for managing AWS Lightsail for Research instances, IAM users, and groups. Simplifies the process of creating and managing multi-user research environments with proper access controls, including streamlined SSH access.

## Features

### 🎯 Current Features

- **User Management**: Create, remove, and list IAM users with auto-generated passwords
- **Group Management**: Organize users with IAM groups and custom policies
- **Instance Management**: List, start, stop, and monitor Lightsail instances
- **SSH Simplification**: Easy SSH access with automatic key management and tunneling
- **Project Organization**: Tag and organize resources by project
- **Idle Detection**: Built-in idle detection to save costs

### 🚧 Roadmap

#### Phase 1 - Enhanced Core Features
- **EFS Integration**: Mount shared storage via VPC peering (as outlined in Lightsail documentation)
- **Advanced Idle Detection**: More configurable idle thresholds and detection algorithms
- **NICE DCV Integration**: Better DCV settings and direct connection commands for remote desktop access
- **Instance Lifecycle**: Full snapshot, backup, restore, and cloning capabilities

#### Phase 2 - Advanced Management
- **Cost Optimization**: Automated recommendations, scheduling, and budget alerts
- **Multi-region Support**: Deploy and manage across multiple AWS regions
- **Bulk Operations**: Mass user/instance operations with progress tracking and rollback
- **Usage Analytics**: Detailed reporting on instance utilization and costs

#### Phase 3 - Enterprise Features
- **Auth Integration**: External authentication systems (LDAP, SAML, OAuth)
- **RBAC**: Role-based access control with custom permission sets
- **Audit Logging**: Comprehensive activity logging and compliance reporting
- **API & SDK**: RESTful API and language-specific SDKs for integration

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap scttfrdmn/lfr-tools
brew install lfr
```

### GitHub Releases

Download the latest binary from the [releases page](https://github.com/scttfrdmn/lfr-tools/releases).

### Docker

```bash
docker run --rm -it ghcr.io/scttfrdmn/lfr-tools:latest --help
```

### From Source

```bash
go install github.com/scttfrdmn/lfr-tools@latest
```

## Quick Start

1. **Configure AWS credentials** (AWS CLI or environment variables)

2. **Create users and instances:**
   ```bash
   lfr users create \
     --project "my-research" \
     --blueprint "ubuntu_22_04" \
     --bundle "nano_2_0" \
     --region "us-east-1" \
     --users "alice,bob,charlie"
   ```

3. **Connect via SSH:**
   ```bash
   lfr ssh connect alice --project "my-research"
   ```

4. **Monitor instances:**
   ```bash
   lfr instances monitor --project "my-research"
   ```

## Usage

### User Management

```bash
# Create users with instances
lfr users create -p myproject -b ubuntu_22_04 --bundle nano_2_0 -r us-east-1 -u alice,bob

# List users
lfr users list -p myproject

# Remove users
lfr users remove -p myproject -u alice,bob
lfr users remove -p myproject --all
```

### Group Management

```bash
# Create group with policies
lfr groups create -n researchers -p arn:aws:iam::aws:policy/ReadOnlyAccess

# List groups
lfr groups list

# Remove group
lfr groups remove -n researchers
```

### Instance Management

```bash
# List instances
lfr instances list -p myproject

# Start/stop instances
lfr instances start -u alice,bob
lfr instances stop -u alice,bob

# Monitor usage
lfr instances monitor -p myproject --idle-threshold 60

# Snapshot and restore
lfr instances snapshot alice-ubuntu_22_04
lfr instances restore snapshot-name new-instance-name

# Clone instances
lfr instances clone alice-ubuntu_22_04 bob-ubuntu_22_04

# Reboot instances
lfr instances reboot alice-ubuntu_22_04 bob-ubuntu_22_04
```

### SSH Management

```bash
# Connect to instance
lfr ssh connect alice -p myproject

# Download SSH keys
lfr ssh keys download alice -o ~/.ssh/

# Generate SSH config
lfr ssh config -p myproject -o ~/.ssh/config.d/lfr-tools

# Create SSH tunnel
lfr ssh tunnel alice 8888:8888 -p myproject
```

## Configuration

Create `~/.lfr-tools.yaml`:

```yaml
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
```

### NICE DCV Management

```bash
# Connect via NICE DCV
lfr dcv connect alice -p myproject

# Configure DCV settings
lfr dcv config -p myproject --quality high

# Check DCV status
lfr dcv status -p myproject

# List active sessions
lfr dcv sessions list
```

## Development

### Prerequisites

- Go 1.20+
- Make
- golangci-lint
- pre-commit

### Setup

```bash
git clone https://github.com/scttfrdmn/lfr-tools.git
cd lfr-tools
make deps
pre-commit install
```

### Build and Test

```bash
make build
make test
make lint
make check  # Run all checks
```

### Release

Tags are automatically released via GoReleaser when pushed:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

- 📖 [Documentation](https://github.com/scttfrdmn/lfr-tools/wiki)
- 🐛 [Issues](https://github.com/scttfrdmn/lfr-tools/issues)
- 💬 [Discussions](https://github.com/scttfrdmn/lfr-tools/discussions)
- ☕ [Ko-Fi](https://ko-fi.com/scttfrdmn)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Credits

Built with ❤️ by [Scott Friedman](https://github.com/scttfrdmn)

Based on the original [lightsail-multiuser](https://github.com/scttfrdmn/lightsail-multiuser) bash script.
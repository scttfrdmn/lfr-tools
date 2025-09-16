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

- **EFS Integration**: Mount shared storage via VPC peering
- **Advanced Monitoring**: Enhanced idle detection and usage analytics
- **Cost Optimization**: Automated recommendations and scheduling
- **Multi-region Support**: Deploy across multiple AWS regions
- **Auth Integration**: External authentication system support

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap scttfrdmn/lfr-tools
brew install lfr-tools
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
   lfr-tools users create \
     --project "my-research" \
     --blueprint "ubuntu_22_04" \
     --bundle "nano_2_0" \
     --region "us-east-1" \
     --users "alice,bob,charlie"
   ```

3. **Connect via SSH:**
   ```bash
   lfr-tools ssh connect alice --project "my-research"
   ```

4. **Monitor instances:**
   ```bash
   lfr-tools instances monitor --project "my-research"
   ```

## Usage

### User Management

```bash
# Create users with instances
lfr-tools users create -p myproject -b ubuntu_22_04 --bundle nano_2_0 -r us-east-1 -u alice,bob

# List users
lfr-tools users list -p myproject

# Remove users
lfr-tools users remove -p myproject -u alice,bob
lfr-tools users remove -p myproject --all
```

### Group Management

```bash
# Create group with policies
lfr-tools groups create -n researchers -p arn:aws:iam::aws:policy/ReadOnlyAccess

# List groups
lfr-tools groups list

# Remove group
lfr-tools groups remove -n researchers
```

### Instance Management

```bash
# List instances
lfr-tools instances list -p myproject

# Start/stop instances
lfr-tools instances start -u alice,bob
lfr-tools instances stop -u alice,bob

# Monitor usage
lfr-tools instances monitor -p myproject --idle-threshold 60
```

### SSH Management

```bash
# Connect to instance
lfr-tools ssh connect alice -p myproject

# Download SSH keys
lfr-tools ssh keys download alice -o ~/.ssh/

# Generate SSH config
lfr-tools ssh config -p myproject -o ~/.ssh/config.d/lfr-tools

# Create SSH tunnel
lfr-tools ssh tunnel alice 8888:8888 -p myproject
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
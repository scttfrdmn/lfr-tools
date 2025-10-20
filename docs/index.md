---
layout: default
title: Home
---

# LFR Tools

A comprehensive CLI tool and GUI application for managing AWS Lightsail for Research instances, IAM users, and groups. Designed specifically for educational and research environments.

[![CI](https://github.com/scttfrdmn/lfr-tools/actions/workflows/ci.yml/badge.svg)](https://github.com/scttfrdmn/lfr-tools/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scttfrdmn/lfr-tools)](https://goreportcard.com/report/github.com/scttfrdmn/lfr-tools)
[![Release](https://img.shields.io/github/release/scttfrdmn/lfr-tools.svg)](https://github.com/scttfrdmn/lfr-tools/releases)

## Quick Links

- 📚 [Getting Started](getting-started.md)
- 📖 [Documentation](documentation.md)
- 🎓 [Tutorials](tutorials.md)
- 💬 [Discussions](https://github.com/scttfrdmn/lfr-tools/discussions)
- 🐛 [Report a Bug](https://github.com/scttfrdmn/lfr-tools/issues/new?template=bug_report.yml)
- ✨ [Request a Feature](https://github.com/scttfrdmn/lfr-tools/issues/new?template=feature_request.yml)

## Key Features

### 🎯 For Professors
- **Bulk User Management**: Create and manage student accounts in minutes
- **Project Organization**: Tag and organize resources by course or project
- **Cost Control**: Built-in idle detection to minimize AWS costs
- **SSH Simplification**: Automatic key management and easy connections

### 👨‍🏫 For Teaching Assistants
- **Instance Monitoring**: Track student instance usage and status
- **Group Management**: Organize students with IAM groups and policies
- **Quick Actions**: Start, stop, and manage instances easily

### 🎓 For Students
- **Simple Access**: Easy SSH connections to research instances
- **DCV Support**: Remote desktop access via NICE DCV
- **User-Friendly**: Clear error messages and helpful documentation

### 🔬 For Researchers
- **Flexible Configuration**: Customize instances for your workflow
- **EFS Integration**: Shared storage via VPC peering
- **Software Packs**: Pre-configured software environments

## Installation

### Homebrew (macOS/Linux)
```bash
brew tap scttfrdmn/lfr-tools
brew install lfr
```

### GitHub Releases
Download the latest binary from [releases](https://github.com/scttfrdmn/lfr-tools/releases).

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

## Documentation

- **Quick Start**
  - 📚 [Getting Started Guide](getting-started.md)
  - 🎓 [Student Guide](tutorials/student-guide.md)
  - 👨‍🏫 [Teacher Guide](tutorials/teacher-guide.md)

- **Detailed Guides**
  - 🔧 [Installation](installation.md)
  - ⚙️ [Configuration](configuration.md)
  - 🧪 [Testing](testing.md)

- **Tutorials**
  - 📖 [Common Tasks](tutorials/common-tasks.md)
  - 🎯 [Real Examples](tutorials/examples.md)
  - 🏫 [Educational Workflows](educational-workflows.md)
  - 🆘 [Troubleshooting](troubleshooting.md)

## Community

- **Get Help**: [GitHub Discussions](https://github.com/scttfrdmn/lfr-tools/discussions)
- **Report Issues**: [GitHub Issues](https://github.com/scttfrdmn/lfr-tools/issues)
- **Contribute**: [Contributing Guide](https://github.com/scttfrdmn/lfr-tools/blob/main/CONTRIBUTING.md)

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/scttfrdmn/lfr-tools/blob/main/LICENSE) file for details.

## Credits

Built with ❤️ by [Scott Friedman](https://github.com/scttfrdmn)

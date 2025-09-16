# Configuration Reference

lfr-tools supports configuration through YAML files, environment variables, and command-line flags.

## Configuration File

The configuration file is loaded from the following locations (in order of precedence):

1. `--config` flag specified file
2. `$PWD/.lfr-tools.yaml`
3. `$HOME/.lfr-tools.yaml`

### Example Configuration

```yaml
# ~/.lfr-tools.yaml
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

debug: false
```

## Configuration Sections

### AWS Configuration

Controls AWS SDK behavior and default settings.

```yaml
aws:
  profile: "production"     # AWS profile to use
  region: "us-west-2"       # Default AWS region
```

**Environment Variables:**
- `LFR_AWS_PROFILE` / `AWS_PROFILE`
- `LFR_AWS_REGION` / `AWS_DEFAULT_REGION`

### Defaults

Default values for command parameters to reduce repetitive typing.

```yaml
defaults:
  blueprint: "ubuntu_22_04"      # Default Lightsail blueprint
  bundle: "nano_2_0"             # Default Lightsail bundle
  idle_threshold: 120            # Default idle threshold (minutes)
```

**Environment Variables:**
- `LFR_DEFAULTS_BLUEPRINT`
- `LFR_DEFAULTS_BUNDLE`
- `LFR_DEFAULTS_IDLE_THRESHOLD`

### SSH Configuration

SSH-related settings for key management and connections.

```yaml
ssh:
  key_path: "~/.ssh/lfr-tools"           # Directory for SSH keys
  config_path: "~/.ssh/config.d/lfr-tools" # SSH config file path
```

**Environment Variables:**
- `LFR_SSH_KEY_PATH`
- `LFR_SSH_CONFIG_PATH`

### Debug

Enable debug logging and verbose output.

```yaml
debug: true
```

**Environment Variables:**
- `LFR_DEBUG`

## Command-Line Flags

Global flags that override configuration file settings:

| Flag | Description | Environment Variable |
|------|-------------|---------------------|
| `--config` | Configuration file path | `LFR_CONFIG` |
| `--debug` | Enable debug logging | `LFR_DEBUG` |
| `--profile` | AWS profile to use | `LFR_AWS_PROFILE` |
| `--region` | AWS region to use | `LFR_AWS_REGION` |

## Available Blueprints

Common Lightsail for Research blueprints:

- `ubuntu_22_04` - Ubuntu 22.04 LTS
- `ubuntu_20_04` - Ubuntu 20.04 LTS
- `amazon_linux_2` - Amazon Linux 2
- `centos_7` - CentOS 7
- `debian_11` - Debian 11

## Available Bundles

Common Lightsail for Research bundles:

- `nano_2_0` - 512 MB RAM, 1 vCPU, 20 GB SSD
- `micro_2_0` - 1 GB RAM, 1 vCPU, 40 GB SSD
- `small_2_0` - 2 GB RAM, 1 vCPU, 60 GB SSD
- `medium_2_0` - 4 GB RAM, 2 vCPUs, 80 GB SSD
- `large_2_0` - 8 GB RAM, 2 vCPUs, 160 GB SSD
- `xlarge_2_0` - 16 GB RAM, 4 vCPUs, 320 GB SSD
- `2xlarge_2_0` - 32 GB RAM, 8 vCPUs, 640 GB SSD

## Available Regions

Common AWS regions for Lightsail:

- `us-east-1` - US East (N. Virginia)
- `us-east-2` - US East (Ohio)
- `us-west-1` - US West (N. California)
- `us-west-2` - US West (Oregon)
- `eu-west-1` - Europe (Ireland)
- `eu-central-1` - Europe (Frankfurt)
- `ap-southeast-1` - Asia Pacific (Singapore)
- `ap-northeast-1` - Asia Pacific (Tokyo)

## Configuration Validation

lfr-tools validates configuration on startup:

```bash
# Test configuration
lfr-tools --config ~/.lfr-tools.yaml instances list --dry-run
```

## Environment Variable Precedence

Configuration is loaded in the following order (highest to lowest precedence):

1. Command-line flags
2. Environment variables
3. Configuration file
4. Default values

## Sensitive Data

**Never store sensitive data in configuration files:**

❌ **Don't do this:**
```yaml
aws:
  access_key: "AKIA..."
  secret_key: "abc123..."
```

✅ **Do this instead:**
```bash
# Use AWS CLI configuration
aws configure

# Or environment variables
export AWS_ACCESS_KEY_ID=AKIA...
export AWS_SECRET_ACCESS_KEY=abc123...
```

## Multiple Environments

Manage different environments using profiles:

```bash
# Development environment
lfr-tools --profile dev --region us-east-1 users list

# Production environment
lfr-tools --profile prod --region us-west-2 users list
```

Or separate configuration files:

```bash
# Development config
lfr-tools --config ~/.lfr-tools-dev.yaml users list

# Production config
lfr-tools --config ~/.lfr-tools-prod.yaml users list
```

## Configuration Examples

### Minimal Configuration

```yaml
aws:
  region: "us-east-1"
```

### Development Environment

```yaml
aws:
  profile: "dev"
  region: "us-east-1"

defaults:
  blueprint: "ubuntu_22_04"
  bundle: "nano_2_0"

debug: true
```

### Production Environment

```yaml
aws:
  profile: "production"
  region: "us-west-2"

defaults:
  blueprint: "ubuntu_22_04"
  bundle: "medium_2_0"
  idle_threshold: 240

ssh:
  key_path: "/etc/lfr-tools/keys"
  config_path: "/etc/ssh/ssh_config.d/lfr-tools"
```

## Troubleshooting Configuration

### View Current Configuration

```bash
lfr-tools --debug users list 2>&1 | head -20
```

### Validate Configuration File

```bash
# Check YAML syntax
yamllint ~/.lfr-tools.yaml

# Test with lfr-tools
lfr-tools --config ~/.lfr-tools.yaml version
```

### Common Issues

1. **YAML Syntax Errors**
   - Ensure proper indentation (2 spaces)
   - Quote string values with special characters
   - No tabs, only spaces

2. **Path Resolution**
   - Use absolute paths or `~` for home directory
   - Ensure directories exist and are writable

3. **AWS Profile Issues**
   - Verify profile exists: `aws configure list-profiles`
   - Check credentials: `aws --profile myprofile sts get-caller-identity`

## Next Steps

- [Installation Guide](installation.md)
- [Command Reference](commands.md)
- [Quick Start](../README.md#quick-start)
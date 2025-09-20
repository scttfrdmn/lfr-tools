# LFR Tools Development Notes

## Project Overview
`lfr-tools` is a comprehensive CLI tool for managing AWS Lightsail for Research instances, IAM users, and groups. The project includes both CLI tools and a GUI application for educational and research environments.

## Architecture

### CLI Tool (Go)
- **Main CLI**: Located in `cmd/` directory
- **Core Libraries**: `internal/` directory with `lfrtypes` and `lfrutils` packages
- **AWS Integration**: `internal/aws/` with IAM and Lightsail services
- **Testing**: Comprehensive unit tests with LocalStack integration

### GUI Application (Wails3 + React)
- **Backend**: Go services in `gui/lfr-gui/pkg/services/`
- **Frontend**: React + TypeScript in `gui/lfr-gui/frontend/`
- **Framework**: **Wails3** (alpha version)
- **UI Components**: Cloudscape Design System
- **Testing**: Playwright e2e tests

## Key Development Information

### GUI Framework: Wails3
This project uses **Wails3** (alpha) for the GUI application:
- **Location**: `gui/lfr-gui/`
- **Installation**: `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`
- **Path**: Add `~/go/bin` to PATH or use `~/go/bin/wails3`
- **Command**: `wails3 dev` (for development)
- **Build**: `wails3 build` (for production)
- **Bindings**: Auto-generated TypeScript bindings in `frontend/bindings/`
- **Regenerate Bindings**: `wails3 generate bindings` (when Go services change)

### Package Structure
- **Core Types**: `internal/lfrtypes/` (renamed from generic `types`)
- **Utilities**: `internal/lfrutils/` (renamed from generic `utils`)
- **AWS Services**: `internal/aws/` with IAM and Lightsail clients

### Testing Strategy

#### Backend Testing
1. **Unit Tests**: `make test`
2. **Integration Tests**: `make integration-test` (LocalStack)
3. **Real AWS Tests**: `make integration-test-real` (requires credentials)

#### Frontend Testing
1. **Component Tests**: `npm run test` (Jest)
2. **E2E Tests (Mocked)**: `make e2e-test`
3. **E2E Tests (LocalStack)**: `make e2e-test-localstack`
4. **E2E Tests (Real AWS)**: `make e2e-test-aws`

### AWS Configuration
- **LocalStack**: Docker-based AWS simulation (IAM works, Lightsail is pro-only)
- **Real AWS**: Requires proper credentials and region configuration
- **AWS Profile**: User's AWS profile is 'aws' (not 'default')
- **Test Commands**: Use `AWS_PROFILE=aws` prefix for real AWS operations
- **Account**: Connected to AWS account 942542972736 with user 'scofri'
- **Region**: Primary testing in us-east-1, 15 regions available

### Build Commands
```bash
# CLI Development
make build          # Build CLI binary
make test          # Run all unit tests
make lint          # Code quality checks

# GUI Development
cd gui/lfr-gui
wails3 dev         # Start GUI in development mode
wails3 build       # Build GUI for production

# Testing
make integration-test      # LocalStack backend tests
make e2e-test-localstack  # GUI + LocalStack tests
make e2e-test-aws         # GUI + Real AWS tests
```

### Known Issues
1. **Wails3 Bindings**: May need regeneration with `wails3 generate` if TypeScript errors occur
2. **React Components**: Some JSX syntax issues in StudentConnect.tsx and StudentWorkspace.tsx
3. **LocalStack Limitations**: Lightsail API not fully supported (IAM works fine)

### Development Workflow
1. **Backend Changes**: Test with `make test` and `make integration-test`
2. **Frontend Changes**: Requires working Wails3 setup and proper bindings
3. **Full E2E**: Use LocalStack for safe testing, Real AWS for final validation

### Security Notes
- **Never commit AWS credentials**
- **Use IAM roles with minimal permissions** for testing
- **LocalStack is safe** for development (no real AWS resources)
- **Real AWS tests** may create actual resources - use with caution
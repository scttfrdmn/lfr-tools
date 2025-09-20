# End-to-End Testing Guide

This document explains how to run comprehensive end-to-end (e2e) tests for LFR Tools, including tests with both mocked services and real AWS integration.

## Overview

LFR Tools supports multiple levels of testing:

1. **Unit Tests** - Fast, isolated component tests
2. **Integration Tests** - AWS SDK integration with LocalStack
3. **E2E Tests (Mocked)** - Full UI testing with mocked AWS responses
4. **E2E Tests (LocalStack)** - UI + LocalStack backend integration
5. **E2E Tests (Real AWS)** - Complete end-to-end with real AWS services

## Quick Start

### Run All Tests
```bash
# Run all test types
make check

# Run specific test types
make test              # Unit tests only
make integration-test  # LocalStack integration
make e2e-test-all     # All e2e configurations
```

## E2E Testing Options

### 1. Mocked E2E Tests (Default)
Fast tests with simulated AWS responses - good for CI/CD.

```bash
make e2e-test
# or
cd gui/lfr-gui/frontend && npm run test:e2e
```

### 2. LocalStack E2E Tests
Tests against LocalStack for safe AWS simulation.

```bash
make e2e-test-localstack
# or manually:
make test-with-localstack
cd gui/lfr-gui/frontend && npm run test:e2e:localstack
make stop-localstack
```

### 3. Real AWS E2E Tests
**⚠️ CAUTION: Uses real AWS resources and may incur costs**

```bash
# Ensure AWS credentials are configured first
aws configure
# or set environment variables:
export AWS_ACCESS_KEY_ID="your-key"
export AWS_SECRET_ACCESS_KEY="your-secret"
export AWS_REGION="us-east-1"

# Run tests
make e2e-test-aws
# or
cd gui/lfr-gui/frontend && npm run test:e2e:aws
```

## AWS Credentials Setup

### Option 1: AWS CLI Configuration
```bash
aws configure
# Enter your AWS credentials when prompted
```

### Option 2: Environment Variables
```bash
export AWS_ACCESS_KEY_ID="AKIA..."
export AWS_SECRET_ACCESS_KEY="..."
export AWS_REGION="us-east-1"
export AWS_PROFILE="your-profile"  # Optional
```

### Option 3: IAM Roles (EC2/ECS)
If running on AWS infrastructure, use IAM roles attached to the compute resource.

## Test Development

### Adding New E2E Tests

1. **For UI-focused tests** - Add to existing `*.spec.ts` files with mocked AWS responses
2. **For AWS integration** - Add to `aws-integration.spec.ts`
3. **For specific workflows** - Create new spec files following the pattern

### Test File Structure
```
gui/lfr-gui/frontend/tests/e2e/
├── end-to-end-workflow.spec.ts     # Complete user workflows (mocked)
├── professor-workflow.spec.ts      # Professor-specific tests (mocked)
├── student-workflow.spec.ts        # Student-specific tests (mocked)
├── aws-integration.spec.ts         # Real AWS integration tests
└── setup/
    ├── aws-setup.ts               # AWS test environment setup
    ├── aws-teardown.ts            # AWS test cleanup
    ├── localstack-setup.ts        # LocalStack initialization
    └── localstack-teardown.ts     # LocalStack cleanup
```

### Configuration Files
```
├── playwright.config.ts          # Default config (mocked)
├── playwright.config.aws.ts      # Real AWS testing
└── playwright.config.localstack.ts # LocalStack testing
```

## Best Practices

### Development Workflow
1. **Start with mocked tests** for rapid development
2. **Use LocalStack** for integration testing
3. **Real AWS tests** only for final validation

### Safety Guidelines
- **Never run real AWS tests in CI/CD** without proper cost controls
- **Use test-specific AWS accounts** for real AWS testing
- **Clean up resources** after real AWS tests
- **Set billing alerts** when using real AWS for testing

### Performance Tips
- **LocalStack tests** are faster than real AWS
- **Mocked tests** are fastest for UI validation
- **Parallel execution** works for mocked/LocalStack, sequential for real AWS

## Troubleshooting

### Common Issues

#### LocalStack Not Starting
```bash
# Check Docker is running
docker ps

# Start LocalStack manually
make test-with-localstack

# Check LocalStack health
curl http://localhost:4566/_localstack/health
```

#### AWS Credentials Issues
```bash
# Test credentials
aws sts get-caller-identity

# Check environment
env | grep AWS

# Verify region
aws configure get region
```

#### GUI Build Errors
```bash
# Install dependencies
cd gui/lfr-gui/frontend && npm install

# Fix any JSX/TypeScript errors
npm run build
```

#### Port Conflicts
- **LocalStack**: Uses port 4566
- **GUI Dev Server**: Uses port 3000
- **Wails App**: May use different ports

```bash
# Check what's using ports
lsof -i :4566
lsof -i :3000

# Kill conflicting processes if needed
killall wails3
```

## Test Data Management

### Test Projects
The tests create temporary projects and resources:
- **Mocked**: No real resources created
- **LocalStack**: Ephemeral, cleaned up automatically
- **Real AWS**: **⚠️ May create actual AWS resources**

### Cleanup
```bash
# LocalStack cleanup (automatic)
make stop-localstack

# Real AWS cleanup (manual verification recommended)
aws iam list-users --query 'Users[?contains(UserName, `test-`)]'
aws lightsail get-instances --query 'instances[?contains(name, `test-`)]'
```

## CI/CD Integration

### GitHub Actions Example
```yaml
name: E2E Tests
on: [push, pull_request]

jobs:
  e2e-mocked:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Install dependencies
        run: cd gui/lfr-gui/frontend && npm install
      - name: Run mocked e2e tests
        run: make e2e-test

  e2e-localstack:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Install dependencies
        run: cd gui/lfr-gui/frontend && npm install
      - name: Run LocalStack e2e tests
        run: make e2e-test-localstack

  # Real AWS tests only on releases or specific branches
  e2e-aws:
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/release' || startsWith(github.ref, 'refs/tags/')
    environment: aws-testing  # Requires environment with AWS credentials
    steps:
      - uses: actions/checkout@v3
      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v2
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Install dependencies
        run: cd gui/lfr-gui/frontend && npm install
      - name: Run real AWS e2e tests
        run: make e2e-test-aws
```

## Monitoring and Reporting

### Test Reports
- **Playwright HTML reports** are generated in `gui/lfr-gui/frontend/playwright-report/`
- **Videos and screenshots** of failures in `gui/lfr-gui/frontend/test-results/`
- **Coverage reports** can be generated with `npm run test:coverage`

### Debugging Failed Tests
```bash
# Run tests in headed mode to see what's happening
cd gui/lfr-gui/frontend && npm run test:e2e:headed

# Use Playwright UI for interactive debugging
npm run test:e2e:ui

# Check the generated reports
open gui/lfr-gui/frontend/playwright-report/index.html
```

## Security Considerations

- **Never commit AWS credentials** to the repository
- **Use IAM roles with minimal permissions** for testing
- **Set up billing alerts** when using real AWS
- **Regularly audit test AWS accounts** for unexpected resources
- **Use separate AWS accounts** for testing vs production

## Next Steps

1. **Fix JSX errors** in React components to enable full e2e testing
2. **Set up AWS test account** with appropriate permissions
3. **Implement test data fixtures** for consistent testing
4. **Add performance testing** for AWS operations
5. **Create automated cleanup** for real AWS test resources
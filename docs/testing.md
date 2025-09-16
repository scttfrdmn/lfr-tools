# Testing Guide

This document covers testing strategies and tools for lfr-tools development.

## Test Types

### Unit Tests
Standard Go unit tests covering individual functions and components.

```bash
# Run all unit tests
make test

# Run tests with coverage
make coverage

# Test specific package
go test ./internal/config -v
```

### Integration Tests
Tests that interact with real or mocked AWS services using LocalStack or real AWS.

```bash
# Run integration tests with LocalStack (recommended)
make integration-test

# Run integration tests with real AWS (requires credentials)
make integration-test-real

# Manual development testing with LocalStack
make test-with-localstack
# Run specific tests while LocalStack is running
AWS_ENDPOINT_URL=http://localhost:4566 AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test go test -tags=integration ./internal/aws -v
make stop-localstack
```

## LocalStack Integration

### Overview
LocalStack provides a fully functional local AWS cloud stack for testing without costs or rate limits.

### Supported Services
- **IAM**: ✅ Full support for users, groups, policies, login profiles
- **Lightsail**: ❌ Not supported (returns 501 - Pro feature)
- **EC2**: ✅ Basic support for regions and metadata

### Setup Requirements
- Docker and Docker Compose installed
- No AWS credentials needed for LocalStack tests

### Test Results Summary

#### ✅ **Working with LocalStack:**
- **IAM Policy Creation**: Successfully creates custom policies with proper ARNs
- **IAM Group Management**: Create groups, attach policies, manage memberships
- **IAM User Management**: Create users, set passwords, manage group memberships
- **User Policy Attachment**: Inline user policies work correctly
- **Resource Cleanup**: Proper deletion and resource management

#### ❌ **LocalStack Limitations:**
- **Lightsail Service**: Not implemented in free version (501 error)
- **Complex IAM Features**: Some advanced IAM features may not work
- **Service Limits**: No realistic AWS service limits enforced

### Real AWS Integration Testing

The integration tests automatically detect real AWS credentials and will run against actual AWS services when:
- No `AWS_ENDPOINT_URL` is set
- Valid AWS credentials are available
- **Default Profile**: Uses `aws` profile (can be overridden with `AWS_PROFILE` env var)

**⚠️ Real AWS Testing Notes:**
- Tests create actual AWS resources
- Password policies are enforced (passwords must meet AWS requirements)
- Resources should be cleaned up automatically
- May incur small AWS charges

### Example Test Results

#### LocalStack Test Output:
```
=== RUN   TestIntegrationIAMOperations
    integration_test.go:105: Created policy: arn:aws:iam::000000000000:policy/LFRTestPolicy
    integration_test.go:115: Created group: LFRTestGroup
    integration_test.go:126: Created user: testuser
    integration_test.go:137: Cleanup completed
--- PASS: TestIntegrationIAMOperations (0.16s)

=== RUN   TestIntegrationLightsailOperations
    integration_test.go:150: Get blueprints failed (expected with LocalStack):
    failed to get blueprints: operation error Lightsail: GetBlueprints,
    https response error StatusCode: 501, RequestID: xxx,
    api error InternalFailure: API for service 'lightsail' not yet implemented
--- PASS: TestIntegrationLightsailOperations (0.01s)
```

#### Real AWS Test Output:
```
=== RUN   TestIntegrationIAMOperations
    integration_test.go:105: Created policy: arn:aws:iam::752123829273:policy/LFRTestPolicy
    integration_test.go:115: Created group: LFRTestGroup
    integration_test.go:120: User creation failed: Password should have at least one uppercase letter
--- PASS: TestIntegrationIAMOperations (1.28s)
```

## Test Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `AWS_ENDPOINT_URL` | LocalStack endpoint | `http://localhost:4566` |
| `LOCALSTACK_ENDPOINT` | Alternative LocalStack endpoint | None |
| `AWS_ACCESS_KEY_ID` | AWS access key (use "test" for LocalStack) | None |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key (use "test" for LocalStack) | None |
| `AWS_PROFILE` | AWS profile to use | None |
| `SKIP_INTEGRATION_TESTS` | Skip integration tests | `false` |

### Build Tags

Integration tests use build tags to separate them from unit tests:

```bash
# Run only integration tests
go test -tags=integration ./internal/aws

# Run unit tests (excludes integration)
go test ./internal/aws
```

## Development Workflow

### 1. Quick Testing with LocalStack
```bash
# Start LocalStack in background
make test-with-localstack

# Run specific tests during development
AWS_ENDPOINT_URL=http://localhost:4566 AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test \
  go test -tags=integration ./internal/aws -v -run="TestIntegrationIAMOperations"

# When done
make stop-localstack
```

### 2. Full Integration Testing
```bash
# Complete test cycle with automatic cleanup
make integration-test
```

### 3. Real AWS Validation
```bash
# Test against real AWS (requires valid credentials)
make integration-test-real
```

## Test Coverage Goals

- **Unit Tests**: >80% coverage for business logic
- **Integration Tests**: Cover all major AWS operations
- **End-to-End Tests**: Full CLI workflow testing (future)

## Troubleshooting

### LocalStack Issues

1. **Container won't start**
   ```bash
   docker ps -a
   docker logs lfr-tools-localstack
   ```

2. **Connection refused**
   - Ensure LocalStack is fully started (wait 10+ seconds)
   - Check port 4566 is not in use: `lsof -i :4566`

3. **Permission errors**
   - Ensure Docker daemon is running
   - Check Docker permissions for your user

### Integration Test Issues

1. **Tests fail with real AWS**
   - Check AWS credentials: `aws sts get-caller-identity`
   - Ensure sufficient IAM permissions
   - Check for password policy requirements

2. **Tests timeout**
   - Increase test timeout in setupIntegrationTest()
   - Check network connectivity

## Adding New Tests

### Unit Tests
```go
func TestNewFeature(t *testing.T) {
    fixture := testutils.NewTestFixture(t)

    // Test implementation
    result := NewFeature()

    fixture.AssertNoError(err)
    fixture.AssertEqual(expected, result)
}
```

### Integration Tests
```go
// +build integration

func TestIntegrationNewFeature(t *testing.T) {
    client := setupIntegrationTest(t)
    ctx := testutils.SetupTestContext()

    service := NewService(client)

    // Test against LocalStack/AWS
    result, err := service.DoSomething(ctx, "test-param")
    if err != nil {
        t.Logf("Expected error with LocalStack: %v", err)
        return
    }

    // Validate results and cleanup
}
```

## Continuous Integration

GitHub Actions automatically runs:
- All unit tests on every PR
- Integration tests can be enabled with LocalStack in CI
- Real AWS integration tests only on manual trigger

The testing infrastructure provides comprehensive coverage while keeping costs low through LocalStack simulation.
# Contributing to lfr-tools

Thank you for your interest in contributing to lfr-tools! This document provides guidelines and information for contributors.

## Development Setup

### Prerequisites

- Go 1.20 or higher
- Make
- golangci-lint
- gosec
- pre-commit
- AWS CLI configured
- GitHub CLI (gh)

### Initial Setup

1. **Fork and clone the repository:**
   ```bash
   gh repo fork scttfrdmn/lfr-tools --clone
   cd lfr-tools
   ```

2. **Install development dependencies:**
   ```bash
   make deps
   ```

3. **Set up pre-commit hooks:**
   ```bash
   pre-commit install
   ```

4. **Verify setup:**
   ```bash
   make check
   ```

## Development Workflow

### Making Changes

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes and ensure tests pass:**
   ```bash
   make test
   make lint
   make check
   ```

3. **Commit your changes using conventional commits:**
   ```bash
   git commit -m "feat: add new feature"
   ```

### Commit Message Format

We use [Conventional Commits](https://www.conventionalcommits.org/) for consistent commit messages:

- `feat:` new features
- `fix:` bug fixes
- `docs:` documentation changes
- `style:` formatting changes
- `refactor:` code refactoring
- `test:` adding tests
- `chore:` maintenance tasks

Examples:
- `feat(users): add bulk user creation`
- `fix(ssh): resolve key permission issues`
- `docs: update installation instructions`

### Code Quality

All contributions must:
- Pass all tests (`make test`)
- Pass linting (`make lint`)
- Pass security checks (`make sec`)
- Maintain or improve code coverage
- Follow Go best practices
- Include appropriate documentation

### Pre-commit Hooks

Pre-commit hooks automatically run:
- Go formatting (`go fmt`)
- Import organization (`goimports`)
- Linting (`golangci-lint`)
- Security scanning (`gosec`)
- Tests (`go test`)

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run specific tests
go test ./internal/aws/...
```

### Writing Tests

- Write unit tests for all new functionality
- Use table-driven tests where appropriate
- Mock external dependencies (AWS APIs)
- Aim for high test coverage

Example test structure:
```go
func TestUserCreate(t *testing.T) {
    tests := []struct {
        name     string
        input    CreateUserInput
        expected CreateUserOutput
        wantErr  bool
    }{
        // test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Documentation

### Code Documentation

- Document all public functions and types
- Use clear, concise comments
- Include examples where helpful
- Follow Go documentation conventions

### User Documentation

Update relevant documentation when making changes:
- README.md for feature changes
- CHANGELOG.md following Keep a Changelog format
- Command help text for CLI changes

## Pull Request Process

1. **Ensure your branch is up to date:**
   ```bash
   git checkout main
   git pull upstream main
   git checkout your-branch
   git rebase main
   ```

2. **Push your changes:**
   ```bash
   git push origin your-branch
   ```

3. **Create a pull request:**
   ```bash
   gh pr create --title "feat: your feature description" --body "Detailed description of changes"
   ```

4. **Ensure CI passes:**
   - All tests pass
   - Linting passes
   - Security checks pass
   - Go Report Card grade maintained

### Pull Request Guidelines

- Provide clear description of changes
- Reference any related issues
- Include tests for new functionality
- Update documentation as needed
- Keep changes focused and atomic
- Respond to review feedback promptly

## Issue Reporting

When reporting issues:
- Use the issue templates
- Provide clear reproduction steps
- Include relevant system information
- Add logs or error messages
- Specify expected vs actual behavior

## Feature Requests

When requesting features:
- Check existing issues first
- Provide clear use case description
- Explain why the feature is needed
- Consider implementation complexity
- Be open to alternative solutions

## Code Review

### As a Reviewer

- Focus on code quality and maintainability
- Provide constructive feedback
- Test the changes locally when possible
- Approve when satisfied with changes

### As an Author

- Be open to feedback
- Respond to comments promptly
- Make requested changes
- Ask questions if feedback is unclear

## Release Process

Releases are automated via GitHub Actions when tags are pushed:

1. Update CHANGELOG.md
2. Create and push tag: `git tag v1.0.0 && git push origin v1.0.0`
3. GoReleaser creates release and updates Homebrew tap

## Community

- Be respectful and inclusive
- Help others learn and grow
- Share knowledge and experience
- Follow the [Go Code of Conduct](https://golang.org/conduct)

## Getting Help

- Check existing documentation
- Search closed issues
- Ask questions in discussions
- Reach out to maintainers

Thank you for contributing to lfr-tools! 🚀
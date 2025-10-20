# GitHub Project Setup Guide

This document explains how the lfr-tools project is organized on GitHub and how to use the various features for tracking and managing development.

## Issue Templates

We use structured issue templates to ensure consistent reporting and feature requests:

### 🐛 Bug Report (`bug_report.yml`)
Use this template when something isn't working correctly. It includes:
- Component selection (CLI, GUI, AWS services, etc.)
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, version, AWS region)
- Log output and configuration

### ✨ Feature Request (`feature_request.yml`)
Use this template to suggest new features or enhancements. It includes:
- Persona selection (who benefits?)
- Component affected
- Roadmap phase alignment
- Problem statement and proposed solution
- Example workflows
- Priority assessment

### 🔧 Technical Debt (`technical_debt.yml`)
Use this template for code refactoring and technical improvements. It includes:
- Component and type of improvement
- Current state vs proposed improvement
- Benefits and scope of changes
- Breaking change assessment

### 📚 Documentation (`documentation.yml`)
Use this template for documentation improvements. It includes:
- Documentation type (guide, tutorial, API docs, etc.)
- Target audience (persona)
- Location and issue description
- Suggested improvements
- Reading level consideration

## Labels System

### Label Categories

Our labels are organized into logical categories:

#### Type Labels
- `bug` - Something isn't working correctly
- `enhancement` - New feature or request
- `documentation` - Doc improvements
- `technical-debt` - Code refactoring needed
- `security` - Security-related issues

#### Priority Labels
- `priority: critical` - Blocking work or severe impact
- `priority: high` - Should be addressed soon
- `priority: medium` - Important but not urgent
- `priority: low` - Nice to have

#### Area Labels
- `area: cli` - Command-line interface
- `area: gui` - Desktop GUI application
- `area: aws-iam` - AWS IAM integration
- `area: aws-lightsail` - AWS Lightsail integration
- `area: aws-efs` - AWS EFS integration
- `area: ssh` - SSH and connection management
- `area: dcv` - NICE DCV integration
- `area: idle-detection` - Idle detection system
- `area: config` - Configuration system
- `area: build` - Build system, CI/CD
- `area: tests` - Testing infrastructure
- `area: docs` - Documentation

#### Persona Labels
- `persona: professor` - Benefits professors
- `persona: teaching-assistant` - Benefits TAs
- `persona: student` - Benefits students
- `persona: researcher` - Benefits researchers
- `persona: lab-manager` - Benefits lab managers
- `persona: admin` - Benefits system administrators

#### Status Labels
- `triage` - Needs initial review
- `needs-info` - Waiting for more information
- `blocked` - Blocked by external dependency
- `ready` - Ready to be worked on
- `in-progress` - Currently being worked on
- `in-review` - In code review
- `awaiting-merge` - Approved and ready to merge

#### Resolution Labels
- `duplicate` - Already exists
- `wontfix` - Will not be worked on
- `invalid` - Not applicable
- `works-as-designed` - Intentional behavior

#### Special Labels
- `good first issue` - Good for newcomers
- `help wanted` - Community attention needed
- `breaking-change` - Breaks backward compatibility
- `performance` - Performance optimization
- `dependencies` - Dependency updates

#### Testing Labels
- `test: unit` - Related to unit tests
- `test: integration` - Related to integration tests (LocalStack)
- `test: e2e` - Related to end-to-end tests (GUI)
- `test: aws` - Requires real AWS testing

#### Phase/Roadmap Labels
- `phase: 1-enhanced-core` - Phase 1 features
- `phase: 2-advanced-mgmt` - Phase 2 features
- `phase: 3-enterprise` - Phase 3 features

#### Platform Labels
- `platform: macos` - Specific to macOS
- `platform: linux` - Specific to Linux
- `platform: windows` - Specific to Windows
- `platform: wails` - Related to Wails3 framework

## Applying Labels

### Using GitHub CLI
```bash
# Sync labels from configuration
gh label sync --labels .github/labels.yml --force

# Add labels to an issue
gh issue edit 123 --add-label "bug,area: cli,priority: high"
```

### Using GitHub Web Interface
1. Navigate to the issue or PR
2. Click "Labels" in the right sidebar
3. Search or browse labels
4. Click to add/remove

## Pull Request Template

Our PR template includes:
- Description and related issues
- Type of change checkboxes
- Persona impact assessment
- Comprehensive testing checklist
- Performance and AWS impact sections
- Breaking change documentation

## GitHub Discussions

We use GitHub Discussions for:
- 💬 **Q&A**: Ask questions about using LFR Tools
- 💡 **Ideas**: Discuss potential features before creating formal requests
- 🎉 **Show and Tell**: Share workflows and success stories

## GitHub Projects

(To be implemented)

We will use GitHub Projects for:
- **Roadmap tracking**: Track progress on Phase 1, 2, and 3 features
- **Sprint planning**: Organize work into 2-week sprints
- **Release planning**: Track issues for upcoming releases

### Project Views
- **Roadmap**: Grouped by phase/milestone
- **Sprint Board**: Kanban view for current sprint
- **By Persona**: Grouped by target user persona
- **By Component**: Grouped by technical area

## GitHub Pages

Our documentation is published to GitHub Pages at:
https://scttfrdmn.github.io/lfr-tools/

### Structure
- **Home**: Project overview and quick links
- **Getting Started**: Installation and quick start guide
- **Documentation**: Detailed guides and references
- **Tutorials**: Step-by-step walkthroughs
- **API Reference**: Generated godoc (future)

### Updating Documentation
Documentation is automatically deployed when changes are pushed to `main`:
1. Edit markdown files in `docs/`
2. Commit and push to `main`
3. GitHub Actions builds and deploys to Pages

## Workflow Tips

### Creating a Bug Report
1. Click "Issues" → "New issue"
2. Select "🐛 Bug Report"
3. Fill out all required fields
4. The `bug` and `triage` labels are automatically applied
5. Maintainers will add additional labels during triage

### Proposing a Feature
1. Consider starting a Discussion in "Ideas" first
2. When ready, create a "✨ Feature Request" issue
3. Link to the Discussion if applicable
4. Explain which persona benefits and why
5. Show how it improves current workflows

### Contributing Code
1. Check "good first issue" label for beginner-friendly tasks
2. Comment on an issue to claim it
3. Create a branch: `feature/short-description` or `fix/short-description`
4. Make your changes and test thoroughly
5. Submit a PR using the template
6. Link the related issue with "Closes #123"
7. Address review feedback

### Reporting Security Issues
Do NOT create a public issue. Instead:
- Email: scofri@g.ucla.edu
- Include details and reproduction steps
- We'll acknowledge within 48 hours
- We'll provide a fix timeline

## Maintenance Tasks

### Weekly
- Review and triage new issues
- Update issue status labels
- Close resolved issues
- Merge approved PRs

### Monthly
- Review and update roadmap
- Archive completed milestones
- Update documentation
- Review and update labels if needed

### Quarterly
- Audit issue templates (retire stale fields, add missing ones)
- Review project board structure
- Update GitHub Pages documentation
- Clean up stale branches

## Resources

- [GitHub Issues Best Practices](https://docs.github.com/en/issues/tracking-your-work-with-issues/quickstart)
- [GitHub Projects Guide](https://docs.github.com/en/issues/planning-and-tracking-with-projects/learning-about-projects/about-projects)
- [GitHub Pages Documentation](https://docs.github.com/en/pages)
- [GitHub CLI Manual](https://cli.github.com/manual/)

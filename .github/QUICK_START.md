# GitHub Organization Quick Start

## Immediate Actions to Take

### 1. Sync Labels (2 minutes)
```bash
cd /Users/scttfrdmn/src/lfr-tools
gh label sync --labels .github/labels.yml --force
```

This will create 80+ organized labels for your repository.

### 2. Enable GitHub Discussions (1 minute)
1. Go to https://github.com/scttfrdmn/lfr-tools/settings
2. Scroll to "Features"
3. Check ✅ "Discussions"
4. Click "Set up discussions"
5. Create categories:
   - **Q&A** - Questions and answers
   - **Ideas** - Feature discussions
   - **Show and Tell** - Success stories
   - **General** - General discussions

### 3. Enable GitHub Pages (1 minute)
1. Go to https://github.com/scttfrdmn/lfr-tools/settings/pages
2. Under "Build and deployment"
3. Source: **Deploy from a branch**
4. Branch: **main** → folder: **/docs**
5. Click "Save"
6. Wait 2-3 minutes, then visit: https://scttfrdmn.github.io/lfr-tools/

### 4. Test Issue Templates (2 minutes)
1. Go to https://github.com/scttfrdmn/lfr-tools/issues/new/choose
2. Verify you see:
   - 🐛 Bug Report
   - ✨ Feature Request
   - 🔧 Technical Debt
   - 📚 Documentation
   - Links to Discussions
3. Try creating a test issue with one template

### 5. Create First GitHub Project (5 minutes)
1. Go to https://github.com/scttfrdmn/lfr-tools/projects
2. Click "New project"
3. Choose "Board" template
4. Name it "LFR Tools Development"
5. Add columns:
   - 📋 Backlog
   - 🔍 Triage
   - 📅 Ready
   - 🏗️ In Progress
   - 👀 In Review
   - ✅ Done
6. Link to repository

## Daily Workflow

### Creating Issues
```bash
# Bug report
gh issue create --template bug_report.yml

# Feature request
gh issue create --template feature_request.yml

# Technical debt
gh issue create --template technical_debt.yml

# Documentation
gh issue create --template documentation.yml
```

### Managing Labels
```bash
# Add labels to issue #123
gh issue edit 123 --add-label "bug,area: cli,priority: high"

# List issues by label
gh issue list --label "triage"
gh issue list --label "good first issue"
gh issue list --label "persona: professor"

# Search specific combinations
gh issue list --label "area: gui" --label "priority: high" --state open
```

### Viewing Issues
```bash
# Issues needing triage
gh issue list --label "triage" --state open

# Ready for work
gh issue list --label "ready" --state open

# In progress
gh issue list --label "in-progress" --state open

# By persona
gh issue list --label "persona: student" --state open
```

## Weekly Maintenance (15 minutes)

### Monday: Triage New Issues
```bash
# List untriaged issues
gh issue list --label "triage" --state open

# For each issue:
# 1. Add priority label
# 2. Add area label  
# 3. Add persona label (if applicable)
# 4. Change status to "ready" or "needs-info"
# 5. Remove "triage" label

# Example:
gh issue edit 123 \
  --add-label "priority: high,area: cli,persona: professor,ready" \
  --remove-label "triage"
```

### Friday: Close Resolved Issues
```bash
# List issues in review
gh issue list --label "in-review" --state open

# Close merged issues
gh issue close 123 --reason completed
```

## Monthly Maintenance (30 minutes)

### Review Stale Issues
```bash
# Find old "needs-info" issues (>30 days)
gh issue list --label "needs-info" --json number,title,updatedAt | \
  jq '.[] | select(.updatedAt < "2024-09-01")'

# Close with comment
gh issue close 123 --comment "Closing due to no response. Please reopen if still relevant."
```

### Update Documentation
```bash
# Edit docs
code docs/

# Commit and push (auto-deploys to Pages)
git add docs/
git commit -m "docs: update getting started guide"
git push
```

## Pro Tips

### Creating Issues from CLI
```bash
# Quick bug report
gh issue create \
  --title "[Bug]: SSH connection fails on macOS" \
  --label "bug,area: ssh,priority: high,triage" \
  --body "Description..."

# Quick feature request
gh issue create \
  --title "[Feature]: Add cost reporting" \
  --label "enhancement,area: cli,persona: professor,triage" \
  --body "As a professor..."
```

### Bulk Label Operations
```bash
# Add persona label to multiple issues
for issue in 45 46 47 48; do
  gh issue edit $issue --add-label "persona: student"
done

# Change priority for related issues
gh issue list --label "area: gui" --json number --jq '.[].number' | \
  xargs -I {} gh issue edit {} --add-label "priority: high"
```

### Project Automation
```bash
# Add all "ready" issues to project
gh issue list --label "ready" --json number --jq '.[].number' | \
  xargs -I {} gh project item-add 1 --owner scttfrdmn --issue {}
```

## Common Patterns

### New Feature Workflow
1. User creates feature request or discusses in Ideas
2. Maintainer adds labels: `enhancement`, `area: X`, `persona: Y`, `triage`
3. During triage, add `priority: X` and `phase: X`
4. Move to `ready` when design is approved
5. Developer claims issue, change to `in-progress`
6. PR created, change to `in-review`
7. PR merged, issue auto-closes

### Bug Fix Workflow
1. User reports bug with template
2. Auto-labeled: `bug`, `triage`
3. Maintainer reproduces, adds `priority: X` and `area: X`
4. If can't reproduce, add `needs-info`
5. Move to `ready` when confirmed
6. Developer fixes, PR created
7. Add `test: unit` or `test: integration` label
8. PR merged, issue auto-closes

### Documentation Update Workflow
1. User or maintainer creates documentation issue
2. Add `documentation`, `priority: X`, target persona
3. Writer claims issue
4. Edit in `docs/`
5. PR created (auto-deploys on merge to GitHub Pages)
6. Issue closes

## Useful Queries

### Find Work
```bash
# Good first issues for new contributors
gh issue list --label "good first issue" --state open

# High priority bugs
gh issue list --label "bug" --label "priority: high" --state open

# Documentation needed
gh issue list --label "documentation" --label "ready" --state open
```

### Project Health
```bash
# Count by status
gh issue list --label "in-progress" --json number | jq 'length'
gh issue list --label "ready" --json number | jq 'length'
gh issue list --label "triage" --json number | jq 'length'

# Count by area
for area in cli gui aws-iam aws-lightsail ssh dcv; do
  count=$(gh issue list --label "area: $area" --json number | jq 'length')
  echo "$area: $count"
done
```

## Reference

- **Full Setup Guide**: `.github/PROJECT_SETUP.md`
- **Refactoring Summary**: `.github/GITHUB_REFACTORING_SUMMARY.md`
- **Labels File**: `.github/labels.yml`
- **Issue Templates**: `.github/ISSUE_TEMPLATE/`
- **GitHub Pages**: `docs/`

## Need Help?

- 📚 Read: [PROJECT_SETUP.md](.github/PROJECT_SETUP.md)
- 💬 Ask: [GitHub Discussions](https://github.com/scttfrdmn/lfr-tools/discussions)
- 📖 Docs: [GitHub CLI Manual](https://cli.github.com/manual/)

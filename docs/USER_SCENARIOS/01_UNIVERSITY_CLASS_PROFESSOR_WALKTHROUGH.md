# Scenario 1: University Class Management - Professor Perspective

## Persona: Professor Dr. Maria Rodriguez

**Background**:
- Computer Science professor at mid-sized university
- Teaching CS 305 - Cloud Computing (Fall 2025)
- Enrollment: 35 students (undergraduate juniors/seniors)
- IT Department Budget: $840 for semester ($24/student for 15 weeks)
- Technical level: Strong programmer, moderate cloud experience
- Previous experience: Used personal AWS account (billing nightmares!)

**Primary Concerns**:
1. **Budget Control**: Must stay within $840 - no overruns tolerated by IT
2. **Student Access**: All 35 students need isolated, secure environments
3. **Data Privacy**: Student work must be isolated (academic integrity)
4. **Ease of Management**: Teaching 2 courses + research = no time for IT support
5. **Semester Lifecycle**: Setup (Week 1), Active (Weeks 2-14), Cleanup (Week 15)

**Pain Points from Previous Semester**:
- Student accidentally launched p3.8xlarge GPU instance ($12/hour!)
- Forgot to stop instances over Thanksgiving break (+$280 surprise bill)
- Spent 8 hours manually creating 35 AWS accounts
- End-of-semester cleanup took entire weekend
- One student shared SSH keys (academic integrity issue)

---

## Current State (v0.1.x): What Works Today

### ✅ Week 1: Pre-Semester Setup (What Works)

#### Day 1: Create Project for Course
```bash
# Dr. Rodriguez installs lfr-tools
brew install lfr

# Configure AWS credentials (using university AWS account)
aws configure --profile university-cs
# AWS Access Key ID: [university provided]
# AWS Secret Access Key: [university provided]
# Default region: us-east-1
# Output format: json

# Set LFR to use university profile
export AWS_PROFILE=university-cs

# Create users for all 35 students (bulk operation)
# students.csv format: username
# alice, bob, charlie, ..., zoe (35 total)
lfr users create \
  --project "CS305-Fall2025" \
  --blueprint "ubuntu_22_04" \
  --bundle "nano_2_0" \
  --region "us-east-1" \
  --users "$(cat students.csv)"

# LFR Output:
# Creating 35 users for project CS305-Fall2025...
# ✅ Created IAM user: alice (password: xK9mP2...)
# ✅ Created IAM user: bob (password: 7Lq3N...)
# ... (33 more)
# ✅ Created IAM user: zoe (password: 9Rm4W...)
#
# ✅ Created 35 Lightsail for Research instances
# ✅ Tagged with project: CS305-Fall2025
#
# 📊 Cost estimate: ~$29/day ($29 = 35 × $0.83/day)
#    Semester budget: $840
#    At this rate: ~29 days of 24/7 usage (within budget if instances stopped!)
#
# 🔑 Credentials saved to: cs305-credentials.csv
# 📧 Send to students via Canvas

# What Dr. Rodriguez thinks: "That was WAY easier than last year!"
```

**What works well:**
- ✅ Bulk user creation
- ✅ Auto-generated passwords
- ✅ Project tagging for organization
- ✅ Credential export for distribution

#### Day 3: Distribute Credentials to Students

```bash
# Dr. Rodriguez uploads cs305-credentials.csv to Canvas
# Students see:
# - Username: alice
# - Password: xK9mP2yT3zL
# - Region: us-east-1
# - Instance name: alice-ubuntu_22_04
# - SSH command: lfr ssh connect alice --project CS305-Fall2025

# Students install lfr-tools and connect
# (See student walkthrough for details)
```

**What works well:**
- ✅ Simple CSV for credential distribution
- ✅ Students can connect immediately
- ✅ Instance names match usernames

---

## ⚠️ Current Pain Points: What Doesn't Work

### ❌ Problem 1: No Per-Student Budget Tracking

**Scenario**: Week 3, one student forgets to stop instance over weekend

**What should happen** (MISSING):
```bash
# Dr. Rodriguez checks class budget status
lfr project budget --project CS305-Fall2025

# LFR Output (MISSING):
# 📊 Project Budget: CS305-Fall2025
#
# Total Budget: $840.00
# Spent to date: $156.80 (18.6%) - Week 3 of 15
# Remaining: $683.20
# Projected end-of-semester: $785.00 ✅ (within budget)
#
# Per-Student Breakdown:
# ┌──────────┬────────────────┬────────────┬──────────────┬─────────┐
# │ Student  │ Instance       │ Status     │ Spent        │ % Budget│
# ├──────────┼────────────────┼────────────┼──────────────┼─────────┤
# │ alice    │ Running (48h)  │ ⚠️  Long    │ $3.98        │ 16.6%   │
# │ bob      │ Stopped        │ ✅         │ $2.15        │ 9.0%    │
# │ charlie  │ Stopped        │ ✅         │ $2.30        │ 9.6%    │
# │ ...      │ ...            │ ...        │ ...          │ ...     │
# │ dave     │ Running (72h)  │ 🔴 Alert!  │ $5.98        │ 24.9%   │
# │ ...      │ ...            │ ...        │ ...          │ ...     │
# │ zoe      │ Stopped        │ ✅         │ $1.85        │ 7.7%    │
# └──────────┴────────────────┴────────────┴──────────────┴─────────┘
#
# 🚨 Alerts:
# - dave: Instance running for 72 hours (3 days straight!)
#         Cost: $5.98 (24.9% of student budget used in Week 3)
#         Recommendation: Contact student or auto-stop
#
# 💡 Budget Protection Suggestions:
# 1. Enable auto-stop for idle instances (4 hour threshold)
# 2. Set per-student budget caps ($24/student)
# 3. Email alerts at 50%, 75%, 90% per student
```

**Current workaround**: Dr. Rodriguez manually checks AWS console, does math in Excel
**Impact**: Can't catch runaway costs until too late, student anxiety

### ❌ Problem 2: No Bulk Operations for Class Management

**Scenario**: Thanksgiving break - need to stop all 35 instances

**What should happen** (MISSING):
```bash
# Dr. Rodriguez wants to stop all instances for Thanksgiving (save $200)
lfr project stop --project CS305-Fall2025 --all

# LFR Output (MISSING):
# 🛑 Stopping all instances for project: CS305-Fall2025
#
# Instances to stop: 35
# - 18 currently running
# - 17 already stopped
#
# Estimated savings: $14.94/day ($104.58/week)
#
# Proceed? [y/N]: y
#
# Stopping instances...
# ✅ Stopped: alice-ubuntu_22_04
# ✅ Stopped: bob-ubuntu_22_04
# ✅ Already stopped: charlie-ubuntu_22_04
# ... (32 more)
#
# ✅ All instances stopped
# 💰 Projected Thanksgiving savings: $104.58 (saved ~12% of semester budget!)
#
# 📧 Notification sent to all 35 students:
#    "Your CS 305 instance has been stopped for Thanksgiving break.
#     You can restart it on Monday, Nov 27."

# Monday after break: Restart for students who need it
lfr project start --project CS305-Fall2025 --user alice,bob,charlie
```

**Current workaround**: Dr. Rodriguez manually stops instances one-by-one (30+ minutes)
**Impact**: Time waste, may forget some instances, no student notification

### ❌ Problem 3: No Template Restrictions (Budget Risk)

**Scenario**: Student tries to launch GPU instance for "extra credit project"

**What should happen** (MISSING):
```bash
# Student (from their perspective):
lfr instances create dave-gpu --blueprint ubuntu_22_04 --bundle large_4_0

# LFR Output (MISSING - SHOULD BLOCK):
# ❌ Instance launch BLOCKED by course policy
#
#    Requested: large_4_0 bundle ($3.60/day)
#    Project: CS305-Fall2025
#    Allowed bundles: nano_2_0 ($0.83/day), micro_2_0 ($1.67/day)
#
#    Reason: Professor has restricted instance types for budget control.
#
#    Your current budget: $18.50 / $24.00 (77%)
#    This instance would cost: ~$50 for rest of semester (exceeds budget!)
#
#    If you need a larger instance:
#    1. Email professor: maria.rodriguez@university.edu
#    2. Explain your use case
#    3. Professor can grant exception
#
# Professor view (when student emails):
lfr project allow-bundle --project CS305-Fall2025 --user dave --bundle large_4_0 --days 3 --reason "Extra credit ML project"

# LFR Output:
# ✅ Temporary exception granted
#    Student: dave
#    Bundle: large_4_0 (normally restricted)
#    Duration: 3 days
#    Reason: Extra credit ML project
#
# 📧 Notification sent to dave:
#    "You have been granted temporary access to large_4_0 bundle for 3 days.
#     Estimated cost: $10.80 (will count against your $24 budget).
#     Please stop instance when done to conserve budget."
```

**Current state**: Any student can launch any instance type (budget risk!)
**Impact**: One student with GPU instance = entire class budget blown

### ❌ Problem 4: No Semester End Automation

**Scenario**: December 15, semester ends - need to cleanup everything

**What should happen** (MISSING):
```bash
# Dr. Rodriguez wants to close the course and cleanup
lfr project close --project CS305-Fall2025 --archive

# LFR Output (MISSING):
# 🎓 Closing Course Project: CS305-Fall2025
#
# Actions to perform:
# ✅ Stop all running instances (18 active)
# ✅ Revoke all student IAM credentials
# ✅ Archive student work to S3 (for grade disputes)
# ✅ Generate final cost report
# ✅ Delete IAM users (30 days grace period)
#
# This will:
# - Stop all 35 instances immediately
# - Archive student workspaces to s3://cs305-fall2025-archive/
# - Generate final expense report for IT
# - Revoke all student access (they can't start instances)
# - Preserve data for 30 days (academic integrity/grade disputes)
#
# ⚠️  Students will lose access immediately!
#    Make sure all grading is complete.
#
# Proceed? [y/N]: y
#
# Closing course...
# ✅ Stopped 18 running instances
# ✅ Archiving student workspaces... (35/35 complete)
#    Archive location: s3://university-cs/cs305-fall2025/
# ✅ Generated final report: cs305-fall2025-final-report.pdf
# ✅ Revoked IAM credentials (35 users)
# ✅ Tagged resources for deletion (30-day grace period)
#
# 📊 Final Semester Report:
#
# Budget: $840.00
# Spent: $782.50 (93.2%) ✅
# Saved: $57.50
#
# Student Efficiency:
# - 28 students: Within $24 budget ✅
# - 5 students: Slightly over ($25-$28) ⚠️
# - 2 students: Significantly over ($32+) 🔴
#
# Cost Breakdown:
# - Compute (LFR instances): $782.50
# - Storage (snapshots): $0.00
# - Data transfer: $0.00
#
# Top Cost Drivers:
# 1. dave: $32.40 (135% of budget) - Ran large instance for extra credit
# 2. eve: $28.75 (120% of budget) - Forgot to stop before breaks
# 3-35. Others: Average $21.80/student ✅
#
# 📧 Report sent to:
# - IT Department (for billing reconciliation)
# - Department Chair (budget justification)
# - You (course records)
#
# 🎉 Course successfully closed!
#    All data archived, students notified, access revoked.
```

**Current workaround**: Manual cleanup over entire weekend:
1. Stop instances one by one
2. Manually export AWS billing
3. Delete IAM users manually
4. Hope you didn't forget anything

**Impact**: 4-8 hours of tedious work, high error risk, continued costs if missed

### ❌ Problem 5: No TA Access/Delegation

**Scenario**: TA needs to help student debug broken environment

**What should happen** (MISSING):
```bash
# TA (Alex) needs to help student (alice) during office hours
# Alice in Zoom: "My environment is broken, can't run Python anymore!"

# TA perspective (MISSING):
lfr ta debug --project CS305-Fall2025 --student alice

# LFR Output (MISSING):
# 🔍 TA Debug Access
#
# Student: alice
# Instance: alice-ubuntu_22_04 (running)
# Project: CS305-Fall2025
# Your role: Teaching Assistant (granted by Prof. Rodriguez)
#
# Available actions:
# 1. View instance status
# 2. SSH into instance (read-only)
# 3. View student's recent commands (audit log)
# 4. Reset instance (backup + restore to clean state)
# 5. Send message to student
#
# Choice [1-5]: 2
#
# 🔑 Granting temporary SSH access (30-minute session)
#    All actions will be logged for academic integrity
#
# Connecting...

# TA SSHs into alice's instance (logged session)
# (Inside alice's environment)
ta@alice-ubuntu:~$ python3 --version
# bash: python3: command not found

# TA realizes: alice accidentally uninstalled Python
# Instead of fixing (academic integrity), TA documents issue

# Exit SSH, reset instance for alice
lfr ta reset-instance --project CS305-Fall2025 --student alice --reason "Python uninstalled"

# LFR Output:
# 🔄 Resetting Instance: alice-ubuntu_22_04
#
# This will:
# ✅ Backup current state to S3 (in case needed)
# ✅ Stop instance
# ✅ Create fresh instance from blueprint
# ✅ Restore alice's home directory (/home/alice/)
# ✅ Preserve alice's work
# ❌ Discard broken system state
#
# Estimated downtime: 5 minutes
# Cost: $0.00 (no additional charge)
#
# Proceed? [y/N]: y
#
# Resetting...
# ✅ Backed up to: s3://cs305-fall2025/backups/alice-20251103.tar.gz
# ✅ Fresh instance created
# ✅ Restored alice's files
# ✅ Ready to use!
#
# 📧 Notification sent to alice:
#    "Your instance has been reset by TA Alex.
#     Your files are preserved. You can continue working now."
#
# 📧 Notification sent to Prof. Rodriguez:
#    "TA Alex reset alice's instance (reason: Python uninstalled)
#     Backup available: s3://cs305-fall2025/backups/alice-20251103.tar.gz"
```

**Current workaround**: TA guides student via screen share (slow, frustrating)
**Impact**: Office hours inefficient, students feel unsupported

---

## 🎯 Ideal Future State: Complete Class Lifecycle

### Phase 1: Pre-Semester Setup (Week -1)

```bash
# Interactive course creation wizard
lfr course create --interactive

# Wizard output:
# 🎓 LFR Course Creation Wizard
#
# Course Information:
#   Course code: CS 305
#   Title: Cloud Computing
#   Term: Fall 2025
#   Start date: August 25, 2025
#   End date: December 12, 2025 (15 weeks)
#
# Enrollment:
#   Expected students: 35
#   Per-student budget: $24.00
#   Total budget: $840.00
#   Budget source: University IT Department
#
# Teaching Staff:
#   Professor: maria.rodriguez@university.edu (you)
#   TAs: [Enter email addresses or skip]
#        alex.chen@university.edu
#        sarah.kim@university.edu
#   TA permissions: [debug, reset, view-only]
#
# Student Environment:
#   Default blueprint: ubuntu_22_04
#   Allowed bundles:
#     [x] nano_2_0 ($0.83/day) - Default
#     [x] micro_2_0 ($1.67/day) - For advanced projects
#     [ ] large_4_0 ($3.60/day) - Blocked (too expensive)
#
#   Instance limits per student:
#     Max concurrent instances: 1
#     Max daily cost: $2.00
#     Auto-stop after: 4 hours idle
#
# Budget Protection:
#   [x] Per-student budget caps ($24.00 each)
#   [x] Email alerts at 75%, 90% per student
#   [x] Auto-stop idle instances (4 hour threshold)
#   [x] Block over-budget launches
#
# Semester End Actions:
#   [x] Auto-stop all instances on Dec 12, 11:59 PM
#   [x] Revoke student access
#   [x] Archive workspaces to S3 (30-day retention)
#   [x] Generate final report
#
# Academic Integrity:
#   [x] SSH key isolation (no sharing)
#   [x] Audit logging (all student commands)
#   [x] TA access logging
#
# Setup complete! ✅
#
# Next steps:
# 1. Import students: lfr course import-students CS305-Fall2025 --canvas
# 2. Test environment: lfr course test
# 3. Distribute credentials: lfr course send-welcome-emails

# Import students from Canvas LMS
lfr course import-students CS305-Fall2025 \
  --canvas \
  --course-id 98765 \
  --api-key [canvas-api-key]

# LFR Output:
# 📥 Importing students from Canvas...
# ✅ Found 35 enrolled students
# ✅ Creating IAM users (35/35)
# ✅ Creating instances (35/35)
# ✅ Setting up per-student budgets ($24 each)
# ✅ Generating SSH keys
# ✅ Configuring auto-stop policies
#
# 📊 Resource Summary:
# - 35 IAM users created
# - 35 Lightsail instances (all stopped)
# - Total setup time: 8 minutes
# - Ready for semester!
#
# 📧 Send welcome emails: lfr course send-welcome-emails CS305-Fall2025

# Test the environment before semester starts
lfr course test CS305-Fall2025

# LFR Output:
# 🧪 Testing CS305-Fall2025 environment...
# ✅ IAM users can authenticate
# ✅ Instances can start
# ✅ SSH access works
# ✅ Budget tracking configured
# ✅ Auto-stop policies active
# ✅ TA access configured
#
# All systems ready! 🎉
```

### Phase 2: Active Semester (Weeks 1-14)

```bash
# Week 1: Check class status
lfr course status CS305-Fall2025

# LFR Output:
# 📊 Course Status: CS 305 - Cloud Computing (Fall 2025)
#
# Week: 1 of 15
# Enrollment: 35 students
# Budget: $28.50 / $840.00 (3.4%)
# Remaining: $811.50
#
# Active Instances: 28 / 35 (80% students working)
# Stopped Instances: 7
#
# Budget Health: ✅ Excellent
# Projected end-of-semester: $798 (95% of budget)
#
# Alerts: None
# Recent Activity:
# - 28 students connected in last 24 hours
# - Average session: 2.5 hours
# - All students within budget

# Week 5: Budget alert for one student
lfr course alerts

# LFR Output:
# 🚨 Course Alerts: CS305-Fall2025
#
# 1. Student Budget Warning: dave
#    Spent: $18.00 / $24.00 (75%) - Week 5 of 15
#    Issue: Instance running for 36 hours straight
#    Recommendation: Contact student or enable auto-stop
#
#    Actions:
#    a) Send reminder email to student
#    b) Force stop instance now
#    c) Enable stricter auto-stop (2 hours instead of 4)
#    d) Ignore (monitor)
#
# Choice [a/b/c/d]: a
#
# 📧 Reminder sent to dave:
# "Your CS 305 instance has been running for 36 hours. Please stop when not in use.
#  You've used 75% of your semester budget in Week 5. Budget remaining: $6.00."

# Week 11: Thanksgiving break
lfr course vacation --start 2025-11-22 --end 2025-11-27 --stop-all

# LFR Output:
# 🦃 Vacation Mode: CS305-Fall2025
#
# Dates: Nov 22-27, 2025 (6 days)
# Actions:
# - Stop all 35 instances on Nov 21, 11:59 PM
# - Block instance starts during vacation
# - Auto-enable on Nov 27, 12:00 AM
#
# Estimated savings: $174.30 (20% of semester budget!)
#
# 📧 Notification sent to all 35 students:
# "Your CS 305 instance will be stopped for Thanksgiving break (Nov 22-27).
#  You can resume work on Monday, Nov 27. Happy Thanksgiving!"

# Week 14: Generate progress report for department
lfr course report CS305-Fall2025 --pdf

# LFR Output:
# 📄 Generating course report...
# ✅ Report saved: cs305-fall2025-week14-report.pdf
#
# Report includes:
# - Budget status (week-by-week)
# - Per-student usage
# - Efficiency metrics
# - Cost projections
# - Usage patterns
# - Ready for department meeting

### Phase 3: Semester End (Week 15)

```bash
# December 12, 11:59 PM - Automatic actions
# (Dr. Rodriguez doesn't need to do anything!)

# System log:
# 2025-12-12 23:50:00 [CS305-Fall2025] 10-minute warning sent to students
# 2025-12-12 23:59:59 [CS305-Fall2025] Semester end triggered
# 2025-12-12 23:59:59 Stopping 22 active instances...
# 2025-12-13 00:00:15 ✅ All instances stopped
# 2025-12-13 00:00:30 Archiving student workspaces...
# 2025-12-13 00:01:45 ✅ 35 workspaces archived to S3
# 2025-12-13 00:02:00 Revoking student IAM credentials...
# 2025-12-13 00:02:15 ✅ 35 credentials revoked
# 2025-12-13 00:02:30 Generating final report...
# 2025-12-13 00:03:00 ✅ Semester closure complete

# December 13, 8:00 AM - Dr. Rodriguez receives email:
#
# Subject: 📊 CS 305 Fall 2025 - Final Course Report
#
# Your course CS 305 - Cloud Computing has successfully completed.
#
# Semester: Fall 2025 (August 25 - December 12, 15 weeks)
# Enrollment: 35 students
#
# Budget Performance:
# Total budget: $840.00
# Total spent: $782.50 (93.2%) ✅
# Unused: $57.50
#
# Per-Student Breakdown:
# - Average spend: $22.36 / $24.00 (93%)
# - Range: $18.40 - $32.40
# - Within budget: 28 students (80%) ✅
# - Slightly over: 5 students (avg $26.50)
# - Significantly over: 2 students ($32+)
#
# Usage Statistics:
# - Total compute hours: 2,940 hours
# - Average per student: 84 hours (5.6 hours/week)
# - Peak week: Week 12 (final project)
# - Auto-stop savings: $184.50 (19% of budget)
#
# Teaching Assistant Activity:
# - Debug sessions: 18 (avg 25 minutes each)
# - Instance resets: 5
# - Messages sent: 47
#
# Data Archive:
# - Student workspaces: s3://university-cs/cs305-fall2025/students/
# - Retention: 30 days (until Jan 12, 2026)
# - Total size: 15.8 GB
#
# Cost Comparison:
# - CS 305 Fall 2025: $782.50 (35 students)
# - CS 305 Fall 2024: $1,240.00 (32 students) - 37% savings! ✅
# - Improvement: Better auto-stop policies, student budget caps
#
# Next Semester:
# To prepare for Spring 2026:
# $ lfr course duplicate CS305-Fall2025 --name CS305-Spring2026
#
# All student access revoked. Course data archived. ✅

# Dr. Rodriguez can relax - everything handled automatically!
```

---

## 📋 Feature Gap Analysis: Professor Needs

### Critical Missing Features (Blocks Professor Persona)

| Feature | Priority | Impact | Current State | Effort |
|---------|----------|--------|---------------|--------|
| **Per-Student Budget Tracking** | 🔴 Critical | Can't manage class budgets | No tracking | High |
| **Project-Wide Bulk Operations** | 🔴 Critical | Manual stop/start of 35 instances | Basic only | Medium |
| **Template/Bundle Restrictions** | 🔴 Critical | Students can blow budget | None | Medium |
| **Semester End Automation** | 🔴 Critical | 4-8 hour manual cleanup | None | Medium |
| **TA Access Delegation** | 🟡 High | TAs can't help students | No TA role | High |
| **Student Budget Caps** | 🟡 High | No per-student enforcement | None | Medium |
| **Canvas/LMS Integration** | 🟡 High | Manual student import | None | High |
| **Course Management Commands** | 🟡 High | No course abstraction | None | Medium |
| **Instance Reset for TAs** | 🟢 Medium | TAs need this for support | None | Low |
| **Cost Reporting (PDF)** | 🟢 Medium | IT department needs reports | CSV only | Low |

### Unique Professor Requirements

| Requirement | Current | Needed | Priority |
|-------------|---------|--------|----------|
| **Bulk student import (35+ at once)** | ✅ CSV import | ✅ Canvas API integration | High |
| **Per-student budget isolation ($24 each)** | ❌ None | ✅ Budget tracking per user | Critical |
| **Template whitelisting** | ❌ None | ✅ Restrict bundles per project | Critical |
| **Auto-stop all instances (break/holiday)** | ⚠️ Manual | ✅ Project-wide stop | Critical |
| **Semester lifecycle (start/end dates)** | ❌ None | ✅ Auto-cleanup at semester end | Critical |
| **TA delegation (debug, reset)** | ❌ None | ✅ TA role with permissions | High |
| **Student isolation (no SSH key sharing)** | ⚠️ Assumed | ✅ Enforce isolation | High |
| **Cost reports for IT** | ⚠️ Basic | ✅ PDF reports with breakdowns | Medium |
| **Grade correlation data** | ❌ None | ⚠️ Nice-to-have | Low |

---

## 🎯 Priority Recommendations: Professor Persona

### Phase 1: Enable Basic Class Management (Milestone 1)

**Target**: Professors can safely run classes of 10-50 students

1. **Per-Student Budget Tracking** (2 weeks)
   - Store budget per student in IAM tags
   - Track spend per student via CloudWatch/Cost Explorer
   - Display per-student budget in `lfr project budget`

2. **Template/Bundle Restrictions** (1 week)
   - Add `--allowed-bundles` to project creation
   - Block launches of non-allowed bundles
   - Show educational error messages

3. **Project Bulk Operations** (1 week)
   - `lfr project stop --all`
   - `lfr project start --all`
   - `lfr project status` for entire class

4. **Semester End Automation** (1 week)
   - Add `--end-date` to project creation
   - Auto-stop all instances at end date
   - Auto-revoke IAM credentials
   - Generate final report

**Estimated effort**: 5 weeks

### Phase 2: TA Support (Milestone 2)

**Target**: TAs can efficiently help students

5. **TA Access Model** (2 weeks)
   - Add TA role to IAM policies
   - `lfr ta debug` command for SSH access
   - Audit logging for TA sessions
   - TA permissions per project

6. **Instance Reset** (3 days)
   - `lfr ta reset-instance` command
   - Backup to S3 before reset
   - Restore student files
   - Notify student of reset

7. **TA Dashboard** (1 week)
   - `lfr ta status` - view all students
   - Show which students need help
   - Budget warnings for TAs

**Estimated effort**: 4 weeks

### Phase 3: Advanced Class Features (Milestone 3)

**Target**: Full educational platform integration

8. **Canvas LMS Integration** (2 weeks)
   - Import students from Canvas API
   - Sync rosters automatically
   - Push grades back to Canvas (optional)

9. **Vacation Mode** (3 days)
   - `lfr course vacation --start --end`
   - Auto-stop all instances
   - Block starts during vacation
   - Student notifications

10. **Advanced Reporting** (1 week)
    - PDF reports for IT department
    - Cost breakdown per student
    - Usage analytics
    - Grade correlation (optional)

**Estimated effort**: 4 weeks

---

## Success Metrics: Professor Perspective

### Time Savings
- ✅ **Setup Time**: 8 hours → 15 minutes (97% reduction)
- ✅ **Weekly Management**: 2 hours → 10 minutes (92% reduction)
- ✅ **Semester End Cleanup**: 4-8 hours → 0 minutes (100% automated)
- **Total semester savings**: ~40 hours

### Budget Control
- ✅ **Budget Overruns**: 40% → 5% (better forecasting, auto-stop)
- ✅ **Cost Surprises**: Eliminated (real-time per-student tracking)
- ✅ **Average Cost per Student**: $28 → $22 (22% savings via auto-stop)

### Student Support
- ✅ **TA Efficiency**: Office hours 50% more productive
- ✅ **Student Satisfaction**: "Easier than my own AWS account" (92%)
- ✅ **Academic Integrity**: Full audit trail for investigations

### IT Department
- ✅ **Budget Predictability**: Classes stay within 95% of budget
- ✅ **Cost Reporting**: Automated PDF reports monthly
- ✅ **Security**: Student isolation enforced, audit logs available

---

## Next Steps

1. **Validate with Real Professors**: Interview 2-3 CS professors currently using AWS
2. **Prototype Budget Tracking**: Mock up per-student budget display
3. **Design TA Access Model**: IAM policies, audit logging architecture
4. **Implementation Plan**: Break into 2-week sprints

**Estimated Timeline**:
- Phase 1 (Basic Class Management): 5 weeks
- Phase 2 (TA Support): 4 weeks
- Phase 3 (Advanced Features): 4 weeks
- **Total**: ~13 weeks (3 months) to full professor support

---

**Status**: Draft Walkthrough
**Persona**: Professor (Course Management)
**Priority**: 🔴 Critical (blocks educational use case)
**Dependencies**: Core LFR functionality (exists), budget tracking (missing), TA roles (missing)

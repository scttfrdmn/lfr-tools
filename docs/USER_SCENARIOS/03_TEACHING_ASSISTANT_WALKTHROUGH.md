# Scenario 3: Teaching Assistant - Office Hours Support

## Persona: Alex Chen - Head Teaching Assistant

**Background**:
- PhD student, Computer Science (3rd year)
- Head TA for CS 305 - Cloud Computing (manages 2 other TAs)
- Technical level: Expert (builds distributed systems for research)
- TA experience: 4 semesters (various courses)
- Time commitment: 20 hours/week (10 hrs office hours, 10 hrs grading)

**Responsibilities**:
1. **Office Hours**: 6 hours/week helping students debug
2. **Student Support**: Primary contact for technical issues
3. **Environment Management**: Help students with broken setups
4. **Grading**: Access student work for assignment grading
5. **TA Coordination**: Manage 2 junior TAs (Sarah and Mike)

**Primary Concerns**:
1. **Efficiency**: 35 students, limited office hours - must help quickly
2. **Access**: Need to see student environments to debug issues
3. **Academic Integrity**: All help must be logged (no cheating accusations)
4. **Empowerment**: Junior TAs need tools to help without constant supervision
5. **Scale**: Can't spend 30 minutes per student - need 5-10 minute solutions

**Pain Points from Previous Semester**:
- Spent 45 minutes debugging via screen share (student's slow WiFi)
- Couldn't access student instance to see actual error
- Student accidentally deleted their environment, lost all work
- No way to prove TA didn't help student cheat
- Junior TAs constantly asked "How do I help with X?"

---

## Current State (v0.1.x): What Doesn't Work

### ❌ Problem: No TA Access at All

**Scenario**: Monday office hours, Emily needs help with broken Python environment

**Current reality**:
```bash
# Emily (student) in Zoom office hours:
# "Alex, my code crashes with some weird error. Can you help?"

# Alex (TA) tries to help:
alex@laptop:~$ lfr users list --project CS305-Fall2025

# Output:
# Error: You do not have permission for this operation.
# Required role: Project Administrator
# Your role: None (you are not a member of this project)

# Alex thinks: "Wait, I'm the TA but I have no access?!"

# Alex asks Emily to share screen (slow, painful debugging)
# Emily's WiFi keeps dropping
# 45 minutes later, Alex finds the issue: typo in import statement
# Could have been 2 minutes if Alex could SSH into Emily's instance

# Alex emails Prof. Rodriguez:
# "Can you add me to the CS305-Fall2025 project? I can't help students without access."

# Prof. Rodriguez replies:
# "I can make you an admin, but then you'd have full control (including billing).
#  I don't think that's appropriate for a TA. Not sure what to do here."

# Alex thinks: "There should be a TA role between 'no access' and 'full admin'!"
```

**Current workaround**: Screen sharing, talking students through fixes, very slow
**Impact**: Office hours inefficient, students frustrated, TAs frustrated

---

## 🎯 Ideal Future State: Complete TA Workflow

### Pre-Semester: TA Onboarding (Week -1)

```bash
# Prof. Rodriguez adds TAs to course project
lfr course add-ta --project CS305-Fall2025 \
  --email alex.chen@university.edu \
  --role head-ta \
  --permissions debug,reset,view

lfr course add-ta --project CS305-Fall2025 \
  --email sarah.kim@university.edu \
  --role ta \
  --permissions view,message

lfr course add-ta --project CS305-Fall2025 \
  --email mike.wong@university.edu \
  --role ta \
  --permissions view,message

# LFR Output:
# ✅ Added Teaching Assistants:
#
# Head TA: alex.chen@university.edu
#   Permissions:
#   - ✅ View all student instances
#   - ✅ Debug access (temporary SSH, logged)
#   - ✅ Reset student instances
#   - ✅ Send messages to students
#   - ✅ View budget status
#   - ❌ No billing access (professor only)
#   - ❌ No user creation (professor only)
#
# TAs: sarah.kim@university.edu, mike.wong@university.edu
#   Permissions:
#   - ✅ View all student instances
#   - ✅ Send messages to students
#   - ⚠️  No debug access (must escalate to Head TA)
#   - ⚠️  No reset access (must escalate to Head TA)
#
# 📧 Welcome emails sent to all TAs

# Alex receives email:
# Subject: 🎓 You're a TA for CS 305 - Your Access is Ready
#
# Hi Alex,
#
# Prof. Rodriguez has added you as Head TA for CS 305 - Cloud Computing.
#
# Your Role: Head Teaching Assistant
#   - Manage 35 student environments
#   - Debug student issues
#   - Reset broken environments
#   - Coordinate with TAs Sarah and Mike
#
# TA Tools:
#   Install: brew install lfr
#   View students: lfr ta students --project CS305-Fall2025
#   Help a student: lfr ta debug --student <name>
#   TA Dashboard: lfr ta dashboard
#
# Quick Start:
#   1. Install LFR Tools (if you haven't)
#   2. Configure AWS with your TA credentials
#   3. Run: lfr ta dashboard
#
# TA Training:
#   - Date: August 20, 2025 (Week before classes)
#   - Location: Zoom link or in-person
#   - Duration: 1 hour
#   - Topics: TA commands, debugging workflow, academic integrity
#
# Questions? Contact Prof. Rodriguez.
#
# Let's have a great semester!

# Alex installs and tests TA access
alex@laptop:~$ lfr ta dashboard --project CS305-Fall2025

# LFR Output:
# 🎓 CS 305 TA Dashboard (Head TA: Alex Chen)
#
# Course: CS 305 - Cloud Computing (Fall 2025)
# Professor: Dr. Rodriguez
# Enrollment: 35 students
# Other TAs: Sarah Kim, Mike Wong
#
# Week -1 (Pre-semester setup)
# Office Hours: Not started yet
#
# Student Status: (All instances created, none active yet)
#   35 students enrolled
#   0 instances running
#   0 students need help
#
# Quick Actions:
#   View students: lfr ta students
#   Test student access: lfr ta test-access
#   Office hours mode: lfr ta office-hours start
#
# 💡 Tip: Run "lfr ta help" for all TA commands
#
# Ready for semester! 🎉
```

### Week 2: Monday Office Hours - Student Needs Help

```bash
# Alex starts office hours mode
alex@laptop:~$ lfr ta office-hours start

# LFR Output:
# 🕐 Office Hours Mode: ACTIVE
#
# Time: Monday, Sep 9, 2025 - 2:00 PM
# Duration: 2 hours (until 4:00 PM)
# Location: Virtual (Zoom)
#
# TA Dashboard will auto-refresh every 30 seconds.
# Students can request help via: lfr student request-help
#
# Waiting for students...

# Meanwhile, Emily (student) has an issue
# emily@ubuntu:~$ python3 analysis.py
# ImportError: No module named 'pandas'

emily@laptop:~$ lfr student request-help "ImportError with pandas module"

# Alex's dashboard updates:
# 🔔 New Help Request!
#
# Student: Emily Chen (emily)
# Issue: "ImportError with pandas module"
# Time: 2:15 PM (just now)
# Instance: emily-ubuntu_22_04 (running)
# Assignment: Homework 2
#
# Actions:
#   [1] Message student
#   [2] Debug session (SSH access)
#   [3] View instance logs
#   [4] Assign to another TA
#   [5] Ignore
#
# Choice [1-5]: 2

# Alex selects debug session
alex@laptop:~$ lfr ta debug --student emily

# LFR Output:
# 🔍 TA Debug Session Starting...
#
# Student: Emily Chen (emily@university.edu)
# Instance: emily-ubuntu_22_04
# Your role: Head TA (Alex Chen)
# Project: CS305-Fall2025
#
# ⚠️  Academic Integrity Notice:
#    - This session will be fully logged
#    - All commands recorded
#    - Student will be notified of your access
#    - Session ID: debug-20250909-001
#    - Log: /var/log/ta-sessions/debug-20250909-001.log
#
# Session Details:
#   Duration: 30 minutes (auto-expires)
#   Access: Read-write (be careful!)
#   Student notification: Sent
#
# Available Commands (in session):
#   ta-view-only-mode    Toggle read-only mode
#   ta-exit              End session early
#   ta-annotate "msg"    Leave note for student
#
# Connecting to emily-ubuntu_22_04...
# ✅ Connected!

# Alex is now in Emily's environment (logged)
┌─────────────────────────────────────────────────────────┐
│ ⚠️  TA DEBUG SESSION ACTIVE                              │
│ Student: Emily Chen (emily@university.edu)              │
│ TA: Alex Chen (alex.chen@university.edu)                │
│ Instance: emily-ubuntu_22_04                            │
│ Session ID: debug-20250909-001                          │
│ Time: 2:16 PM - Expires 2:46 PM (30 min)                │
│ Mode: READ-WRITE ⚠️  (be careful!)                       │
│ All actions logged for academic integrity               │
└─────────────────────────────────────────────────────────┘

emily@ubuntu:~$ cd homework2
emily@ubuntu:~/homework2$ python3 analysis.py
# ImportError: No module named 'pandas'

# Alex checks if pandas is installed
emily@ubuntu:~/homework2$ pip3 list | grep pandas
# (no output - not installed)

# Alex realizes: Student forgot to install dependencies
emily@ubuntu:~/homework2$ ls
# analysis.py  data.csv  requirements.txt

emily@ubuntu:~/homework2$ cat requirements.txt
# pandas==1.5.0
# numpy==1.24.0
# matplotlib==3.6.0

# Alex sees the issue: requirements.txt exists but not installed
# Instead of fixing it (academic integrity), Alex documents

emily@ubuntu:~/homework2$ ta-annotate "Found the issue: pandas not installed. Check requirements.txt and install dependencies."

# LFR Output:
# ✅ Annotation saved for student
# Student will see this note next time they connect

# Alex exits debug session
emily@ubuntu:~/homework2$ ta-exit

# Back on Alex's laptop
alex@laptop:~$
# 🔚 Debug Session Ended
#
# Session Summary:
#   Duration: 3 minutes
#   Student: Emily Chen
#   Issue: Missing pandas installation
#   Resolution: Guided student (annotation left)
#   Academic integrity: ✅ No direct code changes
#
# 📧 Notifications sent:
#   - Emily: "TA Alex found your issue - see annotation"
#   - Prof. Rodriguez: Debug session summary (weekly digest)
#
# Log saved: /var/log/ta-sessions/debug-20250909-001.log

# Alex messages Emily to explain
alex@laptop:~$ lfr ta message emily "Found the issue! You need to install dependencies from requirements.txt. Run: pip3 install -r requirements.txt"

# LFR Output:
# 📧 Message sent to Emily Chen
#
# Emily will receive:
#   - In-app notification (next lfr command)
#   - Email (if not seen in 10 minutes)

# Emily sees the message
emily@laptop:~$ lfr ssh connect emily --project CS305-Fall2025

# LFR Output:
# ┌─────────────────────────────────────────────────┐
# │ 📨 Message from TA Alex Chen (3 minutes ago)     │
# │                                                 │
# │ Found the issue! You need to install            │
# │ dependencies from requirements.txt.             │
# │                                                 │
# │ Run: pip3 install -r requirements.txt           │
# └─────────────────────────────────────────────────┘

emily@ubuntu:~$ cd homework2
emily@ubuntu:~/homework2$ pip3 install -r requirements.txt
# Installing pandas, numpy, matplotlib...
# ✅ Done!

emily@ubuntu:~/homework2$ python3 analysis.py
# (Works now!)

emily@ubuntu:~/homework2$ exit

emily@laptop:~$ lfr student thank-ta alex

# LFR Output:
# 📧 Thank you sent to TA Alex Chen!
#
# Message: "Thanks for helping me with the pandas issue! 🙏"

# Alex receives:
# 📧 Emily Chen thanked you for helping with their issue!

# Alex thinks: "3 minutes to solve vs 45 minutes screen sharing. This is amazing!"
```

### Week 5: Student Broke Environment - Needs Reset

```bash
# David (student) emails: "Help! I deleted Python and my assignment is due tomorrow!"

# Alex checks David's status
alex@laptop:~$ lfr ta student-info david

# LFR Output:
# 👤 Student: David Kim (david@university.edu)
#
# Instance: david-ubuntu_22_04
# Status: Running (been running 6 hours)
# Budget: $16.50 / $24.00 (69%) ✅
#
# Recent Activity:
#   - 6 hours ago: Started instance
#   - No stop in last 6 hours (unusual for David)
#   - Average session: 2.5 hours
#
# Assignments:
#   - Homework 3: Due tomorrow (9/18)
#   - Homework 2: Submitted ✅ (graded: 95%)
#
# Previous Help Requests: 1 (Week 2, pandas issue, resolved)
#
# Quick Actions:
#   Debug: lfr ta debug --student david
#   Reset: lfr ta reset-instance david
#   Message: lfr ta message david "..."

# Alex debugs to see what happened
alex@laptop:~$ lfr ta debug --student david

# (In David's instance)
emily@ubuntu:~$ python3 --version
# bash: python3: command not found

# Alex confirms: Python is gone
# This is beyond simple fix - needs full reset

emily@ubuntu:~$ ta-exit

# Alex resets David's instance
alex@laptop:~$ lfr ta reset-instance david --reason "Python deleted by student"

# LFR Output:
# 🔄 Resetting Instance: david-ubuntu_22_04
#
# ⚠️  This will:
#   ✅ Backup current state to S3 (just in case)
#   ✅ Stop current instance
#   ✅ Launch fresh instance from course blueprint
#   ✅ Restore David's home directory (/home/david/)
#   ✅ Preserve all student work
#   ❌ Discard broken system state
#
# Student Notification:
#   David will receive email explaining the reset
#   His files will be intact
#
# Professor Notification:
#   Prof. Rodriguez will receive summary (weekly digest)
#
# Academic Integrity Log:
#   - TA: Alex Chen
#   - Student: David Kim
#   - Reason: Python deleted by student
#   - Files preserved: Yes
#   - Backup location: s3://cs305-fall2025/backups/david-20250917.tar.gz
#
# Estimated downtime: 5 minutes
# Cost: $0.00 (no additional charge)
#
# Proceed? [y/N]: y
#
# Resetting...
# ✅ Backup created: s3://cs305-fall2025/backups/david-20250917.tar.gz
# ✅ Stopped old instance
# ✅ Launched fresh instance
# ✅ Restored /home/david/
# ✅ Verified Python installation
# ✅ Reset complete!
#
# 📧 Notifications sent:
#   - David: "Your instance has been reset by TA Alex"
#   - Prof. Rodriguez: "TA Alex reset david's instance (weekly summary)"
#
# David can now continue working on his assignment.

# Alex messages David with guidance
alex@laptop:~$ lfr ta message david "Your instance has been reset. Python is back! Your homework files are safe in ~/homework3/. Please be careful with sudo commands. Let me know if you need help!"

# LFR Output:
# 📧 Message sent to David Kim

# David connects and sees his files are there
david@laptop:~$ lfr ssh connect david --project CS305-Fall2025
david@ubuntu:~$ cd homework3
david@ubuntu:~/homework3$ ls
# assignment.py  data/  README.md (all there!)

david@ubuntu:~/homework3$ python3 --version
# Python 3.10.12 (back!)

# David thinks: "Wow, that was fast! And my work is saved!"

# Alex thinks: "5 minutes to reset vs telling David to re-create everything. Perfect!"
```

### Week 8: Coordinating with Junior TAs

```bash
# Sarah (junior TA) has office hours, student needs help
# Sarah doesn't have reset permissions, must escalate

sarah@laptop:~$ lfr ta debug --student alice
# Error: You do not have 'debug' permission
# Your role: TA (view, message only)
# For debug access, ask Head TA: alex.chen@university.edu

# Sarah messages Alex
sarah@laptop:~$ lfr ta escalate --student alice --issue "Student broke environment, needs reset"

# LFR Output:
# 📧 Help request escalated to Head TA Alex Chen
#
# Student: Alice Johnson
# Issue: "Student broke environment, needs reset"
# TA: Sarah Kim
# Time: Wednesday, 3:15 PM
#
# Alex will be notified and can assist.

# Alex receives notification
alex@laptop:~$ lfr ta notifications

# LFR Output:
# 🔔 TA Notifications
#
# 1. [3:15 PM] Sarah Kim escalated: Alice needs environment reset
#    Issue: "Student broke environment, needs reset"
#    Quick action: lfr ta reset-instance alice
#
# 2. [Earlier] David thanked you for helping with pandas issue
#
# 3. [Yesterday] Prof. Rodriguez: "Great job on quick student support!"

# Alex handles the escalation
alex@laptop:~$ lfr ta reset-instance alice --reason "Escalated from TA Sarah - broken environment"

# (Reset happens as before)

# LFR sends notification to Sarah:
# 📧 Your escalation for Alice has been resolved by Alex Chen
#    Action taken: Instance reset
#    Student can continue working

# Sarah messages Alice
sarah@laptop:~$ lfr ta message alice "Alex reset your instance. You're good to go!"

# Sarah thinks: "Love that I can escalate quickly and students get help fast!"
```

### Week 15: End of Semester - TA Summary Report

```bash
# Alex checks semester summary
alex@laptop:~$ lfr ta report --project CS305-Fall2025

# LFR Output:
# 📊 TA Report: CS 305 Fall 2025 (Alex Chen - Head TA)
#
# Semester: August 25 - December 12, 2025 (15 weeks)
# Your Role: Head Teaching Assistant
#
# Office Hours Efficiency:
#   Total office hours: 90 hours (6 hrs/week × 15 weeks)
#   Students helped: 47 unique help requests
#   Average resolution time: 8 minutes ⚡
#   Students helped per hour: 3.1 (vs ~1.3 without TA tools)
#   Time saved: ~60 hours (vs screen sharing method)
#
# Debug Sessions:
#   Total sessions: 32
#   Average duration: 6 minutes
#   Shortest: 2 minutes (quick fix)
#   Longest: 18 minutes (complex debugging)
#   All sessions logged: ✅ (academic integrity)
#
# Instance Resets:
#   Total resets: 8 students
#   Reasons:
#     - Deleted system packages: 5
#     - Corrupted environment: 2
#     - Requested fresh start: 1
#   Average downtime: 5 minutes
#   All files preserved: ✅
#
# Messages Sent: 87
#   - Help responses: 47
#   - Proactive guidance: 28
#   - Assignment clarifications: 12
#
# Student Feedback:
#   "TA was very helpful": 95% (45/47 students)
#   Average help rating: 4.8/5.0
#   "TA response was fast": 98%
#
# Junior TA Performance:
#   Sarah Kim:
#     - Students helped: 18
#     - Escalations: 5 (appropriate)
#     - Messages: 32
#   Mike Wong:
#     - Students helped: 12
#     - Escalations: 3 (appropriate)
#     - Messages: 21
#
# Academic Integrity:
#   - All debug sessions logged: ✅ 32/32
#   - No academic integrity incidents involving TAs
#   - No student complaints about TA access
#   - All TA actions tracked and auditable
#
# Top Student Issues:
#   1. Missing dependencies (pip install): 15 cases
#   2. Environment corruption: 8 cases
#   3. File permissions: 7 cases
#   4. Git conflicts: 6 cases
#   5. Budget confusion: 4 cases
#
# Time Investment:
#   Office hours: 90 hours
#   Grading: 60 hours
#   TA coordination: 10 hours
#   Total: 160 hours (vs estimated 180 without TA tools)
#
# Cost to Course:
#   TA access: $0.00 (no additional AWS cost)
#   Instance resets: $0.00 (no additional cost)
#   Debug sessions: $0.00 (use student instances)
#
# Recommendations for Next Semester:
#   1. Create "common issues" guide (top 5 issues repeat)
#   2. Proactive message at Week 3: "Remember pip install -r requirements.txt"
#   3. Add permissions workshop in Week 2
#   4. Consider +1 junior TA (class growing to 50 students)
#
# Professor Feedback:
#   "Alex and TAs did an excellent job this semester. Student satisfaction
#    with support was the highest I've ever seen. The TA tools made a huge
#    difference in efficiency and academic integrity tracking."
#    - Prof. Rodriguez
#
# Export this report: lfr ta report --export pdf

# Alex exports for his CV/portfolio
alex@laptop:~$ lfr ta report --export pdf --output ~/Desktop/cs305-ta-report.pdf

# LFR Output:
# ✅ Report exported: ~/Desktop/cs305-ta-report.pdf
#    Ready for your CV/portfolio!
#
# Thank you for being a great TA! 🎉
```

---

## 📋 Feature Gap Analysis: TA Needs

### Critical Missing Features (TA Persona Completely Missing)

| Feature | Priority | Impact | Current State | Effort |
|---------|----------|--------|---------------|--------|
| **TA Role & Permissions** | 🔴 Critical | No TA access at all | ❌ None | High |
| **Debug Session (SSH Access)** | 🔴 Critical | Can't help students remotely | ❌ None | High |
| **Instance Reset** | 🔴 Critical | Broken envs = hours to fix | ❌ None | Medium |
| **TA Dashboard** | 🔴 Critical | No visibility into students | ❌ None | Medium |
| **Student Messaging** | 🟡 High | No communication channel | ❌ None | Low |
| **Office Hours Mode** | 🟡 High | Manual coordination | ❌ None | Low |
| **Academic Integrity Logging** | 🟡 High | No proof of help | ❌ None | Medium |
| **TA Escalation** | 🟡 High | Junior TAs can't escalate | ❌ None | Low |
| **TA Reports** | 🟢 Medium | No performance tracking | ❌ None | Low |
| **Help Request System** | 🟢 Medium | Students email/zoom only | ❌ None | Medium |

### TA Requirements (All Missing)

| Requirement | Needed | Priority |
|-------------|--------|----------|
| **TA IAM role (between student and admin)** | ✅ Role with scoped permissions | Critical |
| **Temporary SSH to student instances** | ✅ Time-limited, logged access | Critical |
| **View all student instances** | ✅ Dashboard with status | Critical |
| **Reset student instance (preserve files)** | ✅ One-command reset | Critical |
| **Send messages to students** | ✅ In-app + email notifications | High |
| **Academic integrity audit logs** | ✅ All TA actions logged | High |
| **Escalation workflow (junior → head TA)** | ✅ Junior TAs can request help | High |
| **Office hours coordination** | ✅ Help queue, active mode | Medium |
| **TA performance reports** | ✅ Semester summary | Medium |

---

## 🎯 Priority Recommendations: TA Persona

### Phase 1: Enable Basic TA Support (Milestone 1)

**Target**: TAs can help students remotely with logged sessions

1. **TA IAM Role** (2 weeks)
   - Create TA role in IAM policies
   - Permissions: view instances, read-only by default
   - No billing access, no user creation
   - Scope to project only

2. **Debug Session (Logged SSH)** (2 weeks)
   - `lfr ta debug --student <name>` command
   - Time-limited sessions (30 min default)
   - All commands logged to S3
   - Student notification when TA connects
   - Read-write access with logging

3. **Instance Reset** (1 week)
   - `lfr ta reset-instance <student>` command
   - Backup to S3 before reset
   - Restore student home directory
   - Fresh instance from blueprint
   - Notify student and professor

4. **TA Dashboard** (1 week)
   - `lfr ta dashboard` command
   - List all students
   - Instance status
   - Budget warnings
   - Recent issues

**Estimated effort**: 6 weeks

### Phase 2: TA Coordination (Milestone 2)

**Target**: Multiple TAs can coordinate efficiently

5. **TA Permission Levels** (1 week)
   - Head TA: full debug, reset
   - Junior TA: view, message, escalate
   - Configurable per TA

6. **Student Messaging** (3 days)
   - `lfr ta message <student> "..."` command
   - In-app notifications
   - Email fallback
   - Message history

7. **TA Escalation** (3 days)
   - `lfr ta escalate --student <name>` command
   - Notify head TA
   - Track escalations
   - Resolution notifications

**Estimated effort**: 2 weeks

### Phase 3: Advanced TA Tools (Milestone 3)

**Target**: Professional TA experience with analytics

8. **Office Hours Mode** (1 week)
   - `lfr ta office-hours start/stop`
   - Help request queue
   - Auto-refresh dashboard
   - Student request system

9. **TA Reports** (1 week)
   - `lfr ta report` command
   - Semester summary
   - Performance metrics
   - Export to PDF

10. **Academic Integrity Audit** (1 week)
    - All TA sessions logged
    - Searchable audit trail
    - Export for investigations
    - Retention policy (1 year)

**Estimated effort**: 3 weeks

---

## Success Metrics: TA Perspective

### Efficiency
- ✅ **Average help time**: 45 min → 8 min (82% reduction)
- ✅ **Students helped per hour**: 1.3 → 3.1 (138% increase)
- ✅ **Office hours utilization**: 60% → 95%
- ✅ **Time saved per semester**: ~60 hours (TA can help more students)

### Student Satisfaction
- ✅ **"TA was helpful"**: 95% (up from ~70% with screen sharing)
- ✅ **Help response time**: 30 min → 5 min
- ✅ **Issue resolution rate**: 85% → 98%

### Academic Integrity
- ✅ **TA actions logged**: 0% → 100%
- ✅ **Academic integrity incidents involving TAs**: 0
- ✅ **Audit trail available**: 100% of sessions

### TA Experience
- ✅ **TA satisfaction**: "Tools made TA job much easier" (100%)
- ✅ **Junior TA empowerment**: Can help without constant supervision
- ✅ **TA retention**: More TAs want to TA again (anecdotal)

---

## Next Steps

1. **Interview TAs**: Talk to 2-3 current TAs about pain points
2. **Design TA IAM Role**: Security model, permissions scoping
3. **Prototype Debug Session**: Mock up logged SSH access
4. **Academic Integrity Review**: Consult with university policies

**Estimated Timeline**:
- Phase 1 (Basic TA Support): 6 weeks
- Phase 2 (TA Coordination): 2 weeks
- Phase 3 (Advanced Tools): 3 weeks
- **Total**: ~11 weeks (3 months) to comprehensive TA support

**Dependencies**:
- TA role creation (new IAM policies)
- Audit logging infrastructure (S3, retention)
- Professor project setup (must add TAs to projects)
- Student awareness (notifications when TA accesses)

---

**Status**: Draft Walkthrough
**Persona**: Teaching Assistant (Office Hours Support)
**Priority**: 🔴 Critical (enables scalable student support)
**Note**: This persona is currently **completely missing** from LFR Tools - highest priority after Professor support

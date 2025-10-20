# Scenario 2: Student First-Time Cloud User

## Persona: Emily Chen - Undergraduate Student

**Background**:
- Junior, Computer Science major
- First cloud computing course (CS 305 - Cloud Computing)
- Technical level: Comfortable with Python, Git, command line basics
- Never used AWS or any cloud service before
- Lives off-campus, uses personal laptop (MacBook Air M1)
- Part-time job (15 hours/week) - budget conscious

**Primary Concerns**:
1. **Cost Anxiety**: "Will this cost me money? I can't afford surprise bills!"
2. **Complexity**: "AWS seems really complicated, what if I break something?"
3. **Access**: "Can I work from home? Coffee shops? My parents' house during breaks?"
4. **Grades**: "If my environment breaks, will it affect my grade?"
5. **Time**: "I have 4 other classes, I can't spend hours on infrastructure"

**Previous Experience**:
- Used university computer labs (free, simple, "just works")
- Ran Python locally on laptop (small datasets only)
- Heard horror stories: "My friend got a $500 AWS bill in one month!"

---

## Current State (v0.1.x): What Works Today

### ✅ Day 1: Receiving Credentials from Professor

```bash
# Emily receives email from Prof. Rodriguez:
#
# Subject: CS 305 - Your Cloud Computing Environment
#
# Hi Emily,
#
# Welcome to CS 305! You have been assigned a personal cloud environment
# for this course.
#
# Your credentials:
# - Username: emily
# - Password: xK9mP2yT3zL (change this first time)
# - Instance: emily-ubuntu_22_04
# - Region: us-east-1
#
# Budget: $24.00 for entire semester (covered by course fee)
# ⚠️  This is a FIXED budget - no additional charges to you!
#
# Getting Started:
# 1. Install LFR Tools: https://github.com/scttfrdmn/lfr-tools#installation
# 2. Follow the student guide: https://github.com/scttfrdmn/lfr-tools/docs/tutorials/student-guide.md
# 3. First assignment: Due September 5
#
# Need help? Contact TAs during office hours (schedule on Canvas)
#
# Prof. Rodriguez

# Emily's first thought: "Okay, $24 for the semester, that's fine.
# But what if I mess up and go over? I can't afford that!"
```

### ✅ Day 1: Installation and First Login

```bash
# Emily installs LFR Tools
# (Using Homebrew on macOS)
emily@laptop:~$ brew install lfr

# Output:
# ==> Downloading lfr-tools...
# ==> Installing lfr-tools
# ✅ lfr-tools installed successfully
#
# Get started:
# $ lfr --help

# Emily checks if it's installed
emily@laptop:~$ lfr --help

# Output:
# LFR Tools - Lightsail for Research Management
#
# Usage:
#   lfr [command]
#
# Available Commands:
#   users       Manage IAM users
#   groups      Manage IAM groups
#   instances   Manage LFR instances
#   ssh         SSH connection management
#   dcv         NICE DCV remote desktop
#   help        Help about any command
#
# Use "lfr [command] --help" for more information

# Emily thinks: "Okay, lots of commands. I just need to connect to my instance.
# Let me check the SSH command..."

# Configure AWS credentials (Emily uses professor-provided credentials)
emily@laptop:~$ aws configure --profile cs305
# AWS Access Key ID: [from professor's Canvas page]
# AWS Secret Access Key: [from professor's Canvas page]
# Default region: us-east-1
# Output format: json

export AWS_PROFILE=cs305

# Connect to instance
emily@laptop:~$ lfr ssh connect emily --project CS305-Fall2025

# Output:
# 🔍 Finding instance: emily-ubuntu_22_04
# ⚠️  Instance is stopped (not running)
# 💰 Cost: $0.00/day (no charges while stopped!)
#
# Start instance? [y/N]: y
#
# ✅ Starting instance...
# ⏳ Waiting for instance to be ready (usually 30-60 seconds)...
# ✅ Instance ready!
# 🔑 Downloading SSH keys...
# 🔗 Connecting...

# Emily is now connected!
# Welcome to Ubuntu 22.04 LTS
#
# Instance: emily-ubuntu_22_04
# Course: CS 305 - Cloud Computing
# Budget: $24.00 for semester
#
# 💡 Tip: Your instance will cost $0.83/day if left running 24/7.
#         Stop it when you're done to save your budget!
#
# To stop this instance: lfr instances stop emily-ubuntu_22_04
# To disconnect: exit or Ctrl+D

emily@ubuntu:~$ ls
# (Empty home directory, ready for her work)

emily@ubuntu:~$ python3 --version
# Python 3.10.12

emily@ubuntu:~$ git --version
# git version 2.34.1

# Emily thinks: "Okay, this is actually pretty easy! And it's just Linux."
```

**What works well:**
- ✅ Simple installation (Homebrew)
- ✅ Automatic instance start when connecting
- ✅ Clear cost information ($0.83/day)
- ✅ Standard Ubuntu environment (familiar)
- ✅ Pre-installed tools (Python, Git)

### ✅ Week 1: Working on First Assignment

```bash
# Emily disconnects and reconnects the next day
emily@laptop:~$ lfr ssh connect emily --project CS305-Fall2025

# Output:
# 🔍 Finding instance: emily-ubuntu_22_04
# ⚠️  Instance is stopped
# 💰 Yesterday's cost: $0.83 (you stopped it - good job!)
# 💰 Budget remaining: $23.17 / $24.00 (96%)
#
# Start instance? [y/N]: y
#
# ✅ Starting...
# 🔗 Connected!

# Emily works on her assignment for 3 hours
emily@ubuntu:~$ mkdir hw1
emily@ubuntu:~$ cd hw1
emily@ubuntu:~$ git clone https://github.com/cs305/homework1.git
emily@ubuntu:~$ cd homework1
emily@ubuntu:~$ python3 analysis.py
# (Works on assignment)

# Emily finishes and disconnects
emily@ubuntu:~$ exit

# Back on her laptop
emily@laptop:~$ lfr instances list --project CS305-Fall2025

# Output:
# Instances for project CS305-Fall2025:
# ┌──────────────────────┬──────────┬──────────┬──────────┐
# │ Name                 │ Status   │ Uptime   │ Cost/day │
# ├──────────────────────┼──────────┼──────────┼──────────┤
# │ emily-ubuntu_22_04   │ Running  │ 3h 15m   │ $0.83    │
# └──────────────────────┴──────────┴──────────┴──────────┘
#
# 💡 Your instance is still running!
#    Cost: $0.11 so far today
#    Stop to save budget: lfr instances stop emily-ubuntu_22_04

# Emily stops her instance
emily@laptop:~$ lfr instances stop emily-ubuntu_22_04

# Output:
# 🛑 Stopping instance: emily-ubuntu_22_04
# ✅ Instance stopped
# 💰 Today's cost: $0.11 (3h 15m of usage)
# 💰 Budget remaining: $23.06 / $24.00 (96%)
#
# 💡 Your work is saved! Next time you start, everything will be there.

# Emily thinks: "Okay, I need to remember to stop it. Only cost me $0.11 for 3 hours!"
```

**What works well:**
- ✅ Clear budget feedback when connecting
- ✅ Warning when instance is left running
- ✅ Simple stop command
- ✅ Cost tracking per session

---

## ⚠️ Current Pain Points: What Doesn't Work

### ❌ Problem 1: No Budget Impact Preview Before Starting

**Scenario**: Emily wants to check how much budget she has before starting instance

**What should happen** (MISSING):
```bash
# Emily checks her budget status
lfr student budget

# LFR Output (MISSING):
# 📊 Your CS 305 Budget
#
# Course: CS 305 - Cloud Computing (Fall 2025)
# Professor: Dr. Rodriguez
#
# Budget:
#   Total for semester: $24.00
#   Spent so far: $4.50 (18.8%) - Week 3 of 15
#   Remaining: $19.50
#   Days left: 84 days (12 weeks)
#
# Your Instance: emily-ubuntu_22_04
#   Status: Stopped (safe - no cost!)
#   Type: nano_2_0 ($0.83/day if running 24/7)
#   Average session: 3.2 hours
#   Effective cost: ~$0.11/session
#
# 💡 Budget Forecast:
#    At your current usage rate (~3 sessions/week):
#    - Projected semester total: $21.60 ✅ (within budget!)
#    - Buffer remaining: $2.40
#
# What-if Scenarios:
#   If you work 4 hours/day, 5 days/week:
#     → Projected: $28.00 ⚠️  ($4 over budget)
#
#   If you always stop after sessions (like now):
#     → Projected: $21.60 ✅ (within budget!)
#
# 💡 Tips to stay within budget:
#   1. Always stop instance when done (lfr instances stop emily-ubuntu_22_04)
#   2. Your instance auto-stops after 4 hours idle (professor's policy)
#   3. Don't leave running overnight or over weekends
#
# Next assignment due: September 12 (4 days)

# Now Emily tries to start her instance
lfr ssh connect emily --project CS305-Fall2025

# LFR Output (BETTER - MISSING):
# 🔍 Finding instance: emily-ubuntu_22_04
# ⚠️  Instance is stopped
#
# 💰 Budget Impact:
#    Cost to start: ~$0.11 (estimated 3-hour session, your average)
#    Budget remaining: $19.50 → $19.39 after session
#    Still safe: ✅ You'll have $19.39 left (81% of budget)
#
# Start instance? [y/N]: y
```

**Current state**: Emily has to guess or calculate manually
**Impact**: Constant budget anxiety, doesn't know if it's safe to start

### ❌ Problem 2: No Automatic Stop (Forgotten Instances)

**Scenario**: Emily forgets to stop instance, goes to bed, wakes up to cost surprise

**What should happen** (PARTIALLY WORKS):
```bash
# Emily starts instance, works for 2 hours, gets distracted, closes laptop
# (Instance is still running!)

# 4 hours later (auto-stop should trigger - PROFESSOR CONFIGURED THIS)
# System log:
# 2025-09-08 22:30:00 [emily-ubuntu_22_04] Idle detected (4 hours)
# 2025-09-08 22:30:00 [emily-ubuntu_22_04] Auto-stopping per course policy
# 2025-09-08 22:30:05 ✅ Instance stopped

# Next morning, Emily receives email (MISSING):
# Subject: 💰 CS 305 - Your instance auto-stopped (budget protection)
#
# Hi Emily,
#
# Your instance (emily-ubuntu_22_04) was automatically stopped last night
# after 4 hours of inactivity.
#
# This is normal! The course has an auto-stop policy to protect your budget.
#
# Yesterday's cost: $0.83 (ran from 6:30 PM to 10:30 PM)
# Budget remaining: $18.67 / $24.00 (78%)
#
# 💡 Reminder: Always stop your instance when you're done working:
#    $ lfr instances stop emily-ubuntu_22_04
#
# If you need help, contact TAs during office hours.
#
# Best,
# LFR Tools (automated)

# Emily checks status next day
lfr student budget

# Output (MISSING):
# 📊 Your CS 305 Budget
#
# ⚠️  Notice: Your instance auto-stopped yesterday (9/8 at 10:30 PM)
#     Reason: 4-hour idle timeout (course policy)
#     Cost: $0.83 (you forgot to stop manually)
#     💡 Remember to stop when done to save budget!
#
# Budget remaining: $18.67 / $24.00 (78%)
# On track: ✅ (projected: $22.10 for semester)
```

**Current state**: Auto-stop works (professor configured) but no student notification
**Impact**: Emily doesn't learn from the incident, may repeat

### ❌ Problem 3: No "I Broke My Environment" Recovery

**Scenario**: Emily accidentally deletes Python or breaks her system

**What should happen** (MISSING):
```bash
# Emily accidentally runs: sudo apt remove python3 (oops!)
emily@ubuntu:~$ python3 --version
# bash: python3: command not found

# Emily panics: "OH NO! My assignment is due tomorrow!"
# She exits SSH in panic

emily@laptop:~$ lfr student help

# LFR Output (MISSING):
# 🆘 LFR Student Help
#
# Common Issues:
#
# 1. "I broke my environment" (Python missing, can't run code, etc.)
#    → Solution: Ask TA to reset your instance
#    → Your files will be preserved!
#    → Email: ta-alex@university.edu or ta-sarah@university.edu
#    → Office hours: Mon/Wed 2-4 PM, Fri 1-3 PM
#
# 2. "I'm running out of budget"
#    → Check: lfr student budget
#    → Make sure you stop instances: lfr instances stop <name>
#    → Contact professor if you need budget extension
#
# 3. "I can't connect to my instance"
#    → Try: lfr ssh connect emily --project CS305-Fall2025
#    → If still fails, email TA or professor
#
# 4. "I forgot my password"
#    → Contact professor: maria.rodriguez@university.edu
#
# 5. "My assignment is due and I can't work!"
#    → Email TA immediately with error details
#    → Take screenshot of error
#    → TAs can reset your instance in < 5 minutes

# Emily emails TA: "Help! I deleted Python and my assignment is due tomorrow!"

# TA perspective (see TA walkthrough):
# lfr ta reset-instance emily-ubuntu_22_04 --project CS305-Fall2025

# 5 minutes later, Emily receives email:
# Subject: ✅ Your CS 305 instance has been reset
#
# Hi Emily,
#
# TA Alex has reset your instance to fix the Python issue.
#
# What was done:
# - Your instance was backed up (in case we need it)
# - Fresh Ubuntu environment created
# - Your files restored from /home/emily/
# - Ready to use now!
#
# What you need to do:
# 1. Connect: lfr ssh connect emily --project CS305-Fall2025
# 2. Verify your homework files are there: cd ~/hw2
# 3. Continue working!
#
# Cost: $0.00 (no charge for reset)
# Downtime: 5 minutes
#
# 💡 Tip: Be careful with sudo commands! If unsure, ask first.
#
# You can now continue working on your assignment. Good luck!
#
# Best,
# TA Alex

# Emily reconnects and sees her files are there
emily@laptop:~$ lfr ssh connect emily --project CS305-Fall2025
emily@ubuntu:~$ cd hw2
emily@ubuntu:~/hw2$ ls
# assignment.py  data.csv  README.md (all there!)

emily@ubuntu:~/hw2$ python3 --version
# Python 3.10.12 (back!)

# Emily thinks: "Wow, that was fast! I love having TAs who can actually help!"
```

**Current state**: No self-service help, no easy TA reset
**Impact**: Students panic, potentially lose work, miss deadlines

### ❌ Problem 4: No Clear Cost Explanation (Budget Anxiety Persists)

**Scenario**: Emily wants to understand costs but AWS pricing is confusing

**What should happen** (MISSING):
```bash
# Emily wants to understand her costs better
lfr student costs explain

# LFR Output (MISSING):
# 💰 Understanding Your CS 305 Costs (In Plain English!)
#
# Your Budget: $24.00 for entire semester (15 weeks)
#   → Paid by course fee, NOT charged to you personally
#   → You will NOT receive a bill from AWS
#   → This is a FIXED budget - no surprises!
#
# Your Instance: emily-ubuntu_22_04
#   Type: nano_2_0 (small, perfect for coursework)
#   Cost: $0.83 per day IF running 24 hours straight
#
# But you DON'T run 24 hours straight! Here's what you actually pay:
#
# Real Cost Breakdown:
#   Running: $0.83 / 24 hours = $0.035 per hour
#   Your average session: 3 hours
#   Your cost per session: ~$0.11
#
# Example Scenarios:
#
# ✅ Good usage (staying within budget):
#    - Work 3 hours, 3 times/week: ~$1.00/week
#    - Stop instance after each session
#    - Total semester: ~$15 (well within $24!)
#    - Buffer: $9 for longer sessions or extra work
#
# ⚠️  Budget risk (going over):
#    - Forget to stop overnight: $0.83 wasted
#    - Forget for weekend (48 hours): $1.66 wasted
#    - Do this 4 times: $6.64 wasted (25% of budget!)
#
# 🔴 Budget disaster (definitely going over):
#    - Never stop instance: $0.83/day × 105 days = $87.15
#    - Result: Blow through $24 in 29 days (Week 4!)
#    - Professor will be sad 😞
#
# 💡 Golden Rule: ALWAYS STOP WHEN DONE!
#    Command: lfr instances stop emily-ubuntu_22_04
#    Or just: exit (then run stop command)
#
# Current Semester Stats (Week 3):
#   Your average cost/week: $0.95 ✅
#   Projected semester total: $14.25 ✅
#   Buffer remaining: $9.75 (plenty!)
#
# You're doing great! Keep stopping your instance and you'll be fine. 🎉

# Emily can also see a simple visualization
lfr student costs graph

# LFR Output (MISSING):
# 📊 Your CS 305 Spending Over Time
#
# Week 1: ████░░░░░░ $2.10  (lots of setup work)
# Week 2: ██░░░░░░░░ $0.85  (short sessions)
# Week 3: ███░░░░░░░ $1.50  (current week, assignment due)
#
# Total: $4.45 / $24.00 (18.5%)
#
# Projection: ████████████████░░░░░ $21.30 (safe!)
#
# 💡 At your current rate, you'll finish with $2.70 to spare!
```

**Current state**: No cost education, AWS billing dashboard is overwhelming
**Impact**: Persistent anxiety, students avoid using instance (hurts learning)

### ❌ Problem 5: No Mobile/Simple Access for Quick Checks

**Scenario**: Emily is at coffee shop, wants to quickly check if instance is stopped

**What should happen** (MISSING):
```bash
# Emily is away from laptop, using phone browser
# Opens: https://lfr-tools-web.university.edu (MISSING)

# Web dashboard (MISSING):
# 🎓 CS 305 - Cloud Computing
# Student: Emily Chen
#
# Your Instance: emily-ubuntu_22_04
# Status: 🛑 Stopped (no cost)
# Budget: $19.50 / $24.00 (81%) ✅
#
# Quick Actions:
# [Start Instance] [View Budget] [Get Help]
#
# Recent Activity:
# - 9/8 10:30 PM: Auto-stopped (4hr idle)
# - 9/8 6:30 PM: Started (3h 15min session)
# - 9/7 2:15 PM: Stopped manually
#
# Next Assignment Due: 9/12 (in 4 days)

# Emily sees it's stopped, relaxes, continues with her day
```

**Current state**: CLI only, no web dashboard, no mobile access
**Impact**: Students can't quickly verify status, anxiety continues

---

## 🎯 Ideal Future State: Student-Friendly Experience

### Day 1: Welcoming Onboarding

```bash
# Emily receives personalized welcome email:
#
# Subject: 🎓 Welcome to CS 305! Your Cloud Computer is Ready
#
# Hi Emily,
#
# Congratulations on enrolling in CS 305 - Cloud Computing!
#
# You have been assigned a personal cloud computer for this course.
# Don't worry - this is WAY easier than it sounds! 😊
#
# 3 Things You Need to Know:
#
# 1. 💰 Cost: $0 to you! Covered by course fee.
#    - Budget: $24 for entire semester
#    - You will NEVER get a personal bill
#    - No credit card needed
#    - We'll help you stay within budget!
#
# 2. 🖥️ Your Computer: A personal Linux server in the cloud
#    - Name: emily-ubuntu_22_04
#    - Always available (from anywhere!)
#    - More powerful than your laptop
#    - Your files are safe (backed up)
#
# 3. 📚 Getting Started (15 minutes):
#    - Install: brew install lfr (Mac) or download for Windows
#    - Connect: lfr student join CS305-Fall2025
#    - Follow interactive tutorial
#    - Done!
#
# Need Help?
#   - Student guide: [link to friendly guide]
#   - Video tutorial: [link to 5-minute video]
#   - TAs: Office hours Mon/Wed/Fri
#   - Professor: maria.rodriguez@university.edu
#
# First Assignment: Due September 5 (on Canvas)
#
# We're here to help! Don't hesitate to ask questions. 🙂
#
# Prof. Rodriguez

# Emily installs and runs onboarding
emily@laptop:~$ lfr student join CS305-Fall2025

# Interactive onboarding (MISSING):
# 🎓 Welcome to CS 305, Emily!
#
# Let's get you set up. This will take about 5 minutes.
#
# Step 1/5: AWS Configuration
#   I'll configure your AWS access automatically.
#   Your credentials: [automatically applied from course]
#   ✅ Configured!
#
# Step 2/5: Your Instance
#   Your cloud computer: emily-ubuntu_22_04
#   Type: Small Linux server (perfect for coursework)
#   Cost: $24 for entire semester (already paid!)
#   ✅ Found!
#
# Step 3/5: Budget Protection
#   We'll help you stay within budget:
#   - Show budget every time you connect
#   - Remind you to stop when done
#   - Auto-stop after 4 hours idle
#   - Email alerts if running too long
#   ✅ Configured!
#
# Step 4/5: Test Connection
#   Let's make sure everything works...
#   Starting your instance...
#   Connecting...
#   ✅ Connected!
#
# You're now logged into your cloud computer!
# Let's do a quick tutorial...

# (Inside instance - interactive tutorial)
# 🎓 Quick Tutorial (2 minutes)
#
# This is a Linux command line. Don't worry if you're new to this!
#
# Try these commands:
#
# 1. See where you are:
#    $ pwd
#
emily@ubuntu:~$ pwd
# /home/emily
#
# Great! This is your home directory. All your files go here.
#
# 2. Create a file:
#    $ echo "Hello CS 305!" > test.txt
#
emily@ubuntu:~$ echo "Hello CS 305!" > test.txt
#
# 3. View the file:
#    $ cat test.txt
#
emily@ubuntu:~$ cat test.txt
# Hello CS 305!
#
# Perfect! You're a natural. 🎉
#
# 4. Check Python (you'll use this for assignments):
#    $ python3 --version
#
emily@ubuntu:~$ python3 --version
# Python 3.10.12
#
# Excellent! Python is ready for your assignments.
#
# Step 5/5: Disconnecting
#   When you're done working, ALWAYS stop your instance to save budget!
#
#   Two ways:
#   1. Disconnect: exit (or Ctrl+D)
#   2. Stop from laptop: lfr instances stop emily-ubuntu_22_04
#
#   Let's practice. Type: exit

emily@ubuntu:~$ exit

# Back on laptop
emily@laptop:~$
# 🎉 Tutorial Complete!
#
# You're all set for CS 305!
#
# Quick Reference:
#   Connect: lfr ssh connect emily --project CS305-Fall2025
#   Stop: lfr instances stop emily-ubuntu_22_04
#   Check budget: lfr student budget
#   Get help: lfr student help
#
# 💡 Pro Tip: Your instance is still running! Stop it now:
#    $ lfr instances stop emily-ubuntu_22_04
#
# Need help? Run: lfr student help
#
# Good luck with your first assignment! 🚀

emily@laptop:~$ lfr instances stop emily-ubuntu_22_04
# ✅ Instance stopped!
# Today's cost: $0.06 (10 minutes for tutorial)
# Budget remaining: $23.94 / $24.00 (99.8%)
```

### Week 3: Confident Student Using Budget Tools

```bash
# Emily before starting work
emily@laptop:~$ lfr student status

# Output (MISSING):
# 📊 Emily's CS 305 Dashboard
#
# Instance: emily-ubuntu_22_04
# Status: 🛑 Stopped (no cost)
# Budget: $18.50 / $24.00 (77%) ✅
#
# Next Assignment: Homework 3 - Due 9/15 (in 3 days)
#
# Quick Actions:
#   Start working: lfr ssh connect emily --project CS305-Fall2025
#   Check costs: lfr student costs
#   Get help: lfr student help
#
# Recent Activity:
#   ✅ 9/11: Worked 2.5 hours, stopped manually ($0.09)
#   ✅ 9/10: Worked 4 hours, auto-stopped ($0.14)
#   ⚠️  9/9: Forgot to stop, auto-stopped after 4h idle ($0.83)
#
# 💡 Tip: You're on track! At current rate, you'll finish semester at $21.80 (within budget).

# Emily starts working
emily@laptop:~$ lfr ssh connect emily --project CS305-Fall2025

# Better output (MISSING):
# 🔍 Starting your CS 305 instance...
#
# 💰 Budget Check:
#    Current: $18.50 / $24.00 (77%)
#    This session: ~$0.11 (estimated)
#    After session: $18.61 / $24.00 ✅
#
# ⏳ Starting... (usually 30 seconds)
# ✅ Ready!
# 🔗 Connecting...

# (Works for 3 hours)

# Emily finishes work
emily@ubuntu:~$ exit

# Back on laptop - friendly reminder (MISSING):
emily@laptop:~$
# 👋 Session ended!
#
# ⚠️  Don't forget: Your instance is still running!
#
# Stop now to save budget:
#   $ lfr instances stop emily-ubuntu_22_04
#
# Or it will auto-stop in 4 hours if idle.
# (Auto-stop saves budget but uses more $ than manual stop)

emily@laptop:~$ lfr instances stop emily-ubuntu_22_04

# Better output (MISSING):
# 🛑 Stopping emily-ubuntu_22_04...
# ✅ Stopped!
#
# Today's session:
#   Duration: 3h 12min
#   Cost: $0.11
#
# Budget update:
#   Spent today: $0.11
#   Spent this week: $0.34
#   Semester total: $18.61 / $24.00 (78%)
#   Remaining: $5.39
#
# 🎉 Great job stopping manually! You saved ~$0.17 vs auto-stop.
#
# See you next time! 👋
```

### Week 10: Emily Helps a Friend

```bash
# Emily's friend (new to class, confused):
# "Emily, I'm scared to use this AWS thing. What if I get charged?"

emily@laptop:~$ lfr student share-guide

# Output (MISSING):
# 📧 Sharing Student Guide
#
# I've created a personalized guide for you to share:
#
# Link: https://lfr-tools.dev/student-guide/emily-cs305
#
# This guide includes:
#   - Your personal experience and tips
#   - Actual cost numbers from your usage
#   - Screenshots (optional)
#   - "I was scared too, but..." testimonial
#
# Share this link with classmates to help them get started!
#
# Example message (copy-paste):
#
# "Hey! I was nervous about AWS too, but it's actually really easy and safe.
#  I've been using it for 10 weeks and spent $18 out of my $24 budget.
#  Check out this guide I made: [link]
#
#  Key things:
#  1. You won't get a bill - it's covered by course fee
#  2. Always stop your instance when done (one command)
#  3. TAs can reset if you break something
#
#  Let me know if you need help! - Emily"

# Emily shares with friend, helps them get started
# Emily thinks: "I'm basically a cloud expert now!" 😊
```

---

## 📋 Feature Gap Analysis: Student Needs

### Critical Missing Features (Blocks Student Persona)

| Feature | Priority | Impact | Current State | Effort |
|---------|----------|--------|---------------|--------|
| **Budget Impact Preview** | 🔴 Critical | Reduces anxiety | None | Low |
| **Student Dashboard** | 🔴 Critical | Quick status check | CLI only | Medium |
| **Cost Education/Explanation** | 🔴 Critical | Eliminates anxiety | None | Low |
| **Interactive Onboarding** | 🟡 High | First-time user success | Basic docs | Medium |
| **Self-Service Help** | 🟡 High | Reduces TA burden | None | Low |
| **Email Notifications** | 🟡 High | Awareness of costs | None | Medium |
| **Instance Reset Request** | 🟡 High | Broken environment fix | Requires TA | Low |
| **Mobile/Web Dashboard** | 🟢 Medium | Convenience | None | High |
| **Peer Help Tools** | 🟢 Low | Student collaboration | None | Low |

### Unique Student Requirements

| Requirement | Current | Needed | Priority |
|-------------|---------|--------|----------|
| **"Will this cost me money?" clarity** | ⚠️ Unclear | ✅ Clear, repeated messaging | Critical |
| **Budget preview before starting** | ❌ None | ✅ Show cost impact | Critical |
| **Automatic notifications** | ❌ None | ✅ Email when running, when stopped | High |
| **Plain-English cost explanations** | ❌ None | ✅ "In your terms" cost breakdown | Critical |
| **Recovery from broken environment** | ⚠️ Manual TA | ✅ Self-service or fast TA reset | High |
| **Quick status check (mobile)** | ❌ None | ✅ Web dashboard | Medium |
| **Interactive learning** | ⚠️ Docs only | ✅ Tutorial, guided setup | High |

---

## 🎯 Priority Recommendations: Student Persona

### Phase 1: Eliminate Budget Anxiety (Milestone 1)

**Target**: Students feel safe and confident using LFR

1. **Budget Impact Preview** (3 days)
   - Show budget before every start
   - "This session will cost ~$X, you'll have $Y left"
   - Clear "safe" vs "warning" indicators

2. **Cost Education Tool** (1 week)
   - `lfr student costs explain` command
   - Plain-English breakdown
   - Real examples from student's usage

3. **Better Stop Reminders** (2 days)
   - Reminder after SSH disconnect
   - Email if running > 6 hours
   - Friendly, not scary

4. **Student Dashboard** (1 week)
   - `lfr student status` command
   - Budget, instance status, recent activity
   - Next assignment reminder (from Canvas)

**Estimated effort**: 3 weeks

### Phase 2: Empower Self-Service (Milestone 2)

**Target**: Students can solve common problems themselves

5. **Interactive Onboarding** (1 week)
   - `lfr student join` command
   - Step-by-step tutorial
   - Test connection, verify setup

6. **Self-Service Help** (3 days)
   - `lfr student help` command
   - Common issues and solutions
   - When to contact TA

7. **Instance Reset Request** (3 days)
   - `lfr student request-reset` command
   - Explains what will happen
   - Notifies TA, preserved files

**Estimated effort**: 2 weeks

### Phase 3: Enhanced Experience (Milestone 3)

**Target**: Best-in-class student cloud experience

8. **Email Notifications** (1 week)
   - Instance started/stopped
   - Budget milestones (50%, 75%, 90%)
   - Auto-stop notifications

9. **Web Dashboard** (2 weeks)
   - Simple web UI for quick checks
   - Mobile-friendly
   - No CLI required for basic tasks

10. **Cost Visualization** (3 days)
    - Week-over-week graphs
    - Projection to semester end
    - Comparison with classmates (anonymous)

**Estimated effort**: 4 weeks

---

## Success Metrics: Student Perspective

### Anxiety Reduction
- ✅ **"Will this cost me money?" question**: Eliminated (100% clear)
- ✅ **Students checking budget**: From 0% to 95% (visible in dashboard)
- ✅ **Budget surprises**: Eliminated (proactive notifications)

### Usage Confidence
- ✅ **Students avoiding usage due to fear**: 40% → 5%
- ✅ **Average sessions per week**: 2.5 → 4.5 (more confident to use)
- ✅ **Students within budget**: 75% → 95%

### Self-Sufficiency
- ✅ **TA help requests**: 15/week → 5/week (67% reduction)
- ✅ **"I broke my environment" incidents**: 10/semester → 2/semester
- ✅ **Time to recover from issues**: 2 hours → 10 minutes

### Satisfaction
- ✅ **Student satisfaction**: "Easier than own AWS" (92%)
- ✅ **Would recommend to friend**: 95%
- ✅ **Felt supported**: 90%

---

## Next Steps

1. **Survey Students**: Ask current CS students about cloud anxiety
2. **Prototype Budget Tools**: Mock up budget preview and cost explanations
3. **Design Onboarding**: Create step-by-step student join flow
4. **Implementation Plan**: Prioritize based on anxiety reduction

**Estimated Timeline**:
- Phase 1 (Eliminate Anxiety): 3 weeks
- Phase 2 (Self-Service): 2 weeks
- Phase 3 (Enhanced Experience): 4 weeks
- **Total**: ~9 weeks (2 months) to comprehensive student support

---

**Status**: Draft Walkthrough
**Persona**: Student (First-Time Cloud User)
**Priority**: 🔴 Critical (blocks educational adoption)
**Dependencies**: Core LFR functionality (exists), professor budget tracking (needed), TA reset (needed)

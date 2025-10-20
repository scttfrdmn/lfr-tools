# Scenario 6: AWS Credits Manager - University IT

## Persona: David Kim - Cloud Infrastructure Manager

**Background**:
- Cloud Infrastructure Manager, University IT
- 8 years experience in cloud operations (AWS, Azure, GCP)
- Manages university's AWS Organizations (200+ accounts)
- Reports to: Director of IT Infrastructure
- Technical level: Expert (AWS Certified Solutions Architect, FinOps practitioner)
- Budget responsibility: $2M/year cloud spending across university

**Primary Responsibilities**:
1. **AWS Credits Management**: Track, allocate, optimize $150K/year in credits
2. **Cost Optimization**: Keep university cloud spending within budget
3. **Account Management**: Provision/deprovision AWS accounts for grants/courses
4. **Budget Forecasting**: Project cloud spending, prevent overages
5. **Vendor Relations**: Negotiate with AWS, manage enterprise agreements

**Primary Concerns**:
1. **Credits Expiration**: $20K in credits expired last year (waste!)
2. **Allocation Efficiency**: Credits sitting unused while some projects pay cash
3. **Visibility**: Hard to see which projects have credits, which need them
4. **Accountability**: Need to prove credits were used for intended purposes
5. **Forecasting**: Will credits last? Should we request more?

**Pain Points from Last Year**:
- $20K in AWS Educate credits expired (nobody realized until too late)
- Biology department overspent by $15K (should have used credits first)
- Spent 3 days manually tracking which account used which credits
- Dean asked "Where are our credits going?" - took a week to compile report
- No way to move credits from underutilized projects to high-demand ones

---

## Current State: What Doesn't Work

### ❌ Problem: No Credits Visibility

**Current reality**:
```bash
# David logs into AWS Organizations console
# Navigates to Billing > Credits
# Sees: "Total Credits: $75,000"
# But... which accounts? What types? When do they expire?!

# He clicks through 50+ accounts one by one
# Account 1: $2,000 (AWS Educate, expires 2026-06-30)
# Account 2: $0 (no credits)
# Account 3: $5,000 (AWS Research, expires 2027-12-31)
# ... (48 more accounts to check manually)

# 2 hours later, David has an incomplete spreadsheet
# He thinks: "There MUST be a better way!"
```

**What should happen** (MISSING):
```bash
# David uses LFR Credits Manager
david@laptop:~$ lfr credits dashboard

# LFR Output (MISSING):
# 💰 AWS Credits Dashboard
#
# University: University of XYZ
# Total Credits: $75,000
#
# By Source:
#   - AWS Educate: $50,000 (expires 2026-06-30) - 18 months remaining
#   - AWS Research: $20,000 (expires 2027-12-31) - 30 months remaining
#   - Promotional: $5,000 (expires 2025-12-31) - 60 days remaining ⚠️
#
# By Status:
#   - Allocated: $45,000 (60%)
#   - Unallocated: $30,000 (40%)
#   - Used: $18,000 (24%)
#   - Remaining: $57,000 (76%)
#
# Burn Rate:
#   - Last 30 days: $6,000 ($200/day)
#   - Projected: $73,000/year
#   - Utilization: 97% (excellent!) ✅
#
# ⚠️  Urgent Actions:
#   1. $5,000 promotional credits expire in 60 days
#      → Allocate to high-usage projects immediately
#
#   2. Biology dept paying cash while CS has unused credits
#      → Reallocate $3,000 from CS to Biology
#
# Quick Actions:
#   View by account: lfr credits by-account
#   Allocate credits: lfr credits allocate
#   Export report: lfr credits report --month september
```

---

## 🎯 Ideal Future State: Complete Credits Management Workflow

### Initial Setup: Credits Inventory

```bash
# David receives notification: $50,000 AWS Educate credits awarded
# He adds them to LFR tracking system

david@laptop:~$ lfr credits import \
  --source "AWS Educate" \
  --amount 50000 \
  --expiration 2026-06-30 \
  --program "Educational Credit Program" \
  --restrictions "Teaching and research only"

# LFR Output:
# 💰 Importing AWS Credits
#
# Source: AWS Educate
# Amount: $50,000
# Expiration: 2026-06-30 (18 months from now)
# Restrictions: Teaching and research only
#
# Querying AWS Organizations for credit application...
# ✅ Credits confirmed in AWS account 111111111111 (Management)
#
# Credit Details:
#   - Program ID: EDU-2025-12345
#   - Award Date: 2025-01-15
#   - Expiration: 2026-06-30
#   - Type: Shared credits (all member accounts)
#   - Restrictions: No production workloads
#
# 📊 Credit Portfolio Updated:
#   - Total credits: $75,000 (was $25,000)
#   - AWS Educate: $50,000 (66.7%)
#   - Other sources: $25,000 (33.3%)
#
# 📧 Notifications sent:
#   - Research Admin: "New credits available for allocation"
#   - Department Chairs: "AWS Educate credits ready"
#   - Finance: "Credit asset recorded"
#
# Next steps:
#   1. Set allocation strategy: lfr credits strategy
#   2. Allocate to projects: lfr credits allocate
#   3. Monitor burn rate: lfr credits monitor

# David sets up expiration monitoring
david@laptop:~$ lfr credits set-alerts \
  --expiration-warning 90days,60days,30days,7days \
  --underutilization-threshold 50% \
  --email david.kim@university.edu,finance@university.edu

# Output:
# 🔔 Credit Alert Configuration
#
# Expiration Warnings:
#   ✅ 90 days before: Email to david.kim + finance team
#   ✅ 60 days before: Email + weekly reminders
#   ✅ 30 days before: Email + daily reminders + escalate to CIO
#   ✅ 7 days before: Email + Slack alert + emergency meeting
#
# Underutilization Alerts:
#   ✅ If < 50% utilized at halfway point: Suggest reallocation
#   ✅ Weekly burn rate reports
#   ✅ Monthly utilization analysis
#
# Saved! You'll never miss expiring credits again. 🎉
```

### Allocation Strategy: Prioritize High-Value Use

```bash
# David creates allocation strategy for credits
david@laptop:~$ lfr credits strategy create \
  --name "University-Wide-2025" \
  --interactive

# Interactive Strategy Builder:
# 💡 AWS Credits Allocation Strategy
#
# Let's define how to allocate $50,000 in AWS Educate credits:
#
# 1. Priority Levels (which projects get credits first?)
#    [x] Active research grants (high compute needs)
#    [x] Educational courses (aligned with AWS Educate mission)
#    [ ] Administrative/IT infrastructure (not allowed by AWS Educate)
#    [x] Student research projects
#
# 2. Allocation Method:
#    ( ) Equal split (everyone gets same amount)
#    (x) Need-based (high spenders get more)
#    ( ) First-come-first-served
#    ( ) Committee approval
#
# 3. Reallocation Policy:
#    [x] Allow reallocation from underutilized to high-demand
#    [x] Automatic reallocation if < 50% used at midpoint
#    [x] Emergency reallocation for grant overspend
#
# 4. Expiration Management:
#    [x] Prioritize expiring credits first
#    [x] Alert at 90/60/30/7 days before expiration
#    [x] Auto-allocate expiring credits to high-usage projects
#
# 5. Reporting:
#    [x] Monthly utilization reports to Finance
#    [x] Quarterly reports to Provost
#    [x] Real-time dashboard for department chairs
#
# Strategy saved: University-Wide-2025
#
# Initial Allocations (recommended):
#   - CS Department (teaching): $10,000 (20%)
#   - Biology Department (research): $15,000 (30%)
#   - Engineering (research): $10,000 (20%)
#   - General pool (new grants): $10,000 (20%)
#   - Emergency reserve: $5,000 (10%)
#
# Apply these allocations? [y/N]: y

# David applies the strategy
david@laptop:~$ lfr credits allocate-batch \
  --strategy University-Wide-2025 \
  --dry-run

# Output:
# 📊 Credits Allocation Preview (Dry Run)
#
# Total to allocate: $50,000 (AWS Educate)
#
# Allocations:
#   1. CS305-Fall2025 (Prof. Rodriguez): $2,000
#   2. CS405-Spring2026 (Prof. Martinez): $3,000
#   3. Bio-Research-Pool: $15,000
#      → NSF-2024-12345 (Martinez): $10,000
#      → NIH-R01-2025 (Chen): $5,000
#   4. Engineering-Research: $10,000
#      → DOE-2025-11111 (Kim): $8,000
#      → Private-2024-88888 (Lee): $2,000
#   5. General-Pool: $10,000 (for new grants)
#   6. Emergency-Reserve: $5,000
#
# Impact Analysis:
#   - CS Department: $5,000 credits → save $5,000 cash
#   - Biology: $15,000 credits → save $15,000 cash
#   - Engineering: $10,000 credits → save $10,000 cash
#   - Total cash savings: $40,000 ✅
#
# Looks good? [y/N]: y
#
# Applying allocations...
# ✅ Allocated $2,000 to CS305-Fall2025
# ✅ Allocated $3,000 to CS405-Spring2026
# ✅ Allocated $10,000 to NSF-2024-12345
# ... (7 more)
# ✅ All allocations complete!
#
# 📧 Notifications sent to all PIs and department chairs
```

### Monthly Operations: Monitor & Optimize

```bash
# It's September 15, David checks monthly status
david@laptop:~$ lfr credits status --month september

# Output:
# 💰 Credits Status - September 2025
#
# Total Credits: $75,000
# Used This Month: $6,000
# Used YTD: $18,000 (24%)
# Remaining: $57,000 (76%)
#
# Burn Rate:
#   - September: $200/day (projected $6,000/month)
#   - Average: $200/day
#   - Trend: ➡️  Stable
#
# By Source:
#   AWS Educate ($50,000, exp 2026-06-30):
#     - Allocated: $45,000 (90%)
#     - Used: $12,000 (24%)
#     - Remaining: $38,000
#     - Burn rate: $133/day
#     - Projected expiration use: 95% ✅
#
#   AWS Research ($20,000, exp 2027-12-31):
#     - Allocated: $15,000 (75%)
#     - Used: $4,000 (20%)
#     - Remaining: $16,000
#     - Burn rate: $44/day
#     - Projected expiration use: 90% ✅
#
#   Promotional ($5,000, exp 2025-12-31):
#     - Allocated: $0 (0%) ⚠️
#     - Used: $2,000 (40%)
#     - Remaining: $3,000
#     - Days remaining: 107
#     - 🚨 URGENT: Only 40% used, expires in 3.5 months!
#
# ⚠️  Actions Required:
#   1. Promotional credits underutilized
#      → Allocate $3,000 to high-usage projects this week
#      → Recommended: NSF-2024-12345 (Martinez, high GPU use)
#
# Top Credit Users (This Month):
#   1. NSF-2024-12345 (Martinez): $2,400 (40%)
#   2. NIH-R01-2025 (Chen): $1,200 (20%)
#   3. CS305-Fall2025 (Rodriguez): $800 (13%)
#   4. DOE-2025-11111 (Kim): $600 (10%)
#   5. Others: $1,000 (17%)
#
# Quick Actions:
#   View details: lfr credits details --source promotional
#   Reallocate: lfr credits reallocate
#   Export report: lfr credits report --format pdf

# David checks which projects are paying cash (should use credits)
david@laptop:~$ lfr credits cash-opportunity

# Output:
# 💡 Cash-to-Credits Opportunity Analysis
#
# Projects currently paying CASH that could use CREDITS:
#
# 1. NSF-2023-99999 (Rodriguez, CS)
#    - Cash spending: $1,367/month (very high!)
#    - Available credits: $3,000 (promotional, expiring soon)
#    - Opportunity: Allocate expiring credits → save $3,000 cash
#    - Impact: Extend grant runway by 2+ months
#
# 2. Bio-Teaching (Wilson)
#    - Cash spending: $500/month
#    - Available credits: $5,000 (AWS Educate, not expiring soon)
#    - Opportunity: Shift to credits → save $500/month
#    - Impact: Preserve departmental budget
#
# 3. Engineering-General (Lee)
#    - Cash spending: $300/month
#    - Available credits: $2,000 (AWS Research)
#    - Opportunity: Use research credits → save cash
#
# Total Cash Savings Opportunity: $8,000/month → $96,000/year! 💰
#
# Execute recommended reallocations? [y/N]: y
#
# Reallocating...
# ✅ Allocated $3,000 promotional credits to NSF-2023-99999
# ✅ Allocated $5,000 AWS Educate to Bio-Teaching
# ✅ Allocated $2,000 AWS Research to Engineering-General
#
# 📧 Notifications sent to PIs: "Your project now using credits instead of cash"
#
# Expected cash savings: $8,000/month ✅
```

### Expiration Crisis: Urgent Reallocation

```bash
# November 1, David receives urgent alert
# Email: "⚠️ URGENT: $5,000 in credits expire in 60 days!"

david@laptop:~$ lfr credits expiring --days 60

# Output:
# ⚠️  Credits Expiring Soon
#
# Promotional Credits: $5,000
# Expiration: 2025-12-31 (60 days)
# Currently allocated: $2,000 (40%)
# Currently used: $2,000 (40%)
# At risk: $3,000 (60%) 🚨
#
# Burn rate: $50/day (too slow!)
# Projected expiration waste: $3,000 💸
#
# Recommended Actions:
#
# 1. IMMEDIATE: Reallocate to high-usage projects
#    Top candidates (can burn $3,000 in 60 days):
#      - NSF-2024-12345 (Martinez): Burns $200/day (GPU heavy)
#      - NIH-R01-2025 (Chen): Burns $150/day
#      - DOE-2025-11111 (Kim): Burns $100/day
#
# 2. ALTERNATIVE: Convert to discount vouchers
#    Some AWS programs allow converting expiring credits to
#    future discount vouchers (50% value retention)
#
# 3. LAST RESORT: Emergency allocation
#    Offer to all active projects: "Use it or lose it"
#
# Recommended: Allocate all $3,000 to Martinez (can definitely use it)
#
# Execute? [y/N]: y

david@laptop:~$ lfr credits allocate \
  --source promotional \
  --amount 3000 \
  --project NSF-2024-12345 \
  --reason "Expiring credits - high GPU usage project" \
  --expires-priority

# Output:
# 💰 Allocating Expiring Credits
#
# Credits: $3,000 (Promotional, expires 60 days)
# Project: NSF-2024-12345 (Prof. Martinez, Biology)
# Reason: High GPU usage, can utilize before expiration
#
# Impact Analysis:
#   - Martinez cash burn: $200/day → $0/day (credits cover)
#   - Credits will be used: 15 days ($200/day × 15 = $3,000)
#   - Credits saved from expiration: $3,000 ✅
#   - Cash saved: $3,000 ✅
#   - Win-win! 🎉
#
# ✅ Allocation complete
# ✅ Martinez notified: "You've been allocated $3,000 in credits"
# ✅ Credits marked as "expires-priority" (use these first)
#
# 📊 Updated Status:
#   - Promotional credits: $5,000 → $2,000 used, $0 at risk ✅
#   - Expiration waste prevented: $3,000
#
# Great save! 🎉

# Two weeks later, David checks if credits were actually used
david@laptop:~$ lfr credits verify-usage \
  --project NSF-2024-12345 \
  --allocation-id ALLOC-2025-11-001

# Output:
# ✅ Credit Usage Verification
#
# Allocation: ALLOC-2025-11-001
# Project: NSF-2024-12345 (Martinez)
# Amount: $3,000
# Allocated: 2025-11-01
#
# Usage Status:
#   - Used: $2,100 (70%) in 10 days
#   - Remaining: $900 (30%)
#   - Burn rate: $210/day (on track!)
#   - Projected full utilization: 2025-11-15 (4 days) ✅
#
# Resources using these credits:
#   - GPU instances: $1,800 (86%)
#   - Storage: $200 (9%)
#   - Data transfer: $100 (5%)
#
# Compliance:
#   ✅ All usage tagged with grant number
#   ✅ All usage for research (no personal use)
#   ✅ Credits being used as intended
#
# Conclusion: Credits are being used efficiently! 🎉
# No waste expected. ✅
```

### Quarterly Reporting: Prove Value to Leadership

```bash
# End of Q3, David generates credits report for CFO
david@laptop:~$ lfr credits report \
  --period Q3-2025 \
  --audience executive \
  --format pdf

# Output:
# 📊 Generating Executive Credits Report - Q3 2025
#
# Report Contents:
#   ✅ Credits inventory (sources, amounts, expirations)
#   ✅ Allocation strategy and execution
#   ✅ Utilization metrics (burn rate, efficiency)
#   ✅ Cash savings (credits used instead of budget)
#   ✅ Waste prevention (expired vs saved)
#   ✅ ROI analysis (time saved, cash saved)
#
# Generated: AWS-Credits-Report-Q3-2025.pdf
#
# Key Findings (for CFO):
#
# 1. Credits Portfolio
#    - Total credits received: $75,000
#    - Total credits used: $54,000 (72%)
#    - Total credits remaining: $21,000
#    - Expiration waste: $0 (vs $20K last year!) 🎉
#
# 2. Cash Savings
#    - Projects using credits instead of cash: $54,000
#    - Cash budget preserved: $54,000
#    - Effective discount rate: 100% (credits are free money!)
#    - ROI: Infinite (no cost, all benefit)
#
# 3. Utilization Efficiency
#    - Target: 95% utilization before expiration
#    - Actual: 97% on track ✅
#    - Improvement vs last year: 72% → 97% (+25%)
#    - Waste reduction: $20K → $0 (100% improvement!)
#
# 4. Time Savings
#    - David's time on credits management:
#      * Last year: 20 hours/month (manual tracking)
#      * This year: 2 hours/month (LFR automation)
#      * Time saved: 18 hours/month = 216 hours/year
#      * Value: $21,600/year (assuming $100/hour)
#
# 5. Impact on Research
#    - Grants supported: 25 active grants
#    - Researchers supported: 45 (PIs, postdocs, students)
#    - Publications enabled: 8 papers (credits acknowledged)
#    - Teaching: 3 courses (400 students)
#
# 6. Recommendations for Next Year
#    - Request additional $50K in AWS Research credits
#    - Justify: 97% utilization, high demand, enables more research
#    - Expected impact: Support 15 more grants
#
# This report demonstrates excellent stewardship of AWS credits! 🎉
#
# CFO will be pleased. ✅

# David also generates technical report for IT team
david@laptop:~$ lfr credits report \
  --period Q3-2025 \
  --audience technical \
  --include-details

# Output includes:
# - Per-account credit breakdown
# - Burn rate trends (charts)
# - Allocation history (who got what, when)
# - Reallocation decisions (why moved credits)
# - API usage statistics
# - System health metrics
```

---

## 📋 Feature Gap Analysis: Credits Manager Needs

### Critical Missing Features

| Feature | Priority | Impact | Current State | Effort |
|---------|----------|--------|---------------|--------|
| **Credits Dashboard** | 🔴 Critical | No visibility | Manual tracking (2hrs) | Medium |
| **Expiration Monitoring** | 🔴 Critical | $20K waste last year | None | Low |
| **Allocation Management** | 🔴 Critical | Manual, error-prone | None | Medium |
| **Burn Rate Tracking** | 🔴 Critical | Can't forecast | None | Medium |
| **Cash Opportunity Analysis** | 🟡 High | Suboptimal allocation | None | Medium |
| **Reallocation Tools** | 🟡 High | Waste prevention | Manual | Low |
| **Utilization Reports** | 🟡 High | Prove value to leadership | Manual (8hrs) | Medium |
| **Alert System** | 🟡 High | Proactive management | Email chaos | Low |
| **AWS API Integration** | 🟢 Medium | Auto-sync credit balance | None | High |

### Unique Credits Manager Requirements

| Requirement | Needed | Priority |
|-------------|--------|----------|
| **Real-time credit balance** | ✅ AWS Billing API integration | Critical |
| **Expiration tracking** | ✅ Alerts 90/60/30/7 days before | Critical |
| **Allocation commands** | ✅ Assign credits to accounts/projects | Critical |
| **Burn rate forecasting** | ✅ Project expiration usage | High |
| **Cash-to-credits opportunities** | ✅ Identify misallocated spend | High |
| **Reallocation flexibility** | ✅ Move credits between accounts | High |
| **Executive reporting** | ✅ Prove ROI to leadership | Medium |

---

## 🎯 Priority Recommendations: Credits Manager Persona

### Phase 1: Credits Visibility (Milestone 1)

**Target**: Credits manager can see entire portfolio

1. **Credits Dashboard** (2 weeks)
   - Query AWS Billing API for credits
   - Display by source, account, expiration
   - Burn rate calculations
   - Allocation status

2. **Expiration Monitoring** (1 week)
   - Track expiration dates
   - Alert at 90/60/30/7 days
   - Highlight at-risk credits
   - Recommend actions

**Estimated effort**: 3 weeks

### Phase 2: Allocation Management (Milestone 2)

**Target**: Efficient allocation and reallocation

3. **Allocation Commands** (2 weeks)
   - Allocate credits to accounts/projects
   - Batch allocation
   - Allocation strategies
   - Track allocation history

4. **Reallocation Tools** (1 week)
   - Move credits between accounts
   - Emergency reallocation
   - Expiring credits → high-usage projects
   - Approval workflows

5. **Cash Opportunity Analysis** (1 week)
   - Identify cash spend that could use credits
   - Recommend reallocations
   - Calculate savings potential

**Estimated effort**: 4 weeks

### Phase 3: Reporting & Optimization (Milestone 3)

**Target**: Prove value, optimize usage

6. **Utilization Reports** (2 weeks)
   - Executive summaries (CFO)
   - Technical details (IT)
   - ROI analysis
   - Year-over-year comparison

7. **Forecasting** (1 week)
   - Burn rate projections
   - Expiration risk assessment
   - Budget planning assistance
   - "Should we request more credits?" tool

**Estimated effort**: 3 weeks

---

## Success Metrics: Credits Manager Perspective

### Financial Impact
- ✅ **Credits waste**: $20K/year → $0/year (100% reduction)
- ✅ **Utilization**: 73% → 97% (+24%)
- ✅ **Cash savings**: $54K/year preserved
- **Total value**: $74K/year benefit

### Operational Efficiency
- ✅ **Time spent managing**: 20 hrs/month → 2 hrs/month (90% reduction)
- ✅ **Allocation decisions**: Days → Minutes
- ✅ **Reporting time**: 8 hours → 5 minutes (99% reduction)

### Strategic Value
- ✅ **Leadership visibility**: None → Real-time dashboards
- ✅ **Waste prevention**: Reactive → Proactive (90+ day alerts)
- ✅ **ROI proof**: Anecdotal → Data-driven

---

## Next Steps

1. **Research AWS Billing API**: Credits query endpoints
2. **Prototype Dashboard**: Credits visibility mockup
3. **Design Allocation Model**: How to assign credits to accounts
4. **Interview Credits Managers**: 2-3 at universities
5. **Implementation Plan**: Prioritize based on waste prevention

**Estimated Timeline**:
- Phase 1 (Credits Visibility): 3 weeks
- Phase 2 (Allocation Management): 4 weeks
- Phase 3 (Reporting & Optimization): 3 weeks
- **Total**: ~10 weeks (2.5 months) to comprehensive credits management

---

**Status**: Draft Walkthrough
**Persona**: AWS Credits Manager (University IT)
**Priority**: 🔴 Critical (prevents waste, enables efficient use of free resources)
**Dependencies**: AWS Billing API (credit balance query), AWS Organizations (account management), cost allocation (track credit usage)
**ROI**: $74K/year benefit vs ~$50K implementation cost = positive ROI in first year

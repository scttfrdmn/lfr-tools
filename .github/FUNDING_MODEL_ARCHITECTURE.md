# LFR Tools - Funding Model & Account Architecture

**Date**: October 20, 2025
**Status**: Research & Design Phase

## Executive Summary

This document defines how LFR Tools maps funding sources (grants, departments, credits) to AWS account structures, addressing the fundamental requirement that **someone must pay for AWS usage**. Unlike Google Cloud's built-in Projects concept, AWS requires careful architecture around Organizations, accounts, billing, and credits.

---

## The Fundamental Problem

### "Who Pays for AWS?"

Every AWS usage must be funded by one of these sources:
1. **Grant Funding** - NSF, NIH, DOE, private foundations (time-bounded, restricted use)
2. **Departmental Budgets** - University IT allocations, course fees
3. **AWS Credits** - Educational credits, research credits, promotional credits (expiring)
4. **Personal Funding** - Individual researcher credit cards (rare, discouraged)

### Key Constraints

1. **Grants Have Restrictions**:
   - Specific budget amounts
   - Specific time periods (start/end dates)
   - Specific allowable uses
   - Audit requirements
   - Cannot commingle funds from different grants

2. **AWS Credits Expire**:
   - Typical lifespan: 1-2 years
   - Use-it-or-lose-it (no rollover)
   - Cannot be transferred (usually)
   - Apply at org or account level

3. **AWS Doesn't Have "Projects"**:
   - Google Cloud: Projects are first-class billing entities
   - AWS: Only Organizations > Accounts > Resources
   - No built-in grant tracking
   - Tagging is optional and manual

4. **University Systems Are Complex**:
   - Workday (financial system)
   - Research.gov / eRA Commons (grant systems)
   - Petri (research organization - needs integration)
   - Internal chartstring/cost center codes

---

## AWS Account Structure Options

### Option 1: Single Account with Tagging (Current Implied Model)

```
University AWS Account (123456789012)
├─ IAM Users: alice, bob, charlie, ...
├─ LFR Instances: Tagged with "Project"
│  ├─ alice-ubuntu (Project: CS305-Fall2025)
│  ├─ bob-ubuntu (Project: CS305-Fall2025)
│  ├─ researcher1-ubuntu (Project: NSF-Grant-2024)
│  └─ researcher2-ubuntu (Project: NIH-R01-2025)
└─ Billing: One consolidated bill
   ├─ AWS Credits applied at account level
   └─ Cost allocation via tags (manual tracking)
```

**Pros:**
- ✅ Simple to set up
- ✅ Credits apply to all usage
- ✅ Volume discounts apply across all usage

**Cons:**
- ❌ No hard budget separation (grants commingle)
- ❌ Manual cost allocation (error-prone)
- ❌ No per-project billing enforcement
- ❌ Audit challenges (tags can be changed/deleted)
- ❌ Risk of one project overspending entire account

**Verdict**: ⚠️ Only suitable for small-scale or single-funding-source scenarios

---

### Option 2: Multi-Account with AWS Organizations (Recommended)

```
University Organization (Management Account: 111111111111)
│
├─ Billing Account (no workloads, billing only)
│  └─ Consolidated billing for all member accounts
│
├─ Credits Pool
│  ├─ Educational credits (AWS Educate): $50,000/year
│  ├─ Research credits (AWS Research): $25,000/grant
│  └─ Promotional credits: $5,000
│
├─ Organizational Units (OUs)
│  │
│  ├─ OU: Computer Science Department
│  │  ├─ Account: CS305-Fall2025 (222222222222)
│  │  │  └─ Funding: Course fees ($840)
│  │  ├─ Account: CS405-Spring2026 (333333333333)
│  │  │  └─ Funding: Course fees ($1,200)
│  │  └─ Account: CS-Research-General (444444444444)
│  │     └─ Funding: Departmental budget
│  │
│  ├─ OU: Biology Department
│  │  ├─ Account: NSF-Grant-2024-12345 (555555555555)
│  │  │  └─ Funding: NSF Grant $50,000 (2024-01-01 to 2026-12-31)
│  │  ├─ Account: NIH-R01-2025-67890 (666666666666)
│  │  │  └─ Funding: NIH R01 $250,000 (2025-07-01 to 2030-06-30)
│  │  └─ Account: Bio-Teaching (777777777777)
│  │     └─ Funding: Departmental budget + AWS Educate credits
│  │
│  └─ OU: Engineering Department
│     └─ Account: DOE-Grant-2025-11111 (888888888888)
│        └─ Funding: DOE Grant $100,000 (2025-01-01 to 2027-12-31)
│
└─ Cost Allocation
   ├─ Per-account billing (hard separation)
   ├─ Credits applied per account or shared
   └─ Budgets enforced via AWS Budgets API
```

**Pros:**
- ✅ **Hard budget separation** (grants don't commingle)
- ✅ **Audit-friendly** (each grant = separate account = separate bill)
- ✅ **Enforceable limits** (AWS Budgets can stop resources)
- ✅ **Credits can be targeted** (assign to specific accounts)
- ✅ **Compliance-ready** (NSF/NIH audits easier)
- ✅ **Hierarchical organization** (mirrors university structure)

**Cons:**
- ⚠️ **More complex setup** (multiple accounts to create)
- ⚠️ **Volume discounts less effective** (smaller per-account usage)
- ⚠️ **Credit management complexity** (which account gets which credits?)

**Verdict**: ✅ **Recommended for institutional scale**

---

## Funding Source Mapping

### 1. Grant-Funded Research

**Scenario**: Prof. Sarah Chen has NSF grant for computational biology research

```yaml
Grant Details:
  Sponsor: National Science Foundation
  Grant Number: NSF-2024-12345
  PI: Sarah Chen
  Award Amount: $50,000
  Period: 2024-01-01 to 2026-12-31
  Chartstring: 12-34567-890 (university internal code)

AWS Mapping:
  Account: NSF-Grant-2024-12345 (555555555555)
  Budget: $15,000 for cloud computing (30% of total)
  OU: Biology Department > Research Grants
  Tags:
    - GrantNumber: NSF-2024-12345
    - PI: sarah.chen@university.edu
    - Chartstring: 12-34567-890
    - Department: Biology
    - EndDate: 2026-12-31

LFR Project:
  Name: NSF-Computational-Biology-2024
  Account: 555555555555
  Budget: $15,000
  Members:
    - sarah.chen (PI, admin)
    - postdoc1 (researcher, $3,000 sub-budget)
    - grad-student1 (researcher, $2,000 sub-budget)
    - grad-student2 (researcher, $2,000 sub-budget)
  End Date: 2026-12-31 (auto-cleanup)

Petri Integration:
  Petri Project ID: PROJ-2024-789
  Sync: Bi-directional
    - LFR pushes AWS costs to Petri daily
    - Petri provides grant balance to LFR
    - Alerts when grant is 80%, 90%, 100% spent
```

**Workflow**:
1. Research admin creates AWS account for grant in AWS Organizations
2. Research admin configures LFR project linked to account
3. PI adds researchers to LFR project
4. LFR tracks spending per researcher (sub-budgets)
5. Daily cost sync to Petri
6. Weekly reports to PI
7. Auto-cleanup at grant end date

---

### 2. Course/Educational Funding

**Scenario**: Prof. Rodriguez teaches CS 305 with 35 students, $840 budget from course fees

```yaml
Course Details:
  Course: CS 305 - Cloud Computing
  Professor: Maria Rodriguez
  Term: Fall 2025
  Students: 35
  Budget: $840 ($24/student from course fee)
  Chartstring: 98-76543-210 (university IT allocation)

AWS Mapping:
  Account: CS305-Fall2025 (222222222222)
  Budget: $840
  OU: Computer Science Department > Teaching
  Credits: AWS Educate credits ($2,000 available)
  Tags:
    - Course: CS305
    - Term: Fall2025
    - Professor: maria.rodriguez@university.edu
    - Chartstring: 98-76543-210
    - Department: ComputerScience
    - EndDate: 2025-12-12

LFR Project:
  Name: CS305-Fall2025
  Account: 222222222222
  Budget: $840 (course fee) + $2,000 (AWS Educate credits)
  Members:
    - maria.rodriguez (professor, admin)
    - alex.chen (head TA, debug/reset)
    - sarah.kim (TA, view only)
    - mike.wong (TA, view only)
    - 35 students ($24 sub-budget each)
  End Date: 2025-12-12 (auto-cleanup)

Petri Integration:
  Petri Course ID: COURSE-CS305-F2025
  Sync: Read-only (costs reported to Petri)
    - LFR pushes daily costs to Petri
    - IT department views costs in Petri dashboard
```

**Workflow**:
1. IT admin creates AWS account for course (or reuses existing)
2. IT admin allocates AWS Educate credits to account
3. Professor creates LFR project (or IT admin does it)
4. Professor imports students from Canvas
5. Students work within budget
6. Auto-stop policies enforce budgets
7. End-of-semester: auto-cleanup, final report to IT

---

### 3. AWS Credits (Educational/Research)

**Scenario**: University has $50,000 in AWS Educate credits, need to allocate across courses

```yaml
Credits Available:
  Source: AWS Educate Program
  Amount: $50,000
  Expiration: 2026-06-30 (18 months)
  Type: Educational use only
  Account: Applied at Organization level

Allocation Strategy:
  Total Credits: $50,000
  Allocation Plan:
    - CS305 (Fall 2025): $2,000
    - CS405 (Spring 2026): $3,000
    - Bio-Teaching: $5,000
    - Research (various): $30,000
    - Reserve: $10,000 (unallocated)

AWS Mapping:
  Credits Pool: Organization level (shared)
  OR
  Credits Per Account: Assigned to specific accounts

LFR Credits Manager:
  View: lfr credits list
    - Total: $50,000
    - Used: $12,000 (24%)
    - Remaining: $38,000
    - Expiring: $0 (no credits expiring in 30 days)

  Allocate: lfr credits allocate --account 222222222222 --amount 2000 --project CS305-Fall2025

  Monitor: lfr credits burn-rate
    - Current: $400/day
    - Projected expiry: 2026-05-15 (45 days before expiration)
    - Alert: On track to use all credits ✅

Petri Integration:
  Credits tracked in Petri
  Real-time balance updates
  Expiration alerts
  Reallocation recommendations
```

**Workflow**:
1. Credits manager receives AWS credits (via AWS Educate, Research program, etc.)
2. Credits manager views balance: `lfr credits list`
3. Credits manager allocates to projects: `lfr credits allocate`
4. Credits automatically apply to project spending
5. LFR tracks burn rate and projects expiration
6. Alerts when credits are underutilized (expiring soon)
7. Credits manager reallocates from low-use to high-use projects

---

### 4. Departmental Budget (Centralized IT)

**Scenario**: University IT has $10,000/year for cloud computing across all departments

```yaml
Department Budget:
  Source: University IT General Fund
  Amount: $10,000/year
  Allocation: Ad-hoc (professors request, IT approves)
  Chartstring: 11-11111-111 (central IT)

AWS Mapping:
  Account: University-IT-General (999999999999)
  Budget: $10,000/year ($833/month)
  OU: IT Services > General Fund

LFR Projects:
  Multiple projects share this account:
    - Quick research projects (< $500)
    - Proof-of-concepts
    - Student capstone projects
    - IT infrastructure testing

Approval Workflow:
  1. Researcher requests access: lfr request-access --budget 500 --reason "ML testing"
  2. IT admin reviews: lfr admin review-requests
  3. IT admin approves: lfr admin approve --request 12345 --budget 500 --duration 30days
  4. Researcher gets access for 30 days with $500 budget
  5. Auto-revoke after 30 days

Petri Integration:
  IT dashboard shows all active allocations
  Approval history
  Budget burn by requester
```

---

## Mapping Grants to AWS Accounts

### Grant Lifecycle in LFR Tools

```bash
# 1. Research Admin: Create grant account
lfr admin grant create \
  --grant-number NSF-2024-12345 \
  --pi sarah.chen@university.edu \
  --amount 50000 \
  --cloud-allocation 15000 \
  --start-date 2024-01-01 \
  --end-date 2026-12-31 \
  --chartstring 12-34567-890 \
  --sponsor NSF \
  --aws-account 555555555555

# LFR creates:
# - AWS account in Organizations (if needed)
# - Budget alerts at 80%, 90%, 100%
# - Cost allocation tags
# - Scheduled termination at grant end date
# - Petri project link

# 2. PI: Create LFR project for grant
lfr project create NSF-Computational-Biology-2024 \
  --grant NSF-2024-12345 \
  --budget 15000 \
  --account 555555555555

# 3. PI: Add researchers with sub-budgets
lfr project member add NSF-Computational-Biology-2024 \
  --email postdoc1@university.edu \
  --role researcher \
  --budget 3000

# 4. Ongoing: Track spending
lfr grant status NSF-2024-12345

# Output:
# Grant: NSF-2024-12345 (Computational Biology)
# PI: Sarah Chen
# Period: 2024-01-01 to 2026-12-31 (18 months remaining)
#
# Budget:
#   Total award: $50,000
#   Cloud allocation: $15,000 (30%)
#   Spent to date: $8,500 (57%)
#   Remaining: $6,500 (43%)
#   Projected end: 2026-10-15 (on track) ✅
#
# Researchers:
#   - postdoc1: $1,800 / $3,000 (60%)
#   - grad-student1: $1,200 / $2,000 (60%)
#   - grad-student2: $900 / $2,000 (45%)
#
# Compliance:
#   - All usage tagged with grant number ✅
#   - No non-research usage detected ✅
#   - Audit logs available ✅

# 5. End of grant: Auto-cleanup
# (December 31, 2026 at 11:59 PM)
# - Stop all instances
# - Archive data to S3
# - Revoke researcher access
# - Generate final report
# - Submit to research admin + Petri
```

---

## AWS Credits Management

### Credits Lifecycle

```bash
# 1. Credits Manager: View available credits
lfr credits list

# Output:
# AWS Credits Summary
#
# Total Credits: $75,000
#   - AWS Educate (expires 2026-06-30): $50,000
#   - AWS Research (expires 2027-12-31): $20,000
#   - Promotional (expires 2025-12-31): $5,000
#
# Allocated: $45,000 (60%)
#   - CS305-Fall2025: $2,000 (used $800)
#   - CS405-Spring2026: $3,000 (used $0)
#   - NSF-Grant-2024: $10,000 (used $8,500)
#   - Bio-Teaching: $5,000 (used $3,200)
#   - Research-Pool: $25,000 (used $12,000)
#
# Unallocated: $30,000 (40%)
#
# Burn Rate: $400/day
# Projected Expiration: 2026-05-15 (45 days before credits expire)
#
# ⚠️  Alert: $5,000 in promotional credits expire in 60 days!
#     Recommendation: Allocate to active projects

# 2. Credits Manager: Allocate credits to project
lfr credits allocate \
  --project CS305-Fall2025 \
  --amount 2000 \
  --source "AWS Educate" \
  --reason "Course cloud computing budget"

# 3. Credits Manager: Reallocate unused credits
lfr credits reallocate \
  --from CS405-Spring2026 \
  --to NSF-Grant-2024 \
  --amount 1000 \
  --reason "CS405 cancelled, NSF project needs more"

# 4. Credits Manager: Monitor burn rate
lfr credits burn-rate

# Output:
# Credits Burn Rate Analysis
#
# Last 30 days: $12,000 ($400/day)
# Last 7 days: $3,500 ($500/day) ⚠️  Increasing
#
# By Project:
#   - NSF-Grant-2024: $200/day (highest)
#   - Bio-Teaching: $150/day
#   - Research-Pool: $100/day
#   - CS305-Fall2025: $50/day
#
# Expiration Risk:
#   - AWS Educate ($50,000, expires 2026-06-30):
#     → 18 months remaining
#     → Current burn: $400/day
#     → Will use: $216,000 projected (need more projects!)
#     → ⚠️  Underutilized! Consider more allocations
#
#   - Promotional ($5,000, expires 2025-12-31):
#     → 60 days remaining
#     → Current burn: $50/day
#     → Will use: $3,000 projected
#     → ⚠️  $2,000 will expire! Reallocate ASAP!

# 5. Credits Manager: Generate report
lfr credits report --month september

# Output: PDF report for finance department
# - Credits received
# - Credits allocated
# - Credits used
# - Credits expired (waste)
# - Credits remaining
# - Efficiency: 95% utilization ✅
```

---

## Petri Integration Architecture

### What is Petri?

**Petri** is a university research management system that:
- Tracks research projects across departments
- Manages grant funding and budgets
- Integrates with financial systems (Workday, Oracle, etc.)
- Provides dashboards for research admins, PIs, deans
- Generates compliance reports for sponsors

### Integration Points

```yaml
LFR Tools ←→ Petri Integration:

1. Authentication:
   - SSO via Petri (SAML/OAuth)
   - User provisioning from Petri
   - Role mapping (Petri roles → LFR roles)

2. Project/Grant Import:
   - Petri creates research project
   - LFR imports project metadata:
     * Grant number
     * PI information
     * Budget allocation
     * Start/end dates
     * Chartstring
   - Auto-create LFR project

3. Cost Export:
   - LFR pushes daily AWS costs to Petri
   - Tagged with grant number, chartstring
   - Petri aggregates with other costs (servers, software, etc.)
   - Unified cost view for research admin

4. Budget Alerts:
   - LFR monitors AWS spending
   - Alert Petri when 80%, 90%, 100% of allocation
   - Petri notifies PI and research admin
   - Petri can trigger budget reallocation approval workflow

5. Compliance Reporting:
   - LFR exports usage logs to Petri
   - Petri generates sponsor reports (NSF, NIH format)
   - Audit trail: Who used what, when, for what
   - Petri stores for 7 years (compliance requirement)

6. User Management:
   - Petri manages user affiliations (who works on which grant)
   - LFR syncs user access
   - When user leaves grant in Petri → auto-revoke in LFR

API Spec (Hypothetical):
  Petri Exposes:
    - GET /api/projects (list all research projects)
    - GET /api/projects/{id} (project details)
    - GET /api/projects/{id}/members (who has access)
    - POST /api/costs (LFR posts costs)
    - POST /api/alerts (LFR posts budget alerts)

  LFR Exposes:
    - POST /api/sync-project (Petri triggers sync)
    - POST /api/revoke-user (Petri revokes access)
    - GET /api/costs (Petri pulls costs)
```

---

## Best Practices & Recommendations

### For Small Institutions (< 100 users)
- Start with Option 1 (single account + tagging)
- Use LFR projects for logical organization
- Manual cost tracking acceptable
- Transition to multi-account when > 5 grants

### For Medium Institutions (100-1000 users)
- Use Option 2 (multi-account Organizations)
- One account per grant
- One account per course
- Central credits management

### For Large Institutions (1000+ users)
- Multi-account with OU hierarchy
- Automated account creation via Petri
- Dedicated credits manager role
- Full Petri integration

### Credits Management
- Track expiration dates religiously
- Alert 90/60/30 days before expiry
- Reallocate underutilized credits
- Aim for 95%+ utilization (avoid waste)

### Grant Compliance
- One AWS account = One grant (hard rule)
- Tag ALL resources with grant number
- Export audit logs to Petri
- Retain logs for 7 years

---

## Implementation Roadmap

### Phase 1: Multi-Account Foundation (8 weeks)
1. AWS Organizations setup commands
2. Account creation automation
3. Cost allocation tag enforcement
4. Budget alerts per account

### Phase 2: Credits Management (6 weeks)
5. Credits balance query (AWS Billing API)
6. Credits allocation commands
7. Burn rate tracking
8. Expiration alerts

### Phase 3: Grant Integration (8 weeks)
9. Grant lifecycle commands
10. Chartstring tracking
11. PI/researcher management
12. Compliance reporting

### Phase 4: Petri Integration (12 weeks)
13. Petri API research and design
14. SSO authentication
15. Project import/sync
16. Cost export (daily sync)
17. User management sync

**Total Estimated Effort**: 34 weeks (~8 months)

---

**Status**: Architecture Design
**Next Steps**:
1. Validate with university research admins
2. Prototype credits management
3. Research Petri API specifications
4. Create GitHub issues for implementation

**Dependencies**:
- AWS Organizations access
- Petri system access (if exists)
- University IT approvals
- Research admin input

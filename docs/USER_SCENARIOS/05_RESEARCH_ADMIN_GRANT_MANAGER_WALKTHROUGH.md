# Scenario 5: Research Administrator - Grant & Budget Manager

## Persona: Jennifer Park - Senior Research Administrator

**Background**:
- Senior Research Administrator, Office of Sponsored Programs
- 15 years experience in research administration
- Manages 50+ active grants across 5 departments
- Total research portfolio: $25M/year (NSF, NIH, DOE, private foundations)
- Technical level: Moderate (Excel expert, knows Workday/SAP, not a programmer)
- Reports to: Associate Dean of Research

**Primary Responsibilities**:
1. **Grant Lifecycle**: Pre-award through closeout
2. **Budget Management**: Ensure spending aligns with grant restrictions
3. **Compliance**: NSF/NIH reporting, audits, cost accounting
4. **Systems Integration**: Workday (finance), Research.gov, eRA Commons, Petri
5. **PI Support**: Help faculty manage grant budgets
6. **Audit Preparation**: Prepare for federal audits (A-133, single audit)

**Primary Concerns**:
1. **Compliance**: One mistake = audit finding = university reputation at risk
2. **Visibility**: Need real-time view of spending across all grants
3. **Audit Trail**: Must prove every dollar spent is allowable and allocable
4. **PI Communication**: PIs need budget alerts before overspending
5. **Systems Integration**: Too many disconnected systems (Workday, AWS, Petri, etc.)

**Pain Points from Previous Year**:
- Discovered $15K AWS overspend on NSF grant during quarterly review (too late!)
- Spent 40 hours manually reconciling AWS bills with Workday chartstrings
- Federal auditor questioned cloud spending - took 2 weeks to compile documentation
- PI asked "How much budget do I have left?" - took 3 days to calculate across all systems
- AWS credits expired unused ($20K wasted) - nobody was tracking

---

## Current State (v0.1.x): What Doesn't Work

### ❌ Problem: No Central Visibility Across Grants

**Scenario**: Jennifer needs to see all active grants and their AWS spending

**Current reality**:
```bash
# Jennifer logs into AWS Console
# Sees one massive bill: $45,000 for the month
# But... which grant paid for what?!

# She exports CSV, opens Excel, starts manual tagging
# - Matches instance names to PIs (if tagged correctly)
# - Looks up PI → Grant mapping in Workday
# - Allocates costs to chartstrings (manual data entry)
# - 8 hours later... still not confident in the numbers

# Jennifer thinks: "There MUST be a better way!"
```

**What should happen** (MISSING):
```bash
# Jennifer uses LFR Tools research admin dashboard
jennifer@laptop:~$ lfr admin dashboard

# LFR Output (MISSING):
# 📊 Research Administration Dashboard
#
# University: University of XYZ
# Period: September 2025
# Total AWS Spending: $45,230
#
# Active Grants: 52
# Departments: 5 (Biology, Chemistry, Engineering, CS, Physics)
# Principal Investigators: 28
#
# 💰 Top Spending Grants:
# 1. NSF-2024-12345 (Martinez, Biology): $8,500 / $15,000 (57%)
# 2. NIH-R01-2025-67890 (Chen, Chemistry): $6,200 / $25,000 (25%)
# 3. DOE-2025-11111 (Kim, Engineering): $5,800 / $10,000 (58%)
# 4. NSF-2023-99999 (Rodriguez, CS): $4,100 / $5,000 (82%) ⚠️
# 5. Private-2024-88888 (Lee, Physics): $3,900 / $8,000 (49%)
# ... (47 more)
#
# 🚨 Alerts (Require Action):
# 1. NSF-2023-99999 (Rodriguez): 82% spent, only 3 months into 36-month grant!
#    → Recommendation: Contact PI about spending rate
# 2. NIH-R01-2024-55555 (Wilson): Grant ends in 30 days, $12K unspent
#    → Recommendation: Notify PI to use remaining budget
# 3. AWS Credits expiring: $8,000 in 45 days (AWS Educate)
#    → Recommendation: Allocate to active projects
#
# 📈 Spending Trends:
# - This month: $45,230 (↑ 15% vs last month)
# - YTD: $420,000
# - Projected EOY: $560,000 (on track with budget)
#
# Quick Actions:
#   View grant: lfr admin grant-details NSF-2024-12345
#   Generate report: lfr admin report --month september
#   Check compliance: lfr admin compliance-check --all
#   Export to Workday: lfr admin export-workday
```

---

## 🎯 Ideal Future State: Complete Research Admin Workflow

### Pre-Award: Grant Account Setup

```bash
# Jennifer receives notice: New NSF grant awarded to Prof. Martinez
# Grant Number: NSF-2024-12345
# Award Amount: $500,000 over 3 years
# Cloud computing: $15,000 (3% of total)
# Chartstring: 12-34567-890 (from Workday)

# Jennifer creates AWS account for the grant
jennifer@laptop:~$ lfr admin grant create \
  --grant-number NSF-2024-12345 \
  --sponsor NSF \
  --pi martinez@university.edu \
  --department Biology \
  --total-award 500000 \
  --cloud-allocation 15000 \
  --start-date 2024-01-01 \
  --end-date 2026-12-31 \
  --chartstring 12-34567-890 \
  --workday-project PRJ-2024-789 \
  --petri-project PETRI-BIO-2024-123

# LFR Output:
# 🎓 Creating Grant Account
#
# Grant: NSF-2024-12345
# PI: Prof. Martinez (Biology)
# Period: 2024-01-01 to 2026-12-31 (36 months)
# Cloud budget: $15,000 (3% of $500,000 total)
#
# AWS Setup:
#   ✅ Created AWS account: 555555555555
#   ✅ Added to Organization: University-XYZ
#   ✅ Organizational Unit: Biology Department > Research Grants
#   ✅ Consolidated billing: Enabled
#   ✅ Cost allocation tags:
#      - GrantNumber: NSF-2024-12345
#      - PI: martinez@university.edu
#      - Department: Biology
#      - Chartstring: 12-34567-890
#      - Sponsor: NSF
#      - EndDate: 2026-12-31
#
# Budget Alerts:
#   ✅ 80% threshold: $12,000
#   ✅ 90% threshold: $13,500
#   ✅ 100% threshold: $15,000
#   ✅ Notifications: PI + Research Admin + Department Chair
#
# Compliance:
#   ✅ Auto-tagging: All resources tagged with grant number
#   ✅ Audit logging: S3 bucket created (7-year retention)
#   ✅ Grant end date: Auto-cleanup scheduled for 2026-12-31
#
# Integrations:
#   ✅ Workday: Linked to project PRJ-2024-789
#   ✅ Petri: Synced with PETRI-BIO-2024-123
#   ✅ Research.gov: Grant number validated
#
# 📧 Notifications sent:
#   - PI (martinez@university.edu): "Your grant AWS account is ready"
#   - Dept Chair: "New grant account created for Martinez"
#   - IT Admin: "AWS account 555555555555 provisioned"
#
# Next steps:
#   1. PI adds researchers: lfr project member add ...
#   2. PI allocates sub-budgets to team
#   3. Researchers start work
#   4. LFR tracks spending automatically
#
# Setup complete! Total time: 3 minutes (vs 2 hours manual setup)

# Jennifer configures Petri integration for daily sync
jennifer@laptop:~$ lfr admin petri-sync configure \
  --grant NSF-2024-12345 \
  --sync-frequency daily \
  --push-costs yes \
  --pull-budget yes

# Output:
# 🔄 Petri Integration Configured
#
# Grant: NSF-2024-12345
# Petri Project: PETRI-BIO-2024-123
#
# Sync Schedule:
#   - Daily at 2:00 AM (after AWS billing updates)
#   - Push AWS costs to Petri
#   - Pull budget balance from Petri
#   - Alert if discrepancies detected
#
# Data Flow:
#   LFR → Petri: Daily AWS costs (by researcher, by resource type)
#   Petri → LFR: Budget balance, spending alerts
#   Bi-directional: User access (when user added/removed in Petri)
#
# ✅ First sync scheduled for tonight at 2:00 AM
```

### Post-Award: Monitoring & Compliance

```bash
# 6 months later, Jennifer checks on all grants
jennifer@laptop:~$ lfr admin grants list --status active

# Output:
# 📊 Active Grants (52 total)
#
# ┌─────────────────────┬─────────────┬──────────┬─────────┬─────────┬─────────┐
# │ Grant Number        │ PI          │ Dept     │ Budget  │ Spent   │ Status  │
# ├─────────────────────┼─────────────┼──────────┼─────────┼─────────┼─────────┤
# │ NSF-2024-12345      │ Martinez    │ Biology  │ $15,000 │ $8,500  │ ✅ OK   │
# │ NIH-R01-2025-67890  │ Chen        │ Chem     │ $25,000 │ $6,200  │ ✅ OK   │
# │ DOE-2025-11111      │ Kim         │ Engr     │ $10,000 │ $5,800  │ ✅ OK   │
# │ NSF-2023-99999      │ Rodriguez   │ CS       │ $5,000  │ $4,100  │ ⚠️  82% │
# │ NIH-R01-2024-55555  │ Wilson      │ Biology  │ $20,000 │ $8,000  │ ⚠️  End  │
# │ ... (47 more)                                                               │
# └─────────────────────┴─────────────┴──────────┴─────────┴─────────┴─────────┘
#
# Alerts:
#   2 grants require attention (see above)

# Jennifer investigates the Rodriguez grant (82% spent)
jennifer@laptop:~$ lfr admin grant-details NSF-2023-99999

# Output:
# 📊 Grant Details: NSF-2023-99999
#
# Grant Information:
#   Number: NSF-2023-99999
#   Title: "Scalable Machine Learning Infrastructure"
#   PI: Prof. Maria Rodriguez (CS Department)
#   Sponsor: National Science Foundation
#   Total Award: $450,000
#   Cloud Allocation: $5,000 (1.1% of total)
#   Period: 2023-07-01 to 2026-06-30 (36 months)
#   Elapsed: 3 months (8% of grant period)
#
# Budget Status:
#   Allocated: $5,000
#   Spent: $4,100 (82%)
#   Remaining: $900 (18%)
#   Burn rate: $1,367/month (last 3 months)
#
# ⚠️  RED FLAG: Spending too fast!
#   - 82% spent in only 8% of grant period
#   - At current rate: Budget exhausted in 21 days
#   - Projected total need: $49,200 (vs $5,000 budgeted)
#
# Spending Breakdown:
#   - GPU instances: $3,200 (78% - very expensive!)
#   - Storage (S3): $600 (15%)
#   - Standard compute: $300 (7%)
#
# Researchers:
#   - Prof. Rodriguez: $1,200 (29%)
#   - Grad Student 1: $1,500 (37%)
#   - Grad Student 2: $1,400 (34%)
#
# Recommendations:
#   1. Contact PI immediately about overspending
#   2. Review budget - may need reallocation from other line items
#   3. Consider spot instances for GPU work (70% savings)
#   4. Implement stricter auto-stop policies
#
# Compliance Status:
#   ✅ All spending tagged correctly
#   ✅ No non-grant usage detected
#   ✅ Audit logs complete
#   ⚠️  Budget overage risk
#
# Actions:
#   [1] Email PI warning
#   [2] Request budget reallocation
#   [3] Implement spending cap
#   [4] Schedule meeting with PI

# Jennifer sends alert to PI
jennifer@laptop:~$ lfr admin alert-pi NSF-2023-99999 \
  --subject "Urgent: Cloud budget 82% spent" \
  --include-recommendations

# Email sent to Prof. Rodriguez:
# Subject: ⚠️  Urgent: NSF-2023-99999 Cloud Budget Alert
#
# Dear Prof. Rodriguez,
#
# Your cloud computing budget for NSF grant 2023-99999 is critically low:
#
# Budget: $5,000
# Spent: $4,100 (82%)
# Remaining: $900 (18%)
# Burn rate: $1,367/month
# Estimated depletion: 21 days
#
# This is concerning because:
# - Only 8% of grant period has elapsed (3 of 36 months)
# - At current rate, you'll need $49,200 over 3 years (vs $5,000 budgeted)
#
# Primary cost driver: GPU instances ($3,200 = 78% of spending)
#
# Recommendations:
# 1. Consider spot instances for GPU work (70% savings = $2,240 saved!)
# 2. Implement aggressive auto-stop policies (currently minimal)
# 3. Request budget reallocation from other grant line items
# 4. Review whether GPU instances are necessary for all workloads
#
# Please contact me this week to discuss budget reallocation options.
#
# Best regards,
# Jennifer Park
# Office of Sponsored Programs
#
# cc: Department Chair, Grant Manager

### Quarterly Reporting: Export to Workday

```bash
# End of quarter, Jennifer needs to update Workday with actual spending
jennifer@laptop:~$ lfr admin export-workday \
  --period 2025-07-01:2025-09-30 \
  --format csv \
  --output Q3-AWS-Spending.csv

# Output:
# 📊 Exporting to Workday Format
#
# Period: Q3 2025 (July 1 - September 30)
# Grants: 52 active grants
# Total spending: $135,600
#
# Generating Workday-compatible CSV...
# ✅ Export complete: Q3-AWS-Spending.csv
#
# File format:
#   - Chartstring
#   - Grant Number
#   - PI
#   - Date
#   - Amount
#   - Description
#   - Cost Category (AWS-Compute, AWS-Storage, AWS-Data-Transfer)
#
# Rows: 1,247 transactions
# Ready for Workday import
#
# 💡 Tip: This file can be imported directly to Workday via:
#    Financial Management > Actuals > Import > AWS Cloud Costs

# Jennifer also generates compliance report for NSF
jennifer@laptop:~$ lfr admin compliance-report NSF-2024-12345 \
  --format nsf-annual \
  --year 2025

# Output:
# 📊 NSF Annual Report: Cloud Computing Section
#
# Grant: NSF-2024-12345 (Martinez, Biology)
# Report Year: 2025 (Year 2 of 3)
#
# Generating sections:
#   ✅ Budget vs. Actual (cloud allocation: $15,000)
#   ✅ Primary use cases (RNA-seq, ML training)
#   ✅ Researchers supported (1 postdoc, 3 grad students)
#   ✅ Publications (2 submitted, 1 accepted)
#   ✅ Data management (storage, retention, sharing)
#   ✅ Compliance (tagging, audit logs, cost accounting)
#
# Generated: NSF-2024-12345-Annual-Report-2025.pdf
#
# Key Findings:
#   - Cloud spending: $8,500 (57% of allocation) ✅ On track
#   - Spending is allowable (research-related) ✅
#   - No cost overruns ✅
#   - All spending properly tagged ✅
#   - Audit trail complete (7 years retention) ✅
#
# This report is ready for NSF submission.
# Estimated time saved: 6 hours (vs manual compilation)
```

### Federal Audit: Documentation Ready

```bash
# Federal auditor arrives, requests cloud spending documentation
# for NSF grant NSF-2024-12345

auditor: "I need to see all cloud spending for this grant, including
         proof that spending was allowable and allocable."

# Jennifer generates audit package
jennifer@laptop:~$ lfr admin audit-package NSF-2024-12345 \
  --format federal-audit

# Output:
# 📋 Generating Federal Audit Package
#
# Grant: NSF-2024-12345 (Martinez, Biology)
# Period: 2024-01-01 to present (21 months)
#
# Package Contents:
#
# 1. Budget Documentation
#    ✅ Original budget ($15,000 for cloud computing)
#    ✅ Budget modifications (none)
#    ✅ Current status ($8,500 spent, $6,500 remaining)
#
# 2. Spending Summary
#    ✅ Month-by-month breakdown
#    ✅ Researcher-by-researcher breakdown
#    ✅ Cost category breakdown (compute, storage, data transfer)
#
# 3. Allowability Evidence
#    ✅ All resources tagged with grant number
#    ✅ All usage by authorized grant personnel
#    ✅ All usage for research purposes (no personal use detected)
#    ✅ Publications citing grant (2 papers)
#
# 4. Allocability Evidence
#    ✅ No commingling with other grants (separate AWS account)
#    ✅ Cost allocation tags (100% compliance)
#    ✅ No shared resources across grants
#
# 5. Audit Logs
#    ✅ Complete access logs (who used what, when)
#    ✅ Resource creation logs
#    ✅ Cost allocation logs
#    ✅ All logs retained for 7 years (federal requirement)
#
# 6. Compliance Certifications
#    ✅ PI certification (proper use of funds)
#    ✅ Research admin review (quarterly)
#    ✅ IT security review (annual)
#
# 📦 Package created: NSF-2024-12345-Audit-Package.zip
#
# Contents:
#   - Summary-Report.pdf (executive summary)
#   - Detailed-Spending.xlsx (all transactions)
#   - Audit-Logs/ (directory with 21 months of logs)
#   - Certifications/ (signed documents)
#   - Publications/ (papers citing grant)
#
# Size: 45 MB
# Ready for auditor review
#
# Estimated time saved: 20 hours (vs manual compilation)

# Auditor reviews and is satisfied
auditor: "This is excellent documentation. Everything is properly
         tagged and traceable. Audit finding: Passed ✅"

# Jennifer thinks: "That was SO much easier than last time!"
```

### AWS Credits Management (Coordination)

```bash
# Jennifer works with Credits Manager to ensure optimal use
jennifer@laptop:~$ lfr admin credits-status

# Output:
# 💰 AWS Credits Status
#
# Total Credits: $75,000
#   - AWS Educate (expires 2026-06-30): $50,000
#   - AWS Research (expires 2027-12-31): $20,000
#   - Promotional (expires 2025-12-31): $5,000
#
# Allocated to Grants:
#   - CS305-Fall2025: $2,000 (teaching)
#   - Bio-Teaching: $5,000 (teaching)
#   - NSF-2024-12345: $10,000 (research)
#   - NIH-R01-2025: $15,000 (research)
#   - Research-Pool: $13,000 (various)
#
# Unallocated: $30,000
#
# ⚠️  Alert: $5,000 in promotional credits expire in 60 days
#    Recommendation: Allocate to active research projects
#
# Coordination with Credits Manager:
#   - Weekly sync on credit allocation
#   - Jennifer prioritizes grants needing credits
#   - Credits Manager executes allocations
#
# Grant Priority for Credits:
#   1. Grants ending soon (use budget before closeout)
#   2. Grants with high spending rates (relieve cash pressure)
#   3. New grants (help PIs get started)

# Jennifer recommends credit allocation
jennifer@laptop:~$ lfr admin recommend-credits-allocation

# Output:
# 💡 Credit Allocation Recommendations
#
# Based on grant portfolio analysis:
#
# 1. NSF-2023-99999 (Rodriguez) - URGENT
#    Current: Cash spending $1,367/month (too fast)
#    Recommendation: Allocate $5,000 in credits
#    Reason: High GPU usage, budget exhaustion risk
#    Impact: Reduces cash burn, extends budget runway
#
# 2. NIH-R01-2024-55555 (Wilson)
#    Grant ends: 30 days
#    Unspent budget: $12,000
#    Recommendation: Convert to credits for follow-on grant
#    Reason: Won't spend all cash before closeout
#
# 3. New Grant: DOE-2025-22222 (Jackson)
#    Starting: Next month
#    Recommendation: Allocate $3,000 in credits
#    Reason: Help PI get started quickly
#
# Send recommendations to Credits Manager? [y/N]: y
# ✅ Recommendations sent to credits@university.edu
```

---

## 📋 Feature Gap Analysis: Research Admin Needs

### Critical Missing Features (Institutional Adoption Blockers)

| Feature | Priority | Impact | Current State | Effort |
|---------|----------|--------|---------------|--------|
| **Multi-Grant Dashboard** | 🔴 Critical | No visibility across portfolio | None | High |
| **Workday Integration** | 🔴 Critical | Manual reconciliation (40 hrs/qtr) | None | High |
| **Petri Integration** | 🔴 Critical | University system disconnect | None | High |
| **Grant Lifecycle Management** | 🔴 Critical | Manual account setup | Basic | Medium |
| **Compliance Reporting** | 🔴 Critical | Federal audit prep (20+ hours) | None | High |
| **Budget Alert System** | 🟡 High | Late discovery of overspend | None | Medium |
| **Audit Package Generation** | 🟡 High | Manual doc compilation | None | Medium |
| **Cost Allocation Automation** | 🟡 High | Manual tagging, errors | Partial | Low |
| **PI Communication Tools** | 🟢 Medium | Email overload | None | Low |
| **Credits Coordination** | 🟢 Medium | Waste prevention | None | Low |

### Unique Research Admin Requirements

| Requirement | Needed | Priority |
|-------------|--------|----------|
| **Cross-grant visibility** | ✅ Dashboard for 50+ grants | Critical |
| **Workday/SAP export** | ✅ Direct integration | Critical |
| **Petri sync** | ✅ Bi-directional data flow | Critical |
| **Audit compliance** | ✅ Federal audit packages | Critical |
| **Grant tagging enforcement** | ✅ 100% compliance | High |
| **Budget alerts** | ✅ Proactive overspend detection | High |
| **PI reporting** | ✅ Self-service for PIs | Medium |

---

## 🎯 Priority Recommendations: Research Admin Persona

### Phase 1: Core Grant Management (Milestone 1)

**Target**: Research admins can manage grant portfolio efficiently

1. **Grant Lifecycle Commands** (2 weeks)
   - `lfr admin grant create` (setup new grant accounts)
   - `lfr admin grant close` (closeout at end)
   - Auto-tagging enforcement
   - Budget alert configuration

2. **Multi-Grant Dashboard** (2 weeks)
   - Portfolio view (all grants)
   - Spending by grant, dept, PI
   - Alert highlights
   - Drill-down to details

3. **Budget Monitoring** (1 week)
   - Burn rate calculations
   - Overspend alerts
   - PI notifications
   - Department rollups

**Estimated effort**: 5 weeks

### Phase 2: Compliance & Reporting (Milestone 2)

**Target**: Federal audit ready, compliance automated

4. **Compliance Reporting** (3 weeks)
   - NSF annual report generation
   - NIH progress report format
   - Federal audit packages
   - Export audit logs

5. **Cost Allocation Automation** (2 weeks)
   - Enforce tagging at resource creation
   - Validate tags against Workday
   - Detect untagged resources
   - Auto-remediation

**Estimated effort**: 5 weeks

### Phase 3: Systems Integration (Milestone 3)

**Target**: Seamless data flow with university systems

6. **Workday Integration** (3 weeks)
   - Export actual spending (CSV/API)
   - Map to chartstrings
   - Quarterly reporting automation
   - Budget vs actual comparison

7. **Petri Integration** (4 weeks)
   - SSO authentication
   - Daily cost sync
   - Budget balance pull
   - User access sync

**Estimated effort**: 7 weeks

---

## Success Metrics: Research Admin Perspective

### Time Savings
- ✅ **Grant account setup**: 2 hours → 3 minutes (98% reduction)
- ✅ **Quarterly reconciliation**: 40 hours → 30 minutes (99% reduction)
- ✅ **Federal audit prep**: 20 hours → 1 hour (95% reduction)
- **Total time saved**: ~250 hours/year per admin

### Compliance
- ✅ **Audit findings**: Historical issues → Zero findings
- ✅ **Tagging compliance**: 60% → 100%
- ✅ **Documentation completeness**: Partial → Complete

### Financial Management
- ✅ **Overspend detection time**: 30 days → Real-time
- ✅ **Budget variance**: ±15% → ±2%
- ✅ **Credits waste**: $20K/year → $0

### Institutional Impact
- ✅ **Faculty satisfaction**: "Easier to manage grants" (90%)
- ✅ **Audit pass rate**: 85% → 100%
- ✅ **Grant portfolio growth**: Enabled by efficient admin

---

## Next Steps

1. **Interview Research Admins**: 3-5 at different universities
2. **Map Workday APIs**: Integration architecture
3. **Define Petri Specs**: Data flow requirements
4. **Prototype Dashboard**: Multi-grant view mockup
5. **Implementation Plan**: Prioritize based on compliance requirements

**Estimated Timeline**:
- Phase 1 (Core Grant Management): 5 weeks
- Phase 2 (Compliance & Reporting): 5 weeks
- Phase 3 (Systems Integration): 7 weeks
- **Total**: ~17 weeks (4 months) to comprehensive research admin support

---

**Status**: Draft Walkthrough
**Persona**: Research Administrator (Grant & Budget Manager)
**Priority**: 🔴 Critical (enables institutional adoption at scale)
**Dependencies**: Multi-account architecture (AWS Orgs), Workday/Petri APIs, compliance frameworks (NSF, NIH), federal audit requirements
**Note**: This persona is essential for universities to adopt LFR Tools institutionally (vs individual PI adoption)

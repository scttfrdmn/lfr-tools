# LFR Tools - User Requirements

**Date**: October 20, 2025
**Status**: Authoritative
**Purpose**: Define core user personas, their requirements, and constraints

---

## Document Authority

This document is **AUTHORITATIVE** and defines:
1. **Who** uses LFR Tools (personas)
2. **What** they need (requirements)
3. **Why** they need it (pain points)
4. **How** we prioritize (critical vs nice-to-have)

All feature development should reference this document. Any feature not serving these personas should be questioned.

---

## Core Personas (Priority Order)

### 1. Professor (University Class Management) 🔴 CRITICAL
**Who**: Faculty managing 10-50 student cloud environments for courses
**Scale**: 35 students typical, 1-2 TAs, $840 semester budget ($24/student)
**Example**: Prof. Maria Rodriguez, CS 305 - Cloud Computing

**Primary Needs**:
- Bulk student management (import 35 students in minutes, not hours)
- Per-student budget tracking and enforcement ($24/student caps)
- Template restrictions (prevent GPU overspend by students)
- Semester lifecycle automation (setup → active → auto-cleanup)
- TA delegation (TAs need debug access, not full admin)

**Critical Pain Points**:
- ❌ Manual student setup (8 hours → should be 15 minutes)
- ❌ No budget visibility per student (discover overspend too late)
- ❌ No semester end automation (4-8 hour manual cleanup)
- ❌ TAs can't help students remotely (office hours inefficient)

**Success Criteria**:
- Setup time: 8 hours → 15 minutes (97% reduction)
- Budget overruns: 40% → 5% (enforcement + alerts)
- Semester end: 4-8 hours → 0 minutes (fully automated)

---

### 2. Student (First-Time Cloud User) 🔴 CRITICAL
**Who**: Undergraduate/graduate students using AWS for first time
**Scale**: Budget-conscious, anxious about costs, non-experts
**Example**: Emily Chen, junior CS major in CS 305

**Primary Needs**:
- **Budget confidence**: "Will this cost me money?" answered clearly
- Budget preview before starting (show impact: "This will cost $X, you'll have $Y left")
- Cost education (plain English: "You'll only pay for hours used, not 24/7")
- Self-service help (common issues: broken environment, budget questions)
- Auto-stop reminders (prevent forgotten instances)

**Critical Pain Points**:
- ❌ Budget anxiety (persistent fear of surprise bills)
- ❌ No cost preview (don't know if safe to start instance)
- ❌ No broken environment recovery (TAs must fix)
- ❌ No clear cost explanations (AWS pricing overwhelming)

**Success Criteria**:
- Budget anxiety: High → Eliminated ("I feel safe")
- Students within budget: 75% → 95%
- TA help requests: 15/week → 5/week (better self-service)

---

### 3. Teaching Assistant (Office Hours Support) 🔴 CRITICAL
**Who**: Graduate students helping 35 students debug environments
**Scale**: 10 hours/week office hours, need 5-10 minute resolutions
**Example**: Alex Chen, PhD student, Head TA for CS 305

**Primary Needs**:
- **TA role** (between student and admin - THIS DOESN'T EXIST YET!)
- Debug sessions (temporary SSH to student instances, logged)
- Instance reset (backup + restore to clean state, preserve files)
- Student visibility (dashboard of all 35 students)
- Academic integrity (all TA actions logged, auditable)

**Critical Pain Points**:
- ❌ **NO TA ACCESS AT ALL** (biggest gap in entire system!)
- ❌ Screen sharing only (45 min/student vs could be 5-10 min)
- ❌ Can't reset broken environments (students lose work)
- ❌ No visibility into student status (reactive, not proactive)

**Success Criteria**:
- Help time: 45 min → 8 min (82% reduction)
- Students helped/hour: 1.3 → 3.1 (138% increase)
- TA satisfaction: "Tools made TA job much easier" (100%)

**NOTE**: This persona is **completely missing** from current implementation. Highest priority after Professor support.

---

### 4. Researcher (Grant-Funded Computational Work) 🟡 HIGH
**Who**: Postdocs/faculty conducting grant-funded research
**Scale**: $3K-$50K cloud budgets, 6-36 month grants, paper deadlines
**Example**: Dr. Sarah Chen, postdoc on PI's NSF grant

**Primary Needs**:
- **Reproducibility** (snapshots for paper methods sections)
- Budget forecasting ("Will I have enough for final experiments?")
- Grant compliance (auto-tagging, NSF/NIH reports)
- Cost optimization (spot instances = 70% savings)
- Secure collaboration (share with co-authors, peer reviewers)

**Critical Pain Points**:
- ❌ No reproducibility tools (8 hours to document environment)
- ❌ No budget forecasting (surprise budget exhaustion)
- ❌ No spot instance support (waste 3x budget)
- ❌ Manual grant compliance (8+ hours for NSF report)

**Success Criteria**:
- Reproducibility package time: 8 hours → 5 minutes
- Budget forecast accuracy: N/A → ±5%
- Cost optimization: 0% → 70% savings (via spot)

---

### 5. Research Administrator (Grant & Budget Manager) 🟡 HIGH
**Who**: University staff managing 50+ grants, $25M portfolio
**Scale**: Institutional, cross-department, federal compliance
**Example**: Jennifer Park, Office of Sponsored Programs

**Primary Needs**:
- **Multi-grant dashboard** (visibility across 50+ active grants)
- **Workday/Petri integration** (university financial systems)
- Grant lifecycle automation (setup → monitor → closeout)
- Federal audit compliance (NSF/NIH documentation)
- Budget alerts (detect overspend before it's too late)

**Critical Pain Points**:
- ❌ No central visibility (manual tracking, 40 hrs/quarter)
- ❌ No Workday integration (manual reconciliation)
- ❌ No Petri integration (disconnected systems)
- ❌ Federal audit prep (20+ hours manual compilation)

**Success Criteria**:
- Quarterly reconciliation: 40 hours → 30 minutes (99% reduction)
- Federal audit prep: 20 hours → 1 hour (95% reduction)
- Overspend detection: 30 days → Real-time

**NOTE**: This persona enables **institutional adoption** at scale (vs individual PI adoption).

---

### 6. AWS Credits Manager (University IT) 🔴 CRITICAL
**Who**: IT staff managing $150K/year in AWS credits
**Scale**: 200+ AWS accounts, 3-5 credit sources, expiration dates
**Example**: David Kim, Cloud Infrastructure Manager

**Primary Needs**:
- **Credits visibility** (dashboard of all credits, sources, expirations)
- Expiration monitoring (alerts 90/60/30/7 days before waste)
- Allocation management (assign credits to grants/courses)
- Burn rate tracking (forecast if credits will be fully used)
- Cash-to-credits opportunities (find projects paying cash when credits available)

**Critical Pain Points**:
- ❌ $20K/year in expired credits (complete waste!)
- ❌ Manual tracking (20 hrs/month across 200 accounts)
- ❌ No burn rate forecasting (can't tell if underutilized)
- ❌ Cash spending while credits sit unused (suboptimal)

**Success Criteria**:
- Credits waste: $20K/year → $0/year (100% reduction)
- Time spent: 20 hrs/month → 2 hrs/month (90% reduction)
- Utilization: 73% → 97% (+24%)

**ROI**: $74K/year benefit (waste prevention + time savings)

---

## Cross-Cutting Requirements

### Funding Model (ALL PERSONAS)
**Constraint**: "Who pays for AWS?" must be answered for every usage

**Funding Sources**:
1. **Grant funding** (NSF, NIH, DOE) - Time-bounded, restricted use, audit requirements
2. **Departmental budgets** (IT allocations, course fees) - Annual budgets
3. **AWS credits** (Educate, Research, Promotional) - Expiring, free money
4. **Personal funding** (researcher credit cards) - Rare, discouraged

**Account Structure** (Recommended):
- Multi-account AWS Organizations
- One AWS account per grant (hard budget separation, audit-friendly)
- One AWS account per course
- Consolidated billing at org level
- Credits applied at org level, allocated to accounts

**System Integrations** (Required):
- **Workday/SAP**: University financial system (chartstrings)
- **Petri**: Research management portal (grant tracking)
- **Research.gov / eRA Commons**: Federal grant systems (NSF, NIH)
- **Canvas/Blackboard**: Learning management systems (student rosters)

---

## Compliance & Security (ALL PERSONAS)

### Technical Compliance
- **GDPR**: The ONLY technical compliance framework supported by Lightsail for Research
  - Data protection and privacy controls
  - User consent management
  - Right to erasure (data deletion)
  - Data portability
  - Security controls

### Academic Integrity
- **SSH key isolation**: Students can't share keys
- **TA access logging**: All debug sessions recorded
- **Audit trail**: 7-year retention for grant financial records
- **Plagiarism detection**: Activity logs for investigations

### Grant Financial Accountability (NOT Technical Compliance)
**Important**: NSF/NIH requirements are about **financial accounting**, not technical compliance frameworks.

- **Allowability**: All spending is research-related
- **Allocability**: No commingling of grant funds
- **Cost allocation**: 100% tagging with grant numbers
- **Financial documentation**: Audit-ready reports (NSF, NIH cost accounting formats)
- **Retention**: 7-year log retention for financial records

### Budget Controls
- **Hard caps**: Stop resources at budget limit (configurable)
- **Soft alerts**: Warn at 80%, 90%, 100%
- **Per-user budgets**: Sub-budgets for team members
- **Auto-stop policies**: Idle detection (save budget)

---

## Persona Priority Matrix

| Persona | Priority | Blocks | Effort (weeks) | ROI |
|---------|----------|--------|----------------|-----|
| **Professor** | 🔴 Critical | Educational use at scale | 13 | High - enables courses |
| **Student** | 🔴 Critical | User adoption, retention | 9 | High - user experience |
| **TA** | 🔴 Critical | Office hours efficiency | 11 | High - support quality |
| **Researcher** | 🟡 High | Research productivity | 12 | Medium - paper quality |
| **Research Admin** | 🟡 High | Institutional adoption | 17 | High - scaling |
| **Credits Manager** | 🔴 Critical | Waste prevention | 10 | Very High - $74K/year |

**Total Estimated Effort**: ~72 weeks (~18 months) for comprehensive support of all personas

**Phased Approach** (Recommended):
- **Phase 1** (6 months): Professor + Student + TA (core educational use case)
- **Phase 2** (6 months): Researcher + Credits Manager (research & efficiency)
- **Phase 3** (6 months): Research Admin (institutional scale)

---

## Feature Prioritization Framework

### 🔴 Critical (Blockers)
Features without which a persona **cannot use the system**:
- TA role & permissions (TA persona completely blocked)
- Budget tracking per student (Professor can't manage class)
- Budget impact preview (Student won't use due to anxiety)
- Credits dashboard (Credits Manager can't prevent waste)

### 🟡 High (Major Pain Points)
Features that cause significant inefficiency or frustration:
- Semester end automation (Professor wastes 4-8 hours)
- Instance reset (TA/Student struggle with broken environments)
- Reproducibility snapshots (Researcher wastes 8 hours/paper)
- Workday integration (Research Admin wastes 40 hours/quarter)

### 🟢 Medium (Nice to Have)
Features that improve experience but have workarounds:
- Jupyter integration (can SSH manually)
- Mobile dashboard (can use laptop)
- Cost visualization graphs (can read tables)

---

## Constraints & Assumptions

### Technical Constraints
1. **AWS Lightsail for Research** is the primary compute service
2. **AWS Organizations** is available for multi-account structure
3. **IAM** is used for user management (not Cognito)
4. **No AWS SSO** (university uses own SSO/Shibboleth)

### Organizational Constraints
1. **University IT approval** required for account creation
2. **Research admin involvement** required for grant accounts
3. **Federal compliance** must be maintained (NSF, NIH, A-133)
4. **FERPA compliance** for student data

### Budget Constraints
1. **Limited IT budgets** - cost efficiency critical
2. **Grant restrictions** - spending must be allowable/allocable
3. **Credits expiration** - use-it-or-lose-it pressure
4. **No surprise bills** - students/PIs need predictability

### Time Constraints
1. **Semester deadlines** - setup must happen in 1-2 weeks before classes
2. **Paper deadlines** - researchers need quick reproducibility
3. **Grant closeout** - end-of-grant cleanup must be automated
4. **Office hours** - TAs need 5-10 min resolutions, not 45 min

---

## Non-Goals (Out of Scope)

### What LFR Tools Is NOT
1. ❌ Full AWS management tool (only LFR-focused)
2. ❌ Multi-cloud (GCP, Azure) - AWS only
3. ❌ On-premises infrastructure management
4. ❌ Learning management system (not Canvas replacement)
5. ❌ Grant writing tool (not proposal software)
6. ❌ Financial system (not Workday replacement)

### What LFR Tools Does NOT Do
1. ❌ Create AWS accounts automatically (requires IT approval)
2. ❌ Manage university user authentication (SSO handled externally)
3. ❌ Handle non-research AWS usage (production apps, websites)
4. ❌ Replace Petri (integrates with, doesn't replace)

---

## Success Metrics (Overall)

### Adoption
- ✅ Professors using LFR: 0 → 50+ (2 years)
- ✅ Students onboarded: 0 → 2,000+ (2 years)
- ✅ Grants managed: 0 → 100+ (2 years)

### Efficiency
- ✅ Setup time: Hours → Minutes (98% reduction)
- ✅ Management overhead: 20% → 5% (75% reduction)
- ✅ Support burden: High → Low (TAs empowered)

### Financial
- ✅ Budget control: 60% → 95% within budget
- ✅ Credits waste: $20K+/year → $0 (100% reduction)
- ✅ Cost optimization: 0% → 40% (spot, auto-stop)

### Quality
- ✅ Reproducibility: 0% → 100% (all research)
- ✅ Compliance: Risky → Audit-ready
- ✅ User satisfaction: N/A → 90%+ ("would recommend")

---

## Revision History

- **2025-10-20**: Initial version, 6 personas documented
- **Future**: Update based on user interviews, pilot deployments

---

**For Questions**: Refer to persona walkthroughs in `docs/USER_SCENARIOS/`
**For Architecture**: See `.github/FUNDING_MODEL_ARCHITECTURE.md`
**For Development**: See `.github/PROJECT_ALIGNMENT_ANALYSIS.md`

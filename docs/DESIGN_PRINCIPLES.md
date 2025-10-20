# LFR Tools - Design Principles

**Date**: October 20, 2025
**Status**: Authoritative
**Purpose**: Define core philosophy and decision-making framework

---

## Mission Statement

**LFR Tools makes AWS Lightsail for Research accessible, affordable, and accountable for educational and research use.**

We believe that:
- Cloud computing should be accessible to students and researchers of all skill levels
- Budget anxiety should never prevent learning or research
- Compliance and accountability should be automatic, not burdensome
- Infrastructure management should take minutes, not hours

---

## Core Principles

### 1. Educational Mission First 🎓

**Principle**: Optimize for teaching and learning, not production infrastructure.

**What this means**:
- **Students are not DevOps engineers** - interfaces must be beginner-friendly
- **Professors are not cloud admins** - setup should be automated
- **TAs are educators, not sysadmins** - support should be efficient
- **Budget is limited** - waste prevention is critical

**Design Decisions**:
- ✅ Plain-English error messages (not AWS jargon)
- ✅ Budget preview before every action (reduce anxiety)
- ✅ Auto-stop policies by default (prevent waste)
- ✅ Student-friendly documentation (14-year-old reading level)
- ❌ No assumption of cloud expertise
- ❌ No complex IAM policy editing
- ❌ No manual cost tracking required

**Example**:
```bash
# Bad (production tool):
$ aws ec2 run-instances --image-id ami-0c55b159cbfafe1f0 \
  --instance-type t3.large --count 1 --key-name mykey \
  --security-group-ids sg-0123456789abcdef0

# Good (educational tool):
$ lfr ssh connect emily --project CS305-Fall2025
# → Auto-starts instance, connects, shows budget
```

---

### 2. Budget Accountability 💰

**Principle**: Every dollar must be tracked, justified, and optimized.

**What this means**:
- **Grant funding** - Compliance with NSF, NIH, DOE requirements
- **Course budgets** - IT departments need accountability
- **Credits** - Use-it-or-lose-it, track expiration
- **Personal budgets** - No surprise bills for students

**Design Decisions**:
- ✅ Mandatory tagging (grant number, project, user)
- ✅ Per-user sub-budgets (enforce limits)
- ✅ Real-time budget tracking (not monthly surprises)
- ✅ Audit trails (7-year retention for federal compliance)
- ✅ Budget forecasting (proactive, not reactive)
- ❌ No untagged resources allowed
- ❌ No "pay later, track later" approach
- ❌ No hidden costs

**Example**:
```bash
# Before launching instance:
💰 Budget Impact:
   Cost: $0.83/day ($108/month if 24/7)
   Your budget: $18.50 / $24.00 (77%)
   After session: $18.61 / $24.00 ✅

Proceed? [y/N]:
```

---

### 3. Reproducible Science 🔬

**Principle**: Research must be reproducible. Period.

**What this means**:
- **Papers** - Methods sections must be accurate and complete
- **Peer review** - Reviewers must be able to verify results
- **Collaboration** - Co-authors must access exact environments
- **Compliance** - Federal grants require reproducibility

**Design Decisions**:
- ✅ Environment snapshots (capture everything)
- ✅ Software inventory (all versions documented)
- ✅ DOI minting (citable via Zenodo)
- ✅ One-click reproduction (reviewers can launch identical env)
- ✅ Methods section generation (auto-generate from logs)
- ❌ No "it worked on my machine" excuses
- ❌ No manual environment documentation
- ❌ No lost reproducibility information

**Example**:
```bash
$ lfr research snapshot rnaseq-analysis \
  --name "Chen2025-RNAseq" \
  --paper "Nature Communications submission"

# Generates:
# - AMI snapshot
# - Software manifest (all versions)
# - Dockerfile
# - Methods section text
# - Zenodo DOI: 10.5281/zenodo.XXXXXX
```

---

### 4. Minimal Cognitive Load 🧠

**Principle**: Users should think about research/teaching, not infrastructure.

**What this means**:
- **Students** - Focus on learning, not AWS billing
- **Professors** - Focus on teaching, not user management
- **Researchers** - Focus on science, not DevOps
- **Admins** - Focus on strategy, not manual tracking

**Design Decisions**:
- ✅ Sensible defaults (auto-stop, tagging, regions)
- ✅ Progressive disclosure (advanced features hidden until needed)
- ✅ Contextual help (right information, right time)
- ✅ Automated workflows (setup, monitoring, cleanup)
- ✅ Self-healing (auto-recovery from common failures)
- ❌ No 20-flag commands
- ❌ No "read the manual" requirements for basic use
- ❌ No manual checklist maintenance

**Example**:
```bash
# Minimal command, maximum automation:
$ lfr ssh connect emily --project CS305-Fall2025

# Behind the scenes (automatic):
# - Checks budget
# - Starts instance if stopped
# - Configures SSH keys
# - Opens tunnel if needed
# - Tags with project
# - Logs activity
# - Reminds to stop when done
```

---

### 5. Institutional Integration 🏛️

**Principle**: LFR Tools must fit into existing university systems, not replace them.

**What this means**:
- **Workday/SAP** - University financial systems are authoritative
- **Petri** - Research management systems track grants
- **Canvas/Blackboard** - LMS systems manage student rosters
- **SSO/Shibboleth** - Authentication is centralized
- **AWS Organizations** - University IT controls accounts

**Design Decisions**:
- ✅ API integrations (push/pull data, don't duplicate)
- ✅ SSO authentication (no separate user database)
- ✅ Chartstring mapping (university codes)
- ✅ Export formats (Workday CSV, NSF reports)
- ❌ No custom user management (use university SSO)
- ❌ No duplicate financial tracking (sync with Workday)
- ❌ No proprietary lock-in (open APIs)

**Example**:
```bash
# Export to Workday for quarterly reconciliation:
$ lfr admin export-workday \
  --period Q3-2025 \
  --format csv

# Output: Q3-AWS-Spending.csv
# Columns: Chartstring, Grant, Date, Amount, Description
# → Direct import to Workday Financial Management
```

---

### 6. Accountability by Default ✅

**Principle**: Financial accountability and data protection should be automatic, not a checkbox.

**What this means**:
- **GDPR compliance** - The ONLY technical compliance framework supported by Lightsail for Research
- **Grant financial accountability** (NSF, NIH) - Cost accounting and audit trails for federal grants
- **Academic integrity** - All actions logged, auditable
- **Data retention** - 7 years automatic (for grant financial records)
- **Cost allocation** - 100% tagged and traceable
- **FERPA** - Student data protected

**Important Distinction**:
- **Technical Compliance Framework**: GDPR only (data protection, privacy, security controls)
- **Financial Accountability**: NSF/NIH grant requirements (cost accounting, allowability, allocability)
- LFR Tools provides **financial tracking and reporting**, not technical compliance certification

**Design Decisions**:
- ✅ Auto-tagging enforced (can't create untagged resources)
- ✅ Audit logs automatic (S3, encrypted, retained)
- ✅ Financial reports generated (NSF, NIH formats for cost accounting)
- ✅ Separation of concerns (grants don't commingle)
- ✅ GDPR-compliant data handling
- ❌ No manual tagging required
- ❌ No "remember to document this" prompts
- ❌ No accountability as afterthought

**Example**:
```bash
# Grant financial documentation request:
$ lfr admin audit-package NSF-2024-12345

# Generates (automatic):
# - Budget summary
# - Spending breakdown (allowability evidence)
# - Cost allocation tags (all tagged correctly)
# - Allocability evidence (no grant commingling)
# - Access logs (7 years retention)
# - Publications citing grant
# → Ready for federal grant audit in 1 minute
```

---

### 7. Progressive Complexity 📈

**Principle**: Simple things should be simple, complex things should be possible.

**What this means**:
- **Beginners** - Can get started in minutes with defaults
- **Intermediate** - Can customize as they learn
- **Advanced** - Can access full power when needed
- **Experts** - Can integrate with other tools

**Complexity Layers**:
1. **Student**: `lfr ssh connect emily` (1 command, everything else automatic)
2. **Professor**: Bulk operations, budget management
3. **Researcher**: Snapshots, spot instances, collaboration
4. **Admin**: Multi-grant dashboard, Petri integration, compliance

**Design Decisions**:
- ✅ Guided workflows (interactive wizards)
- ✅ Sensible defaults (works out of box)
- ✅ Progressive disclosure (advanced flags hidden)
- ✅ Multiple interfaces (CLI, GUI, API)
- ❌ No "expert mode only"
- ❌ No forcing simple users to understand complexity
- ❌ No hiding power from advanced users

---

### 8. Fast Feedback Loops ⚡

**Principle**: Users should know immediately if something is wrong.

**What this means**:
- **Budget overages** - Alert before it happens, not after
- **Credits expiration** - Warn 90/60/30/7 days in advance
- **Instances running** - Remind to stop, don't wait for bill
- **Errors** - Clear messages with solutions, not cryptic codes

**Design Decisions**:
- ✅ Real-time budget tracking (not monthly surprises)
- ✅ Proactive alerts (prevent issues)
- ✅ Actionable error messages (tell user how to fix)
- ✅ Immediate feedback (confirm actions)
- ❌ No "check back in 30 days"
- ❌ No silent failures
- ❌ No error codes without explanations

**Example**:
```bash
# Proactive alert (before problem):
⚠️  Your instance has been running for 6 hours.
   Did you forget to stop it?

   Cost so far: $0.45
   If left running: $0.83/day

   Stop now? [y/N]: y

# vs Reactive (too late):
💸 Your bill last month: $450
   (You left instance running for 30 days)
```

---

### 9. Graceful Degradation 🛡️

**Principle**: Failures should be safe and recoverable.

**What this means**:
- **Data loss** - Never lose student/researcher work
- **Broken environments** - Easy reset, files preserved
- **Budget overages** - Soft warnings before hard stops
- **System failures** - Degraded service better than no service

**Design Decisions**:
- ✅ Auto-backup before destructive operations
- ✅ Soft budget limits (warn at 80%, 90%, hard stop at 100%)
- ✅ Instance reset preserves user files
- ✅ API retries with exponential backoff
- ❌ No "delete without backup"
- ❌ No hard budget stops without warning
- ❌ No data loss from system failures

**Example**:
```bash
# TA resets student's broken environment:
$ lfr ta reset-instance emily-ubuntu

# Before reset:
# ✅ Backup to S3: s3://backups/emily-20251015.tar.gz
# ✅ Preserve /home/emily/ (student work)
# ❌ Discard broken system state
# ✅ Launch fresh from template
# ✅ Restore student files
# → Student's work safe, environment fixed
```

---

### 10. Measurable Impact 📊

**Principle**: Everything should be measurable and improvable.

**What this means**:
- **Time savings** - Quantify efficiency gains
- **Cost savings** - Track optimization impact
- **User satisfaction** - Measure experience
- **Compliance** - Prove audit readiness

**Design Decisions**:
- ✅ Metrics built-in (usage, costs, efficiency)
- ✅ Dashboards for all personas
- ✅ ROI calculable (time + cost saved)
- ✅ A/B testing capability (try new approaches)
- ❌ No "feels faster" - measure it
- ❌ No anecdotal evidence only
- ❌ No unmeasured features

**Example**:
```bash
# Research admin quarterly report:
$ lfr admin report Q3-2025

# Metrics:
# - Time saved: 40 hrs/quarter (reconciliation) → 30 min
# - Cost saved: $20K (credits waste) → $0
# - Audit readiness: 20 hrs → 1 hr (95% reduction)
# - User satisfaction: 4.6/5.0 (90% would recommend)
# → Clear ROI: $50K time savings + $20K waste prevention
```

---

## Decision Framework

When making design decisions, ask:

### 1. Does this serve our personas?
- ✅ If yes: Which persona? What pain point?
- ❌ If no: Why are we building it?

### 2. Does this reduce cognitive load?
- ✅ If yes: How much simpler is it?
- ❌ If no: Can we automate it instead?

### 3. Does this improve accountability?
- ✅ If yes: Is it audit-ready?
- ❌ If no: What compliance risk does it introduce?

### 4. Does this scale institutionally?
- ✅ If yes: Can research admin use it for 50+ grants?
- ❌ If no: Does it only work for 1-2 users?

### 5. Does this have measurable impact?
- ✅ If yes: What metrics improve? By how much?
- ❌ If no: How will we know if it works?

---

## Anti-Patterns (What NOT to Do)

### ❌ Don't Build Production Tools
LFR Tools is for **education and research**, not production workloads.
- No high-availability guarantees
- No 24/7 SLA requirements
- No enterprise support tiers

### ❌ Don't Assume Expertise
Users are **students and researchers**, not AWS experts.
- No IAM policy editing required
- No VPC configuration knowledge assumed
- No CloudFormation templates to write

### ❌ Don't Hide Costs
Budget transparency is **non-negotiable**.
- No "contact us for pricing"
- No hidden fees
- No surprise bills

### ❌ Don't Create Lock-In
Data and workflows should be **portable**.
- Export capabilities everywhere
- Standard formats (CSV, JSON, PDF)
- Open APIs
- No proprietary data formats

### ❌ Don't Duplicate Systems
Integrate, don't replace **existing university systems**.
- Not a financial system (Workday exists)
- Not an LMS (Canvas exists)
- Not a grant system (Research.gov exists)
- Not an SSO provider (university SSO exists)

---

## Trade-offs & Tensions

### Simplicity vs Power
**Tension**: Beginners want simple, experts want control.
**Resolution**: Progressive complexity - simple default, power available.

### Automation vs Control
**Tension**: Automation helps beginners, experts want manual override.
**Resolution**: Automate with escape hatches - can override when needed.

### Security vs Convenience
**Tension**: Security wants restrictions, users want ease.
**Resolution**: Secure by default, educational exceptions with logging.

### Cost vs Performance
**Tension**: Budget limits vs research needs.
**Resolution**: Budget visibility first, then cost optimization (spot, etc).

### Speed vs Compliance
**Tension**: Fast iteration vs audit requirements.
**Resolution**: Compliance built-in, not bolted-on - no compromise.

---

## Success Stories (Examples)

### Before LFR Tools:
- Prof. Rodriguez: 8 hours to set up 35 students
- Emily (student): Constant budget anxiety, avoided using AWS
- Alex (TA): 45 minutes per student debug (screen sharing)
- Jennifer (admin): 40 hours/quarter reconciliation
- David (IT): $20K in expired credits annually

### After LFR Tools:
- Prof. Rodriguez: 15 minutes to set up 35 students (97% faster)
- Emily: "I feel safe" - budget preview before every action
- Alex: 8 minutes per student debug (82% faster, SSH access)
- Jennifer: 30 minutes/quarter reconciliation (99% faster)
- David: $0 in expired credits (100% waste prevention)

---

## Revision History

- **2025-10-20**: Initial version
- **Future**: Update based on user feedback, pilot results

---

**These principles guide all development decisions. When in doubt, refer back to this document.**

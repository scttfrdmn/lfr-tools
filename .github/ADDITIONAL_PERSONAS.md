# Additional Critical Personas - Research Enablement

**Date**: October 20, 2025
**Status**: Identified, walkthroughs pending

## Newly Identified Personas

Based on user feedback, these additional personas are critical for research enablement and have NOT been captured by existing tools:

### 1. University Research Administrator / Grant Manager

**Role**: Manages grants, allocates budgets, integrates with institutional systems

**Key Responsibilities**:
- **Grant Management**: Track grant-funded research budgets across departments
- **Budget Allocation**: Distribute AWS credits/budgets to PIs and labs
- **Petri Integration**: Hookup to university's Petri research portal/system
- **Compliance**: Ensure spending aligns with grant restrictions
- **Reporting**: Generate reports for sponsored programs office, funding agencies

**Pain Points** (Hypothesis):
- No visibility into AWS spending across multiple grants
- Manual tracking of which AWS resources belong to which grant
- Can't reallocate budgets between grants/projects
- No integration with Petri or institutional research systems
- Difficult to prove grant spending to auditors

**Critical Features Needed**:
- Grant budget tracking (multiple grants, multiple PIs)
- Cross-project visibility (university-wide or college-wide)
- Budget reallocation (move unused budget between grants)
- Petri integration (research portal connection)
- Audit reports (grant-compliant spending reports)
- Multi-tenant hierarchy (university → college → department → lab)

**Priority**: 🟡 High (enables institutional adoption)

**Dependencies**:
- Understanding of Petri system architecture
- University IT integration requirements
- Grant compliance requirements (NSF, NIH, etc.)

---

### 2. AWS Credits Manager

**Role**: Manages AWS credits, tracks usage, reallocates between accounts/projects

**Key Responsibilities**:
- **Credit Tracking**: View available AWS credits across org
- **Credit Allocation**: Assign credits to specific accounts/projects
- **Credit Monitoring**: Track burn rate, expiration dates
- **Credit Reallocation**: Move credits between accounts (if top-level)
- **Credit Reporting**: Who used what credits, when

**Pain Points** (Hypothesis):
- AWS credits expire if not used (wasteful)
- No easy way to see credit balance
- Can't move credits from underutilized projects to active ones
- No alerts when credits are about to expire
- Difficult to track which project used which credits

**Critical Features Needed**:
- View AWS credits balance (total, per account, per project)
- Credit allocation commands (assign credits to projects)
- Credit reallocation (move between accounts if permitted by AWS)
- Credit burn rate tracking ($ per day, projected expiration)
- Credit expiration alerts (30 days, 7 days, 1 day warnings)
- Credit usage reports (per project, per grant, per PI)

**Priority**: 🔴 Critical (enables efficient use of educational/research credits)

**Dependencies**:
- AWS Credits API understanding
- AWS Organizations credit sharing model
- Credit transfer limitations (what AWS allows)

---

## Integration Points

### Petri Integration
**What is Petri**: [Need to research - appears to be university research management system]

**Potential Integration Points**:
- User authentication (SSO via Petri)
- Grant budget import (pull grant allocations from Petri)
- Spending export (push AWS spending to Petri for reporting)
- Project creation (auto-create LFR projects from Petri research projects)
- Compliance tagging (tag AWS resources with Petri grant IDs)

### AWS Credits System
**AWS Credits Background**:
- Educational credits: AWS Educate, AWS Academy credits
- Research credits: AWS Cloud Credits for Research program
- Promotional credits: Various AWS promotional programs
- Credits have expiration dates (typically 1-2 years)
- Credits can be at org level or account level
- Some credits can be reallocated, others cannot

**Integration Needs**:
- Query AWS Billing API for credit balance
- Track credit expiration dates
- Monitor credit usage vs. remaining balance
- Alert on expiration risk
- Report credit usage by project/grant

---

## Updated Persona Priority Matrix

| Persona | Current Status | Priority | Blocks | Estimated Effort |
|---------|----------------|----------|--------|------------------|
| **Professor** | Partial walkthrough | 🔴 Critical | Educational use at scale | 13 weeks |
| **TA** | Walkthrough complete | 🔴 Critical | Office hours efficiency | 11 weeks |
| **Student** | Walkthrough complete | 🔴 Critical | User adoption | 9 weeks |
| **Research Admin** | Not started | 🟡 High | Institutional adoption | 8-12 weeks |
| **AWS Credits Manager** | Not started | 🔴 Critical | Credit waste prevention | 6-8 weeks |
| **Researcher** | Not started | 🟡 High | Individual research use | 6 weeks |

---

## Next Steps

1. **Research Petri System**: Understand architecture, APIs, integration points
2. **Research AWS Credits API**: Understand credit query, allocation, limitations
3. **Create Research Admin Walkthrough**: Detailed scenario with Petri integration
4. **Create AWS Credits Manager Walkthrough**: Detailed credit management workflow
5. **Create Researcher Walkthrough**: Individual PI or grad student use case
6. **Update GitHub Issues**: Add issues for Petri integration, credit management

---

## Questions to Answer

### Petri Integration
- What is Petri exactly? (Research portal? Budget system? Both?)
- What APIs does Petri expose?
- How do other tools integrate with Petri?
- What data needs to flow between LFR Tools and Petri?
- Is Petri university-specific or multi-tenant?

### AWS Credits
- Can credits be queried via API? (AWS Billing API)
- Can credits be reallocated between accounts?
- What credit management tools exist?
- How do we track credit expiration?
- Can we set alerts for expiring credits?

### Research Admin Use Cases
- What reports do sponsored programs need?
- How are grants typically structured? (PI → Co-PI → Students?)
- What compliance requirements exist? (NSF, NIH, DOE, etc.)
- How do admins currently track AWS spending?
- What other tools do research admins use?

---

**Status**: Draft - Personas Identified
**Next**: Create detailed walkthroughs for Research Admin, AWS Credits Manager, and Researcher
**Owner**: Project maintainer

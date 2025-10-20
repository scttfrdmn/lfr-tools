# LFR Tools - GitHub Issues Summary

**Created**: 2025-10-20
**Total New Issues**: 21 (issues #8-28)
**Status**: Ready for implementation

---

## Critical Priority (Blocks Persona Usage) - 5 issues

### Phase 1 Critical
1. **#8** - TA role and permissions system
   - **BLOCKS**: Entire TA persona (no TA access at all)
   - Persona: Teaching Assistant
   - Estimate: 4 weeks

2. **#9** - Real-time budget impact preview for students
   - **BLOCKS**: Student adoption (budget anxiety barrier)
   - Persona: Student
   - Estimate: 2 weeks

3. **#10** - Per-student budget tracking and enforcement
   - **BLOCKS**: Professor course management at scale
   - Persona: Professor
   - Estimate: 3 weeks

4. **#11** - AWS credits dashboard and expiration tracking
   - **BLOCKS**: $20K/year waste prevention
   - Persona: Admin (Credits Manager)
   - Estimate: 3 weeks
   - ROI: $74K/year benefit

5. **#23** - Funding model architecture (AWS accounts + grants)
   - **BLOCKS**: All financial tracking (foundational)
   - Persona: Admin
   - Estimate: 4 weeks
   - **NOTE**: Should be implemented FIRST (dependency for #10, #11, others)

---

## High Priority (Major Pain Points) - 11 issues

### Phase 1 High
6. **#12** - Bulk student management operations
   - Pain: 8 hours → 15 minutes (97% reduction)
   - Persona: Professor
   - Estimate: 2 weeks

7. **#26** - Template restrictions for courses
   - Pain: Students overspend on GPU instances
   - Persona: Professor
   - Estimate: 1 week

8. **#14** - TA student dashboard
   - Depends on: #8 (TA role)
   - Persona: Teaching Assistant
   - Estimate: 2 weeks

9. **#15** - TA instance reset with backup
   - Depends on: #8 (TA role)
   - Persona: Teaching Assistant, Student
   - Estimate: 2 weeks

10. **#17** - Auto-stop reminders for students
    - Pain: Prevents biggest cause of budget overruns
    - Persona: Student
    - Estimate: 1 week

11. **#16** - Semester end automation
    - Pain: 4-8 hours → 0 minutes (100% reduction)
    - Persona: Professor
    - Estimate: 2 weeks

12. **#4** - Idle detection policy templates (existing)
    - Persona: Professor
    - Estimate: TBD

13. **#3** - EFS integration documentation (existing)
    - Persona: Professor
    - Estimate: TBD

### Phase 2 High
14. **#13** - Research reproducibility snapshots
    - Pain: 8 hours → 5 minutes (documentation)
    - Persona: Researcher
    - Estimate: 3 weeks

15. **#18** - Grant financial reporting (NSF/NIH formats)
    - Pain: 8+ hours → 5 minutes
    - Persona: Researcher, Admin
    - Estimate: 2 weeks

16. **#19** - Research admin multi-grant dashboard
    - Enables: Institutional adoption at scale
    - Persona: Admin (Research Administrator)
    - Estimate: 4 weeks

17. **#20** - Workday/SAP financial system integration
    - Pain: 40 hrs/quarter → 30 min (99% reduction)
    - Persona: Admin
    - Estimate: 3 weeks

18. **#28** - Federal grant audit package generation
    - Pain: 20+ hours → 1 hour (95% reduction)
    - Persona: Admin
    - Estimate: 2 weeks

---

## Medium Priority (Nice to Have) - 5 issues

### Phase 1 Medium
19. **#25** - Cost education and plain-English explanations
    - Improves: Student confidence
    - Persona: Student
    - Estimate: 1 week

### Phase 2 Medium
20. **#21** - Spot instance support for cost optimization
    - Benefit: 70% cost savings
    - Persona: Researcher
    - Estimate: 2 weeks

21. **#22** - Budget forecasting
    - Benefit: Proactive budget management
    - Persona: Researcher, Professor
    - Estimate: 2 weeks

22. **#24** - Petri research management integration
    - Benefit: University-wide research coordination
    - Persona: Admin
    - Estimate: 3 weeks

23. **#27** - Secure research collaboration and sharing
    - Benefit: Safer co-author/reviewer access
    - Persona: Researcher
    - Estimate: 2 weeks

### Phase 2/3 Medium
24. **#5** - AWS Cost Explorer integration (existing)
    - Persona: Professor
    - Estimate: TBD

25. **#6** - Hardware-tied educational tokens (existing)
    - Phase: 3
    - Persona: Student
    - Estimate: TBD

---

## Recommended Implementation Order (Phase 1 Focus)

### Sprint 1 (4 weeks) - FOUNDATIONAL
1. **#23** - Funding model architecture (4 weeks)
   - **MUST be first** - blocks all financial features

### Sprint 2 (4 weeks) - CRITICAL BLOCKERS (Part 1)
2. **#8** - TA role and permissions (4 weeks)
   - Unblocks TA persona completely

### Sprint 3 (6 weeks) - CRITICAL BLOCKERS (Part 2)
3. **#11** - Credits dashboard (3 weeks)
   - High ROI ($74K/year)
4. **#10** - Per-student budget tracking (3 weeks, parallel with #11)

### Sprint 4 (2 weeks) - CRITICAL BLOCKERS (Part 3)
5. **#9** - Student budget preview (2 weeks)

### Sprint 5 (6 weeks) - HIGH PRIORITY PROFESSOR FEATURES
6. **#12** - Bulk student operations (2 weeks)
7. **#26** - Template restrictions (1 week)
8. **#16** - Semester end automation (2 weeks)
9. **#4** - Idle detection templates (1 week)

### Sprint 6 (4 weeks) - HIGH PRIORITY TA FEATURES
10. **#14** - TA dashboard (2 weeks)
11. **#15** - TA instance reset (2 weeks)

### Sprint 7 (2 weeks) - HIGH PRIORITY STUDENT FEATURES
12. **#17** - Auto-stop reminders (1 week)
13. **#25** - Cost education (1 week, can parallelize)

**Total Phase 1 Estimate**: ~24 weeks (6 months)

After Phase 1, move to Phase 2 (researcher, research admin features).

---

## Key Insights from Persona Analysis

1. **TA persona is completely missing** - highest gap (#8 blocks everything)
2. **Budget anxiety blocks student adoption** - critical to address (#9)
3. **Funding model is foundational** - must implement first (#23)
4. **High ROI opportunities**:
   - Credits management: $74K/year benefit
   - Bulk operations: 8 hours → 15 minutes (97% reduction)
   - Audit packages: 20 hours → 1 hour (95% reduction)
   - Quarterly reconciliation: 40 hours → 30 minutes (99% reduction)
5. **Institutional adoption blocked** without research admin features (#19, #20, #28)

---

## Documentation Created

All documentation committed to main branch:

1. **Persona Walkthroughs** (~6,000 lines total):
   - `docs/USER_SCENARIOS/01_UNIVERSITY_CLASS_PROFESSOR_WALKTHROUGH.md`
   - `docs/USER_SCENARIOS/02_STUDENT_FIRST_TIME_USER_WALKTHROUGH.md`
   - `docs/USER_SCENARIOS/03_TEACHING_ASSISTANT_WALKTHROUGH.md`
   - `docs/USER_SCENARIOS/04_RESEARCHER_GRANT_FUNDED_WALKTHROUGH.md`
   - `docs/USER_SCENARIOS/05_RESEARCH_ADMIN_GRANT_MANAGER_WALKTHROUGH.md`
   - `docs/USER_SCENARIOS/06_AWS_CREDITS_MANAGER_WALKTHROUGH.md`

2. **Core Documentation**:
   - `docs/USER_REQUIREMENTS.md` - Authoritative persona/requirements doc
   - `docs/DESIGN_PRINCIPLES.md` - Decision framework and philosophy

3. **Architecture Documentation**:
   - `.github/FUNDING_MODEL_ARCHITECTURE.md` - AWS account structure design
   - `.github/PROJECT_ALIGNMENT_ANALYSIS.md` - Gap analysis
   - `.github/ADDITIONAL_PERSONAS.md` - Persona discovery notes

4. **This Summary**:
   - `.github/GITHUB_ISSUES_SUMMARY.md` - Issue prioritization and roadmap

---

## Next Steps

1. **Review**: Validate priorities with stakeholders
2. **Refine**: Adjust estimates based on team capacity
3. **Implement**: Start with Sprint 1 (#23 - Funding Model)
4. **Iterate**: Reassess after each sprint based on feedback

---

**Note**: All compliance references have been corrected to clarify that Lightsail for Research supports GDPR as the ONLY technical compliance framework. NSF/NIH requirements are about financial accountability (cost accounting), not technical compliance.

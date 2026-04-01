# framework-review

## Purpose
Review current migration changes against architecture, phase boundaries, and compatibility requirements.

This skill is for review only.
It should not make code edits unless the caller explicitly asks for fixes after the review.

## Required context
Before reviewing, read:
- `AGENTS.md`
- `docs/framework-migration-plan.md`
- the active phase/task objective from the current thread or repo file

## Review checklist
Always check:
1. phase discipline
   - did the changes stay inside the requested phase?
   - did execution spill into adjacent tasks/phases?

2. framework boundary discipline
   - does `framework/pkg` expose only contracts/abstractions?
   - are implementations kept in `framework/internal`?
   - are implementation types leaking across the boundary?

3. dependency discipline
   - does framework-facing code still depend directly on Fiber, Redis, SQL/Ent, or other internal libraries?
   - were adapters introduced where needed?

4. repository-role discipline
   - is `api/` still functioning as application/reference implementation?
   - was `fe/` kept compatible or were risks reported clearly?

5. compatibility discipline
   - are routes, payloads, auth semantics, and runtime behavior preserved unless explicitly planned otherwise?

## Output format
Use exactly these headings:
- Review Scope
- Passed Checks
- Blocking Issues
- Non-Blocking Issues
- Phase Violations
- Boundary Violations
- Compatibility Risks
- Final Verdict

## Verdict rules
Use one of only three verdicts:
- PASS
- PASS WITH RISKS
- FAIL

If FAIL:
- list the smallest blocking fixes needed
- do not propose unrelated improvements

## Never do
- do not rewrite the plan during review unless explicitly asked
- do not suggest large redesigns unless the current changes are fundamentally incompatible with the migration plan
- do not fix code in the same response unless explicitly requested

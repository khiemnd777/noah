# framework-orchestrator

## Purpose
Run the full repository transformation from application-monorepo shape toward extracted framework shape with minimal supervision.

This skill is the orchestration entrypoint.
It should plan, analyze, execute, review, fix safe issues, checkpoint progress, and continue phase by phase.

## Required context
Before orchestration, read:
- `AGENTS.md`
- `docs/framework-migration-plan.md`
- any existing checkpoint files in `docs/` if they exist

## Core operating model
Run the migration as a staged pipeline:
1. load migration source of truth
2. determine current phase status
3. plan current phase
4. use subagents for analysis where useful
5. execute one task at a time
6. review after each meaningful task or small task group
7. auto-fix only safe, local issues
8. checkpoint progress
9. continue until the phase is complete
10. move to next phase only after phase validation passes

## Subagent policy
You may delegate bounded analysis work to these roles when useful:
- analyzer
- architect
- migration-planner
- contract-guard
- reviewer
- tester

Subagents should primarily analyze and report.
Do not let multiple subagents perform overlapping code edits.
Prefer centralized execution in the main thread.

## Auto-fix policy
You may auto-fix only when all are true:
- the fix is local and clearly scoped
- it does not introduce a new abstraction beyond the approved phase
- it does not expand the execution scope materially
- it does not create a breaking FE/API contract change

If the issue is architectural, ambiguous, or cross-phase, checkpoint and stop instead of improvising.

## Checkpoint policy
After each phase, write or update a checkpoint file in `docs/`.
Recommended file name pattern:
- `docs/framework-migration-checkpoint.md`

The checkpoint should include:
- completed phases/tasks
- current repo state summary
- validations passed/failed
- remaining risks
- next intended phase/task

## Final report policy
At the end of the full run, produce a final report covering:
- completed phases
- major structural changes
- public API surface created
- adapters introduced
- compatibility verification summary
- remaining risks or deferred items

## Repo-specific migration policy
Respect the transformation strategy for this repo:
- extract a framework layer without destroying the current application too early
- treat `api/` as the consumer/reference implementation during migration
- keep `fe/` compatible throughout
- avoid simultaneous broad refactors across `api/` and `fe/`

## Stop conditions
Stop and checkpoint instead of guessing when:
- the migration plan is missing or contradictory
- a phase boundary becomes ambiguous
- a change would require destructive rewrite without validation
- compatibility risk to `fe/` cannot be assessed safely
- repeated review failures indicate the plan needs adjustment

## Output style during orchestration
Keep interim messages concise and operational.
For each completed task, record:
- task completed
- files changed
- validation result
- blockers if any

## Never do
- do not merge phases into one large implementation burst
- do not silently bypass reviews
- do not delete working application code before replacement is proven
- do not expose implementation details through `framework/pkg`

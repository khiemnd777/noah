# framework-plan

## Purpose
Turn one migration phase into an execution-ready task list for Codex.

This skill is for planning only.
It must not edit code, create implementation files, or move packages.

## Required inputs
The caller should provide:
- the target phase name
- optional subphase or focus area
- any additional constraints for that run

## Required context
Before planning, read:
- `AGENTS.md`
- `docs/framework-migration-plan.md`

If those files conflict, follow `AGENTS.md` first.

## Planning contract
For the requested phase, produce:
1. phase objective
2. in-scope packages and files
3. out-of-scope packages and files
4. ordered atomic tasks
5. dependencies between tasks
6. expected risks per task
7. validation gate per task
8. phase-completion validation gate

## Atomic task rules
A task is valid only if it:
- touches one coherent concern only
- stays within one phase
- can be reviewed independently
- can be reverted independently
- keeps the system runnable after completion

Do not combine multiple concerns into one task.

## Repo-specific migration policy
This repository currently contains:
- `api/`
- `fe/`
- `docs/`

The transformation goal is to extract a reusable framework while preserving the application during migration.

For planning purposes, treat the repo as:
- `framework/` = target reusable core
- `api/` = current backend app that will become a consumer/reference implementation
- `fe/` = frontend consumer that must remain compatible

## Phase policy
Use only these phases unless the caller explicitly narrows one further:
1. Reframe / Repository boundary definition
2. Framework skeleton creation
3. Contract extraction
4. Adapter construction
5. API migration to framework abstractions
6. Framework runtime introduction
7. Stabilization and compatibility verification

If `docs/framework-migration-plan.md` uses a different label set, map your output to that file but preserve the same staged discipline.

## Output format
Use exactly these headings:
- Phase
- Objective
- In Scope
- Out of Scope
- Ordered Tasks
- Task Dependencies
- Risks
- Validation Gates
- Stop Conditions

## Stop conditions
Planning must stop and report a blocker instead of inventing assumptions when:
- the requested phase depends on missing prior phases
- the phase boundary is ambiguous
- the requested work would require cross-phase execution
- the repo structure no longer matches the migration plan

## Never do
- do not edit code
- do not create files other than plan/checkpoint files if explicitly requested
- do not mix planning with implementation
- do not silently expand scope beyond the requested phase

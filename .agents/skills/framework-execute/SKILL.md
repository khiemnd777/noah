# framework-execute

## Purpose
Execute one approved migration task safely.

This skill is for implementation of exactly one task at a time.
It must not continue to the next task automatically unless the caller explicitly asks.

## Required inputs
The caller should provide:
- phase name
- task identifier or exact task description
- any execution constraints

## Required context
Before executing, read:
- `AGENTS.md`
- `docs/framework-migration-plan.md`
- the approved phase/task plan if one exists in the current thread or a repo file

## Execution contract
Before editing code, restate:
1. the exact task being executed
2. the exact in-scope files/packages
3. the exact out-of-scope files/packages
4. the validation checks that must pass before the task is considered done

Then execute only that task.

## Safety rules
- do not cross phase boundaries
- do not silently execute adjacent tasks
- do not perform broad cleanup unrelated to the task
- keep `api/` runnable
- do not break FE/API contract unless the task explicitly includes a compatible contract update
- prefer introducing framework code before replacing application code

## Repo-specific migration policy
Treat the repo as temporarily dual-mode:
- `framework/` is being extracted as reusable core
- `api/` remains the running application/reference implementation during migration
- `fe/` remains a consumer and compatibility target

Never prematurely delete or invalidate application behavior just because the framework target exists.

## Expected implementation style
- extract contracts before moving implementations
- move implementations behind adapters before removing direct dependency usage
- keep changes small and reviewable
- preserve compile/runtime continuity wherever practical

## Required output after execution
Use exactly these headings:
- Executed Task
- Files Changed
- What Changed
- Validation Performed
- Risks / Unverified Areas
- Recommended Next Task

## Validation expectations
At minimum, perform the nearest practical validation for the task, such as:
- compile/build for affected package(s)
- targeted tests for affected package(s)
- route/contract sanity checks if HTTP code changed
- import-boundary checks if `pkg/` or `internal/` changed

If validation cannot be run, state exactly what remains unverified.

## Never do
- do not continue automatically into the next task
- do not rewrite unrelated packages
- do not move to later phases because it seems convenient
- do not expose implementation types through `framework/pkg`

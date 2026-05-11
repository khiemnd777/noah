---
name: noah-identity-contract
description: Use when a Noah task touches staff, user, department, department administrator assignment, staff route params, or identity mapping so users.id and staffs.id semantics stay correct.
---

# Noah Identity Contract

Use this skill when a task touches staff, user, department, department administrator assignment or unassignment, staff route params, staff DTOs, or identity mapping.

## Required reading

1. `/AGENTS.md`
2. `docs/identity-contract.md`
3. the owning handler, service, repository, frontend API wrapper, mapper, schema, table, or widget for the affected flow

## Identity invariants

- `users.id` is the account identity.
- `staffs.id` is the staff record identity.
- `staffs.user_staff` points to `users.id`.
- `departments.administrator_id` stores `users.id`, never `staffs.id`.
- frontend staff DTO identity currently uses `users.id` unless an explicitly approved contract change says otherwise.

## No shortcut patch rule

- trace the real identity owner before editing
- do not accept both `users.id` and `staffs.id` in one endpoint unless the user explicitly approves compatibility mode
- do not hide identity drift with mapper-only fallbacks, dual lookup fallbacks, ambiguous route params, or hardcoded aliases
- if the correct fix crosses handler, service, repository, DTO, mapper, permission, or UI boundaries, update those layers coherently

## Required checks

- verify whether each route param is `users.id` or `staffs.id`; never infer from a param named `id`
- verify department admin assignment and unassignment read and write `users.id`
- verify DTOs and mappers preserve the canonical ID contract
- verify permission checks use the correct caller and target identity domains
- verify variable names make identity explicit:
  - use `userID` for `users.id`
  - use `staffRecordID` for `staffs.id`
  - avoid `staffID` unless it truly means `staffs.id`

## Output expectations

Before editing, state:

- which identity domain each touched flow uses
- whether any route param is ambiguous
- whether the task changes the contract or only fixes usage of the existing contract
- the validation plan for assignment, unassignment, mapper, and permission behavior touched by the change

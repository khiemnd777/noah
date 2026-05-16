---
name: noah-fast-triage
description: Use when the user provides a Noah issue, keyword, screenshot, UI text, route, API endpoint, log line, field name, or suspected module and wants the agent to quickly find the owning files, classify the boundary, and propose a verified change plan before editing.
---

# Noah Fast Triage

Use this skill to quickly turn user-provided signals into an evidence-backed file map and change proposal.

This skill is for triage before implementation. It does not override the root `AGENTS.md` confirmation gate: do not edit code, config, tests, workflows, runtime assets, or UI until the user approves the proposed change set.

## Source Of Truth

For Noah-managed agent artifacts, use `agents/**` as the editable source of truth.

- Skills live under `agents/skills/<skill-name>/`.
- Subagents live under `agents/subagents/<agent-name>.toml`.
- `.codex/**` is repo-local runtime or local repo-specific artifact space, not the canonical source for Noah-managed artifacts.
- Before creating or editing Noah-managed skills/subagents, read `agents/skills/README.md`.

## Input Signals

Classify all user-provided signals before searching:

- issue title/body or review comment
- screenshot/image
- UI text, table column, button label, error message, or i18n copy
- route/page URL
- API endpoint or network request
- log line, error text, stack trace, or status code
- DTO/model/schema/field name
- migration/database table or column
- suspected module, file, or feature name

If a signal is ambiguous, label it as ambiguous. Do not present guesses as facts.

## Search Order

Prefer exact, high-signal searches before broad exploration.

1. Search exact quoted strings from the user.
2. Search visible UI text and i18n/copy keys.
3. Search route paths and route params.
4. Search API endpoint paths and HTTP method names.
5. Search DTO/model/schema/field names.
6. Search component, table, form, widget, mapper, or API-client names.
7. Search backend handler/service/repository names.
8. Search migration/database table or column names when persistence may be involved.
9. Search nearby tests only after the likely owner is identified.

Use `rg` and `rg --files` first. Keep searches scoped once ownership becomes clear.

## Screenshot Triage

When the user provides a screenshot or image:

- Extract visible text, labels, buttons, table headers, route clues, and error messages.
- Identify the likely surface type: page, modal, form, table, detail view, nav, or toast.
- Search extracted text before searching visual guesses.
- Treat visual identification as a hypothesis unless a file, route, or copy string confirms it.
- Do not infer hidden business rules from the screenshot.

For UI changes, still follow root `AGENTS.md`: read `DESIGN.md`, draw the intended layout in text/ASCII, and wait for explicit approval before editing.

## Ownership Mapping

Map the likely owning layer before proposing any edit.

For frontend signals, trace:

```text
route -> feature index/registration -> page -> widget/component -> api client -> mapper/model/schema
```

For backend signals, trace:

```text
route/endpoint -> handler/controller -> service/use-case -> repository -> DTO/schema/model
```

For FE/API contract signals, inspect both sides:

```text
backend DTO/handler/service/repository
frontend API wrapper/mapper/model/UI consumer
permission/error/invalidation behavior
```

For unclear ownership, read `docs/module-inventory.md` before expanding the search.

For staff/user/department identity signals, read `docs/identity-contract.md` before inspecting implementation files.

## Required Output Before Editing

Before any implementation, produce a concise triage report:

```text
Signals provided:
- ...

Evidence found:
- keyword/route/file matches with paths

Likely owner:
- frontend/backend/cross-boundary/platform-impact
- owning module/feature

File map:
- files inspected
- files likely requiring changes
- files intentionally not touched

Root-cause hypothesis:
- label as hypothesis unless proven by code path

Proposed change set:
- exact files/layers to update
- why each layer is included

Verification plan:
- behavior-specific checks
- nearest tests/build/lint/typecheck if relevant
- known unverified gaps or residual risk
```

Ask for explicit user approval before editing.

## Guardrails

- Start from evidence, not guesses.
- Do not scan the whole repo without a concrete reason.
- Do not stop at the first plausible match if the owning layer is not proven.
- Do not propose shortcut patches, guard-only suppressions, hardcoded fallbacks, or symptom-only fixes.
- If the correct fix crosses handler/service/repository, mapper/API/UI, auth, registry, cache, migration, or job boundaries, classify it as cross-boundary and include each affected layer.
- If the task is only a question or review, do not propose edits unless the user asks for implementation.

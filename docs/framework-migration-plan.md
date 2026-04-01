# Noah Framework Migration Plan

## Purpose
This document is the source of truth for transforming this repository from an application-oriented monorepo into a framework-oriented repository while preserving the running application during migration.

## Current Repository Shape
- `api/` — current backend application
- `fe/` — current frontend application
- `docs/` — planning, checkpoints, reports

## Target Repository Shape
- `framework/` — reusable framework core
  - `framework/pkg/` — stable public contracts and abstractions
  - `framework/internal/` — concrete implementations and runtime wiring
- `api/` — backend consumer / reference implementation during and after migration
- `fe/` — frontend consumer that must remain compatible

## Migration Rules
- Migrate incrementally, never as a destructive rewrite.
- Keep `api/` runnable throughout the migration.
- Preserve FE/API compatibility unless an explicit compatible contract migration is planned.
- Extract contracts before moving implementations.
- Move implementations behind adapters before removing direct dependency usage.
- Validate each phase before moving to the next one.

---

## Phase 1 — Reframe / Boundary Definition
### Goals
- Establish clear repository roles for framework, api, and fe.
- Prevent accidental treatment of `api/` as the framework itself.

### Scope
- Planning and documentation only.
- Optional creation of top-level `framework/` directory skeleton if explicitly approved.

### Risks
- Misaligned mental model causing destructive refactors later.

### Validation Gate
- The migration framing is documented and unambiguous.
- `framework` target role, `api` role, and `fe` role are explicit.

---

## Phase 2 — Framework Skeleton Creation
### Goals
- Create `framework/` directory structure.
- Establish `framework/pkg` and `framework/internal` boundaries.

### Scope
- Create empty or minimal package scaffolding only.
- No broad migration of existing implementation code.

### Risks
- Premature code movement before contracts exist.

### Validation Gate
- `framework/` exists with a clear boundary structure.
- No existing runtime behavior is changed.

---

## Phase 3 — Contract Extraction
### Goals
- Define the initial public abstractions for the framework.
- Identify the smallest stable contracts needed first.

### Candidate Contracts
- App
- Module
- Router
- Cache
- DB
- Auth
- Lifecycle

### Scope
- Contract/interface creation.
- Small supporting types if required.
- No broad implementation moves.

### Risks
- Over-designing contracts too early.
- Leaking implementation details into public packages.

### Validation Gate
- Initial framework contracts exist in `framework/pkg/`.
- Public contracts do not expose Fiber, Redis, SQL/Ent, or similar implementation details.

---

## Phase 4 — Adapter Construction
### Goals
- Wrap existing implementation technologies behind framework abstractions.

### Candidate Adapters
- Fiber → router/http adapter
- Redis → cache adapter
- DB access layer → db adapter
- Auth runtime → auth adapter

### Scope
- Introduce concrete implementations under `framework/internal/`.
- Keep migration local and reversible.

### Risks
- Partial wrapping that still leaves direct dependency leaks.
- Creating adapters that are too thin to be useful or too broad to remain stable.

### Validation Gate
- Target adapters exist under `framework/internal/`.
- Adapters map cleanly to public contracts.
- Direct dependency leakage is reduced in the targeted areas.

---

## Phase 5 — API Migration to Framework Abstractions
### Goals
- Convert `api/` from direct-implementation usage to framework-consumer usage.
- Migrate one concern at a time.

### Recommended Order
1. HTTP/router boundary
2. Cache boundary
3. DB boundary
4. Auth boundary
5. Boot/runtime composition boundary

### Scope
- Replace direct technology usage in `api/` with framework abstractions incrementally.
- Preserve runtime behavior.

### Risks
- Cross-layer breakage from broad edits.
- Inconsistent hybrid state if too many concerns are migrated at once.

### Validation Gate
- Targeted `api/` layer uses framework contracts instead of direct implementation dependency.
- `api/` still builds and remains operationally coherent.

---

## Phase 6 — Framework Runtime Introduction
### Goals
- Move core application boot/runtime behavior into framework-managed abstractions.
- Introduce `App` and module system in usable form.

### Scope
- Framework runtime composition.
- Controlled migration of boot logic from `api/` into framework runtime.

### Risks
- Boot/runtime regressions.
- Leaking old app-specific assumptions into framework runtime.

### Validation Gate
- The framework can construct and run the application through its abstractions.
- `api/` acts as a consumer/reference implementation rather than a hand-wired core.

---

## Phase 7 — Stabilization and Compatibility Verification
### Goals
- Ensure framework boundary discipline.
- Confirm FE/API compatibility.
- Confirm public API surface is coherent and stable enough for continued use.

### Scope
- Review, cleanup, documentation, compatibility checks, targeted fixes.

### Risks
- Undetected contract drift.
- Hidden implementation leaks in `framework/pkg/`.
- FE breakage from backend contract changes.

### Validation Gate
- `framework/pkg/` exposes only public contracts and abstractions.
- `framework/internal/` owns implementations.
- `api/` remains a working consumer/reference implementation.
- `fe/` compatibility risks are either resolved or explicitly documented.

---

## Execution Policy
- Execute one task at a time.
- Review after each meaningful task or small task group.
- Auto-fix only local, low-risk issues.
- Write checkpoints in `docs/framework-migration-checkpoint.md`.

## Final Deliverables
- `framework/` directory with stable initial contract layer
- internal adapters and runtime components
- migrated `api/` consuming framework abstractions
- compatibility review for `fe/`
- checkpoint and final migration report

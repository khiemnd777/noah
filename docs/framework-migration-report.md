# Framework Migration Report

## Summary
This run established a reusable `framework/` layer and migrated the lowest-risk backend runtime boundaries into it while keeping `api/` as the runnable application and `fe/` unchanged. The migration is intentionally incremental: the gateway now boots through framework application abstractions, cache access flows through framework cache contracts, runtime registry persistence flows through framework lifecycle contracts, and auth token pairs now come from framework public contracts.

## Structural Changes
- Added root workspace wiring with `go.work`.
- Added a new Go module at `framework/`.
- Added public framework contracts:
  - `framework/pkg/app`
  - `framework/pkg/auth`
  - `framework/pkg/cache`
  - `framework/pkg/context`
  - `framework/pkg/db`
  - `framework/pkg/http`
  - `framework/pkg/lifecycle`
  - `framework/pkg/module`
- Added internal framework implementations:
  - `framework/internal/http/fiber`
  - `framework/internal/cache/redis`
  - `framework/internal/db/factory`
  - `framework/internal/db/driver`
  - `framework/internal/lifecycle/jsonstore`
- Added public runtime composition helpers:
  - `framework/runtime/app.go`
  - `framework/runtime/cache.go`
  - `framework/runtime/db.go`
  - `framework/runtime/lifecycle.go`
- Updated `api/go.mod` to consume the local framework module through `replace`.
- Migrated backend entrypoints and shims:
  - `api/main.go`
  - `api/gateway/main.go`
  - `api/gateway/runtime/start.go`
  - `api/shared/auth/model.go`
  - `api/shared/cache/cache.go`
  - `api/shared/cache/invalidate.go`
  - `api/shared/runtime/registry.go`
- Migrated HTTP-facing API seams:
  - `api/shared/app/*`
  - `api/shared/middleware/*`
  - `api/gateway/proxy/*`
  - `api/modules/*/handler/*` signatures

## Phase Validation Status
- Phase 1, Extraction: Passed
  - Repository role split is now explicit in code and docs.
- Phase 2, Decoupling: Passed for cache and lifecycle storage; partial for HTTP boot
  - `api/` consumes framework contracts in the migrated areas.
- Phase 3, Internalization: Passed for the introduced framework implementations
  - Implementations now exist under `framework/internal`.
- Phase 4, Framework Core Introduction: Partial
  - Core abstractions exist, but Module/DB/Auth are not yet fully adopted by `api/`.
- Phase 5, Stabilization: Partial
  - Full backend test suite passed, but remaining direct technology coupling still exists outside the migrated boundaries.

## Remaining Risks
- `framework/pkg/db` still cannot be the primary API boundary for `api/` until raw `*sql.DB` dependencies are refactored behind framework-level contracts.
- Handler implementations still reference some Fiber-native response constants and low-level request helpers even though their exposed seams now use framework HTTP contracts.
- Redis is still used directly in some backend packages outside `api/shared/cache`.
- Module bootstrap still exposes concrete runtime dependencies and should be moved behind a framework module/runtime abstraction in a follow-up phase.

## Validation Performed
- `cd framework && go test ./...`
- `cd api && GOCACHE="$PWD/.gocache" go test ./...`

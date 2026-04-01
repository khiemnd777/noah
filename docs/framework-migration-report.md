# Framework Migration Report

## Summary
This run established a reusable `framework/` layer and migrated the lowest-risk backend runtime boundaries into it while keeping `api/` as the runnable application and `fe/` unchanged. The migration is intentionally incremental: the gateway now boots through framework application abstractions, cache access flows through framework cache contracts, runtime registry persistence flows through framework lifecycle contracts, and auth token pairs now come from framework public contracts.
This follow-up consolidation moved remaining framework-worthy DB and Redis/pubsub ownership out of `api/shared/*` infrastructure packages: the public DB contract no longer leaks `database/sql`, Redis instance management now lives in framework runtime/internal adapters, module startup now configures framework cache runtime for framework-owned services, and API wrappers remain only as compatibility shims.
This phase also moved the built-in observability and attribute modules into framework ownership. Their handlers, services, repositories, models, module wiring, and module config now live under `framework/modules/*`, while the corresponding `api/modules/*` trees remain only as composition bridges for auth/middleware, startup, and config mapping.
The latest migration batch removed those composition bridges for the low-complexity built-ins: `attribute`, `folder`, and `profile` now boot directly from `framework/modules/*` through a framework-owned runtime bootstrap path, and their API-side module directories have been deleted.

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
  - `framework/runtime/config.go`
  - `framework/runtime/db.go`
  - `framework/runtime/lifecycle.go`
  - `framework/runtime/module_process.go`
  - `framework/runtime/auth_context.go`
- Added the first framework-owned built-in module:
  - `framework/modules/observability/*`
- Added the second framework-owned built-in module:
  - `framework/modules/attribute/*`
- Added the third and fourth framework-owned built-in modules:
  - `framework/modules/folder/*`
  - `framework/modules/profile/*`
- Updated `api/go.mod` to consume the local framework module through `replace`.
- Migrated backend entrypoints and shims:
  - `api/main.go`
  - `api/gateway/main.go`
  - `api/gateway/runtime/start.go`
  - `api/shared/auth/model.go`
  - `api/shared/cache/cache.go`
  - `api/shared/cache/invalidate.go`
  - `api/shared/module/bootstrapper.go`
  - `api/shared/runtime/registry.go`
- Migrated HTTP-facing API seams:
  - `api/shared/app/*`
  - `api/shared/middleware/*`
  - `api/gateway/proxy/*`
  - `api/modules/*/handler/*` signatures
- Migrated module discovery/startup helpers so `framework/modules/*` is the primary discovery/config root while `api/modules/*` can remain composition entrypoints during transition.
- Removed API module ownership for:
  - `api/modules/attribute`
  - `api/modules/folder`
  - `api/modules/profile`

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
- Built-in module startup still depends on API-owned composition entrypoints for coexistence-mode boot, so module process startup is not fully framework-owned yet.
- Medium, high, and very-high complexity built-ins still depend on API-owned RBAC, pubsub, storage, Ent schema, and search/custom-field helpers that have not yet been extracted into stable framework contracts.
- Observability still relies on API-provided RBAC permission bridging for auth enforcement, which is acceptable for coexistence mode but should become a framework auth contract before broader built-in module migration.
- Attribute still uses module-local raw SQL inside `framework/modules/attribute/repository` because the shared Ent schema is not yet separable from broader application entities; that is acceptable for the current DB boundary, but it should eventually be replaced by a framework-owned persistence adapter or a module-local extracted data layer.

## Validation Performed
- `cd framework && go test ./...`
- `cd api && GOCACHE="$PWD/.gocache" go test ./...`
- `cd framework && GOCACHE="$PWD/.gocache" go test ./runtime ./modules/observability/...`
- `cd api && GOCACHE="$PWD/.gocache" go test ./modules/observability/... ./shared/runtime ./scripts/module_runner/runner ./gateway/...`
- `cd framework && GOCACHE="$PWD/.gocache" go test ./runtime ./modules/attribute/...`
- `cd api && GOCACHE="$PWD/.gocache" go test ./modules/attribute/... ./shared/module`
- `cd framework && GOCACHE="$PWD/.gocache" go test ./runtime ./modules/attribute/... ./modules/folder/... ./modules/profile/...`
- `cd api && GOCACHE="$PWD/.gocache" go test ./shared/runtime ./scripts/module_runner/runner ./gateway/...`
- `cd api && GOCACHE="$PWD/.gocache" go run scripts/module_runner/main.go sync`

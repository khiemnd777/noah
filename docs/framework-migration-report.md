# Framework Migration Report

## Summary
This run established a reusable `framework/` layer and migrated the lowest-risk backend runtime boundaries into it while keeping `api/` as the runnable application and `fe/` unchanged. The migration is intentionally incremental: the gateway now boots through framework application abstractions, cache access flows through framework cache contracts, runtime registry persistence flows through framework lifecycle contracts, and auth token pairs now come from framework public contracts.
This follow-up consolidation moved remaining framework-worthy DB and Redis/pubsub ownership out of `api/shared/*` infrastructure packages: the public DB contract no longer leaks `database/sql`, Redis instance management now lives in framework runtime/internal adapters, module startup now configures framework cache runtime for framework-owned services, and API wrappers remain only as compatibility shims.
This phase also moved the built-in observability and attribute modules into framework ownership. Their handlers, services, repositories, models, module wiring, and module config now live under `framework/modules/*`, while the corresponding `api/modules/*` trees remain only as composition bridges for auth/middleware, startup, and config mapping.
The latest migration batch removed those composition bridges for the low-complexity built-ins: `attribute`, `folder`, and `profile` now boot directly from `framework/modules/*` through a framework-owned runtime bootstrap path, and their API-side module directories have been deleted.
This cleanup batch extracted the remaining legacy shared-layer runtime ownership from `api/shared/*` into `framework/runtime`: module registry/discovery, port/process helpers, circuit breaker, retry-aware HTTP wrappers, HTTP client, JWT/auth context helpers, pubsub, and generic runtime/path/config helpers now live in the framework, while `api/shared/*` keeps only compatibility shims or app-local behavior.
This final built-in migration batch completed the module boundary move: every built-in module now lives under `framework/modules/*`, the corresponding `api/modules/*` built-in trees have been removed, and only `api/modules/main` remains application-owned. To preserve behavior while finishing the cutover, the remaining shared support used by those modules was internalized under `framework/internal/legacy/shared/*`, so framework modules no longer import `api/shared/*` directly.
The final shared-layer cutover is now complete: the entire legacy shared surface moved from `api/shared/*` to `framework/shared/*`, all application imports and generator paths were repointed to framework ownership, and `api/shared` has been deleted.

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
  - `framework/runtime/circuitbreaker.go`
  - `framework/runtime/cache.go`
  - `framework/runtime/config.go`
  - `framework/runtime/db.go`
  - `framework/runtime/http_client.go`
  - `framework/runtime/http_helpers.go`
  - `framework/runtime/http_params.go`
  - `framework/runtime/lifecycle.go`
  - `framework/runtime/module_registry.go`
  - `framework/runtime/module_process.go`
  - `framework/runtime/pubsub.go`
  - `framework/runtime/system.go`
  - `framework/runtime/auth_context.go`
- Added the first framework-owned built-in module:
  - `framework/modules/observability/*`
- Added the second framework-owned built-in module:
  - `framework/modules/attribute/*`
- Added the third and fourth framework-owned built-in modules:
  - `framework/modules/folder/*`
  - `framework/modules/profile/*`
- Added the remaining framework-owned built-in modules:
  - `framework/modules/auditlog/*`
  - `framework/modules/auth/*`
  - `framework/modules/metadata/*`
  - `framework/modules/notification/*`
  - `framework/modules/photo/*`
  - `framework/modules/rbac/*`
  - `framework/modules/realtime/*`
  - `framework/modules/search/*`
  - `framework/modules/token/*`
  - `framework/modules/user/*`
- Added internalized framework-owned legacy support packages required by the final module migration:
  - `framework/internal/legacy/shared/*`
- Added exported framework-owned shared packages for the final consumer cutover:
  - `framework/shared/*`
- Updated `api/go.mod` to consume the local framework module through `replace`.
- Migrated backend entrypoints and shims:
  - `api/main.go`
  - `api/gateway/main.go`
  - `api/gateway/runtime/start.go`
  - `api/gateway/registry/loader.go`
  - `api/gateway/proxy/proxy.go`
  - `api/scripts/module_runner/runner/runner.go`
  - `api/shared/auth/model.go`
  - `api/shared/cache/cache.go`
  - `api/shared/cache/invalidate.go`
  - `api/shared/circuitbreaker/cb.go`
  - `api/shared/pubsub/pubsub.go`
  - `api/shared/module/bootstrapper.go`
  - `api/shared/runtime/registry.go`
  - `api/shared/app/*`
  - `api/shared/utils/auth.go`
  - `api/shared/utils/config.go`
  - `api/shared/utils/fiber_user_utils.go`
  - `api/shared/utils/jwtutil.go`
  - `api/shared/utils/module_utils.go`
  - `api/shared/utils/net.go`
  - `api/shared/utils/param_utils.go`
  - `api/shared/utils/token_context.go`
  - `api/shared/utils/util.go`
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
- Removed API module ownership for the remaining built-ins:
  - `api/modules/auditlog`
  - `api/modules/auth`
  - `api/modules/metadata`
  - `api/modules/notification`
  - `api/modules/observability`
  - `api/modules/photo`
  - `api/modules/rbac`
  - `api/modules/realtime`
  - `api/modules/search`
  - `api/modules/token`
  - `api/modules/user`

## Phase Validation Status
- Phase 1, Extraction: Passed
  - Repository role split is now explicit in code and docs.
- Phase 2, Decoupling: Passed for cache and lifecycle storage; partial for HTTP boot
  - `api/` consumes framework contracts in the migrated areas, and legacy shared runtime ownership has been reduced to compatibility facades for the migrated concerns.
- Phase 3, Internalization: Passed for the introduced framework implementations
  - Implementations now exist under `framework/internal`.
- Phase 4, Framework Core Introduction: Partial
  - Core abstractions exist, and runtime ownership has moved further into `framework/runtime`, but `api/shared/module` and some middleware composition still retain application-specific bootstrap behavior.
- Phase 5, Stabilization: Partial
  - Full backend test suite passed, but remaining direct technology coupling still exists outside the migrated boundaries.
- Final built-in module ownership cutover: Passed
  - Built-in module discovery, compilation, and application-side boot now resolve from `framework/modules/*` only.
- Final shared-layer ownership cutover: Passed
  - `api/shared` no longer exists, `api` imports `framework/shared/*`, and both Go modules still build and test successfully.

## Remaining Risks
- Handler implementations still reference some Fiber-native response constants and low-level request helpers even though their exposed seams now use framework HTTP contracts.
- Redis is still used directly in some backend packages outside `api/shared/cache`.
- Built-in module startup still depends on API-owned composition entrypoints for coexistence-mode boot, so module process startup is not fully framework-owned yet.
- Medium, high, and very-high complexity built-ins still depend on API-owned RBAC, pubsub, storage, Ent schema, and search/custom-field helpers that have not yet been extracted into stable framework contracts.
- Observability still relies on API-provided RBAC permission bridging for auth enforcement, which is acceptable for coexistence mode but should become a framework auth contract before broader built-in module migration.
- Attribute still uses module-local raw SQL inside `framework/modules/attribute/repository` because the shared Ent schema is not yet separable from broader application entities; that is acceptable for the current DB boundary, but it should eventually be replaced by a framework-owned persistence adapter or a module-local extracted data layer.
- `framework/shared/*` and `framework/internal/legacy/shared/*` now overlap intentionally; the next cleanup pass should collapse that duplication as final runtime ownership is normalized.

## Validation Performed
- `cd framework && go test ./...`
- `cd framework && GOCACHE="$PWD/.gocache" go test ./shared/...`
- `cd framework && GOCACHE="$PWD/.gocache" go test ./runtime/...`
- `cd api && GOCACHE="$PWD/.gocache" go test ./...`
- `cd api && GOCACHE="$PWD/.gocache" go run scripts/module_runner/main.go sync`
- `cd api && GOCACHE="$PWD/.gocache" go test ./shared/... ./gateway/... ./scripts/module_runner/runner`
- `cd framework && GOCACHE="$PWD/.gocache" go test ./runtime ./modules/observability/...`
- `cd api && GOCACHE="$PWD/.gocache" go test ./modules/observability/... ./shared/runtime ./scripts/module_runner/runner ./gateway/...`
- `cd framework && GOCACHE="$PWD/.gocache" go test ./runtime ./modules/attribute/...`
- `cd api && GOCACHE="$PWD/.gocache" go test ./modules/attribute/... ./shared/module`
- `cd framework && GOCACHE="$PWD/.gocache" go test ./runtime ./modules/attribute/... ./modules/folder/... ./modules/profile/...`
- `cd api && GOCACHE="$PWD/.gocache" go test ./shared/runtime ./scripts/module_runner/runner ./gateway/...`
- `cd api && GOCACHE="$PWD/.gocache" go run scripts/module_runner/main.go sync`
- `GOCACHE="$PWD/.gocache" go test ./framework/...`
- `GOCACHE="$PWD/.gocache" go test ./api/...`
- `GOCACHE="$PWD/.gocache" go run ./api/scripts/module_runner/main.go sync`

# Framework Migration Checkpoint

## 2026-04-01

### Phase Status
- Phase 1: Completed
- Phase 2: Completed
- Phase 3: Completed
- Phase 4: Completed
- Phase 5: Partial
- Phase 6: Partial
- Phase 7: Partial

### Completed This Run
- Removed the API-owned `attribute`, `folder`, and `profile` module trees entirely after migrating runtime ownership to `framework/modules/*`.
- Added a framework-owned module bootstrap path in `framework/runtime/*` for built-in module config loading, DB/cache wiring, runtime registry binding, and JWT auth context setup so framework modules can run without API-side `main.go` shims.
- Converted `framework/modules/attribute/main.go` from a placeholder into a runnable framework-owned entrypoint using the new runtime bootstrap and framework auth middleware.
- Migrated the built-in `folder` module into `framework/modules/folder/*` with framework-owned handler, service, SQL repository, model, module wiring, and config files.
- Migrated the built-in `profile` module into `framework/modules/profile/*` with framework-owned handler, service, SQL repository, model, module wiring, and config files.
- Moved the built-in attribute module implementation out of `api/modules/attribute/*` and into `framework/modules/attribute/*`, including handler, service, repository, model, module wiring, and framework-owned module config files.
- Replaced the attribute module's API-owned Ent repository with a framework-owned SQL repository that uses `framework/pkg/db` plus the runtime SQL bridge only inside `framework/modules/attribute/repository`.
- Reduced `api/modules/attribute` to a composition-only shim: framework module config bridging and route/auth composition through the existing module bootstrap path.
- Extended `api/shared/module` bootstrapping to expose the framework DB client to built-in module shims and to configure framework cache runtime for module processes before framework-owned services use cache contracts.
- Moved the built-in observability module implementation out of `api/modules/observability/*` and into `framework/modules/observability/*`, including handler, service, repository, model, module wiring, and framework-owned module config files.
- Reduced `api/modules/observability` to a composition-only shim: app config to framework config mapping, Ent bootstrap for RBAC permission checks, and route/middleware composition.
- Updated module discovery/startup helpers so `framework/modules/*` owns module discovery and config lookup, while `api/modules/*` can remain composition entrypoints during coexistence.
- Moved the framework DB raw-SQL bridge out of `framework/pkg/db` and into `framework/runtime/db.go`, so the public DB contract no longer exposes `database/sql`.
- Consolidated Redis/cache/pubsub instance ownership under `framework/runtime/cache.go` plus `framework/internal/cache/redis`, leaving `api/shared/redis` and `api/shared/pubsub` as framework-backed compatibility shims only.
- Removed unused app-local DB driver ownership in `api/shared/db/driver/*` and `api/shared/db/interface/client.go`.
- Established a top-level `framework/` Go module and root `go.work`.
- Created initial framework contract packages for App, Module, Context, Router, Cache, DB, Auth, and Lifecycle.
- Added initial internal adapters for Fiber app/router, Redis cache, DB client factory, and JSON-backed lifecycle storage.
- Migrated `api/main.go`, `api/gateway/main.go`, and `api/gateway/runtime/start.go` to boot the gateway through framework application abstractions while preserving the Fiber-backed runtime during migration.
- Migrated `api/shared/cache` to consume framework cache abstractions through compatibility shims.
- Migrated `api/shared/db` to construct database clients through framework DB runtime factories while preserving the existing `GetSQL()` bridge for application-mode compatibility.
- Migrated `api/shared/runtime` storage operations to consume framework lifecycle abstractions through compatibility shims.
- Migrated `api/shared/auth` token pair type to the framework contract.
- Migrated HTTP seam helpers in `api/shared/app` and `api/shared/middleware` to `framework/pkg/http` contracts while keeping Fiber-specific behavior behind compatibility bridges.
- Converted module handler and registry signatures from `fiber.Router` / `*fiber.Ctx` to `framework/pkg/http.Router` / `framework/pkg/http.Context` without changing business logic.
- Moved gateway reverse-proxy registration onto the framework router abstraction and kept websocket upgrade bridging inside adapter code only.

### Validation Snapshot
- No frontend files changed.
- `framework/`: `GOCACHE="$PWD/.gocache" go test ./runtime ./modules/attribute/... ./modules/folder/... ./modules/profile/...` passed.
- `framework/`: `go test ./...` passed.
- `api/`: `GOCACHE="$PWD/.gocache" go test ./shared/runtime ./scripts/module_runner/runner ./gateway/...` passed.
- `api/`: `GOCACHE="$PWD/.gocache" go test ./...` passed.
- `api/`: `GOCACHE="$PWD/.gocache" go run scripts/module_runner/main.go sync` resolved `attribute`, `folder`, and `profile` through module discovery after the API module directories were removed.
- `framework/`: `GOCACHE="$PWD/.gocache" go test ./runtime ./modules/attribute/...` passed.
- `api/`: `GOCACHE="$PWD/.gocache" go test ./modules/attribute/... ./shared/module` passed.
- `framework/`: `GOCACHE="$PWD/.gocache" go test ./runtime ./modules/observability/...` passed.
- `api/`: `GOCACHE="$PWD/.gocache" go test ./modules/observability/... ./shared/runtime ./scripts/module_runner/runner ./gateway/...` passed.
- `framework/`: `go test ./...` passed.
- `api/`: `GOCACHE="$PWD/.gocache" go test ./...` passed.
- Backend compatibility is preserved through `api/shared/*` wrappers plus framework-backed gateway boot.
- `framework/pkg/` no longer exposes raw `*sql.DB` or Redis client types.
- Observability is now discovered from `framework/modules/observability` through the existing multi-root loader, with `api/modules/observability/main.go` retained only as a composition/startup bridge.
- Attribute is now discovered from `framework/modules/attribute` through the existing multi-root loader, with `api/modules/attribute/main.go` retained only as a composition/startup bridge.
- Folder is now discovered from `framework/modules/folder`; `api/modules/folder` no longer exists.
- Profile is now discovered from `framework/modules/profile`; `api/modules/profile` no longer exists.

### Remaining Work
- Complete API-side migration for module boot composition and replace remaining raw DB/Fiber bridging with stable framework-native contracts.
- Reduce remaining Fiber-native response/body helper usage inside handler implementations by introducing stable HTTP response helpers where useful.
- Reduce remaining direct Redis and raw SQL usage outside the newly introduced framework boundaries.
- Move the next built-in module into framework ownership, then remove more API-side composition assumptions from module startup over time.

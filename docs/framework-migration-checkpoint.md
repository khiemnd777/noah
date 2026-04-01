# Framework Migration Checkpoint

## 2026-04-01

### Phase Status
- Phase 1: Completed
- Phase 2: Completed
- Phase 3: Completed
- Phase 4: Completed
- Phase 5: Completed
- Phase 6: Completed
- Phase 7: Completed

### Completed This Run
- Removed the remaining built-in module directories from `api/modules/*` and left only `api/modules/main` as the application-owned feature area.
- Moved the remaining built-in module implementations into `framework/modules/*`: `auditlog`, `auth`, `metadata`, `notification`, `photo`, `rbac`, `realtime`, `search`, `token`, and `user`; `observability` remained framework-owned and its API bridge was removed.
- Internalized the remaining shared support surface required by those modules into `framework/internal/legacy/shared/*` so framework modules no longer import `api/shared/*` directly while preserving current runtime behavior during the final migration.
- Confirmed framework module discovery, entrypoint resolution, and module-runner sync work with built-ins owned only by `framework/modules/*`.
- Extracted remaining legacy shared runtime ownership into `framework/runtime` for circuit breaker, retry-aware HTTP wrappers, HTTP client, JWT/auth context helpers, pubsub, module discovery/registry, port/process helpers, and generic HTTP param parsing.
- Reduced `api/shared/app` framework helpers (`fiber_listener`, `fiber_parser`, `fiber_router`, `fiber_wrapper`, `http_client`, `native_bridge`) to compatibility shims backed by `framework/runtime`.
- Reduced `api/shared/runtime/registry`, `api/shared/circuitbreaker`, `api/shared/pubsub`, and framework-worthy `api/shared/utils/*` helpers (`auth`, `config`, `fiber_user_utils`, `jwtutil`, `module_utils`, `net`, `param_utils`, `token_context`, `util`) to framework-backed compatibility shims.
- Switched gateway/module-runner runtime consumers to framework-owned implementations for module registry generation, module discovery/config resolution, Fiber context bridging, and internal-token lookup.
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
- Completed the final system extraction by moving gateway runtime, module-runner orchestration, reusable commands, SQL migrations, observability assets, and sample deployment tooling from `api/*` into framework-owned paths.
- Deleted the remaining API-owned system directories: `api/gateway`, `api/starter`, `api/scripts`, `api/migrations`, `api/observability`, and `api/docker`.
- Reduced `api` to the sample entrypoint, sample config, and `api/modules/main/*`.

### Validation Snapshot
- No frontend files changed.
- `framework/`: `GOCACHE="$PWD/.gocache" go test ./...` passed with all built-in modules present under `framework/modules/*`.
- `api/`: `GOCACHE="$PWD/.gocache" go test ./...` passed after removing all built-in `api/modules/*` directories except `api/modules/main`.
- `api/`: `GOCACHE="$PWD/.gocache" go run scripts/module_runner/main.go sync` discovered and synchronized built-ins from framework-owned module roots after API-side built-in module removal.
- `framework/`: `GOCACHE="$PWD/.gocache" go test ./runtime/...` passed.
- `framework/`: `GOCACHE="$PWD/.gocache" go test ./...` passed.
- `api/`: `GOCACHE="$PWD/.gocache" go test ./shared/... ./gateway/... ./scripts/module_runner/runner` passed.
- `api/`: `GOCACHE="$PWD/.gocache" go test ./...` passed.
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
- `framework/`: `go test ./...` passed after the final gateway/tooling/migration asset extraction.
- `api/`: `go test ./...` passed after deleting API-owned system packages.
- `framework/`: `go run ./cmd/module_runner sync` completed successfully and wrote synchronized runtime state into `api/tmp/modules.json` from framework-owned discovery/orchestration code.

### Final Shared Layer Cutover
- Migrated the entire `api/shared/**` tree into exported framework ownership at `framework/shared/**`.
- Repointed all live `api/**` imports, scripts, templates, and Ent generation paths from `api/shared` to `framework/shared`.
- Removed `api/shared` completely; the backend no longer has an API-owned shared system layer.
- Preserved `framework/pkg/*` as contract-only packages; Fiber, SQL, Redis, websocket, and Ent concrete types remain outside the public contract surface.

### Remaining Work
- System extraction is complete. Remaining work is cleanup-only: reduce duplication between `framework/shared/*` and `framework/internal/legacy/shared/*`, and generalize the remaining sample-app filesystem assumptions inside framework runtime helpers.

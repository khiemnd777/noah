# Noah Module Inventory

## Purpose

This file is a navigation index for agents and engineers. Use it to find the likely owning runtime, module, feature, and source files before making a change.

This inventory is not a substitute for code inspection. Before editing, verify the nearest source files listed here and follow the owning module's handler, service, repository, API, mapper, route, and widget patterns.

## Status Legend

| Status | Meaning |
| --- | --- |
| `Registered` | Backend module has `config.yaml`, frontend feature has `index.tsx`, or main subfeature has `registry.go`. |
| `Main-owned area` | Owned by `api/modules/main`, but not under `api/modules/main/features/*`. |
| `Platform` | Shared/runtime module or infrastructure used across features. |
| `Verify before edit` | Evidence exists, but exact ownership or counterpart must be confirmed from nearby source before changing behavior. |

## Runtime Index

| Runtime / Boundary | Entrypoint | Registry / Loader Source | Evidence |
| --- | --- | --- | --- |
| `api` | `api/main.go` | `api/gateway/runtime/start.go`, `api/modules/*/config.yaml`, `api/tmp/runtime.json` | `api/main.go`, `api/gateway/runtime/start.go`, `api/shared/runtime/registry.go`, `api/scripts/module_runner` |
| `fe` | `fe/src/main.tsx` | `fe/src/core/index.ts`, `fe/src/features/*/index.tsx` | `fe/src/main.tsx`, `fe/src/core/index.ts` |
| `.github` / `deploy` | `.github/workflows/deploy.yml` | `deploy/scripts/provision-and-deploy.sh`, `deploy/scripts/render-production-config.sh` | `.github/workflows/deploy.yml`, `deploy/scripts/provision-and-deploy.sh` |

Backend module lifecycle control uses `api/tmp/runtime.json` as the runtime source of truth for `start`, `stop`, `restart`, `sync`, and `status`. `api/tmp/modules.json` is deprecated and is not part of lifecycle control.

## Backend Top-Level Modules

| Module | Route / Boundary Evidence | Primary Folders | Status | Notes |
| --- | --- | --- | --- | --- |
| `api/modules/attribute` | `api/modules/attribute/config.yaml` (`/api/attribute`) | `api/modules/attribute` | `Registered` | Verify handler/service/repository files before edit. |
| `api/modules/auditlog` | `api/modules/auditlog/config.yaml` (`/api/audit`) | `api/modules/auditlog` | `Platform` | Audit/event infrastructure; verify consumers before edit. |
| `api/modules/auth` | `api/modules/auth/config.yaml` (`/api/auth`) | `api/modules/auth` | `Platform` | Auth-sensitive; verify middleware and token behavior before edit. |
| `api/modules/folder` | `api/modules/folder/config.yaml` (`/api/folder`) | `api/modules/folder` | `Registered` | Verify owning handlers before edit. |
| `api/modules/i18n` | `api/modules/i18n/config.yaml` (`/api/i18n`) | `api/modules/i18n` | `Registered` | FE counterpart: `fe/src/features/languages`. |
| `api/modules/main` | `api/modules/main/config.yaml` (`/api/department`) | `api/modules/main`, `api/modules/main/department`, `api/modules/main/features` | `Registered` | Contains main-owned areas and feature registry. Identity-sensitive areas must use `docs/identity-contract.md`. |
| `api/modules/metadata` | `api/modules/metadata/config.yaml` (`/api/metadata`) | `api/modules/metadata` | `Platform` | Metadata/custom-field infrastructure; verify downstream consumers before edit. |
| `api/modules/notification` | `api/modules/notification/config.yaml` (`/api/notification`) | `api/modules/notification` | `Registered` | FE counterpart: `fe/src/features/notification`. |
| `api/modules/observability` | `api/modules/observability/config.yaml` (`/api/observability`) | `api/modules/observability` | `Platform` | Log-query module; verify observability config before runtime changes. |
| `api/modules/photo` | `api/modules/photo/config.yaml` (`/api/photo`) | `api/modules/photo` | `Registered` | Storage-sensitive; verify storage config and service paths before edit. |
| `api/modules/profile` | `api/modules/profile/config.yaml` (`/api/profile`) | `api/modules/profile` | `Registered` | Verify auth/profile ownership before edit. |
| `api/modules/rbac` | `api/modules/rbac/config.yaml` (`/api/rbac`) | `api/modules/rbac` | `Platform` | Permission-sensitive; verify RBAC middleware and FE route permissions before edit. |
| `api/modules/realtime` | `api/modules/realtime/config.yaml` (`/ws`) | `api/modules/realtime` | `Platform` | WebSocket/pubsub runtime; verify event consumers before edit. |
| `api/modules/search` | `api/modules/search/config.yaml` (`/api/search`) | `api/modules/search` | `Platform` | Search infrastructure; verify guards and index/event assumptions before edit. |
| `api/modules/token` | `api/modules/token/config.yaml` (`/api/token`) | `api/modules/token` | `Registered` | Token lifecycle; verify cron/cleanup behavior before edit. |
| `api/modules/user` | `api/modules/user/config.yaml` (`/api/user`) | `api/modules/user` | `Registered` | Identity-sensitive; read `docs/identity-contract.md` before staff/user/department changes. |

## Main-Owned Backend Areas

| Area | Evidence | Handler / Service / Repository | Status | Notes |
| --- | --- | --- | --- | --- |
| `api/modules/main/department` | `api/modules/main/department/department.go` | `handler/`, `service/`, `repository/`, `model/` | `Main-owned area` | Identity-sensitive; read `docs/identity-contract.md` before edit. FE counterpart: `fe/src/features/department`. |
| `api/modules/main/features/staff` | `api/modules/main/features/staff/registry.go` | `handler/`, `service/`, `repository/` | `Registered` | Identity-sensitive; read `docs/identity-contract.md` before edit. FE counterpart: `fe/src/features/staff`. |
| `api/modules/main/features/__relation` | `api/modules/main/features/__relation/registry.go` | `handler/`, `service/`, `repository/`, `policy/`, `registrar/` | `Verify before edit` | Internal relation feature evidence exists; confirm call sites and ownership before changing behavior. |

## Frontend Features

| Feature | Registration Evidence | Nearby Folders | Status | Notes |
| --- | --- | --- | --- | --- |
| `fe/src/features/auth` | `fe/src/features/auth/index.tsx` | `schemas/`, `widgets/` | `Registered` | Auth/account route; verify auth state and route guards before edit. |
| `fe/src/features/department` | `fe/src/features/department/index.tsx` | `api/`, `model/`, `schemas/`, `tables/`, `widgets/` | `Registered` | Identity-sensitive; read `docs/identity-contract.md`. API counterpart: `api/modules/main/department`. |
| `fe/src/features/languages` | `fe/src/features/languages/index.tsx` | `api/`, `mapper/`, `model/`, `schemas/`, `tables/`, `widgets/` | `Registered` | API counterpart: `api/modules/i18n`. |
| `fe/src/features/metadata` | `fe/src/features/metadata/index.tsx` | feature root | `Registered` | API counterpart: `api/modules/metadata`; verify exact API files before edit. |
| `fe/src/features/notification` | `fe/src/features/notification/index.tsx` | `components/`, `widgets/` | `Registered` | API counterpart: `api/modules/notification`. |
| `fe/src/features/observability_logs` | `fe/src/features/observability_logs/index.tsx` | `api/`, `components/`, `model/`, `pages/`, `tables/` | `Registered` | API counterpart: `api/modules/observability`. |
| `fe/src/features/rbac` | `fe/src/features/rbac/index.tsx` | `api/`, `components/`, `mapper/`, `model/`, `schemas/`, `tables/`, `widgets/` | `Registered` | Permission-sensitive; API counterpart: `api/modules/rbac`. |
| `fe/src/features/search` | `fe/src/features/search/index.tsx` | `components/`, `widgets/` | `Registered` | API counterpart: `api/modules/search`; verify guards before behavior changes. |
| `fe/src/features/settings` | `fe/src/features/settings/index.tsx` | `api/`, `components/`, `mapper/`, `model/`, `schemas/`, `widgets/` | `Registered` | Verify backend counterpart from API calls before edit. |
| `fe/src/features/staff` | `fe/src/features/staff/index.tsx` | `api/`, `mapper/`, `model/`, `schemas/`, `tables/`, `widgets/` | `Registered` | Identity-sensitive; read `docs/identity-contract.md`. API counterpart: `api/modules/main/features/staff`. |

## FE / API Ownership Hints

| Frontend Area | Backend Area | Confidence | Notes |
| --- | --- | --- | --- |
| `fe/src/features/staff` | `api/modules/main/features/staff` | `Registered` | Identity-sensitive; verify `users.id` vs `staffs.id` flow before edit. |
| `fe/src/features/department` | `api/modules/main/department` | `Registered` | Identity-sensitive; verify department administrator uses `users.id`. |
| `fe/src/features/languages` | `api/modules/i18n` | `Registered` | Verify request/response DTOs and mapper before contract changes. |
| `fe/src/features/rbac` | `api/modules/rbac` | `Registered` | Permission-sensitive; verify FE route permissions and backend middleware. |
| `fe/src/features/metadata` | `api/modules/metadata` | `Verify before edit` | Counterpart exists; verify exact route/API ownership before changing behavior. |
| `fe/src/features/notification` | `api/modules/notification` | `Registered` | Verify realtime or chip consumers before behavior changes. |
| `fe/src/features/observability_logs` | `api/modules/observability` | `Registered` | Verify Loki/config assumptions before runtime changes. |
| `fe/src/features/search` | `api/modules/search` | `Verify before edit` | Verify guards, scopes, and indexed entities before behavior changes. |

## Reading Rules

- Start from the runtime or feature row, then inspect the listed `config.yaml`, `index.tsx`, or `registry.go`.
- For backend behavior, continue through handler, service, repository, registry, middleware, and config files inside the owning module only.
- For frontend behavior, continue through the feature's `api/`, `mapper/`, `model/`, `schemas/`, `tables/`, `widgets/`, pages, and route metadata files as applicable.
- For FE/API work, inspect both sides of the mapping and verify request/response payloads, permissions, mappers, invalidation, and route ownership before editing.
- For staff, user, or department identity work, read `docs/identity-contract.md` before inspecting implementation files.

## Maintenance Rules

- Update this file whenever module ownership, feature ownership, registration, route ownership, runtime app boundaries, deploy/runtime boundaries, or FE/API ownership mappings change.
- Keep every row evidence-backed with paths that exist in the repo.
- Do not add inferred modules, future runtimes, or business areas unless repo evidence exists.
- Keep `docs/tech-stack-inventory.md` as the stack/runtime/tooling/deploy infrastructure index; keep this file focused on module and ownership navigation.

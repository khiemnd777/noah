# Framework Architecture

## 1. System Overview

This repository has been transformed into a framework-first architecture:

- `framework/` → owns the entire system (runtime, modules, infra)
- `api/` → sample application only

---

## 2. Directory Structure

### Framework

- `framework/pkg` → public contracts (no infra types)
- `framework/internal` → implementations
- `framework/runtime` → system orchestration
- `framework/modules` → built-in modules
- `framework/cmd` → operational commands
- `framework/scripts` → reusable scripts
- `framework/migrations` → SQL migrations
- `framework/observability` → monitoring stack

---

### API (Sample)

- `api/main.go` → entrypoint
- `api/modules/main/features` → business modules
- `api/config.yaml` → configuration

---

## 3. Module System

### Built-in Modules

Located in:

```text
framework/modules/*
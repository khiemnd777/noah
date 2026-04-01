# Module System

## Module Types

### 1. Built-in Modules
- Located in `framework/modules`
- Provided by framework
- Not reimplemented by users

### 2. Business Modules
- Located in `api/modules/main/features`
- Written by application developers

---

## Module Structure

Each module contains:

- handler
- service
- repository
- models
- config.yaml
- main.go

---

## Module Loading

Modules are discovered via multi-root loader.

---

## Module Rules

- No direct dependency on api/*
- Must use framework contracts
- Must not leak infra types

---

## Lifecycle (future)

- Register
- Start
- Stop
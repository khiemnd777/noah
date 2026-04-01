# Quickstart

## 1. Run sample app

```bash
cd api
go run main.go
```

## 2. Create new business module

Create under:
`
api/modules/main/features/<module_name>
`

## 3. Use built-in modules

Built-in modules are auto-loaded:

- auth
- user
- rbac
- metadata
- etc.

## 4. Configuration

Edit:
`
api/config.yaml
`

## 5. Run module sync

```bash
go run ../framework/cmd/module_runner sync
```

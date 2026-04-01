# API Sample Project

`api/` is a minimal sample application that demonstrates how to start the framework without carrying the previous production module tree.

## What stays here

- `main.go` contains the sample entrypoint.
- `modules/main/features` contains a tiny sample feature registration package.
- Root config files remain as examples for local customization.

## Run the sample

```bash
go run .
```

Optional environment variables:

- `API_HOST` defaults to `127.0.0.1`
- `API_PORT` defaults to `8080`

## Routes

- `GET /`
- `GET /health`

This directory is intentionally lightweight. Framework implementation and reusable runtime code stay in `/framework`.

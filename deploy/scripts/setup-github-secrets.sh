#!/bin/bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
PROJECT_ENV_FILE="${PROJECT_ENV_FILE:-$ROOT_DIR/deploy/config/project.env}"

source "$ROOT_DIR/deploy/scripts/lib-env.sh"

if ! command -v gh >/dev/null 2>&1; then
  echo "GitHub CLI (gh) is required." >&2
  exit 1
fi

if [ ! -f "$PROJECT_ENV_FILE" ]; then
  echo "Missing config file: $PROJECT_ENV_FILE" >&2
  exit 1
fi

if ! gh auth status >/dev/null 2>&1; then
  gh auth login
fi

load_env_file "$PROJECT_ENV_FILE"

require_vars VPS_HOST VPS_PORT VPS_USER VPS_PASSWORD VPS_SUDO_PASSWORD VPS_DEPLOY_PATH LETSENCRYPT_EMAIL

"$ROOT_DIR/deploy/scripts/push-github-secrets.sh" "$PROJECT_ENV_FILE"

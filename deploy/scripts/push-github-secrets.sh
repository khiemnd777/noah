#!/bin/bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
PROJECT_ENV_FILE="${1:-${PROJECT_ENV_FILE:-$ROOT_DIR/deploy/config/project.env}}"

source "$ROOT_DIR/deploy/scripts/lib-env.sh"

if ! command -v gh >/dev/null 2>&1; then
  echo "GitHub CLI (gh) is required." >&2
  exit 1
fi

if [ ! -f "$PROJECT_ENV_FILE" ]; then
  echo "Missing config file: $PROJECT_ENV_FILE" >&2
  exit 1
fi

load_env_file "$PROJECT_ENV_FILE"

require_vars VPS_HOST VPS_PORT VPS_USER VPS_PASSWORD VPS_SUDO_PASSWORD VPS_DEPLOY_PATH LETSENCRYPT_EMAIL

PROJECT_ENV_B64="$(base64 <"$PROJECT_ENV_FILE" | tr -d '\n')"

set_secret() {
  local name="$1"
  local value="$2"

  printf '%s' "$value" | gh secret set "$name"
  echo "Updated secret: $name"
}

set_secret VPS_HOST "$VPS_HOST"
set_secret VPS_PORT "$VPS_PORT"
set_secret VPS_USER "$VPS_USER"
set_secret VPS_PASSWORD "$VPS_PASSWORD"
set_secret VPS_SUDO_PASSWORD "$VPS_SUDO_PASSWORD"
set_secret VPS_DEPLOY_PATH "$VPS_DEPLOY_PATH"
set_secret PROJECT_ENV_B64 "$PROJECT_ENV_B64"
set_secret LETSENCRYPT_EMAIL "$LETSENCRYPT_EMAIL"

smtp_keys=(SMTP_HOST SMTP_PORT SMTP_USERNAME SMTP_PASSWORD SMTP_FROM SMTP_TO)
smtp_filled=0
for key in "${smtp_keys[@]}"; do
  if [ -n "${!key:-}" ]; then
    smtp_filled=$((smtp_filled + 1))
  fi
done

if [ "$smtp_filled" -ne 0 ] && [ "$smtp_filled" -ne "${#smtp_keys[@]}" ]; then
  echo "SMTP config must be either fully empty or fully populated." >&2
  exit 1
fi

if [ "$smtp_filled" -eq "${#smtp_keys[@]}" ]; then
  for key in "${smtp_keys[@]}"; do
    set_secret "$key" "${!key}"
  done
fi

#!/bin/bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

source "$ROOT_DIR/../deploy/scripts/lib-env.sh"

ENV_FILE="$ROOT_DIR/.env"
if [ "${APP_ENV:-development}" = "production" ] && [ -f "$ROOT_DIR/.env.prod" ]; then
  ENV_FILE="$ROOT_DIR/.env.prod"
fi

load_env_file "$ENV_FILE"

LOKI_SCHEME="${LOKI_SCHEME:-http}"
LOKI_HOST="${LOKI_HOST:-loki}"
LOKI_PORT="${LOKI_PORT:-3100}"
LOKI_HOST_PORT="${LOKI_HOST_PORT:-$LOKI_PORT}"
if [ -z "${LOKI_BASE_URL:-}" ]; then
  LOKI_BASE_URL="${LOKI_SCHEME}://${LOKI_HOST}:${LOKI_PORT}"
fi
export LOKI_SCHEME LOKI_HOST LOKI_PORT LOKI_HOST_PORT LOKI_BASE_URL

if docker compose version >/dev/null 2>&1; then
  COMPOSE_CMD=(docker compose)
elif command -v docker-compose >/dev/null 2>&1; then
  COMPOSE_CMD=(docker-compose)
else
  echo "docker compose is required but not found."
  exit 1
fi

if ! docker info >/dev/null 2>&1; then
  echo "docker daemon is not running. Please start Docker Desktop or dockerd first."
  exit 1
fi

mkdir -p \
  "$ROOT_DIR/tmp/observability/logs" \
  "$ROOT_DIR/tmp/observability/loki" \
  "$ROOT_DIR/tmp/observability/promtail" \
  "$ROOT_DIR/tmp/observability/grafana"

touch "$ROOT_DIR/tmp/observability/logs/noah_api.json.log"

"$ROOT_DIR/scripts/render_observability_config.sh" --env-file "$ENV_FILE"

"${COMPOSE_CMD[@]}" -f "$ROOT_DIR/docker-compose.observability.yml" up -d

echo "Observability stack is running."
echo "Grafana: http://127.0.0.1:3001 (admin/admin)"
echo "Loki health: http://127.0.0.1:${LOKI_HOST_PORT}/ready"

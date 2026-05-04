#!/bin/bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$ROOT_DIR/../deploy/scripts/lib-env.sh"

ENV_FILE="${ROOT_DIR}/.env"

while [ $# -gt 0 ]; do
  case "$1" in
    --env-file)
      ENV_FILE="$2"
      shift 2
      ;;
    --env-file=*)
      ENV_FILE="${1#*=}"
      shift
      ;;
    *)
      echo "Unknown argument: $1" >&2
      exit 1
      ;;
  esac
done

if [ ! -f "$ENV_FILE" ]; then
  echo "Env file not found: $ENV_FILE" >&2
  exit 1
fi

load_env_file "$ENV_FILE"

LOKI_SCHEME="${LOKI_SCHEME:-http}"
LOKI_HOST="${LOKI_HOST:-loki}"
LOKI_PORT="${LOKI_PORT:-3100}"
if [ -z "${LOKI_BASE_URL:-}" ]; then
  LOKI_BASE_URL="${LOKI_SCHEME}://${LOKI_HOST}:${LOKI_PORT}"
fi

PROMTAIL_TARGET="${ROOT_DIR}/observability/promtail-config.yaml"
GRAFANA_TARGET="${ROOT_DIR}/observability/grafana/provisioning/datasources/loki.yaml"
PROMTAIL_TEMPLATE="${ROOT_DIR}/observability/promtail-config.yaml.tmpl"
GRAFANA_TEMPLATE="${ROOT_DIR}/observability/grafana/provisioning/datasources/loki.yaml.tmpl"

mkdir -p "$(dirname "$PROMTAIL_TARGET")" "$(dirname "$GRAFANA_TARGET")"

if [ ! -f "$PROMTAIL_TEMPLATE" ] || [ ! -f "$GRAFANA_TEMPLATE" ]; then
  echo "Observability template file is missing." >&2
  exit 1
fi

escaped_loki_base_url="${LOKI_BASE_URL//\\/\\\\}"
escaped_loki_base_url="${escaped_loki_base_url//&/\\&}"

sed "s|\${LOKI_BASE_URL}|$escaped_loki_base_url|g" "$PROMTAIL_TEMPLATE" >"$PROMTAIL_TARGET"
sed "s|\${LOKI_BASE_URL}|$escaped_loki_base_url|g" "$GRAFANA_TEMPLATE" >"$GRAFANA_TARGET"

echo "Rendered observability config:"
echo "- $PROMTAIL_TARGET"
echo "- $GRAFANA_TARGET"

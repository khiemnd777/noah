#!/bin/bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
PROJECT_ENV_FILE="${PROJECT_ENV_FILE:-$ROOT_DIR/deploy/config/project.env}"

source "$ROOT_DIR/deploy/scripts/lib-env.sh"

if [ ! -f "$PROJECT_ENV_FILE" ]; then
  echo "Missing config file: $PROJECT_ENV_FILE" >&2
  exit 1
fi

load_env_file "$ROOT_DIR/.env.sample"
load_env_file "$ROOT_DIR/api/.env.sample"
load_env_file "$PROJECT_ENV_FILE"

require_vars VPS_DEPLOY_PATH APP_DOMAIN LETSENCRYPT_EMAIL JWT_TOKEN_SECRET INTERNAL_AUTH_TOKEN AUTH_LOG_TOKEN PG_USER PG_PASSWORD PG_DBNAME

APP_PROTOCOL="${APP_PROTOCOL:-https}"
APP_FE_ORIGIN="${APP_PROTOCOL}://${APP_DOMAIN}"
APP_ENV="production"
HOST="0.0.0.0"
PG_HOST="postgres"
PG_PORT="5432"
REDIS_HOST="redis"
REDIS_PORT="6379"
FRONTEND_HOST_PORT="${FRONTEND_HOST_PORT:-18080}"
API_HOST_PORT="${API_HOST_PORT:-17999}"
NGINX_SITE_NAME="${NGINX_SITE_NAME:-${APP_DOMAIN}}"
PGDATA_DIR="${PGDATA_DIR:-${VPS_DEPLOY_PATH}/.data/postgres}"
REDISDATA_DIR="${REDISDATA_DIR:-${VPS_DEPLOY_PATH}/.data/redis}"
LOKI_SCHEME="${LOKI_SCHEME:-http}"
LOKI_HOST="${LOKI_HOST:-loki}"
LOKI_PORT="${LOKI_PORT:-3100}"
LOKI_HOST_PORT="${LOKI_HOST_PORT:-$LOKI_PORT}"
if [ -z "${LOKI_BASE_URL:-}" ]; then
  LOKI_BASE_URL="${LOKI_SCHEME}://${LOKI_HOST}:${LOKI_PORT}"
fi

ROOT_ENV_TARGET="$ROOT_DIR/.env.prod"
API_ENV_TARGET="$ROOT_DIR/api/.env.prod"
GENERATED_DIR="$ROOT_DIR/deploy/generated/nginx"
HTTP_NGINX_TARGET="$GENERATED_DIR/${NGINX_SITE_NAME}.http.conf"
HTTPS_NGINX_TARGET="$GENERATED_DIR/${NGINX_SITE_NAME}.https.conf"

mkdir -p "$GENERATED_DIR" "$(dirname "$API_ENV_TARGET")"

write_env_file_from_keys "$ROOT_ENV_TARGET" APP_FE_ORIGIN

api_keys=()
while IFS= read -r key; do
  api_keys+=("$key")
done < <(list_env_keys "$ROOT_DIR/api/.env.sample")
write_env_file_from_keys "$API_ENV_TARGET" "${api_keys[@]}"

cat >"$HTTP_NGINX_TARGET" <<EOF
server {
    listen 80;
    listen [::]:80;
    server_name ${APP_DOMAIN};

    client_max_body_size ${BODY_LIMIT_MB}m;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location /api {
        proxy_pass http://127.0.0.1:${API_HOST_PORT};
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location /ws {
        proxy_pass http://127.0.0.1:${API_HOST_PORT};
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location / {
        proxy_pass http://127.0.0.1:${FRONTEND_HOST_PORT};
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF

cat >"$HTTPS_NGINX_TARGET" <<EOF
server {
    listen 80;
    listen [::]:80;
    server_name ${APP_DOMAIN};

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 301 https://\$host\$request_uri;
    }
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name ${APP_DOMAIN};

    ssl_certificate /etc/letsencrypt/live/${APP_DOMAIN}/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/${APP_DOMAIN}/privkey.pem;

    client_max_body_size ${BODY_LIMIT_MB}m;

    location = /api/ping {
        proxy_pass http://127.0.0.1:${API_HOST_PORT}/ping;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location /api {
        proxy_pass http://127.0.0.1:${API_HOST_PORT};
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location /ws {
        proxy_pass http://127.0.0.1:${API_HOST_PORT};
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location / {
        proxy_pass http://127.0.0.1:${FRONTEND_HOST_PORT};
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF

echo "Rendered production config:"
echo "- $ROOT_ENV_TARGET"
echo "- $API_ENV_TARGET"
echo "- $HTTP_NGINX_TARGET"
echo "- $HTTPS_NGINX_TARGET"

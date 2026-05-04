#!/bin/bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DEPLOY_ENV_FILE="${DEPLOY_ENV_FILE:-$REPO_ROOT/.deploy.env}"

source "$REPO_ROOT/deploy/scripts/lib-env.sh"

if [ ! -f "$DEPLOY_ENV_FILE" ]; then
  echo "Missing deploy env file: $DEPLOY_ENV_FILE" >&2
  exit 1
fi

load_env_file "$DEPLOY_ENV_FILE"

require_vars VPS_DEPLOY_PATH VPS_SUDO_PASSWORD PROJECT_ENV_B64 LETSENCRYPT_EMAIL

PROJECT_ENV_PATH="$REPO_ROOT/deploy/config/project.env"

mkdir -p "$(dirname "$PROJECT_ENV_PATH")"
printf '%s' "$PROJECT_ENV_B64" | base64 -d >"$PROJECT_ENV_PATH"

load_env_file "$PROJECT_ENV_PATH"

sudo_exec() {
  printf '%s\n' "$VPS_SUDO_PASSWORD" | sudo -S "$@"
}

install_packages_apt() {
  sudo_exec apt-get update
  sudo_exec apt-get install -y ca-certificates curl nginx certbot python3-certbot-nginx rsync
  sudo_exec apt-get install -y docker.io || true

  if ! sudo_exec docker compose version >/dev/null 2>&1; then
    sudo_exec apt-get install -y docker-compose-plugin || sudo_exec apt-get install -y docker-compose
  fi
}

install_packages_dnf() {
  sudo_exec dnf install -y ca-certificates curl nginx certbot python3-certbot-nginx rsync docker docker-compose-plugin
}

ensure_runtime() {
  if command -v apt-get >/dev/null 2>&1; then
    install_packages_apt
  elif command -v dnf >/dev/null 2>&1; then
    install_packages_dnf
  else
    echo "Unsupported package manager. Expected apt-get or dnf." >&2
    exit 1
  fi

  sudo_exec systemctl enable --now docker
  sudo_exec systemctl enable --now nginx
  sudo_exec systemctl enable --now certbot.timer || true
}

compose_cmd() {
  if sudo_exec docker compose version >/dev/null 2>&1; then
    sudo_exec docker compose "$@"
    return
  fi

  sudo_exec docker-compose "$@"
}

install_nginx_conf() {
  local source_file="$1"
  local site_name="$2"
  local nginx_target

  sudo_exec mkdir -p /var/www/certbot

  if [ -d /etc/nginx/sites-available ] || [ -d /etc/nginx/sites-enabled ]; then
    sudo_exec mkdir -p /etc/nginx/sites-available /etc/nginx/sites-enabled
    nginx_target="/etc/nginx/sites-available/${site_name}.conf"
    sudo_exec rm -f /etc/nginx/sites-enabled/default
    sudo_exec rm -f "/etc/nginx/sites-enabled/${site_name}.conf"
    printf '%s\n' "$VPS_SUDO_PASSWORD" | sudo -S tee "$nginx_target" >/dev/null <"$source_file"
    sudo_exec ln -sf "$nginx_target" "/etc/nginx/sites-enabled/${site_name}.conf"
  else
    sudo_exec mkdir -p /etc/nginx/conf.d
    nginx_target="/etc/nginx/conf.d/${site_name}.conf"
    printf '%s\n' "$VPS_SUDO_PASSWORD" | sudo -S tee "$nginx_target" >/dev/null <"$source_file"
  fi

  sudo_exec nginx -t
  sudo_exec systemctl reload nginx
}

verify_domain_points_here() {
  local domain="$1"
  local public_ip
  local resolved_ips

  public_ip="$(curl -4fsS https://api.ipify.org)"
  resolved_ips="$(getent ahostsv4 "$domain" | awk '{print $1}' | sort -u || true)"

  if [ -z "$resolved_ips" ]; then
    echo "Domain does not resolve yet: $domain" >&2
    exit 1
  fi

  if ! printf '%s\n' "$resolved_ips" | grep -Fxq "$public_ip"; then
    echo "Domain $domain does not point to this VPS public IP ($public_ip)." >&2
    printf 'Resolved IPs: %s\n' "$resolved_ips" >&2
    exit 1
  fi
}

wait_for_http() {
  local url="$1"
  local attempt

  for attempt in {1..15}; do
    if curl -fsS "$url" >/dev/null 2>&1; then
      return 0
    fi
    sleep 2
  done

  echo "HTTP endpoint did not become ready: $url" >&2
  exit 1
}

run_healthchecks() {
  local public_base_url="$1"
  local frontend_local_port="$2"
  local api_local_port="$3"

  wait_for_http "http://127.0.0.1:${frontend_local_port}/health"
  wait_for_http "http://127.0.0.1:${api_local_port}/ping"
  wait_for_http "${public_base_url}/"
  wait_for_http "${public_base_url}/api/ping"
}

ensure_runtime

"$REPO_ROOT/deploy/scripts/render-production-config.sh"
"$REPO_ROOT/api/scripts/render_observability_config.sh" --env-file "$REPO_ROOT/api/.env.prod"

load_env_file "$PROJECT_ENV_PATH"
APP_PROTOCOL="${APP_PROTOCOL:-https}"
APP_DOMAIN="${APP_DOMAIN:-}"
NGINX_SITE_NAME="${NGINX_SITE_NAME:-${APP_DOMAIN}}"
FRONTEND_HOST_PORT="${FRONTEND_HOST_PORT:-18080}"
API_HOST_PORT="${API_HOST_PORT:-17999}"

require_vars APP_DOMAIN NGINX_SITE_NAME

HTTP_NGINX_FILE="$REPO_ROOT/deploy/generated/nginx/${NGINX_SITE_NAME}.http.conf"
HTTPS_NGINX_FILE="$REPO_ROOT/deploy/generated/nginx/${NGINX_SITE_NAME}.https.conf"

install_nginx_conf "$HTTP_NGINX_FILE" "$NGINX_SITE_NAME"
verify_domain_points_here "$APP_DOMAIN"

sudo_exec certbot certonly \
  --webroot \
  --webroot-path /var/www/certbot \
  --non-interactive \
  --agree-tos \
  --keep-until-expiring \
  --email "$LETSENCRYPT_EMAIL" \
  -d "$APP_DOMAIN"

install_nginx_conf "$HTTPS_NGINX_FILE" "$NGINX_SITE_NAME"

mkdir -p "$REPO_ROOT/api/tmp/observability/logs"

compose_cmd \
  --env-file "$REPO_ROOT/api/.env.prod" \
  -f "$REPO_ROOT/api/docker-compose.prod.yml" \
  up --build -d

PUBLIC_BASE_URL="${APP_PROTOCOL}://${APP_DOMAIN}"
run_healthchecks "$PUBLIC_BASE_URL" "$FRONTEND_HOST_PORT" "$API_HOST_PORT"

echo "Provisioning and deployment completed successfully."

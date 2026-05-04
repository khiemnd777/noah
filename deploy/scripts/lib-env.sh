#!/bin/bash

trim_whitespace() {
  local value="$1"
  value="${value#"${value%%[![:space:]]*}"}"
  value="${value%"${value##*[![:space:]]}"}"
  printf '%s' "$value"
}

strip_matching_quotes() {
  local value="$1"

  if [[ "$value" == \"*\" && "$value" == *\" ]]; then
    value="${value:1:${#value}-2}"
  elif [[ "$value" == \'*\' && "$value" == *\' ]]; then
    value="${value:1:${#value}-2}"
  fi

  printf '%s' "$value"
}

load_env_file() {
  local env_file="$1"

  if [ ! -f "$env_file" ]; then
    return 0
  fi

  while IFS= read -r raw_line || [ -n "$raw_line" ]; do
    raw_line="${raw_line%$'\r'}"

    case "$raw_line" in
      ""|\#*)
        continue
        ;;
    esac

    if [[ "$raw_line" != *=* ]]; then
      echo "Invalid env line in $env_file: $raw_line" >&2
      exit 1
    fi

    local key="${raw_line%%=*}"
    local value="${raw_line#*=}"

    key="$(trim_whitespace "$key")"
    value="$(trim_whitespace "$value")"

    if [[ "$key" == export[[:space:]]* ]]; then
      key="$(trim_whitespace "${key#export}")"
    fi

    value="$(strip_matching_quotes "$value")"

    if [ -z "$key" ]; then
      echo "Invalid env key in $env_file" >&2
      exit 1
    fi

    printf -v "$key" '%s' "$value"
    export "$key"
  done <"$env_file"
}

require_vars() {
  local missing=()
  local name

  for name in "$@"; do
    if [ -z "${!name:-}" ]; then
      missing+=("$name")
    fi
  done

  if [ "${#missing[@]}" -gt 0 ]; then
    printf 'Missing required variables: %s\n' "${missing[*]}" >&2
    exit 1
  fi
}

list_env_keys() {
  local env_file="$1"

  while IFS= read -r raw_line || [ -n "$raw_line" ]; do
    raw_line="${raw_line%$'\r'}"
    case "$raw_line" in
      ""|\#*)
        continue
        ;;
    esac

    if [[ "$raw_line" != *=* ]]; then
      continue
    fi

    local key="${raw_line%%=*}"
    key="$(trim_whitespace "$key")"

    if [[ "$key" == export[[:space:]]* ]]; then
      key="$(trim_whitespace "${key#export}")"
    fi

    if [ -n "$key" ]; then
      printf '%s\n' "$key"
    fi
  done <"$env_file"
}

format_env_value() {
  local value="${1-}"

  if [ -z "$value" ]; then
    printf ''
    return
  fi

  if [[ "$value" =~ [[:space:]#\"] ]]; then
    value="${value//\\/\\\\}"
    value="${value//\"/\\\"}"
    printf '"%s"' "$value"
    return
  fi

  printf '%s' "$value"
}

write_env_file_from_keys() {
  local target="$1"
  shift

  : >"$target"

  local key
  for key in "$@"; do
    printf '%s=%s\n' "$key" "$(format_env_value "${!key-}")" >>"$target"
  done
}

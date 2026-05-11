#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
SKILLS_SOURCE_DIR="${1:-"$ROOT_DIR/agents/skills"}"
SUBAGENTS_SOURCE_DIR="${2:-"$ROOT_DIR/agents/subagents"}"

fail() {
  echo "Validation failed: $*" >&2
  exit 1
}

assert_contains() {
  local file="$1"
  local pattern="$2"
  local message="$3"

  if ! grep -Eq "$pattern" "$file"; then
    fail "$message ($file)"
  fi
}

validate_skills() {
  local source_dir="$1"
  local found_skill=0

  [[ -d "$source_dir" ]] || fail "Skills source directory not found: $source_dir"

  while IFS= read -r -d '' entry; do
    if [[ -d "$entry" ]]; then
      continue
    fi

    if [[ "$(basename "$entry")" == "README.md" ]]; then
      continue
    fi

    fail "Unexpected non-skill artifact in skills source directory: $entry"
  done < <(find "$source_dir" -mindepth 1 -maxdepth 1 -print0 | sort -z)

  while IFS= read -r -d '' skill_dir; do
    found_skill=1

    local skill_name
    local skill_md
    local agents_yaml

    skill_name="$(basename "$skill_dir")"
    skill_md="$skill_dir/SKILL.md"

    [[ -f "$skill_md" ]] || fail "Missing SKILL.md for skill '$skill_name'"

    assert_contains "$skill_md" '^---$' "Missing frontmatter delimiter"
    assert_contains "$skill_md" '^name: [A-Za-z0-9._-]+$' "Missing or invalid frontmatter name"
    assert_contains "$skill_md" '^description: .+$' "Missing frontmatter description"

    agents_yaml="$skill_dir/agents/openai.yaml"
    [[ -f "$agents_yaml" ]] || fail "Missing agents/openai.yaml for skill '$skill_name'"

    assert_contains "$agents_yaml" '^interface:$' "Missing interface block in openai.yaml"
    assert_contains "$agents_yaml" '^  display_name: ".+"$' "Missing display_name in openai.yaml"
    assert_contains "$agents_yaml" '^  short_description: ".+"$' "Missing short_description in openai.yaml"

    if grep -Eq '^  default_prompt: ' "$agents_yaml"; then
      assert_contains "$agents_yaml" "\\\$${skill_name}" "default_prompt must mention \$${skill_name}"
    fi
  done < <(find "$source_dir" -mindepth 1 -maxdepth 1 -type d -print0 | sort -z)

  [[ "$found_skill" -eq 1 ]] || fail "No skills found in $source_dir"

  echo "Validated skills: $source_dir"
}

validate_subagents() {
  local source_dir="$1"

  [[ -d "$source_dir" ]] || fail "Subagents source directory not found: $source_dir"

  python3 - "$source_dir" "$SKILLS_SOURCE_DIR" <<'PY'
import sys
from pathlib import Path

try:
    import tomllib
except ModuleNotFoundError as exc:
    raise SystemExit(f"Validation failed: tomllib is required to validate subagents ({exc})")

source_dir = Path(sys.argv[1])
skills_source_dir = Path(sys.argv[2])
files = sorted(source_dir.glob("*.toml"))

if not files:
    raise SystemExit(f"Validation failed: No subagents found in {source_dir}")

seen_names: dict[str, Path] = {}

for path in files:
    try:
        data = tomllib.loads(path.read_text())
    except Exception as exc:
        raise SystemExit(f"Validation failed: malformed TOML in {path}: {exc}")

    for key in ("name", "description", "developer_instructions"):
        value = data.get(key)
        if not isinstance(value, str) or not value.strip():
            raise SystemExit(f"Validation failed: missing or empty '{key}' in {path}")

    name = data["name"].strip()
    expected_stem = path.stem
    if name != expected_stem:
      raise SystemExit(
          f"Validation failed: subagent name '{name}' must match filename '{expected_stem}' in {path}"
      )

    if name in seen_names:
        raise SystemExit(
            f"Validation failed: duplicate subagent name '{name}' in {path} and {seen_names[name]}"
        )
    seen_names[name] = path

    sandbox_mode = data.get("sandbox_mode")
    if sandbox_mode is not None and sandbox_mode not in {"read-only", "workspace-write", "danger-full-access"}:
        raise SystemExit(f"Validation failed: invalid sandbox_mode '{sandbox_mode}' in {path}")

    instructions = data["developer_instructions"]
    if sandbox_mode == "workspace-write" and (
        "Do not edit files unless the parent agent explicitly states the user approved the edit plan."
        not in instructions
    ):
        raise SystemExit(f"Validation failed: workspace-write subagent missing approval gate in {path}")

    if ("reviewer" in name or "explorer" in name) and sandbox_mode != "read-only":
        raise SystemExit(f"Validation failed: reviewer/explorer subagent must be read-only in {path}")

    if "do not revert edits you did not make" not in instructions.lower():
        raise SystemExit(f"Validation failed: subagent must protect unrelated edits in {path}")

    skills = data.get("skills")
    if skills is not None:
        if not isinstance(skills, dict):
            raise SystemExit(f"Validation failed: 'skills' must be a table in {path}")

        config = skills.get("config")
        if config is not None:
            if not isinstance(config, list) or not config:
                raise SystemExit(f"Validation failed: 'skills.config' must be a non-empty array in {path}")

            for index, entry in enumerate(config, start=1):
                if not isinstance(entry, dict):
                    raise SystemExit(
                        f"Validation failed: skills.config entry #{index} must be a table in {path}"
                    )

                skill_path = entry.get("path")
                if not isinstance(skill_path, str) or not skill_path.strip():
                    raise SystemExit(
                        f"Validation failed: skills.config entry #{index} missing 'path' in {path}"
                    )

                if not (
                    skill_path.startswith("__CODEX_HOME__/skills/")
                    or skill_path.startswith("/")
                ):
                    raise SystemExit(
                        f"Validation failed: skills.config path must start with __CODEX_HOME__/skills/ or / in {path}"
                    )

                if skill_path.startswith("__CODEX_HOME__/skills/"):
                    parts = Path(skill_path.replace("__CODEX_HOME__/skills/", "", 1)).parts
                    if len(parts) != 2 or parts[1] != "SKILL.md":
                        raise SystemExit(
                            f"Validation failed: skills.config path must target __CODEX_HOME__/skills/<skill>/SKILL.md in {path}"
                        )

                    source_skill = skills_source_dir / parts[0] / "SKILL.md"
                    if not source_skill.exists():
                        raise SystemExit(
                            f"Validation failed: skills.config references missing repo skill '{parts[0]}' in {path}"
                        )
                elif not Path(skill_path).exists():
                    raise SystemExit(
                        f"Validation failed: skills.config references missing absolute skill path '{skill_path}' in {path}"
                    )

                enabled = entry.get("enabled")
                if enabled is not None and not isinstance(enabled, bool):
                    raise SystemExit(
                        f"Validation failed: skills.config entry #{index} has non-boolean 'enabled' in {path}"
                    )

print(f"Validated subagents: {source_dir}")
PY
}

validate_skills "$SKILLS_SOURCE_DIR"
validate_subagents "$SUBAGENTS_SOURCE_DIR"

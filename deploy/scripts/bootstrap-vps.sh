#!/bin/bash
set -euo pipefail

echo "bootstrap-vps.sh is deprecated."
echo "Use deploy/scripts/setup-github-secrets.sh locally, then push the deploy branch so GitHub Actions runs deploy/scripts/provision-and-deploy.sh on the VPS."
exit 1

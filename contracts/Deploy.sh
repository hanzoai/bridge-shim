#!/usr/bin/env bash
# Deploy.sh — Hanzo Bridge contract deployment driver.
#
# Wraps @luxfi/standard's DeployTenant.s.sol with this tenant's
# manifest at contracts/tenant.json.
#
# Usage:
#   HANZO_PRIVATE_KEY=0x... ./contracts/Deploy.sh <rpc-url>
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "usage: $0 <rpc-url>" >&2
  exit 2
fi

RPC_URL="$1"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TENANT_FILE="${SCRIPT_DIR}/tenant.json"
STANDARD_DIR="${FOUNDRY_LIB_PATH:-${SCRIPT_DIR}/lib/luxfi-standard}"

if [[ ! -f "${TENANT_FILE}" ]]; then
  echo "error: tenant manifest missing at ${TENANT_FILE}" >&2
  exit 1
fi

if [[ ! -d "${STANDARD_DIR}" ]]; then
  echo "error: @luxfi/standard not found at ${STANDARD_DIR}" >&2
  echo "       git clone --branch v1.7.5 https://github.com/luxfi/standard ${STANDARD_DIR}" >&2
  exit 1
fi

if [[ -z "${HANZO_PRIVATE_KEY:-}" ]]; then
  echo "error: HANZO_PRIVATE_KEY env var required" >&2
  exit 1
fi

cd "${STANDARD_DIR}"
TENANT_FILE="${TENANT_FILE}" \
LUX_PRIVATE_KEY="${HANZO_PRIVATE_KEY}" \
  forge script contracts/script/DeployTenant.s.sol \
    --rpc-url "${RPC_URL}" \
    --broadcast \
    --legacy \
    -vvv

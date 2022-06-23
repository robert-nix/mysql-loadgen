#!/bin/bash
set -euo pipefail

./scripts/run-mysql.sh --version
mkdir -p db/run db/data
./scripts/run-mysql.sh --initialize-insecure

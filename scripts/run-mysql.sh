#!/bin/bash
set -euo pipefail

if [ "$(dirname $0)" != "./scripts" ]; then
	echo "Please run this script from the root directory of the repo; i.e. ./scripts/run-mysql.sh"
	exit 1
fi

if [ -z "$MYSQLD" -o ! -e "$MYSQLD" ]; then
	echo "Please set \$MYSQLD to a path to a mysqld binary."
	exit 1
fi

ulimit -n 1000000
export ICU_DATA=/usr/lib/mysql/private
exec $MYSQLD --defaults-file=./fixtures/my.cnf --user `whoami` "$@"

#!/bin/bash
set -euo pipefail

echo "create user 'root'@'%'; grant all on *.* to 'root'@'%';" | mysql -u root -S ./db/run/mysqld.sock

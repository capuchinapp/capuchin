#!/bin/sh
set -eu

db_path="${SQLITE_DB_PATH:-./capuchin.db}"
data_dir=$(dirname "$db_path")

mkdir -p "$data_dir"
chown -R appuser:appgroup "$data_dir"

if [ ! -e "$db_path" ]; then
	su-exec appuser:appgroup touch "$db_path"
fi

exec su-exec appuser:appgroup sh -c 'goose up && exec ./capuchin'
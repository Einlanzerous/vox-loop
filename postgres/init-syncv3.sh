#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER syncv3 WITH PASSWORD '${SYNCV3_DB_PASSWORD}';
    CREATE DATABASE syncv3 OWNER syncv3;
EOSQL

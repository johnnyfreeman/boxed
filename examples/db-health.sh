#!/bin/bash
set -euo pipefail

BOXED="${BOXED:-./boxed}"

usage() {
    echo "Usage: DB_TYPE=<postgres|mysql> DB_HOST=<host> DB_USER=<user> DB_NAME=<database> $0"
    echo ""
    echo "Example:"
    echo "  DB_TYPE=postgres DB_HOST=localhost DB_USER=admin DB_NAME=mydb $0"
    exit 1
}

DB_TYPE="${DB_TYPE:-}"
DB_HOST="${DB_HOST:-localhost}"
DB_USER="${DB_USER:-}"
DB_NAME="${DB_NAME:-}"

if [ -z "$DB_TYPE" ] || [ -z "$DB_USER" ] || [ -z "$DB_NAME" ]; then
    usage
fi

check_postgres() {
    local conn_count=$(psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT count(*) FROM pg_stat_activity;" 2>/dev/null | tr -d ' ')
    local db_size=$(psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT pg_size_pretty(pg_database_size('$DB_NAME'));" 2>/dev/null | tr -d ' ')
    local max_conn=$(psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -t -c "SHOW max_connections;" 2>/dev/null | tr -d ' ')

    echo "$conn_count:$max_conn:$db_size"
}

check_mysql() {
    local conn_count=$(mysql -h "$DB_HOST" -u "$DB_USER" -D "$DB_NAME" -sN -e "SELECT COUNT(*) FROM information_schema.processlist;" 2>/dev/null)
    local db_size=$(mysql -h "$DB_HOST" -u "$DB_USER" -D "$DB_NAME" -sN -e "SELECT ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) FROM information_schema.tables WHERE table_schema='$DB_NAME';" 2>/dev/null)
    local max_conn=$(mysql -h "$DB_HOST" -u "$DB_USER" -sN -e "SHOW VARIABLES LIKE 'max_connections';" 2>/dev/null | awk '{print $2}')

    echo "$conn_count:$max_conn:${db_size}MB"
}

case "$DB_TYPE" in
    postgres)
        if ! command -v psql &> /dev/null; then
            echo "Error: psql not found"
            exit 1
        fi
        result=$(check_postgres)
        ;;
    mysql)
        if ! command -v mysql &> /dev/null; then
            echo "Error: mysql not found"
            exit 1
        fi
        result=$(check_mysql)
        ;;
    *)
        echo "Error: Unsupported database type: $DB_TYPE"
        usage
        ;;
esac

conn_count=$(echo "$result" | cut -d: -f1)
max_conn=$(echo "$result" | cut -d: -f2)
db_size=$(echo "$result" | cut -d: -f3)

conn_usage=$(awk "BEGIN {printf \"%.1f\", ($conn_count/$max_conn)*100}")

determine_status() {
    local usage=$(echo "$1" | cut -d. -f1)

    if [ "$usage" -gt 90 ]; then
        echo "error"
    elif [ "$usage" -gt 75 ]; then
        echo "warning"
    else
        echo "success"
    fi
}

box_type=$(determine_status "$conn_usage")

case "$box_type" in
    error)
        subtitle="$conn_usage% pool usage"
        ;;
    warning)
        subtitle="$conn_usage% pool usage"
        ;;
    *)
        subtitle="$conn_count/$max_conn"
        ;;
esac

$BOXED "$box_type" \
    --title "Database Health: $DB_NAME" \
    --subtitle "$subtitle" \
    --kv "Type=$DB_TYPE" \
    --kv "Host=$DB_HOST" \
    --kv "Connections=$conn_count/$max_conn ($conn_usage%)" \
    --kv "Size=$db_size" \
    --footer "Checked at $(date '+%Y-%m-%d %H:%M:%S')"

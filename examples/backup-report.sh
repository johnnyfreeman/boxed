#!/bin/bash
set -euo pipefail

BOXED="${BOXED:-./boxed}"

usage() {
    echo "Usage: $0 <success|error> <backup_file> [destination]"
    echo ""
    echo "Example:"
    echo "  $0 success /backups/db-2025-10-19.tar.gz \"s3://my-bucket/backups\""
    echo "  $0 error /backups/db-2025-10-19.tar.gz"
    exit 1
}

if [ $# -lt 2 ]; then
    usage
fi

status="$1"
backup_file="$2"
destination="${3:-local}"

get_file_size() {
    if [ -f "$1" ]; then
        du -h "$1" | cut -f1
    else
        echo "N/A"
    fi
}

get_file_age() {
    if [ -f "$1" ]; then
        local mod_time=$(stat -c %Y "$1" 2>/dev/null || stat -f %m "$1" 2>/dev/null)
        local current_time=$(date +%s)
        local age=$((current_time - mod_time))
        local hours=$((age / 3600))
        local minutes=$(((age % 3600) / 60))

        if [ "$hours" -gt 0 ]; then
            echo "${hours}h ${minutes}m ago"
        else
            echo "${minutes}m ago"
        fi
    else
        echo "N/A"
    fi
}

backup_name=$(basename "$backup_file")
file_size=$(get_file_size "$backup_file")
last_backup=$(get_file_age "$backup_file")

case "$status" in
    success)
        box_type="success"
        title="Backup Complete"
        subtitle="$backup_name"

        $BOXED "$box_type" \
            --title "$title" \
            --subtitle "$subtitle" \
            --kv "Size=$file_size" \
            --kv "Created=$last_backup" \
            --kv "Destination=$destination" \
            --kv "Status=Successfully backed up" \
            --footer "Backup completed at $(date '+%Y-%m-%d %H:%M:%S')"
        ;;
    error)
        box_type="error"
        title="Backup Failed"
        subtitle="$backup_name"

        $BOXED "$box_type" \
            --title "$title" \
            --subtitle "$subtitle" \
            --kv "Size=$file_size" \
            --kv "Attempted=$last_backup" \
            --kv "Destination=$destination" \
            --kv "Status=Backup operation failed" \
            --footer "Failed at $(date '+%Y-%m-%d %H:%M:%S')"
        ;;
    *)
        echo "Error: Invalid status. Use 'success' or 'error'"
        usage
        ;;
esac

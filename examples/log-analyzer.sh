#!/bin/bash
set -euo pipefail

BOXED="${BOXED:-./boxed}"

usage() {
    echo "Usage: $0 <log_file> [lines_to_check]"
    echo ""
    echo "Example:"
    echo "  $0 /var/log/syslog 1000"
    echo "  $0 /var/log/application.log"
    exit 1
}

if [ $# -lt 1 ]; then
    usage
fi

log_file="$1"
lines_to_check="${2:-500}"

if [ ! -f "$log_file" ]; then
    echo "Error: Log file not found: $log_file"
    exit 1
fi

count_pattern() {
    tail -n "$lines_to_check" "$log_file" | grep -ci "$1" 2>/dev/null || echo "0"
}

errors=$(count_pattern "error")
warnings=$(count_pattern "warning")
critical=$(count_pattern "critical\|fatal")
total_lines=$(tail -n "$lines_to_check" "$log_file" | wc -l)

determine_status() {
    if [ "$critical" -gt 0 ]; then
        echo "error"
    elif [ "$errors" -gt 10 ]; then
        echo "warning"
    else
        echo "info"
    fi
}

box_type=$(determine_status)

case "$box_type" in
    error)
        subtitle="Critical Issues Found"
        ;;
    warning)
        subtitle="Multiple Errors Detected"
        ;;
    *)
        subtitle="Log Summary"
        ;;
esac

$BOXED "$box_type" \
    --title "Log Analysis" \
    --subtitle "$subtitle" \
    --kv "File=$(basename "$log_file")" \
    --kv "Analyzed=$total_lines lines" \
    --kv "Critical=$critical occurrences" \
    --kv "Errors=$errors occurrences" \
    --kv "Warnings=$warnings occurrences" \
    --footer "Analyzed at $(date '+%Y-%m-%d %H:%M:%S')"

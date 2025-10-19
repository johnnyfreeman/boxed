#!/bin/bash
set -euo pipefail

BOXED="${BOXED:-./boxed}"

usage() {
    echo "Usage: $0 <passed> <failed> <skipped> [duration] [coverage]"
    echo ""
    echo "Example:"
    echo "  $0 45 0 2 \"3.2s\" \"87.5%\""
    echo "  $0 42 3 1 \"2.8s\" \"82.1%\""
    exit 1
}

if [ $# -lt 3 ]; then
    usage
fi

passed="$1"
failed="$2"
skipped="$3"
duration="${4:-unknown}"
coverage="${5:-N/A}"

total=$((passed + failed + skipped))

determine_status() {
    if [ "$failed" -gt 0 ]; then
        echo "error"
    elif [ "$skipped" -gt "$passed" ]; then
        echo "warning"
    else
        echo "success"
    fi
}

box_type=$(determine_status)

case "$box_type" in
    error)
        title="Test Results"
        subtitle="Tests Failed"
        ;;
    warning)
        title="Test Results"
        subtitle="Many Tests Skipped"
        ;;
    *)
        title="Test Results"
        subtitle="All Tests Passed"
        ;;
esac

$BOXED "$box_type" \
    --title "$title" \
    --subtitle "$subtitle" \
    --kv "Total=$total tests" \
    --kv "Passed=$passed tests" \
    --kv "Failed=$failed tests" \
    --kv "Skipped=$skipped tests" \
    --kv "Duration=$duration" \
    --kv "Coverage=$coverage" \
    --footer "Test run completed at $(date '+%Y-%m-%d %H:%M:%S')" \
    --exit-on-error

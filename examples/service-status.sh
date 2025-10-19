#!/bin/bash
set -euo pipefail

BOXED="${BOXED:-./boxed}"

usage() {
    echo "Usage: $0 <service1> [service2] [service3] ..."
    echo ""
    echo "Example:"
    echo "  $0 docker nginx postgresql"
    exit 1
}

if [ $# -lt 1 ]; then
    usage
fi

check_service() {
    local service="$1"

    if systemctl is-active --quiet "$service" 2>/dev/null; then
        echo "running"
    else
        echo "stopped"
    fi
}

declare -A service_status
total_services=$#
running=0
stopped=0

for service in "$@"; do
    status=$(check_service "$service")
    service_status["$service"]="$status"

    if [ "$status" = "running" ]; then
        ((running++))
    else
        ((stopped++))
    fi
done

determine_status() {
    if [ "$stopped" -gt 0 ]; then
        echo "error"
    else
        echo "success"
    fi
}

box_type=$(determine_status)

case "$box_type" in
    error)
        subtitle="$stopped/$total_services Services Stopped"
        ;;
    *)
        subtitle="All Services Running"
        ;;
esac

for service in "$@"; do
    status="${service_status[$service]}"
    echo "$service=$status"
done | $BOXED "$box_type" \
    --title "Service Status" \
    --subtitle "$subtitle" \
    --stdin-kv \
    --footer "Checked at $(date '+%Y-%m-%d %H:%M:%S')"

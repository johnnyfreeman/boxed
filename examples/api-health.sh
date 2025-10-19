#!/bin/bash
set -euo pipefail

BOXED="${BOXED:-./boxed}"

declare -A endpoints
declare -A statuses
declare -A response_times

usage() {
    echo "Usage: $0 <endpoint1> [endpoint2] [endpoint3] ..."
    echo ""
    echo "Example:"
    echo "  $0 https://api.example.com/health https://api.example.com/status"
    exit 1
}

if [ $# -lt 1 ]; then
    usage
fi

check_endpoint() {
    local url="$1"
    local start=$(date +%s%N)
    local http_code

    http_code=$(curl -s -o /dev/null -w "%{http_code}" --max-time 5 "$url" 2>/dev/null || echo "000")

    local end=$(date +%s%N)
    local duration=$(((end - start) / 1000000))

    echo "$http_code:$duration"
}

total_endpoints=$#
healthy=0
unhealthy=0

echo "Checking $total_endpoints endpoints..."

for url in "$@"; do
    result=$(check_endpoint "$url")
    http_code=$(echo "$result" | cut -d: -f1)
    response_time=$(echo "$result" | cut -d: -f2)

    endpoints["$url"]="$http_code"
    response_times["$url"]="${response_time}ms"

    if [ "$http_code" = "200" ] || [ "$http_code" = "204" ]; then
        statuses["$url"]="healthy"
        ((healthy++))
    else
        statuses["$url"]="down"
        ((unhealthy++))
    fi
done

determine_status() {
    if [ "$unhealthy" -gt 0 ]; then
        echo "error"
    else
        echo "success"
    fi
}

box_type=$(determine_status)

case "$box_type" in
    error)
        title="API Health Check"
        subtitle="$unhealthy/$total_endpoints down"
        ;;
    *)
        title="API Health Check"
        subtitle="$healthy/$total_endpoints healthy"
        ;;
esac

for url in "$@"; do
    short_url=$(echo "$url" | sed 's|https://||' | sed 's|http://||' | cut -c1-40)
    status="${statuses[$url]}"
    response_time="${response_times[$url]}"

    echo "$short_url=$status ($response_time)"
done | $BOXED "$box_type" \
    --title "$title" \
    --subtitle "$subtitle" \
    --stdin-kv \
    --footer "Health check completed at $(date '+%Y-%m-%d %H:%M:%S')"

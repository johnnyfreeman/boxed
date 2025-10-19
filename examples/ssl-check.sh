#!/bin/bash
set -euo pipefail

BOXED="${BOXED:-./boxed}"

usage() {
    echo "Usage: $0 <domain1> [domain2] [domain3] ..."
    echo ""
    echo "Example:"
    echo "  $0 example.com api.example.com"
    exit 1
}

if [ $# -lt 1 ]; then
    usage
fi

check_cert() {
    local domain="$1"
    local expiry_date

    expiry_date=$(echo | openssl s_client -servername "$domain" -connect "$domain:443" 2>/dev/null | \
        openssl x509 -noout -enddate 2>/dev/null | cut -d= -f2)

    if [ -z "$expiry_date" ]; then
        echo "unknown:999"
        return
    fi

    local expiry_epoch=$(date -d "$expiry_date" +%s 2>/dev/null || date -j -f "%b %d %T %Y %Z" "$expiry_date" +%s 2>/dev/null)
    local current_epoch=$(date +%s)
    local days_left=$(( (expiry_epoch - current_epoch) / 86400 ))

    echo "$expiry_date:$days_left"
}

declare -A cert_info
total_domains=$#
expiring_soon=0
expired=0

echo "Checking $total_domains SSL certificates..."

for domain in "$@"; do
    result=$(check_cert "$domain")
    expiry_date=$(echo "$result" | cut -d: -f1)
    days_left=$(echo "$result" | cut -d: -f2)

    cert_info["$domain"]="$days_left days"

    if [ "$days_left" -lt 0 ]; then
        ((expired++))
    elif [ "$days_left" -lt 30 ]; then
        ((expiring_soon++))
    fi
done

determine_status() {
    if [ "$expired" -gt 0 ]; then
        echo "error"
    elif [ "$expiring_soon" -gt 0 ]; then
        echo "warning"
    else
        echo "success"
    fi
}

box_type=$(determine_status)

case "$box_type" in
    error)
        subtitle="$expired expired"
        ;;
    warning)
        subtitle="$expiring_soon expiring <30d"
        ;;
    *)
        subtitle="$total_domains valid"
        ;;
esac

for domain in "$@"; do
    days="${cert_info[$domain]}"
    echo "$domain=$days remaining"
done | $BOXED "$box_type" \
    --title "SSL Certificate Check" \
    --subtitle "$subtitle" \
    --stdin-kv \
    --footer "Checked at $(date '+%Y-%m-%d %H:%M:%S')"

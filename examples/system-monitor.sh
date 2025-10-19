#!/bin/bash
set -euo pipefail

BOXED="${BOXED:-./boxed}"

get_cpu_usage() {
    top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* id.*/\1/" | awk '{print 100 - $1"%"}'
}

get_memory_usage() {
    free | grep Mem | awk '{printf "%.1f%%", $3/$2 * 100.0}'
}

get_disk_usage() {
    df -h / | awk 'NR==2 {print $5}'
}

get_load_average() {
    uptime | awk -F'load average:' '{print $2}' | awk '{print $1}' | tr -d ','
}

get_uptime() {
    uptime -p | sed 's/up //'
}

determine_status() {
    local cpu_val=$(echo "$1" | sed 's/%//')
    local mem_val=$(echo "$2" | sed 's/%//')
    local disk_val=$(echo "$3" | sed 's/%//')

    if (( $(echo "$cpu_val > 90" | bc -l) )) || \
       (( $(echo "$mem_val > 90" | bc -l) )) || \
       (( $(echo "$disk_val > 90" | bc -l) )); then
        echo "error"
    elif (( $(echo "$cpu_val > 75" | bc -l) )) || \
         (( $(echo "$mem_val > 75" | bc -l) )) || \
         (( $(echo "$disk_val > 75" | bc -l) )); then
        echo "warning"
    else
        echo "success"
    fi
}

main() {
    cpu=$(get_cpu_usage)
    memory=$(get_memory_usage)
    disk=$(get_disk_usage)
    load=$(get_load_average)
    uptime=$(get_uptime)

    box_type=$(determine_status "$cpu" "$memory" "$disk")

    case "$box_type" in
        error)
            title="System Monitor"
            subtitle="Critical Usage Detected"
            ;;
        warning)
            title="System Monitor"
            subtitle="High Usage Warning"
            ;;
        *)
            title="System Monitor"
            subtitle="All Systems Normal"
            ;;
    esac

    $BOXED "$box_type" \
        --title "$title" \
        --subtitle "$subtitle" \
        --kv "CPU=$cpu" \
        --kv "Memory=$memory" \
        --kv "Disk=$disk" \
        --kv "Load=$load" \
        --kv "Uptime=$uptime" \
        --footer "Checked at $(date '+%Y-%m-%d %H:%M:%S')"
}

main "$@"

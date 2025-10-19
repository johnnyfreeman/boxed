#!/bin/bash
set -euo pipefail

BOXED="${BOXED:-./boxed}"

check_docker() {
    if ! command -v docker &> /dev/null; then
        echo "Error: docker not found in PATH"
        exit 1
    fi

    if ! docker info &> /dev/null; then
        echo "Error: Cannot connect to Docker daemon"
        exit 1
    fi
}

get_container_counts() {
    local running=$(docker ps -q 2>/dev/null | wc -l)
    local stopped=$(docker ps -aq -f status=exited 2>/dev/null | wc -l)
    local total=$(docker ps -aq 2>/dev/null | wc -l)

    echo "${running}/${total} (${stopped} stopped)"
}

get_image_count() {
    docker images -q 2>/dev/null | wc -l
}

get_volume_count() {
    docker volume ls -q 2>/dev/null | wc -l
}

get_network_count() {
    docker network ls -q 2>/dev/null | wc -l
}

get_disk_usage() {
    docker system df --format "{{.Type}}: {{.Size}}" 2>/dev/null | grep "Total" | awk '{print $2}'
}

determine_status() {
    local containers="$1"
    local stopped=$(echo "$containers" | grep -oP '\(\K[0-9]+' || echo "0")

    if [ "$stopped" -gt 10 ]; then
        echo "warning"
    else
        echo "info"
    fi
}

main() {
    check_docker

    containers=$(get_container_counts)
    images=$(get_image_count)
    volumes=$(get_volume_count)
    networks=$(get_network_count)
    disk=$(get_disk_usage || echo "unknown")

    box_type=$(determine_status "$containers")

    stopped=$(echo "$containers" | grep -oP '\(\K[0-9]+' || echo "0")

    case "$box_type" in
        warning)
            subtitle="$stopped stopped"
            ;;
        *)
            subtitle="$containers"
            ;;
    esac

    $BOXED "$box_type" \
        --title "Docker Status" \
        --subtitle "$subtitle" \
        --kv "Containers=$containers" \
        --kv "Images=$images total" \
        --kv "Volumes=$volumes total" \
        --kv "Networks=$networks total" \
        --kv "Disk=$disk used" \
        --footer "Checked at $(date '+%Y-%m-%d %H:%M:%S')"
}

main "$@"

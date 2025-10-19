#!/bin/bash
set -euo pipefail

BOXED="${BOXED:-./boxed}"

check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        echo "Error: kubectl not found in PATH"
        exit 1
    fi
}

check_jq() {
    if ! command -v jq &> /dev/null; then
        echo "Error: jq not found in PATH"
        exit 1
    fi
}

get_cluster_info() {
    echo "Checking cluster connectivity..."
    if ! kubectl cluster-info &> /dev/null; then
        echo "Error: Cannot connect to cluster"
        exit 1
    fi
}

main() {
    check_kubectl
    check_jq
    get_cluster_info

    echo "Gathering cluster metrics..."

    local total_nodes=$(kubectl get nodes --no-headers 2>/dev/null | wc -l)
    local ready_nodes=$(kubectl get nodes --no-headers 2>/dev/null | grep " Ready" | wc -l)
    local total_pods=$(kubectl get pods --all-namespaces --no-headers 2>/dev/null | wc -l)
    local running_pods=$(kubectl get pods --all-namespaces --no-headers 2>/dev/null | grep "Running" | wc -l)
    local failing_count=$(kubectl get pods --all-namespaces --no-headers 2>/dev/null | grep -vE "Running|Completed|Succeeded" | wc -l)
    local namespace_count=$(kubectl get namespaces --no-headers 2>/dev/null | wc -l)

    local box_type="success"
    local subtitle="$ready_nodes/$total_nodes ready"
    if [ "${failing_count:-0}" -gt 0 ]; then
        box_type="warning"
        subtitle="$failing_count failing"
    fi

    # Generate JSON and pipe to boxed
    jq -n \
        --arg title "Cluster Status" \
        --arg subtitle "$subtitle" \
        --arg nodes "$ready_nodes/$total_nodes ready" \
        --arg pods "$running_pods/$total_pods running" \
        --arg failing "$failing_count pods" \
        --arg namespaces "$namespace_count total" \
        --arg footer "Generated at $(date '+%Y-%m-%d %H:%M:%S')" \
        '{
            title: $title,
            subtitle: $subtitle,
            kv: {
                Nodes: $nodes,
                Pods: $pods,
                Failing: $failing,
                Namespaces: $namespaces
            },
            footer: $footer
        }' | $BOXED "$box_type" --json
}

main "$@"

#!/bin/bash
set -euo pipefail

BOXED="${BOXED:-./boxed}"

check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        echo "Error: kubectl not found in PATH"
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

get_node_status() {
    local total_nodes=$(kubectl get nodes --no-headers 2>/dev/null | wc -l)
    local ready_nodes=$(kubectl get nodes --no-headers 2>/dev/null | grep -c " Ready" || echo "0")
    echo "$ready_nodes/$total_nodes"
}

get_pod_status() {
    local total_pods=$(kubectl get pods --all-namespaces --no-headers 2>/dev/null | wc -l)
    local running_pods=$(kubectl get pods --all-namespaces --no-headers 2>/dev/null | grep -c "Running" || echo "0")
    echo "$running_pods/$total_pods"
}

get_failing_pods() {
    local count=$(kubectl get pods --all-namespaces --no-headers 2>/dev/null | \
        grep -vE "Running|Completed|Succeeded" | wc -l)
    echo "${count:-0}"
}

get_namespace_count() {
    kubectl get namespaces --no-headers 2>/dev/null | wc -l
}

main() {
    check_kubectl
    get_cluster_info

    echo "Gathering cluster metrics..."

    node_status=$(get_node_status)
    pod_status=$(get_pod_status)
    failing_pods=$(get_failing_pods)
    namespace_count=$(get_namespace_count)

    if [ "$failing_pods" -gt 0 ]; then
        box_type="warning"
        title="Cluster Status"
        subtitle="Issues Detected"
    else
        box_type="success"
        title="Cluster Status"
        subtitle="All Systems Operational"
    fi

    $BOXED "$box_type" \
        --title "$title" \
        --subtitle "$subtitle" \
        --kv "Nodes=$node_status ready" \
        --kv "Pods=$pod_status running" \
        --kv "Failing=$failing_pods pods" \
        --kv "Namespaces=$namespace_count total" \
        --footer "Generated at $(date '+%Y-%m-%d %H:%M:%S')"
}

main "$@"

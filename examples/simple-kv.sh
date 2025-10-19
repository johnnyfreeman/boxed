#!/bin/bash
set -euo pipefail

BOXED="${BOXED:-./boxed}"

# Demonstrate the three ways to pass KV pairs:

echo "Method 1: Repeated --kv flags"
$BOXED success \
    --title "Deployment Summary" \
    --kv "Environment=Production" \
    --kv "Version=v2.1.0" \
    --kv "Duration=2m 34s"

echo ""
echo "Method 2: Comma-separated in single --kv flag"
$BOXED success \
    --title "Deployment Summary" \
    --kv "Environment=Production,Version=v2.1.0,Duration=2m 34s"

echo ""
echo "Method 3: Mix of both (most flexible)"
$BOXED success \
    --title "Deployment Summary" \
    --kv "Environment=Production,Version=v2.1.0" \
    --kv "Duration=2m 34s,Status=Complete"

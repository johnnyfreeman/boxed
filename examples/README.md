# Examples

Example scripts demonstrating how to use `boxed` in real-world scenarios.

## deploy-status.sh

A deployment status reporter for CI/CD pipelines.

### Usage

```bash
# Success deployment
./examples/deploy-status.sh success v1.2.3 "2m 34s"

# Failed deployment
./examples/deploy-status.sh error v1.2.3 "1m 12s"

# Set custom environment
ENVIRONMENT=staging ./examples/deploy-status.sh success v1.2.4 "3m 10s"
```

### Features

- Automatically includes git commit hash and branch
- Supports custom environment via `ENVIRONMENT` variable
- Displays deployment duration and timestamp
- Color-coded based on success/failure

---

## k8s-status.sh

A Kubernetes cluster health check script that displays a formatted status summary.

### Prerequisites

- `kubectl` installed and configured
- Active Kubernetes cluster context
- `boxed` binary built in the project root

### Usage

```bash
# From the project root
./examples/k8s-status.sh

# Or specify custom boxed path
BOXED=/usr/local/bin/boxed ./examples/k8s-status.sh
```

### What it checks

- Cluster connectivity
- Node status (ready vs total)
- Pod status (running vs total)
- Failing pods count
- Namespace count

### Output

Displays a success box when all systems are operational, or a warning box when issues are detected.

Example output:
```
╭────────────────────────────────────╮
│╱╱ Cluster Status All Systems Op...│
│                                    │
│   Nodes         3/3 ready          │
│                                    │
│   Pods          45/45 running      │
│                                    │
│   Failing       0 pods             │
│                                    │
│   Namespaces    8 total            │
│                                    │
│╱╱ Generated at 2025-10-19 14:30:42│
╰────────────────────────────────────╯
```

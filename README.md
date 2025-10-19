# boxed

<img width="2081" height="676" alt="image" src="https://github.com/user-attachments/assets/4e9de355-2cc4-49a1-b363-0c46ce8bb23a" />

Beautiful bordered boxes for terminal output with Tokyo Night theme colors.

## Install

```bash
go build -o boxed .
```

## Usage

```bash
boxed <type> [flags]
```

**Types:** `success` (lime green), `error` (pink-red), `info` (sky blue), `warning` (golden orange)

Colors inspired by the Tokyo Night theme.

**Flags:**
- `-t, --title` - Box title (bold, colored)
- `-s, --subtitle` - Box subtitle (italic, gray)
- `-k, --kv` - Key-value pairs (repeatable, format: `key=value` or `key1=value1,key2=value2`)
- `-f, --footer` - Box footer (gray)
- `-b, --border-style` - Border style: `rounded`, `normal`, `thick`, `double` (default: rounded)
- `-w, --width` - Box width (0 for auto-size)
- `--stdin-kv` - Read KV pairs from stdin (one per line)
- `--json` - Read box definition from JSON stdin
- `--json-file` - Read box definition from JSON file
- `--exit-on-error` - Exit with code 1 when rendering an error box (for CI/CD)
- `--exit-on-warning` - Exit with code 2 when rendering a warning box (for CI/CD)

## Examples

```bash
# Basic
./boxed success --title "Deploy Complete" --subtitle "v2.1.0" \
  --kv "Duration=2m 34s" --kv "Commit=abc1234" \
  --footer "Deployed at 2025-10-19"
```

<img width="845" height="431" alt="image" src="https://github.com/user-attachments/assets/ccc7497d-6b22-41bc-bea7-480625dda68f" />

```bash
# Error
./boxed error --title "Build Failed" \
  --kv "File=src/main.go" --kv "Line=142" \
  --footer "Check logs"
```

<img width="845" height="431" alt="image" src="https://github.com/user-attachments/assets/52ada984-d3c6-4c05-8a16-3dd9155c2eb2" />

```bash
# From stdin
echo -e "Region=us-east-1\nEnv=prod" | \
  ./boxed info --title "Config" --stdin-kv
```

<img width="845" height="382" alt="image" src="https://github.com/user-attachments/assets/a3eff49d-0bab-408b-b2a2-793cfb472a14" />

```bash
# Unicode support
./boxed success --title "部署完了 ✅" \
  --kv "環境=本番" --kv "Status=🚀 Deployed"
```

<img width="845" height="379" alt="image" src="https://github.com/user-attachments/assets/bf29c2b3-ae17-4319-858e-a2a9d88017e8" />

```bash
# Border styles
./boxed warning --title "Warning" --border-style thick
```

<img width="845" height="191" alt="image" src="https://github.com/user-attachments/assets/b1a55ac8-a875-471e-aba5-d4a2b9d91bc2" />

```bash
# Error with stack trace
./boxed error --title "Runtime Panic" --subtitle "SIGSEGV: segmentation violation" \
  --kv "Error=runtime error: invalid memory address or nil pointer dereference" \
  --kv "goroutine 1 [running]=main.(*Server).handleRequest(0x0, 0xc00012e000)" \
  --kv "Location=src/server/handler.go:145 +0x2a" \
  --kv "Caller=main.processOrder(0xc00012e000, 0xc0001a4000)" \
  --kv "Trace=src/orders/processor.go:89 +0x15c" \
  --kv "Root Cause=Attempted to access uninitialized database connection" \
  --footer "Process terminated with exit code 2"
```

<img width="2102" height="809" alt="image" src="https://github.com/user-attachments/assets/9abe912e-7a57-47d1-a81b-da0bd21d697c" />

```bash
# Long text example
./boxed warning --title "Database Migration Warning" --subtitle "Schema Changes Detected" \
  --kv "Migration File=/db/migrations/20251019_add_user_authentication_and_session_management_tables_with_foreign_key_constraints.sql" \
  --kv "Affected Tables=users, user_sessions, user_roles, role_permissions, authentication_providers, oauth_tokens, password_reset_tokens, email_verification_tokens" \
  --kv "Estimated Duration=This migration will take approximately 15-20 minutes to complete depending on database size and current load. Large tables may require additional time." \
  --kv "Backup Status=Automatic backup created at /backups/prod_db_20251019_143218.sql.gz (size: 2.4GB, compression ratio: 87%, estimated restore time: 8-10 minutes)" \
  --kv "Breaking Changes=This migration includes breaking changes that may affect the following microservices: auth-service, user-service, session-manager, api-gateway, admin-dashboard" \
  --kv "Rollback Plan=Rollback script available at /db/rollbacks/20251019_revert_authentication_changes.sql - tested in staging environment with zero data loss" \
  --footer "Review migration plan at https://docs.example.com/migrations/20251019 before proceeding • Scheduled maintenance window: 2025-10-19 22:00-23:00 UTC"
```

<img width="2365" height="1196" alt="image" src="https://github.com/user-attachments/assets/93b32e91-9943-46ae-bd55-162bc8ab5533" />

## Advanced Features

### Comma-separated KV syntax

For convenience, you can pass multiple KV pairs in a single flag using comma separation:

```bash
# Traditional syntax
./boxed success --title "Deploy" --kv A=1 --kv B=2 --kv C=3

# Comma-separated (shorter)
./boxed success --title "Deploy" --kv A=1,B=2,C=3

# Mix both styles
./boxed success --title "Deploy" --kv A=1,B=2 --kv C=3,D=4
```

### JSON input

Define your entire box configuration in JSON, perfect for programmatic generation:

```bash
# From stdin
echo '{"title":"Status","subtitle":"All Good","kv":{"CPU":"15%","Memory":"28%"}}' | \
  ./boxed success --json

# From file
./boxed info --json-file status.json

# Combine JSON with CLI overrides (CLI takes precedence)
./boxed success --json-file base.json --title "Override Title"
```

JSON format:
```json
{
  "title": "System Status",
  "subtitle": "All Systems Normal",
  "kv": {
    "CPU": "15%",
    "Memory": "28%",
    "Disk": "45%"
  },
  "footer": "Updated 2025-10-19",
  "width": 50,
  "border_style": "rounded"
}
```

### Exit codes for CI/CD

Use exit codes to integrate with CI/CD pipelines and fail builds based on box type:

```bash
# Exit with code 1 if error box is rendered (fails CI build)
./boxed error --title "Tests Failed" --exit-on-error
echo $?  # 1

# Exit with code 2 if warning box is rendered
./boxed warning --title "Warnings Found" --exit-on-warning
echo $?  # 2

# Success boxes always exit 0
./boxed success --title "All Good" --exit-on-error
echo $?  # 0
```

Example in CI:
```bash
#!/bin/bash
if ./tests/run.sh; then
  ./boxed success --title "Tests Passed" --kv "Duration=2.3s"
else
  ./boxed error --title "Tests Failed" --exit-on-error
fi
```

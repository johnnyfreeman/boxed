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
- `-k, --kv` - Key-value pairs (repeatable, format: `key=value`)
- `-f, --footer` - Box footer (gray)
- `-b, --border-style` - Border style: `rounded`, `normal`, `thick`, `double` (default: rounded)
- `-w, --width` - Box width (0 for auto-size)
- `--stdin-kv` - Read KV pairs from stdin

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

# SecretWatch

Track secret rotation health across **AWS Secrets Manager**, **HashiCorp Vault**, **GCP Secret Manager**, and **GitHub Actions Secrets** — from a single CLI and web dashboard.

```
STATUS    NAME                    STORE          LAST ROTATED  EXPIRY
------    ----                    -----          ------------  ------
✗ overdue  prod-db-password       aws/us-east-1  2026-06-01    5 days overdue
!  warning  jwt-signing-key        vault          2026-07-15    10 days left
✓ ok       stripe-webhook-secret  aws/us-east-1  2026-08-25    75 days left
?  unknown  gcp-service-account    gcp            never         —
```

## Why

Secret rotation is one of those things everyone knows matters but nobody has visibility into across their whole stack. AWS console shows AWS secrets. Vault UI shows Vault secrets. GCP console shows GCP secrets. Nothing shows all of them together, with a clear answer to "what's overdue and who last rotated it?"

SecretWatch does one thing: give you that view.

## Features

- **Multi-store scanning** — AWS Secrets Manager (multi-region), HashiCorp Vault KV v2 (recursive), GCP Secret Manager, GitHub org secrets
- **Rotation policy** — configurable max age per secret or pattern (`prod-* = 30 days`, `dev-* = 180 days`)
- **Status tracking** — `ok`, `warning` (≤14 days until expiry), `overdue`, `unknown`
- **Audit history** — who rotated a secret and when, pulled from CloudTrail (AWS), Cloud Audit Logs (GCP), GitHub Audit Log
- **Web dashboard** — dark-mode UI with filtering by store and status
- **SQLite cache** — `scan` hits your cloud APIs; `list` and `serve` read locally, so the dashboard is always fast
- **JSON API** — `/api/secrets` and `/api/history` for scripting and alerting integrations

## Install

**Go:**
```bash
go install github.com/leylmordor/secretwatch/cmd/secretwatch@latest
```

**Download binary:** see [Releases](../../releases).

## Quick start

```bash
# Copy the example config
mkdir -p ~/.secretwatch
cp config.example.yaml ~/.secretwatch/config.yaml

# Edit — add credentials and configure the stores you use
$EDITOR ~/.secretwatch/config.yaml

# Pull from all stores and cache results
secretwatch scan

# Show what's overdue
secretwatch list --status overdue

# Open the web dashboard
secretwatch serve
# → http://localhost:8080

# Show who rotated a specific secret and when
secretwatch history prod-db-password
```

## Configuration

```yaml
stores:
  # AWS Secrets Manager — one entry per account/region group
  aws:
    - regions:
        - us-east-1
        - us-west-2
      # profile: my-aws-profile  # optional; uses default credential chain if omitted

  vault:
    enabled: true
    address: http://127.0.0.1:8200
    token: ${VAULT_TOKEN}        # env var expansion supported
    paths:
      - secret/

  # GCP Secret Manager — one entry per project
  gcp:
    - project: my-gcp-project    # uses Application Default Credentials

  github:
    enabled: true
    org: my-github-org
    token: ${GITHUB_TOKEN}

rotation_policy:
  default_max_age_days: 90
  overrides:
    - pattern: "prod-*"
      max_age_days: 30
    - pattern: "dev-*"
      max_age_days: 180

web:
  port: 8080
```

Config file location defaults to `~/.secretwatch/config.yaml`. Override with `--config`.

**Enabling stores:** AWS, AWS IAM, GCP, and GCP SA Keys are list-based — include entries to enable them, omit the section entirely to disable. Vault and GitHub are single objects with an `enabled:` flag. `~` is expanded in `db_path`.

## Credentials

SecretWatch is read-only. It never modifies secrets.

| Store | Auth |
|---|---|
| AWS | Default credential chain (`AWS_PROFILE`, instance role, etc.) |
| Vault | `VAULT_TOKEN` env var or `token` in config |
| GCP | Application Default Credentials (`gcloud auth application-default login`) |
| GitHub | `GITHUB_TOKEN` env var or `token` in config (org admin required for audit log) |

## Audit history

```bash
secretwatch history prod-db-password
```

```
TIMESTAMP              ACTOR                        ACTION          STORE             NOTE
---------              -----                        ------          -----             ----
2026-08-15 09:12:44   arn:aws:iam::123:user/alice  RotateSecret    cloudtrail/us-east-1  automated rotation
2026-05-20 14:30:01   arn:aws:iam::123:user/bob    PutSecretValue  cloudtrail/us-east-1
```

**Prerequisites by store:**
- **AWS** — CloudTrail must be enabled (it is by default in most accounts). 90-day window.
- **GCP** — Data Access audit logs must be explicitly enabled for Secret Manager in your project (off by default).
- **GitHub** — org admin token required. 90-day retention on free plans.
- **Vault** — not supported via API; audit log is written to file/syslog on the Vault server.

## Deployment

### GitHub Actions (recommended — zero infra)

Copy [`examples/secretwatch-scheduled.yml`](examples/secretwatch-scheduled.yml) into your repo at `.github/workflows/secretwatch.yml`, fill in your org/project names, add your credentials as GitHub Secrets, and you're done. SecretWatch runs daily and posts to Slack automatically.

Required secrets: `SLACK_WEBHOOK_URL`, plus whichever of `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` / `GOOGLE_CREDENTIALS` / `VAULT_TOKEN` your stores need.

### Docker

```bash
# Pull
docker pull ghcr.io/leylmordor/secretwatch:latest

# Scan and alert (one-off)
docker run --rm \
  -v ~/.secretwatch:/root/.secretwatch \
  -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY \
  -e SLACK_WEBHOOK_URL \
  ghcr.io/leylmordor/secretwatch:latest \
  sh -c "secretwatch scan && secretwatch alert"

# Always-on dashboard
docker compose up dashboard
```

### Kubernetes CronJob

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: secretwatch
spec:
  schedule: "0 8 * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: secretwatch
            image: ghcr.io/leylmordor/secretwatch:latest
            command: ["sh", "-c", "secretwatch scan && secretwatch alert"]
            envFrom:
            - secretRef:
                name: secretwatch-env
          restartPolicy: OnFailure
```

## Commands

```
secretwatch scan                        Scan all configured stores, cache to SQLite
secretwatch list                        List cached secrets
secretwatch list --status overdue       Filter by status (ok, warning, overdue, unknown)
secretwatch list --store aws            Filter by store
secretwatch history <secret-name>       Show rotation history from audit logs
secretwatch serve                       Start web dashboard (default: http://localhost:8080)
```

Global flag: `--config <path>` to use a non-default config file.

## Supported stores

| Store | Scan | History |
|---|---|---|
| AWS Secrets Manager | ✓ | ✓ (CloudTrail) |
| HashiCorp Vault KV v2 | ✓ | — |
| GCP Secret Manager | ✓ | ✓ (Cloud Audit Logs) |
| GitHub Actions Secrets | ✓ | ✓ (Audit Log API) |
| Azure Key Vault | planned | planned |
| 1Password | planned | — |

## Contributing

Pull requests welcome. To add a new store, implement the `store.Store` interface:

```go
type Store interface {
    Type() string
    List(ctx context.Context) ([]*Secret, error)
}
```

And optionally the `history.Fetcher` interface for audit log support:

```go
type Fetcher interface {
    StoreType() string
    Fetch(ctx context.Context, secretName string, limit int) ([]*Entry, error)
}
```

See `internal/store/aws.go` and `internal/history/aws.go` as reference implementations.

## License

MIT

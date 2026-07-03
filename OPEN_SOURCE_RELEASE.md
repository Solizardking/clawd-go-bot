# Open Source Release Checklist

Use this checklist before pushing a public release.

## Required Gate

```bash
make release-check
```

This checks Go formatting, `go vet`, race tests, release entrypoint builds, and
tracked generated artifacts.

## Repository Hygiene

- Keep generated binaries and caches out of git: `.cache/`, `build/`, `dist/`,
  root `clawdbot`, `**/target/`, `**/.next/`, and `*.tsbuildinfo`.
- Keep `.env` local. Commit only `.env.example`.
- Keep live wallets, treasury keypairs, install ledgers, funding receipts, and
  private API keys outside the repository.
- Prefer `CLAWDBOT_SOURCE_MODE=archive` for one-shot installs so
  `.gitattributes export-ignore` keeps downloads small.

## Security Defaults

- The web console binds to `127.0.0.1` unless `--public` is passed.
- CORS only allows same-origin requests unless `CLAWDBOT_CORS_ORIGINS` is set.
- Proxy IP headers are ignored unless `CLAWDBOT_TRUST_PROXY_HEADERS=1` is set.
- `/api/config` returns redacted secrets unless `CLAWDBOT_WEB_EXPOSE_SECRETS=1`
  is set for a trusted local session.

## Publish Notes

Before announcing a release, capture the output of:

```bash
go version
make release-check
```

If `govulncheck` is installed, include the result of:

```bash
govulncheck ./...
```

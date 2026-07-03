# Security Policy

## Supported Surface

Security fixes target the current `main` branch of `clawdbot-go`.

## Reporting

Report vulnerabilities privately by opening a GitHub security advisory for the
repository or by contacting the maintainer listed on the GitHub profile. Do not
publish exploit details until a fix is available.

## Secrets

Never commit real API keys, wallet keypairs, private keys, install ledgers, or
generated runtime state. The repository is configured to ignore `.env`, build
caches, binaries, Rust `target/`, Next `.next/`, and TypeScript incremental
artifacts.

The web console redacts config secrets by default. Only set
`CLAWDBOT_WEB_EXPOSE_SECRETS=1` for a trusted local debugging session.

## Local Verification

Run the hardened verification path before publishing changes:

```bash
make release-check
```

When available, also run:

```bash
govulncheck ./...
```

Trading and funding paths must remain dry-run by default. Any live transfer or
trade must require explicit user opt-in and bounded limits.

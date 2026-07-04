# Security Policy

## Reporting Security Issues

Do not open public issues with secrets, exploit details, production credentials, or sensitive customer data. Report privately to the repository owner.

## Secret Handling

- Never commit real credentials.
- Never store real credentials in `.env` files committed to the repository.
- Keep `.env.example` limited to placeholder values.
- Rotate any credential that may have been exposed.

## Supported Security Baseline

- Review dependency changes before merging.
- Keep administrative changes separate from runtime code changes.
- Require explicit confirmation before destructive operations.

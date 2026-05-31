# Governance

This repository uses documentation-first governance for administrative, security, and delivery decisions.

## Scope

- Keep code changes separate from documentation and repository administration work.
- Prefer small, reversible pull requests.
- Record assumptions, risks, and validation steps in every non-trivial PR.
- Do not claim a feature is complete unless the code and tests support that status.

## Security Rules

- Do not commit API keys, tokens, passwords, credentials, private certificates, or production data.
- Do not create `.env` files with real secrets.
- Use `.env.example` only for variable names and safe placeholder values.
- Consume credentials from the approved secure vault or runtime environment.
- Review logs, screenshots, and generated reports before publishing them.

## Change Control

- Destructive changes require explicit confirmation and a rollback plan.
- Runtime behavior changes require tests or a documented reason when tests are not feasible.
- Dependency changes require a risk note when they affect production execution.
- Documentation-only changes should state that no runtime behavior was changed.

## Audit Requirements

Every PR should include:

- Summary of changes.
- Implementation status.
- Tests or checks run.
- Risk notes.
- Rollback plan.

# Git Workflow

## Branches

- `main`: stable branch. No direct administrative work unless explicitly approved.
- `develop`: optional integration branch for repositories that need staged releases.
- `feature/*`: new functionality.
- `fix/*`: bug fixes.
- `docs/*`: documentation-only changes.
- `chore/*`: maintenance and repository administration.

## Pull Requests

- Open a pull request for every meaningful change.
- Keep code changes separate from repository-organization changes.
- Use the PR template when available.
- Include verification commands and rollback notes.
- Do not merge automatically unless the repository owner requested it.

## Commits

- Prefer focused commits with clear prefixes: `docs:`, `chore:`, `fix:`, `feat:`, `test:`.
- Do not commit secrets, generated caches, local databases, build artifacts, or personal environment files.
- Avoid large mixed commits that combine code, formatting, and documentation.

## Rollback

- Documentation changes should be revertible with `git revert`.
- Runtime changes should describe data migrations and operational side effects.
- Force push, reset, delete, and drop operations require explicit confirmation.

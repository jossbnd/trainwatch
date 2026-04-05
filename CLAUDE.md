# CLAUDE.md

## Git workflow
- Never push directly to main. Always create a feature branch and open a PR.
- Branch naming: `<type>/<short-description>` (e.g. `feat/garmin-app`, `fix/prim-timeout`, `ci/review-on-demand`).
- Write clear, conventional commit messages: `type: short description` (types: feat, fix, refactor, ci, docs, test, chore).

## GitHub interactions
- Always show the full content (title, body) of any PR, issue, or comment and wait for user approval before posting to GitHub.

## Language & project
- Backend is Go (1.25+) using Gin, located in `backend/`.
- Run tests with `go test ./...` from the `backend/` directory.
- Do not commit `.env` files or secrets. Use `.env.example` for templates.

## Code style
- Follow standard Go conventions: `gofumpt`, short variable names, error handling with early returns.
- Keep functions small and focused. Prefer returning errors over panicking.
- Do not add comments that restate the code. Only comment non-obvious logic.

## General behavior
- Respond in the same language the user writes in.
- Be concise. Don't over-explain or add unrequested features.
- Read code before modifying it. Don't guess.

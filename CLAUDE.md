# CLAUDE.md

## Git workflow
- Never push directly to main. Always create a feature branch and open a PR.
- Branch naming: `<type>/<short-description>` (e.g. `feat/garmin-app`, `fix/prim-timeout`, `ci/review-on-demand`).
- Write clear, conventional commit messages: `type: short description` (types: feat, fix, refactor, ci, docs, test, chore).

## GitHub interactions
- Always show the full content (title, body) of any PR, issue, or comment and wait for user approval before posting to GitHub.

## Language & project
- Backend is Go (1.25+) using Gin, located in `backend/`. Docker config is in `backend/docker/`.
- Garmin watch app is Monkey C, located in `garmin/`. Uses Connect IQ SDK 9.x. Target device: **FR165**.
- Run tests with `go test ./...` from the `backend/` directory.
- Do not commit `.env` files or secrets. Use `.env.example` for templates.
- Do not commit `garmin/source/Config.mc` — use `Config.mc.example` as template.
- API endpoint is `GET /departures` with query params `stop_ref`, `line_ref`, `direction`, `limit`.
- JSON response envelope uses key `departures` (not `trains`).

## Observability (Sentry)
- `internal/sentry` is the **single point of contact** with the `sentry-go` SDK. Never import `github.com/getsentry/sentry-go` or `sentry-go/gin` outside of this package.
- Sentry is activated when `SENTRY_DSN` is non-empty. There is no separate `SENTRY_ENABLED` flag.
- Sentry must be initialized **before** the logger so that startup logs are forwarded.
- HTTP 500s are captured in `middleware.SentryCapture()`. Call `c.Error(err)` in handlers to attach the real error; a synthetic fallback is used otherwise.

## Code style
- Follow standard Go conventions: `gofumpt`, short variable names, error handling with early returns.
- Keep functions small and focused. Prefer returning errors over panicking.
- Do not add comments that restate the code. Only comment non-obvious logic.

## General behavior
- Respond in the same language the user writes in.
- Be concise. Don't over-explain or add unrequested features.
- Read code before modifying it. Don't guess.

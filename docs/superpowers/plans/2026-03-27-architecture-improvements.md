# Architecture Improvements — 2026-03-27

## Overview
Ten incremental improvements to the trainwatch backend covering context propagation, package consolidation, enriched responses, variable naming, bug fixes, observability, and tests.

## Tasks

### Task 1 — Propagate context.Context
Add `ctx context.Context` to `prim.Client.FetchStopVisits`, `service.Service.GetNextTrains`, and pass `c.Request.Context()` from the handler.

### Task 2 — Replace dto/ with model/
Create `internal/model/` with `TextValue`, `MonitoredCall`, `MonitoredVehicleJourney`, `NextTrain`, and `NewNextTrain()`. Remove `internal/api/dto/`.

### Task 3 — Enrich API response
Add `destination`, `aimed_at`, `wait_minutes` fields to `model.NextTrain`. Extract destination from PRIM `DestinationName`. Compute `WaitMinutes` in `NewNextTrain()`.

### Task 4 — Fix variable naming
- `cmd/main.go`: `c` → `cfg`, shadow `logger` → `log`
- `service/next_trains.go`: `mvj` → `journey`, `df` → `dirFilter`
- `service/service.go`: `i` → `input`

### Task 5 — Fix handler bugs
- `400`: respond with `gin.H{"error": err.Error()}` not raw error
- `404`: respond with `gin.H{"error": "no trains found"}` not nil
- `500`: respond with `gin.H{"error": "internal server error"}` not raw error

### Task 6 — Graceful shutdown
Use `http.Server` with `Shutdown(ctx)` triggered on SIGTERM/SIGINT.

### Task 7 — GET /health
Register `GET /health` returning `200 {"status":"ok"}`.

### Task 8 — Request logging middleware
Add Gin middleware logging method, path, status code, and latency via slog.

### Task 9 — limit query parameter
Add optional `limit` query param (default 5). Slice result in service before returning.

### Task 10 — Unit tests (service layer)
Create `service/next_trains_test.go` with 7 tests covering direction filtering, time selection, past departure filtering, limit enforcement, and empty results.

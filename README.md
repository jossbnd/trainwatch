# trainwatch

Real-time train departure monitor for Île-de-France, with a Go backend and a Garmin watch app.

## Repo structure

| Path | Description |
|------|-------------|
| `backend/` | Go REST API — fetches live departures from the PRIM API (Île-de-France Mobilités) |
| *(coming soon)* | Garmin watch app — displays next departures on the wrist |

---

## Backend

### Prerequisites

- Go 1.25+
- A PRIM API key — request one at [data.iledefrance-mobilites.fr](https://data.iledefrance-mobilites.fr)

### Setup

```bash
cp backend/.env.example backend/.env
# Edit backend/.env and fill in PRIM_API_KEY and API_KEY
go run ./backend/cmd
```

The server starts on port `8080` by default (set `PORT` in `.env` to change).

### Endpoints

#### `GET /health`

Returns `200 OK` when the server is up.

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

#### `GET /next-train`

Returns upcoming departures for a given stop and line.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `stop` | string | yes | Stop monitoring ref (PRIM stop ID) |
| `line` | string | yes | Line ref (PRIM line ID) |
| `direction` | string | no | Filter by direction ref |
| `limit` | int | no | Max number of results |

```bash
curl "http://localhost:8080/next-train?stop=STIF%3AStopPoint%3AQ%3A473920%3A&line=STIF%3ALine%3A%3AC01742%3A&limit=2" \
  -H 'X-API-Key: <your_key>'
```

```json
[
  {
    "estimated_at": "2024-01-15T08:32:00Z",
    "aimed_at": "2024-01-15T08:30:00Z",
    "destination": "Versailles-Chantiers",
    "status": "delayed",
    "delay_minutes": 2
  },
  {
    "estimated_at": "2024-01-15T08:47:00Z",
    "destination": "Versailles-Chantiers",
    "status": "onTime",
    "delay_minutes": 0
  }
]
```

**HTTP status codes**

| Code | Meaning |
|------|---------|
| `200` | Departures returned |
| `400` | Missing or invalid query parameters |
| `404` | No trains found for the given criteria |
| `500` | Internal server error |

---

## Garmin app

Work in progress.

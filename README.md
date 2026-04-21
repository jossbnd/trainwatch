# trainwatch

Real-time train departure monitor for Île-de-France, with a Go backend and a Garmin watch app.

## Repo structure

| Path | Description |
|------|-------------|
| `backend/` | Go REST API — fetches live departures from the PRIM API (Île-de-France Mobilités) |
| `garmin/` | Garmin watch app — displays next departures on the wrist |

---

## Backend

### Prerequisites

- Docker + Docker Compose
- [mkcert](https://github.com/FiloSottile/mkcert) for local TLS
- A PRIM API key — request one at [data.iledefrance-mobilites.fr](https://data.iledefrance-mobilites.fr)

### Local setup

```bash
# 1. Install mkcert and generate local cert
brew install mkcert
mkcert -install
cd backend/docker/certs
mkcert api.localhost
mv api.localhost.pem cert.pem
mv api.localhost-key.pem key.pem
cd ../../..

# 2. Configure environment
cp backend/.env.example backend/.env
# Edit backend/.env and fill in PRIM_API_KEY, API_KEY

# 3. Start
cd backend && make up
```

The API is available at `https://api.localhost`.

### Makefile commands

| Command | Description |
|---------|-------------|
| `make up` | Start services in background |
| `make reup` | Force recreate and restart services |
| `make fmt` | Format code with gofumpt and goimports |
| `make lint` | Run golangci-lint |
| `make test` | Run tests |

### Project structure

```
backend/
├── cmd/main.go          # entrypoint
├── docker/              # Traefik reverse proxy config
├── docker-compose.yml
└── internal/
    ├── api/             # HTTP handlers
    ├── config/          # config loading
    ├── logger/          # slog-based logger with Sentry forwarding
    ├── middleware/       # Gin middlewares (auth, request ID, logging, Sentry capture)
    ├── model/           # data structs
    ├── prim/            # PRIM HTTP client
    ├── sentry/          # single wrapper around the sentry-go SDK
    └── service/         # business logic
```

### Authentication

All endpoints except `/health` require an `X-API-Key` header matching `API_KEY` in `.env`.

### Endpoints

#### `GET /health`

```bash
curl https://api.localhost/health
# {"status":"ok"}
```

#### `GET /departures`

Returns upcoming departures for a given stop and line.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `stop_ref` | string | yes | Stop ref (PRIM stop ID) |
| `line_ref` | string | yes | Line ref (PRIM line ID) |
| `direction` | string | no | Filter by direction name |
| `limit` | int | no | Max number of results (default 5) |

```bash
curl "https://api.localhost/departures?stop_ref=STIF%3AStopArea%3ASP%3A43198%3A&line_ref=STIF%3ALine%3A%3AC01742%3A&limit=2" \
  -H 'X-API-Key: <your_key>'
```

```json
{
  "departures": [
    {
      "estimated_at": "2024-01-15T08:32:00Z",
      "aimed_at": "2024-01-15T08:30:00Z",
      "destination": "Boissy-Saint-Léger",
      "status": "onTime",
      "delay_minutes": 0
    }
  ]
}
```

| Code | Meaning |
|------|---------|
| `200` | Departures returned |
| `400` | Missing or invalid query parameters |
| `404` | No departures found |
| `500` | Internal server error |

---

## Garmin app

Watch app for the Garmin FR165 displaying next departures.

### Prerequisites

- [Garmin Connect IQ SDK](https://developer.garmin.com/connect-iq/sdk/) (9.x)
- VS Code with the [Monkey C extension](https://marketplace.visualstudio.com/items?itemName=garmin.monkey-c)
- A Garmin developer key (see below)

### Developer key

1. Create a free account at [developer.garmin.com](https://developer.garmin.com)
2. In VS Code, open the command palette (`Shift+Cmd+P`) and run **Monkey C: Generate Developer Key**
3. Save the key as `garmin/developer_key`
4. The `garmin/developer_key` file is gitignored — keep it local

### Setup

```bash
# 1. Generate developer key (see above)

# 2. Configure the app
cp garmin/source/Config.mc.example garmin/source/Config.mc
# Edit Config.mc and fill in API_URL, API_KEY, STOP_REF, LINE_REF

# 3. Configure the SDK path
cp garmin/.env.example garmin/.env
# Edit garmin/.env and set CIQ_SDK_DIR if your SDK is not in the default location
```

Open the workspace in VS Code — the `.vscode/settings.json` and `launch.json` are already configured.

Press `F5` to build and run in the simulator.

### Build for device

```bash
cd garmin
make build   # produces bin/trainwatch.prg, ready to sideload onto the FR165
```

Copy `bin/trainwatch.prg` to the `GARMIN/APPS/` folder of the watch when mounted as USB mass storage.

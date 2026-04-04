# trainwatch

Real-time train departure monitor for Île-de-France, with a Go backend and a Garmin watch app.

## Repo structure

| Path | Description |
|------|-------------|
| `backend/` | Go REST API — fetches live departures from the PRIM API (Île-de-France Mobilités) |
| `docker/` | Traefik reverse proxy configuration |
| *(coming soon)* | Garmin watch app — displays next departures on the wrist |

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
cd docker/certs
mkcert api.localhost
mv api.localhost.pem cert.pem
mv api.localhost-key.pem key.pem
cd ../..

# 2. Configure environment
cp backend/.env.example backend/.env
# Edit backend/.env and fill in PRIM_API_KEY and API_KEY

# 3. Start
make up
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

### Endpoints

#### `GET /health`

Returns `200 OK` when the server is up.

```bash
curl https://api.localhost/health
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
curl "https://api.localhost/next-train?stop=STIF%3AStopArea%3ASP%3A43198%3A&line=STIF%3ALine%3A%3AC01742%3A" \
  -H 'X-API-Key: <your_key>'
```

```json
[
  {
    "estimated_at": "2024-01-15T08:32:00Z",
    "aimed_at": "2024-01-15T08:30:00Z",
    "destination": "Boissy-Saint-Léger",
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

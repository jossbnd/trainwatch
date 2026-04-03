
# trainwatch - backend

Go REST API fetching next departures from the PRIM API (Île-de-France Mobilités).

## Prerequisites

- Go 1.25+
- A PRIM API key — [data.iledefrance-mobilites.fr](https://data.iledefrance-mobilites.fr)

## Setup

```bash
cp .env.example .env
# Fill in PRIM_API_KEY and API_KEY in .env
go run ./cmd
```

The server starts on port `8080` by default.

## Authentication

All endpoints except `/health` require an `X-API-Key` header matching the `API_KEY` value in `.env`.

```bash
curl http://localhost:8080/next-train?... -H 'X-API-Key: <your_key>'
```

## Endpoints

### `GET /health`

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

### `GET /next-train`

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `stop` | string | yes | Stop monitoring ref (PRIM stop ID) |
| `line` | string | yes | Line ref (PRIM line ID) |
| `direction` | string | no | Filter by direction ref |
| `limit` | int | no | Max number of results |

```bash
curl "http://localhost:8080/next-train?stop=STIF%3AStopArea%3ASP%3A43198%3A&line=STIF%3ALine%3A%3AC01742%3A" \
  -H 'X-API-Key: <your_key>'
```

## Project structure

```
backend/
├── cmd/main.go          # entrypoint
└── internal/
    ├── api/             # HTTP handlers
    ├── config/          # config loading
    ├── logger/          # logger abstraction
    ├── middleware/      # Gin middlewares
    ├── model/           # data structs
    ├── prim/            # PRIM HTTP client
    └── service/         # business logic
```

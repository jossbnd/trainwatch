
# trainwatch - backend

Simple API to fetch the next train (RATP) using the public PRIM API.

Prerequisites
- Go 1.25+

Run the server

```bash
go run ./cmd
```

Endpoint

GET `/next-train`

Required query parameters:
- `type`: transport type (`metro`, `rer`, `tram`, ...)
- `line`: line identifier (e.g. `1`)
- `station`: station name (e.g. `Chatelet`)
- `direction`: direction or destination name (e.g. `A` or `La Defense`)

Example

```bash
curl "http://localhost:8080/next-train?type=metro&line=1&station=Chatelet&direction=A"
```

Behavior
- If no upcoming train is found, the API responds with `404` and an error message.

Possible improvements
- Add unit tests for the `service` layer and a mock for `internal/prim`.
- Add normalization for station names (mapping, fuzzy matching).
- Add a healthcheck endpoint.


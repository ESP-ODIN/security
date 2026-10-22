# security

The security pipeline component of the ODIN application. Requires Go 1.27.0 or later.
See [architecture.md](architecture.md) for the full pipeline design.

## Run locally

`DATABASE_URL` must be set (see [Configuration](#configuration)); the API pings the database on startup.

On Windows (PowerShell):

```powershell
$env:DATABASE_URL = "postgresql://user:password@host/dbname?sslmode=require"
go run ./cmd/api
```

In a second terminal:

```powershell
Invoke-RestMethod http://127.0.0.1:8080/health
```

On Linux or macOS (Bash/Zsh):

```sh
export DATABASE_URL=postgresql://user:password@host/dbname?sslmode=require
go run ./cmd/api
```

Check the endpoints with:

```sh
curl http://127.0.0.1:8080/health
curl http://127.0.0.1:8080/ready
```

`/health` is a liveness check. `/ready` pings the database.

Swagger UI is available at [http://127.0.0.1:8080/swagger/index.html](http://127.0.0.1:8080/swagger/index.html).
Regenerate docs after changing annotations with `make swagger`.

## Configuration

By default, the server listens on `127.0.0.1:8080`. To change the port on Windows (PowerShell):

```powershell
$env:HTTP_ADDR = "127.0.0.1:8081"
go run ./cmd/api
```

On Linux or macOS (Bash/Zsh):

```sh
export HTTP_ADDR=127.0.0.1:8081
go run ./cmd/api
```

`DATABASE_URL`, `AUTH_JWKS_URL`, `AUTH_ISSUER` and `AUTH_AUDIENCE` are required.
`ADMIN_USER_IDS` is an optional comma-separated list of Auth user UUIDs allowed to call admin routes.
`.env.exemple` documents the environment variables. `.env` files are not loaded automatically
by `go run`/`go build` — export them yourself, or use `make docker-dev` below, which reads `.env`.
Configuration is read from the process environment.

## Checks and build

```sh
go test ./...
go vet ./...
go build -o bin/ ./cmd/...
```

The API executable is `bin/api.exe` on Windows and `bin/api` on Linux and macOS.
If Make is installed, `make run`, `make test`, `make vet`,
`make build`, and `make fmt` are also available.

## Docker

A development `Dockerfile` runs the API with hot-reload (via [air](https://github.com/air-verse/air)),
meant to be plugged as-is into the project's docker-compose. Fill in `DATABASE_URL` in `.env` first
(see `.env.exemple`).

Build the image:

```sh
make docker-build
```

Run it locally with hot-reload (mounts the source code and exposes the app on the host port set by
`PORT` in `.env`, `8093` by default):

```sh
make docker-dev
```

Then check it with `curl http://127.0.0.1:8093/health`. Editing any `.go` file rebuilds and restarts
the server automatically inside the container. The container itself always listens on `0.0.0.0:8080`
(overridden from `.env`'s `HTTP_ADDR` by the `docker-dev` target); only the host-side port, read from
`PORT` in `.env`, changes.

## Structure

- `cmd/api`: entrypoint HTTP server.
- `cmd/migrate`: reserved for migrations.
- `config`: configuration loading and validation.
- `db`: database connection pool.
- `internal/router`: HTTP route registration, one file per domain.
- `internal/handler`: HTTP handlers.
- `internal/service`: business logic interfaces.
- `internal/service/implementation`: service implementations.
- `internal/repository`: data access interfaces.
- `internal/repository/implementation`: repository implementations.
- `internal/dto`: request/response payloads.
- `internal/model`: domain models.
- `internal/testutil`: JWT fixtures for tests.
- `docs`: generated Swagger/OpenAPI files (`make swagger`).
- `migrations`: SQL migrations (to come).

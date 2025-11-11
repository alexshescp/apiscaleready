# Go Load Lab

Go Load Lab is a flexible load-testing toolkit that can be executed from the command line or driven via an HTTP API. The project now produces structured summaries, exposes an OpenAPI contract, and ships with a C4 context diagram for architecture discussions.

## Features
- Multiple concurrency strategies (wait group, semaphore, channel)
- CLI runner with interactive progress bar
- HTTP API for on-demand load generation
- Structured metrics summaries ready for automation
- OpenAPI documentation (`openapi.yaml`) and C4 context diagram (`docs/c4-context.drawio`)

## Getting Started

### CLI mode
```bash
go mod tidy
go run .
```

Configure the CLI by exporting `GLAB_` environment variables:
```bash
GLAB_TARGET_URL=https://httpbin.org/post \
GLAB_METHOD=POST \
GLAB_JSON_BODY='{"test":42}' \
GLAB_CONCURRENCY=50 \
GLAB_RATE=200 \
GLAB_DURATION=20s \
GLAB_MODE=semaphore \
go run .
```

### HTTP API mode
Set `GLAB_LISTEN_ADDR` to expose the API and run the binary:
```bash
GLAB_LISTEN_ADDR=:8080 go run .
```

Use the API to orchestrate tests programmatically. Refer to `openapi.yaml` for the full contract.

#### List of available curl commands
- `curl http://localhost:8080/healthz`
- `curl -X POST http://localhost:8080/tests -H 'Content-Type: application/json' -d '{"target_url":"https://example.com","duration":"5s"}'`

## Configuration Reference
| Variable | Description |
| --- | --- |
| `GLAB_TARGET_URL` | Target endpoint for the CLI runner. |
| `GLAB_METHOD` | HTTP method (default `GET`). |
| `GLAB_CONCURRENCY` | Number of parallel workers (default `100`). |
| `GLAB_DURATION` | Total test duration (default `10s`). |
| `GLAB_RATE` | Requests per second cap (`0` disables throttling). |
| `GLAB_MODE` | Concurrency strategy (`wg`, `semaphore`, `channel`). |
| `GLAB_NO_KEEP_ALIVE` | Disable HTTP keep-alives when set to `true`. |
| `GLAB_PROXY_FILE` | Path to proxy list file. |
| `GLAB_JSON_BODY` | Request body for non-GET operations. |
| `GLAB_HEADERS_JSON` | JSON map of additional headers. |
| `GLAB_CONTENT_TYPE` | Overrides the `Content-Type` header. |
| `GLAB_LISTEN_ADDR` | When set, the HTTP API listens on this address instead of running the CLI. |

## Testing
Run the unit and integration tests with:
```bash
go test ./...
```

## Documentation Assets
- `openapi.yaml` — OpenAPI 3 specification of the HTTP API.
- `docs/c4-context.drawio` — C4 context diagram outlining primary actors and systems.

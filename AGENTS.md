# Repository Guidelines

## Project Structure & Module Organization
- `cmd/`: Entry points (e.g., `cmd/main.go`).
- `internal/`: Application logic and services (private).
- `pkg/`: Public packages and generated code (`pkg/gen`).
- `proto/`: Protobuf definitions; generates Go and web assets.
- `webapp/_webapp/`: Chrome extension (frontend) source and build.
- `docs/`: Images and documentation assets.
- `helm-chart/` and `hack/`: Deployment templates and scripts.
- Build artifacts land in `dist/`.

## Build, Test, and Development Commands
- `make deps`: Install dev tools (buf, wire, golangci-lint) and web deps.
- `make gen`: Compile protobufs to `pkg/gen` and webapp; run `wire gen`.
- `make fmt`: Format protobufs, Go, and webapp code.
- `make lint`: Run buf + golangci-lint and webapp lint.
- `make build`: Build backend to `dist/pd.exe`.
- `make test`: Run Go tests with coverage; writes `coverage.out/html`.
- Run locally: `./dist/pd.exe` (default `:6060`).
- Webapp build example: `cd webapp/_webapp && npm run build:prd:chrome`.

## Coding Style & Naming Conventions
- Go: idiomatic style; use `go fmt` and `golangci-lint` (gci import ordering).
- Protos: keep packages consistent; run `make gen` after changes.
- Webapp: follow project ESLint/Prettier via `npm run lint` and `npm run format`.
- Naming: Go files `lower_snake_case.go`; tests end with `_test.go`. Packages are short, lower-case.

## Testing Guidelines
- Framework: standard Go testing.
- Naming: `file_test.go`, functions `TestXxx`.
- Coverage: `make test` sets `PD_MONGO_URI` and produces `coverage.html`.
- Scope: prefer unit tests in `internal/...`; mock external services where possible.

## Commit & Pull Request Guidelines
- Conventional Commits for titles/messages: e.g., `feat(api): add project listing`, `fix(webapp): debounce search`.
- Keep PRs focused; include description, linked issues, and screenshots/GIFs for webapp changes.
- Update generated code in the same PR when proto/DI changes (`make gen`).
- CI enforces semantic PR titles.

## Security & Configuration Tips
- Copy `.env.example` to `.env`; do not commit secrets.
- Local dev requires MongoDB (`PD_MONGO_URI`), optional OpenAI key for AI features.
- Publishing/build pipelines use GitHub Actions; avoid committing large build outputs.

# PROJECT.md

## Purpose
This project demonstrates reading and writing a JSON config file stored in the user's home directory (`~/.gatorconfig.json`) using an internal Go package.

## Current Behavior
- `main.go` reads config from disk into app state.
- The CLI currently supports a `login` command that sets `current_user_name` in `~/.gatorconfig.json`.
- Running `gator login <name>` updates the config and prints a confirmation message.
- The CLI also supports `register`, which creates a user in Postgres and sets that user as current.

## Config Package Design (`internal/config`)
- Exported:
  - `type Config` with JSON tags:
    - `db_url`
    - `current_user_name`
  - `func Read() (Config, error)`
  - `func (cfg *Config) SetUser(user string) error`
- Internal helpers:
  - `const configFileName = ".gatorconfig.json"`
  - `getConfigFilePath() (string, error)`
  - `write(cfg Config) error`

## File Format
Expected config JSON:

```json
{
  "db_url": "postgres://example",
  "current_user_name": "copilot"
}
```

`current_user_name` may be absent before first write; it is added by `SetUser`.

## Assumptions
- `~/.gatorconfig.json` exists and is readable before app startup.
- Home directory is resolved using `os.UserHomeDir()`.
- Writing the config should overwrite the existing file contents.

## Operational Notes
- `Makefile` targets:
  - `run`: runs `go run .`
  - `build`: runs `sqlc generate` and `go build`
  - `reset`: runs goose down/up from `sql/schema` against `postgres://postgres:postgres@localhost:5432/gator`; the down step tolerates a fresh database with no current version and the up step is run in a separate shell line
  - `clean`: removes local `gator` binary and generated SQLC Go files under `internal/database/*.go`
- Database schema is managed by goose.
- Migration files live under `sql/schema`.
- From `sql/schema`, use `goose postgres postgres://postgres:postgres@localhost:5432/gator up` to apply migrations.
- From `sql/schema`, use `goose postgres postgres://postgres:postgres@localhost:5432/gator down` to roll back migrations.
- `goose status` shows the current migration state of the database.
- SQLC is configured from the repo root using `sqlc.yaml`.
- SQLC reads schema files from `sql/schema` and queries from `sql/queries`.
- Generated Go code is written to `internal/database`.
- The application opens Postgres using the DB URL from `~/.gatorconfig.json` and stores `*database.Queries` in app state.

## Maintenance Rule
Update this file whenever behavior, architecture, assumptions, or key commands change so future agents have accurate context.

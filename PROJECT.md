# MEMORY.md

## Purpose
This project demonstrates reading and writing a JSON config file stored in the user's home directory (`~/.gatorconfig.json`) using an internal Go package.

## Current Behavior
- `main.go` reads config from disk into app state.
- The CLI currently supports a `login` command that sets `current_user_name` in `~/.gatorconfig.json`.
- Running `gator login <name>` updates the config and prints a confirmation message.

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
  - `clean`: removes local `gator` binary if present
- Database schema is managed by goose.
- Migration files live under `sql/schema`.
- From `sql/schema`, use `goose postgres postgres://postgres:postgres@localhost:5432/gator up` to apply migrations.
- From `sql/schema`, use `goose postgres postgres://postgres:postgres@localhost:5432/gator down` to roll back migrations.
- `goose status` shows the current migration state of the database.

## Maintenance Rule
Update this file whenever behavior, architecture, assumptions, or key commands change so future agents have accurate context.

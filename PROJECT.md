# MEMORY.md

## Purpose
This project demonstrates reading and writing a JSON config file stored in the user's home directory (`~/.gatorconfig.json`) using an internal Go package.

## Current Behavior
- `main.go` reads config from disk through `internal/config.Read()`.
- `main.go` sets `current_user_name` to `copilot` using `Config.SetUser()`.
- `main.go` re-reads config and prints the resulting struct.

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

## Maintenance Rule
Update this file whenever behavior, architecture, assumptions, or key commands change so future agents have accurate context.

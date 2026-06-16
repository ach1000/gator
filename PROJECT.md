# PROJECT.md

## Purpose
This project is a Go CLI RSS reader backed by Postgres. It stores configuration in the user's home directory (`~/.gatorconfig.json`), manages schema and queries with goose and SQLC, and supports user, feed, and follow management commands.

## Current Behavior
- `main.go` reads config from disk into app state.
- The CLI supports these commands:
  - `login`: sets `current_user_name` in `~/.gatorconfig.json`
  - `register`: creates a user in Postgres and sets that user as current
  - `reset`: deletes all users
  - `users`: lists users and marks the current user
  - `agg <time_between_reqs>`: runs a never-ending aggregation loop; fetches the next scheduled feed, prints item titles, and repeats on a ticker interval (e.g. `1s`, `1m`, `1h`)
  - `addfeed`: creates a feed for the current user and auto-follows it
  - `feeds`: lists feeds with owner names
  - `follow`: follows a feed by URL for the current user
  - `unfollow`: unfollows a feed by URL for the current user
  - `following`: lists feed names followed by the current user
- Commands that require an authenticated user are wrapped with `middlewareLoggedIn`, which loads the current `database.User` once and passes it into the handler.

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

## Command Design
- `state` stores:
  - `*database.Queries`
  - `*config.Config`
- `commands` is a string-to-handler registry with handlers of type `func(*state, command) error`.
- Logged-in handlers use the higher-order wrapper `middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error`.
- `follow` and `unfollow` both resolve feeds by URL before creating or deleting rows in `feed_follows`.
- Aggregation flow:
  - `handlerAgg` requires exactly one duration argument and parses it with `time.ParseDuration`
  - It prints `Collecting feeds every <duration>` and starts a `time.Ticker`
  - It runs immediately, then once per tick via a `for ; ; <-ticker.C` loop
  - `scrapeFeeds` selects one feed to fetch (`GetNextFeedToFetch`), marks it fetched (`MarkFeedFetched`), fetches RSS with `fetchFeed`, and prints each item title
  - Scrape errors are logged and the loop continues until process termination (Ctrl+C)

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
- After editing SQL in `sql/queries`, regenerate generated Go code with `sqlc generate` or `make build`.
- Feed scheduling schema/query changes:
  - Migration `sql/schema/004_last_fetched_at.sql` adds nullable `feeds.last_fetched_at`.
  - Query `MarkFeedFetched` sets `last_fetched_at = NOW()` and `updated_at = NOW()` for a feed ID.
  - Query `GetNextFeedToFetch` returns one feed ordered by `last_fetched_at ASC NULLS FIRST`.
  - If code references `last_fetched_at` before migration is applied, Postgres returns `column "last_fetched_at" does not exist (42703)`.
- The application opens Postgres using the DB URL from `~/.gatorconfig.json` and stores `*database.Queries` in app state.

## Maintenance Rule
Update this file whenever behavior, architecture, assumptions, or key commands change so future agents have accurate context.

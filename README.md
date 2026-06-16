# gator

`gator` is a CLI RSS aggregator written in Go with a Postgres backend.

## Prerequisites

You need these installed before running `gator`:

- Go (the project currently uses Go 1.24+)
- PostgreSQL

## Install the CLI with `go install`

If this repository is available at its module path (`github.com/ach1000/gator`), install with:

```bash
go install github.com/ach1000/gator@latest
```

If you are working locally from this repo, you can install from the current directory:

```bash
go install .
```

Make sure your Go bin directory is on your `PATH` (usually `$HOME/go/bin`).

## Database Setup

1. Create a Postgres database (example name: `gator`).
2. Update your connection string as needed, for example:

```text
postgres://postgres:postgres@localhost:5432/gator?sslmode=disable
```

3. Run migrations from the schema directory:

```bash
cd sql/schema
goose postgres postgres://postgres:postgres@localhost:5432/gator up
```

From the repo root, you can also use:

```bash
make reset
```

## Config File Setup

`gator` reads config from:

```text
~/.gatorconfig.json
```

Create this file with your DB URL:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
}
```

After you run `register` or `login`, `current_user_name` is stored automatically.

## Running the Program

From the repo root during development:

```bash
go run . <command> [args...]
```

After installing via `go install`:

```bash
gator <command> [args...]
```

## Common Commands

Register and log in:

```bash
gator register alice
gator login alice
gator users
```

Add and follow feeds:

```bash
gator addfeed "Hacker News" "https://hnrss.org/newest"
gator feeds
gator follow "https://hnrss.org/frontpage"
gator following
```

Run the aggregator loop (fetch/store posts continuously):

```bash
gator agg 1m
```

Browse recent posts from feeds you follow (default limit is 2):

```bash
gator browse
gator browse 10
```

Unfollow a feed:

```bash
gator unfollow "https://hnrss.org/frontpage"
```

## Notes

- `agg` is designed to run continuously; stop it with `Ctrl+C`.
- Duplicate posts (same URL) are ignored while scraping.
- If schema-related errors appear, re-run migrations with goose.

## Troubleshooting

- Error: `column "last_fetched_at" does not exist` or `relation "posts" does not exist`
  - Cause: migrations have not been applied.
  - Fix: run goose up migrations from `sql/schema`.

- Error: `pq: password authentication failed` or DB connection refused
  - Cause: incorrect `db_url` or Postgres is not running.
  - Fix: verify Postgres is running and update `~/.gatorconfig.json` with the correct connection string.

- Error: `command not found: gator`
  - Cause: Go bin path is not on `PATH`.
  - Fix: add `$HOME/go/bin` (or your configured `GOBIN`) to `PATH`, then restart your shell.

- `browse` shows no posts
  - Cause: no scraped posts yet, or current user does not follow any feeds.
  - Fix: run `addfeed`/`follow`, then run `agg` for a bit, then try `browse` again.

- Aggregator runs but inserts fail repeatedly
  - Cause: feed data format issues or unexpected DB errors.
  - Fix: check terminal logs, confirm migrations are current, and test with another RSS feed URL.

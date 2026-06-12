run:
	go run .

build:
	sqlc generate
	go build

reset:
	cd sql/schema && goose postgres postgres://postgres:postgres@localhost:5432/gator down || true
	cd sql/schema && goose postgres postgres://postgres:postgres@localhost:5432/gator up

clean:
	rm -f gator
	rm -f internal/database/*.go

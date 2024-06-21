create-migration:
	@GOOSE_MIGRATION_DIR=./migrations go run github.com/pressly/goose/v3/cmd/goose postgres "user=stream password=stream dbname=stream sslmode=disable" create $(name) sql
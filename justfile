set shell := ['fish', '-c']

default:
    just -lu

db-apply-migrations:
    goose postgres "host=127.0.0.1 port=6666 user=postgres password=postgres dbname=postgres sslmode=disable" up

db-create-migration name:
    goose postgres "host=127.0.0.1 port=6666 user=postgres password=postgres dbname=postgres sslmode=disable" create {{name}}
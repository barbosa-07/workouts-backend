# How to run the project

## Docker compose command
- docker compose up --build


## connect to psql
-  psql -U postgres -h localhost -p 5432


## go project
- go run main.go

# migration
- export PATH=$HOME/go/bin:$PATH                                                                          
- goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
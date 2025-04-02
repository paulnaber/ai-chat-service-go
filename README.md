# AI Chat Service

This project is an AI Chat Proxy Server that handles chat and message functionalities.
The main focus of this project is the architecture. I want to compare it with other technologies and languages.

## Other Languages and Technologies

[AI Chat Proxy Server - Node.js](https://github.com/paulnaber/ai-chat-service-nodejs) <br>
Java...

## Tech Stack

-   Goose - Database migrations
-   sqlc - Type-safe SQL queries
-   Docker Compose - Simplified local database setup
-   Fiber - Fast and minimalist web framework
-   oapi-codegen - Structs and Routes generated from OpenAPI spec
-   Keycloak - Authentication, Authorization

## Requirements

-   Go 1.20 or higher
-   Docker and Docker Compose (for local development)
-   Required tools (for development):
    -   Goose (database migrations)
    -   sqlc (SQL code generation)

### Getting Started

1. Before getting started, make sure to:

```
add .env file (see .env.example)
have oauth2 provider up and running (docker compose)
have postgres up and running (docker compose)
have go installed on your machine
```

2. Install required tools:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

3. Generate SQL code:

```bash
sqlc generate
```

4. Generate Structs and Routes from OpenAPI spec:

```bash
oapi-codegen -generate types,fiber -o ./internal/api/api.gen.go -package api ./api.yml
```

5. Build the application:

```bash
go build -o bin/ai-chat-service-go .
```

6. Start the application:

```bash
./bin/ai-chat-service-go
```

for development, you can run:

```bash
go run .
```

### Running Tests

Run tests with:

```bash
go test -v ./...
```

### Database Migrations

Create a new migration:

```bash
read -p "Enter migration name: " name; goose -dir ./sql/schema create $name sql
```

Run migrations up:

```bash
goose -dir ./sql/schema postgres postgres://postgres:postgres@localhost:5432/aichat up
```

Run migrations down:

```bash
goose -dir ./sql/schema postgres postgres://postgres:postgres@localhost:5432/aichat up
```

Check migration status:

```bash
goose -dir ./sql/schema postgres postgres://postgres:postgres@localhost:5432/aichat status
```

### Docker

Build Docker image:

```bash
docker build -t ai-chat-service-go .
```

Start Docker services:

```bash
docker compose up --build
```

Stop Docker services:

```bash
docker compose down
```

### Swagger

http://localhost:3000/swagger/
http://localhost:3000/api.yml

### TODOs

-   implementation of endpoints
-   better logging
-   prometheus
-   auth, including roles
-   /metrics endpoint

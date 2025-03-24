# AI Chat Service

This project is an AI Chat Proxy Server that handles chat and message functionalities.
The main focus of this project is the architecture. I want to compare it with other technologies and languages.

## Other Languages and Technologies

[AI Chat Proxy Server - Node.js](https://github.com/paulnaber/ai-chat-service-nodejs)
Java...

## Features

-   Create and manage chat sessions
-   Send messages and receive AI-generated responses
-   User authentication via OAuth2 with Keycloak
-   Swagger documentation (generated from code documentation)
-   Dockerized setup with docker-compose

## Features

• Supports chat and message operations
• OpenAPI documentation generated from code documentation
• ~~Authentication with OAuth2 Provider~~
• ~~Downloadable OpenAPI definition (as json)~~
• Easy local setup with Docker Compose

## Tech Stack

• Goose - Database migrations
• sqlc - Type-safe SQL queries
• Docker Compose - Simplified local database setup
• Fiber - Fast and minimalist web framework
• go-swagger - OpenAPI generation from code documentation

## Requirements

-   Go 1.20 or higher
-   Docker and Docker Compose (for local development)
-   Required tools (for development):
    -   Goose (database migrations)
    -   sqlc (SQL code generation)
    -   go-swagger (API documentation)

## Getting Started

### Clone the repository

```bash
git clone https://github.com/yourusername/ai-chat-service-go.git
cd ai-chat-service-go
```

### Configuration

1. Copy the example environment file:

```bash
cp .env.example .env
```

2. Edit the `.env` file to match your environment

### Running with Docker Compose

The easiest way to get started is using Docker Compose, which will set up all the necessary services:

```bash
# Build and start all services
make docker-build
make docker-up

# Stop the services
make docker-down
```

### Running locally (without Docker)

1. Make sure PostgreSQL is running and accessible
2. Make sure Keycloak is running and configured
3. Install required tools:

```bash
# Install development tools
make tools
```

4. Run the database migrations:

```bash
make migrate-up
```

5. Generate SQL code from query files:

```bash
make sqlc
```

6. Generate Swagger documentation:

```bash
make swagger
```

7. Build and run the application:

```bash
make build
make run
```

### Keycloak Setup (Manual)

If you're not using Docker, you'll need to set up Keycloak manually:

1. Create a new realm called `ai-chat`
2. Create a new client called `ai-chat-client`
3. Configure the client:
    - Set access type to `confidential`
    - Enable `Service Accounts`
    - Add redirect URIs for your frontend
4. Get the client secret from the Credentials tab
5. Add the client secret to your `.env` file

### API Documentation

Swagger documentation is available at:

```
http://localhost:3000/swagger/
```

## Development

### Running Tests

```bash
make test
```

### Database Migrations

Migrations are managed using Goose and stored as SQL files in the `migrations` directory.

```bash
# Create a new migration
make migrate-create
# Enter the migration name when prompted

# Run migrations up
make migrate-up

# Roll back migrations
make migrate-down

# Check migration status
make migrate-status
```

### SQL Queries

SQL queries are stored in `sql/queries/` and used to generate Go code with sqlc.

```bash
# Generate SQL code
make sqlc
```

### Swagger Documentation

API documentation is generated automatically from annotations in the code.

```bash
# Generate Swagger documentation
make swagger
```

## Gator

**Gator** is a small Go web service for working with RSS feeds and user follow relationships. It exposes HTTP endpoints for managing users, feeds, and the relationships between them, backed by a Postgres database.

This project is a **lesson from [Boot.dev](https://www.boot.dev/)**, adapted into my own GitHub repository.

### Features

- **User management**: create and manage user accounts.
- **Feed management**: create and store RSS feeds in the database.
- **Follow system**: follow/unfollow feeds per user.
- **Feed ingestion**: pull RSS content into the database for later retrieval.
- **JSON HTTP API**: simple JSON-based endpoints for interacting with the service.

### Tech Stack

- **Language**: Go
- **Database**: PostgreSQL
- **Migrations / SQL**: raw SQL queries (generated Go files in `internal/database`)
- **HTTP**: standard `net/http` with basic routing and middleware

### Getting Started

#### Prerequisites

- Go (1.21+ recommended)
- PostgreSQL

#### Setup

1. **Clone the repo**

   ```bash
   git clone https://github.com/<your-username>/gator.git
   cd gator
   ```

2. **Configure environment**

   Export environment variables as needed, for example:

   ```bash
   export DB_URL="postgres://user:password@localhost:5432/gator?sslmode=disable"
   export PORT="8080"
   ```

   Adjust the variable names and values to match what `main.go` expects in your setup.

3. **Run the server**

   ```bash
   go run ./...
   ```

   The API will start on `http://localhost:${PORT}` (default `8080` if not overridden).

### Project Structure

- `main.go`: application entrypoint and HTTP server setup.
- `user_functions.go`: user-related handlers and helpers.
- `feed_functions.go`: feed-related handlers and helpers.
- `follow_functions.go`: follow/unfollow logic.
- `internal/database/`: generated database access code.
- `sql/queries/`: raw SQL queries and schema pieces.

### Development

Run tests (if present) with:

```bash
go test ./...
```

### Acknowledgements

This project is based on a **Boot.dev** backend curriculum lesson. Thanks to Boot.dev for the clear, hands-on approach to learning backend development in Go.


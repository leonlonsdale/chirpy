# Chirpy

Chirpy is a simple Twitter-like HTTP server written in Go. It provides a RESTful API for user registration, authentication, posting short messages ("chirps"), and user management. The backend uses PostgreSQL for data storage and supports JWT-based authentication and refresh tokens.

## Features

-   User registration and login with secure password hashing (Argon2id)
-   JWT-based authentication for protected endpoints
-   Refresh token support for session management
-   CRUD operations for "chirps" (short messages)
-   Admin endpoints for metrics and user reset (development only)
-   Webhook endpoint for upgrading users to "Chirpy Red"
-   Simple web frontend served from `/web`
-   Environment-based configuration

## Project Structure

```
cmd/chirpy/           # Main application entrypoint
internal/
  auth/               # Authentication and JWT logic
  config/             # API configuration and middleware
  database/           # Database access (generated by sqlc)
  handlers/           # HTTP handlers (API and web)
  server/             # HTTP server setup and routing
  util/               # Utility functions (e.g., JSON responses)
sql/
  schema/             # Database schema migrations
  queries/            # SQL queries for sqlc
web/                  # Static web assets
```

## Getting Started

### Prerequisites

-   Go 1.24+
-   PostgreSQL database
-   [sqlc](https://sqlc.dev/) (for regenerating database code, optional)

### Setup

1. **Clone the repository**

    ```sh
    git clone https://github.com/leonlonsdale/chirpy.git
    cd chirpy
    ```

2. **Configure environment variables**

    Create a `.env` file in the project root with the following variables:

    ```
    DB_URL=postgres://user:password@localhost:5432/chirpydb?sslmode=disable
    JWT_SECRET_KEY=your_jwt_secret
    PLATFORM=dev
    POLKA_KEY=your_polka_webhook_key
    ```

3. **Run database migrations**

    Use your preferred migration tool (e.g., [goose](https://github.com/pressly/goose)) to apply the migrations in `sql/schema/` to your database.

4. **Build and run the server**

    ```sh
    go run ./cmd/chirpy
    ```

    The server will start on port 8080 by default.

## API Endpoints

### User Endpoints

-   `POST /api/users` — Register a new user
-   `PUT /api/users` — Update user email/password (requires JWT)

### Authentication

-   `POST /api/login` — Login and receive JWT + refresh token
-   `POST /api/refresh` — Exchange refresh token for new JWT
-   `POST /api/revoke` — Revoke a refresh token

### Chirps

-   `POST /api/chirps` — Create a new chirp (requires JWT)
-   `GET /api/chirps` — List all chirps (optional `author_id` and `sort` query params, accepting `asc` or `desc`)
-   `GET /api/chirps/{chirpID}` — Get a specific chirp by ID
-   `DELETE /api/chirps/{chirpID}` — Delete a chirp (requires JWT, must be author)

### Admin & Web

-   `GET /admin/metrics` — View server metrics (development only)
-   `POST /admin/reset` — Reset user data (development only)
-   `POST /api/polka/webhooks` — Webhook to upgrade user to Chirpy Red (requires API key)

### Static Files

-   `/app/` — Serves static files from the `web/` directory

## Development

-   Handlers are organized in [`internal/handlers`](internal/handlers/).
-   Database queries are generated by [`sqlc`](https://sqlc.dev/) from SQL files in [`sql/queries`](sql/queries/).
-   JWT and password logic is in [`internal/auth`](internal/auth/).
-   Configuration and middleware are in [`internal/config`](internal/config/).

## Testing

Run all tests with:

```sh
go test ./...
```

## License

MIT License. See [LICENSE](LICENSE)

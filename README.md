# Aarambha Agency

A high-performance real estate discovery and advisory platform built with **Go**, **HTMX**, **Alpine.js**, and **Tailwind CSS v4**.

Features server-side rendering with Go's standard library `net/http` and `html/template`, reactive partial updates via HTMX, lightweight client state with Alpine.js, and self-contained static assets compiled directly into the binary using Go's `embed.FS`.

---

## Tech Stack

- **Backend**: Go 1.24+ (`net/http`, `html/template`, `embed.FS`)
- **Frontend**: HTMX (dynamic swaps), Alpine.js (modals, dark mode, calculators)
- **Styling**: Tailwind CSS v4 CLI
- **Data**: In-memory repository seeded from `seed.json`
- **Deployment**: Multi-stage Docker (Alpine Linux runtime)

---

## Getting Started

### Prerequisites

- [Go 1.24+](https://golang.org/dl/)
- [Node.js 20+](https://nodejs.org/) & `npm` (for Tailwind CSS)
- [Docker](https://www.docker.com/) *(optional)*

### Quick Start

1. **Install frontend dependencies**:
   ```bash
   npm install
   ```

2. **Build stylesheets**:
   ```bash
   npm run build:css
   # or: make css
   ```

3. **Start the server**:
   ```bash
   make dev
   # or: go run ./cmd/server
   ```

4. Open [http://localhost:8080](http://localhost:8080) in your browser.

---

## Development Commands

| Command | Description |
|---|---|
| `make dev` | Run server directly with `go run` |
| `make build` | Compile static binary to `bin/server` |
| `make run` | Build and run the compiled binary |
| `make css` | Compile and minify production Tailwind CSS |
| `npm run watch:css` | Watch and recompile CSS during development |
| `make test` | Run test suite (`go test -v ./...`) |
| `make upload-images-dry` | Preview Supabase Storage image upload plan |
| `make upload-images` | Upload image assets to Supabase Storage & sync URLs |
| `make clean` | Remove compiled binaries and Go cache |
| `make docker-build` | Build the multi-stage Docker image |
| `make docker-run` | Run container locally on port 8080 |

---

## Supabase & Database Setup

The application supports live **Supabase PostgreSQL**.

### 1. Automatic Database Migration & Seeding on Startup
When `DB_URL` is set in `.env`, the server automatically checks if the database is initialized or empty on startup:
- Runs database table migrations (`CREATE TABLE IF NOT EXISTS`, indexes, and RLS policies).
- Automatically seeds all properties, agents, projects, location insights, articles, and testimonials if the properties table is empty.

Alternatively, you can manually execute [`supabase/schema.sql`](supabase/schema.sql) in your **Supabase Dashboard -> SQL Editor** or via `psql`.

### 2. Upload Images to Supabase Storage
Run the automated image uploader to pull external seed assets and host them directly in your Supabase Storage bucket:
```bash
# Preview the 28 unique assets and target storage paths:
make upload-images-dry

# Upload assets to Supabase Storage and generate updated seed files:
make upload-images
```
This produces:
- `supabase/seed.supabase.json`: Seed dataset with Supabase CDN URLs.
- `supabase/seed_supabase_images.sql`: SQL migration updating all database records to Supabase CDN URLs.

---

## Configuration

Configure the server via environment variables or a `.env` file (automatically loaded on startup):

| Variable | Default | Description |
|---|---|---|
| `HOST` | `0.0.0.0` | Bind IP address |
| `PORT` | `8080` | HTTP listening port |
| `APP_ENV` | `development` | Environment mode (`development` or `production`) |
| `DB_URL` | `""` | Supabase / PostgreSQL connection string (`postgres://...`) |
| `SUPABASE_URL` | `""` | Supabase Project API URL (`https://<project-ref>.supabase.co`) |
| `SUPABASE_SECRET_KEY` | `""` | Supabase Service Role / Secret API key |
| `SUPABASE_BUCKET` | `test-bkt` | Supabase Storage bucket name for images |
| `READ_TIMEOUT_SEC` | `10` | Maximum request read duration in seconds |
| `WRITE_TIMEOUT_SEC` | `15` | Maximum response write duration in seconds |
| `IDLE_TIMEOUT_SEC` | `60` | Maximum keep-alive wait time in seconds |

---

## Health & Monitoring

The service exposes production-grade health check endpoints for container orchestrators (Docker, Kubernetes, AWS ECS, GCP Cloud Run) and load balancers:

| Endpoint | Method | Purpose | Response |
|---|---|---|---|
| `/health` | `GET` | Comprehensive diagnostics: uptime, environment, version, and data checks | `200 OK` (JSON) |
| `/healthz` | `GET` | Lightweight liveness probe for orchestrators and load balancers | `200 OK` `{"status":"ok"}` |
| `/livez` | `GET` | Kubernetes standard liveness probe | `200 OK` `{"status":"ok"}` |
| `/readyz` | `GET` | Readiness probe validating core data repository initialization | `200 OK` / `503 Service Unavailable` |

Docker container health is checked automatically via `HEALTHCHECK` in the `Dockerfile`.

---

## Project Structure

```text
├── cmd/server/          # Application entrypoint and server lifecycle
├── configs/             # Environment configuration & company metadata
├── internal/
│   ├── delivery/http/   # HTTP handlers, middleware, and router
│   ├── domain/          # Core entities, models, and interfaces
│   ├── repository/      # In-memory data store and seed.json
│   ├── service/         # Business logic services
│   └── view/            # Template engine and view models
├── web/
│   ├── static/          # Source & compiled CSS, icons, manifest
│   ├── template/        # HTML templates (layouts, pages, partials, admin)
│   └── web.go           # embed.FS declarations for binary embedding
├── Dockerfile           # Multi-stage production container
└── Makefile             # Development and build shortcuts
```

---

## License

Proprietary. Developed for **Aarambha Agency** by **NemeaForge Solutions**.
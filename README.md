# Workspace Overview

This workspace contains the full application stack for the photo and video album platform.

The frontend is managed with Bun in the devcontainer environment, so Bun is the expected package manager for local development and script execution.

## Project structure

- `albumscli/` — Go CLI for album management and imports.
- `clip/` — Python clip service and generated gRPC client/server artifacts.
- `db/` — database schema, migrations, and SQL queries.
- `frontend/` — SvelteKit frontend application.
- `garage/` — garage configuration.
- `grpc/` — protobuf definitions.
- `migration/` — DB migration tooling.
- `worker/` — Go background worker service.

## Common workflows

### Frontend

In the devcontainer, change into the frontend directory first:

```sh
cd /workspaces/frontend
bun install
bun run dev
bun run build
bun run check
```

### Worker

In the devcontainer, change into the worker directory first:

```sh
cd /workspaces/worker
go run .
```

### Clip service

In the devcontainer, change into the clip directory first:

```sh
cd /workspaces/clip
uv server.py
```

### CLI

In the devcontainer, change into the CLI directory first:

```sh
cd /workspaces/albumscli
go run .
```

## Deployment

The application is deployed using prebuilt container images published to Docker Hub:

- `wutipong/albums-db-migrate` — database migration image
- `wutipong/albums-frontend` — frontend service image
- `wutipong/albums-clip:latest` — clip/gRPC service image
- `wutipong/albums-worker:latest` — background worker image

A dedicated `imgproxy` service should be included in the deployment stack to generate resized image URLs for thumbnails, previews, and high-resolution views.

Recommended deployment layout:

- `imgproxy` service for image transformation and signed URL generation
- `albums-frontend` for the SvelteKit application
- `albums-clip` for the Python clipping service
- `albums-worker` for processing jobs and asset pipeline work
- `albums-db-migrate` executed as part of the database bootstrap or migration phase

### Example deployment composition

> Note: this deployment guidance describes the project’s runtime/container deployment stack and is not the same as the local development setup defined in `.devcontainer/docker-compose.yaml`.

The following is an example container layout for a deployment stack:

```yaml
services:
  db:
    image: postgres:16
    environment:
      POSTGRES_DB: albums
      POSTGRES_USER: albums
      POSTGRES_PASSWORD: albums

  db-migrate:
    image: wutipong/albums-db-migrate
    depends_on:
      - db
    restart: "no"

  imgproxy:
    image: darthsim/imgproxy:latest
    environment:
      IMGPROXY_BIND: ":8080"
      IMGPROXY_USE_S3: "true"
      AWS_REGION: "auto"

  frontend:
    image: wutipong/albums-frontend
    depends_on:
      - db
      - imgproxy
      - clip

  clip:
    image: wutipong/albums-clip:latest

  worker:
    image: wutipong/albums-worker:latest
    depends_on:
      - db
      - clip
```

### Start-up order

1. Start the database and storage dependencies.
2. Run `db-migrate` to apply schema changes.
3. Start `imgproxy` so frontend image-generation endpoints are available.
4. Start `clip` and `worker`.
5. Start `frontend` last so the UI can connect to the backend services once they are reachable.

### PostgreSQL and `pgvector`

The database service relies on the PostgreSQL `pgvector` extension for vector-based functionality. If you are provisioning PostgreSQL yourself, install the extension and create it in the target database before running the application migration step.

Example:

```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

On systems where the extension package is not installed by default, install the `pgvector` PostgreSQL extension package for your OS first, then run the statement above inside the application database.

### Environment variables

The services should be wired with deployment-specific environment variables, including at minimum:

- `IMGPROXY_URL` and `IMGPROXY_KEY` / `IMGPROXY_SALT` for frontend image-signing
- `S3_BUCKET`, `AWS_REGION`, and S3 credentials for object storage access
- database connection settings used by the frontend, worker, and migration image
- clip service endpoint settings used by the worker and frontend integration paths
- any runtime secrets required by the deployment environment

### Runtime notes

- The `imgproxy` service is required for thumbnail, preview, and transformed image URL generation.
- `albums-frontend` is the user-facing SvelteKit application.
- `albums-clip` provides the Python clipping/gRPC service.
- `albums-worker` handles queued processing work and asset pipeline jobs.
- `albums-db-migrate` should be treated as a one-time bootstrap or migration step in the deployment flow.

## Development notes

The repository is intended to be developed inside the provided devcontainer environment. The container already includes the languages, tooling, and service wiring needed for local development, including Bun for the frontend workflow.

Use the devcontainer as the primary development environment when working on the project. Within that environment, the expected frontend commands are Bun-based:

```sh
bun install
bun run dev
bun run build
bun run check
```

The devcontainer setup is for development convenience and is separate from the deployment stack described in the previous section.

## Notes

- The project uses a mixed technology stack: Go, Python, Svelte, SQL, gRPC, and Bun for frontend tooling.
- Database migrations live under `db/migrations/`.
- The frontend and worker are the main runtime services for the application.
- Bun-based commands such as `bun run`, `bun install`, and `bunx` are the expected workflow in this environment.

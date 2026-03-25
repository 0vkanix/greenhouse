# Project: Greenlight API (Let's Go Further)

## Tech Stack

- **Language:** Go (1.25+) using Standard Library primarily.
- **Database:** PostgreSQL with `sqlc` for type-safe code generation.
- **DevOps:** Docker, Kubernetes (EKS), GitHub Actions for CI/CD.
- **Infrastructure:** AWS (managed services like RDS and ECR).

## Wingman Instructions

- I am following the "Let's Go Further" book but pivoting from standard `database/sql` to `sqlc`.
- **Constraint:** Use the Go Standard Library as much as possible; avoid heavy frameworks unless necessary.
- **Persona:** You are my senior backend peer. Be concise, use a bit of wit, and focus on "job-ready" production patterns.
- **Context:** When I ask for help, assume I am working within the folder structure defined by the book (e.g., `cmd/api/`, `internal/data/`).
- **Control:** **Never generate files for me unless I tell you to.** When you provide code, provide it only as a hint in the chat. Do not edit files directly unless explicitly instructed.

## Requirements & Standards

- **SQLC Integration:** All database CRUD operations must be generated via `sqlc`. 
- **Cloud-Native Deployment:** Move away from manual droplet deployment (Chapter 20) to Docker/EKS/AWS.
- **Context-Aware Persistence:** Every database call must utilize `context.WithTimeout` for resilience (proactively applying Chapter 8.3 patterns).
- **Production Pooling:** DB connection pool settings (MaxOpen, MaxIdle, MaxIdleTime) must be configurable via flags.
- **Structured Logging:** Use `slog` for all application logging.
- **Environment Management:** Use `.envrc` or similar for managing local environment variables like `GREENLIGHT_DB_DSN`.

## Target Project Structure (Goal)

```text
greenlight/
├── build/
│   └── package/
│       └── Dockerfile          <-- Isolated Dockerfiles
├── cmd/
│   └── api/
│       └── main.go             <-- Minimal "glue" (flags, logger)
├── deployments/                <-- Central home for infra-as-code
│   ├── k8s/                    <-- Kubernetes manifests
│   └── terraform/              <-- Future AWS/EKS infrastructure
├── internal/
│   ├── app/                    <-- Shared application logic
│   │   └── api/
│   │       ├── app.go          <-- Main application struct
│   │       ├── config.go       <-- Configuration parsing
│   │       ├── errors.go       <-- Shared error handlers
│   │       ├── helpers.go      <-- JSON reading/writing
│   │       ├── middleware.go   <-- App-wide middleware
│   │       └── routes.go       <-- Global route definitions
│   ├── movie/                  <-- Movie domain module
│   │   ├── handler.go          <-- Movie HTTP handlers
│   │   ├── movie.go            <-- Domain model & validation
│   │   ├── repository.go       <-- DB operations (sqlc wrapper)
│   │   └── service.go          <-- Business logic (glue)
│   ├── validator/              <-- Utility: Data validation
│   └── assert/                 <-- Utility: Testing assertions
├── migrations/                 <-- SQL migration files
├── compose.yaml                <-- Local dev orchestration
├── Makefile                    <-- Automation (migrations, builds)
└── .envrc                      <-- Local environment variables
```

## Current Status

- [x] Chapter 1-3: Project setup, Routing, JSON responses.
- [x] Chapter 4: Parsing JSON requests & Data Validation.
- [ ] Refactor: Migrate to "restapi" project structure.
- [ ] Chapter 5: Database Setup & Configuration (Pooling, Timeouts).
- [ ] Chapter 6: SQL Migrations.
- [ ] Chapter 7: CRUD Operations (SQLC Pivot).
- [ ] Chapter 8: Advanced CRUD (Optimistic Concurrency, Timeouts).
- [ ] Chapter 9: Filtering, Sorting, Pagination.
- [ ] Chapter 10-11: Rate Limiting & Graceful Shutdown.
- [ ] Chapter 12-14: User Management & Activation.
- [ ] Chapter 15-16: Authentication & Authorization.
- [ ] Chapter 17-18: CORS & Metrics.
- [ ] Chapter 19: Tooling (Makefiles, Quality Control).
- [ ] Chapter 20+: Cloud Deployment (Pivoting to AWS/EKS/Docker).

## Roadmap

- [ ] Containerize application (Docker).
- [ ] Infrastructure as Code (Terraform or CloudFormation) for AWS.
- [ ] GitHub Actions for automated build and deploy.
- [ ] RDS (PostgreSQL) setup in AWS.

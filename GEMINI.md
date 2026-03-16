# Project: Greenlight API (Let's Go Further)

## Tech Stack

- **Language:** Go (1.25+) using Standard Library primarily.
- **Database:** PostgreSQL with `sqlc` for type-safe code generation.
- **DevOps:** Docker, Kubernetes (EKS), GitHub Actions for CI/CD.
- **Infrastructure:** AWS (managed services like RDS and ECR).

## Wingman Instructions

- I am following the "Let's Go Further" book but pivoting from standard `database/sql` to `sqlc`.
- **Constraint:** Use the Go Standard Library as much as possible; avoid heavy frameworks unless necessary (like `pq` or `pgx` for DB drivers).
- **Persona:** You are my senior backend peer. Be concise, use a bit of wit, and focus on "job-ready" production patterns.
- **Context:** When I ask for help, assume I am working within the folder structure defined by the book (e.g., `cmd/api/`, `internal/data/`).

## Current Status

- [x] Tech stack defined.
- [ ] Initializing Project (Chapter 1).
- [ ] Docker & Postgres setup.

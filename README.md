# Resume App (Go + PostgreSQL + Docker + Kubernetes + AWS + CI/CD)

A concise reference project demonstrating a production-style stack:

- Go REST API (gorilla/mux) with PostgreSQL persistence
- Docker & docker-compose for local development
- Kubernetes manifests (Deployment, Service, ConfigMap, Secret)
- AWS Terraform samples (ECR, RDS) + GitHub Actions CI/CD pipeline
- Basic tests for handlers

## Architecture Overview

```
Client -> HTTP (notes CRUD) -> Go API -> PostgreSQL
                       |--> Health probe
Infra: Containerized -> Deployed to K8s -> Image in ECR -> RDS Postgres
CI/CD: GitHub Actions (test -> build -> push -> deploy)
```

## Endpoints
- `GET /health` -> `{"status":"ok"}`
- `POST /notes` JSON `{title, content}` -> created note
- `GET /notes` -> list
- `GET /notes/{id}` -> single note

## Local Development (Docker Compose)
```bash
cp .env.example .env # (optional if you add one later)
docker compose up --build
curl http://localhost:8080/health
```

## Local Development (Direct)
Requires running Postgres yourself.
```bash
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=resumeapp
go build ./cmd/server && ./server
```

## Makefile Quick Commands
```bash
make build
make dev           # docker compose
make test
make k8s-apply     # applies manifests (needs kube context)
```

## Kubernetes Deployment
1. Replace `<ECR_REPO_URI>` in `k8s/deployment.yaml` or let CI pipeline do it.
2. Apply manifests:
```bash
make k8s-apply
kubectl get pods -n resume-app
```
3. (Optional) Add an Ingress or port-forward:
```bash
kubectl port-forward -n resume-app deploy/resume-app 8080:8080
curl http://localhost:8080/health
```

## AWS Terraform Samples
Directory: `infra/terraform`
- `provider.tf`: AWS provider configuration
- `ecr.tf`: ECR repository for container images
- `rds.tf`: Simplified RDS Postgres instance (requires VPC + subnet IDs via variables)

Usage example (after customizing variables):
```bash
cd infra/terraform
terraform init
terraform plan -var 'private_subnet_ids=["subnet-123","subnet-456"]' -var 'db_security_group_id=sg-abc'
terraform apply
```
Capture outputs: `ecr_repository_url`, `rds_endpoint`.

## CI/CD (GitHub Actions)
Workflow: `.github/workflows/ci.yml`
Stages:
1. Build & test
2. Assume AWS role via OIDC (needs `AWS_IAM_ROLE_ARN` secret)
3. Build & push Docker image to ECR
4. Update and apply Kubernetes manifests

Required secrets:
- `AWS_IAM_ROLE_ARN`: Role allowing ECR push + EKS deploy
- Cluster credentials (kubeconfig or other method) made available to the runner (e.g. store KUBECONFIG as a secret or use another action to fetch)

## Folder Structure
```
cmd/server        # main entrypoint
internal/notes    # domain logic & handlers
internal/testutil # stub repository for tests
migrations        # SQL migrations (simple example)
k8s               # Kubernetes manifests
infra/terraform   # AWS infra samples
.github/workflows # CI pipeline
```

## Extending
- Add authentication (JWT) layer
- Replace gorilla/mux with chi or standard net/http for minimalism
- Introduce structured logging (zap / zerolog)
- Add migration tool (golang-migrate) instead of auto-create
- Add integration tests using Testcontainers

## License
Not specified (add one if publishing publicly).

## Resume Tips
Highlight specific responsibilities implemented:
- "Designed containerized Go service with health probes & readiness checks"
- "Implemented automated CI/CD to ECR & EKS with OIDC credentials"
- "Provisioned AWS infrastructure (RDS, ECR) via Terraform"

Enjoy building!

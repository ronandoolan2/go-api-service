# Transaction API - Proof of Concept

This project is a simple proof-of-concept demonstrating:

1. A Golang API that stores transactions in Postgres.  
2. Deployment to Kubernetes using Docker containers.  
3. Basic monitoring and CI/CD workflows.

## Features
- Minimal REST API endpoint:
  - `POST /api/transaction` -> inserts a new transaction into the database.
- Containerized with multi-stage Docker build.
- Kubernetes manifests for deploying the API and Postgres.
- Example GitHub Actions workflow for CI/CD.

## Prerequisites
- Docker
- Kubernetes cluster (e.g., Minikube, KIND, or a cloud-managed service)
- Kubectl
- (Optionally) GitHub Actions set up for CI/CD

## Quickstart (Local)
1. Clone this repository:  
   ```bash
   git clone https://github.com/youruser/transaction-api.git

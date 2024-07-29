#!/bin/bash

# Build app
podman build -f ./deployments/app/Dockerfile . --tag=gofiber-api

# Build worker
podman build -f ./deployments/worker/cron/Dockerfile . --tag=gofiber-api-cron
podman build -f ./deployments/worker/queue/Dockerfile . --tag=gofiber-api-queue

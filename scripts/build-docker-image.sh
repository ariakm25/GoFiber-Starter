#!/bin/bash

# Build app
docker build . -f ./deployments/app/Dockerfile --tag=gofiber-api

# Build worker
docker build . -f ./deployments/worker/cron/Dockerfile --tag=gofiber-api-cron
docker build . -f ./deployments/worker/queue/Dockerfile --tag=gofiber-api-queue

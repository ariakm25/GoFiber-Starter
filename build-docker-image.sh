#!/bin/bash

# Build app
docker build . -f ./docker/app/Dockerfile --tag=gofiber-api

# Build worker
docker build . -f ./docker/worker/cron/Dockerfile --tag=gofiber-api-cron
docker build . -f ./docker/worker/queue/Dockerfile --tag=gofiber-api-queue

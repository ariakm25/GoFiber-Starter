#!/bin/bash

# Build the application
go build -o ./build/app ./cmd/app/main.go

# Build workers
go build -o ./build/cron ./cmd/worker/cron/cron.go
go build -o ./build/queue ./cmd/worker/queue/queue.go

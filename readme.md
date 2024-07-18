<h3 align="center">Go Fiber Starter</h3>

---

## 📝 Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Run](#run)

## 🧐 About <a name = "about"></a>

Go Fiber App Boilerplate

## Prerequisites

### Install:

- [Go](https://go.dev/)
- [Redis](https://redis.io/)
- [Postgres](https://www.postgresql.org/)

## 🏁 Getting Started <a name = "getting_started"></a>

### Install dependencies

```
go get
```

### Set up environment variables

```
cp .env.example .env
```

## 🔧 Running the app <a name = "run"></a>

run:

```
go run main.go
```

or run with live reload using [Air](https://github.com/air-verse/air):

```
air
```

## 🌱 Database Seed

```
go run cmd/database/seeder/seeder.go --table=all --count=100
```

## ⛏️ Worker

### Queue Worker

```
go run cmd/worker/queue/queue.go
```

### Scheduler Worker

```
go run cmd/worker/cron/cron.go
```

## 📊 Monitoring

Open in browser:

```
localhost:3000/monitoring/jobs
```

## 🍃 Built Using <a name = "built_using"></a>

- [GoFiber](https://gofiber.io/)
- [Paseto](https://paseto.io/)
- [Gorm](https://gorm.io/)
- [Casbin](https://casbin.org/)
- [Asynq](https://github.com/hibiken/asynq)
- [Asynqmon](https://github.com/hibiken/asynqmon)
- [Viper](https://github.com/spf13/viper)
- [Testify](https://github.com/stretchr/testify)
- [Zap](https://github.com/uber-go/zap)
- [Gomail](https://github.com/go-gomail/gomail)
- [Validator](https://github.com/go-playground/validator)

---

---

## Checklist

---

### Auth

- [ ] Logout
- [x] Role & Permission
- [x] Reset Password
- [x] Forgot Password
- [x] Register
- [x] Login

### Mailing

- [ ] Template
- [x] Send

### System

- [x] Scheduler / Cron Job
- [x] Queue
- [x] Logging
- [ ] File Upload: S3

### Data

- [ ] Pagination

### Database

- [x] Auto Migration
- [ ] Versioned Migration
- [x] Seeder

### Monitoring

- [x] Job Monitoring (Scheduler & Queue)

### Caching

- [x] Redis Service

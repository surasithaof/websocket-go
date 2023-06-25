# websocket-go

WebSocket using [Go](https://go.dev/), [Gin framework](https://github.com/gin-gonic/gin) and [Gorilla WebSocket](https://github.com/gorilla/websocket)

- [x] Create websocket server and handle connecting clients
- [x] Unit test
- [x] Client for E2E test
- [x] Docker support
- [x] PingPong health check
- [ ] Example chat app

---

## Setup

[Justfile](https://github.com/casey/just) for running command

```bash
brew install just
```

## Set up .env file

Create `.env` file and set up environment variables (you can copy from `.env.example`)

## Create and migrate database

You can start postgresql database by run this command `docker-compose up`,
And do the migration by run `just db-migrate`

## Start service 🚀

```bash
just run
```

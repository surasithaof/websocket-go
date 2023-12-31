set dotenv-load := true
set export

run:
    go run cmd/main.go
build:
    go build -v -o service ./cmd/main.go
test:
    go test ./...
test-cov:
    go test ./... -coverprofile=coverage.out
show-cov:
    go tool cover -html=coverage.out
generate:
    go generate ./...

env name:
	ln -sf .env.{{name}} .env
envcreate name:
    cp .env.example .env.{{name}}

# release:
#     # autotag --scheme=conventional
release:
    git sv next-version

# https://github.com/pantheon-systems/autotag
#  go install  github.com/git-chglog/git-chglog/cmd/git-chglog
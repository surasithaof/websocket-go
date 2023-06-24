set dotenv-load := true
set export

run:
    go run cmd/main.go
build:
    go build
test:
    go test .../.
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
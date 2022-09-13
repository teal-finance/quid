help:
	# make all          Build both frontend and backend.
	# make front        Build the frontend UI.
	# make quid         Build the backend.
	#
	# make run          Run the backend (serves the frontend static files).
	# make run-dev      Run the backend in dev mode (also serves the frontend).
	# make run-front    Run the frontend in dev mode (NodeJS serves the frontend).
	#
	# make compose-up   Run Quid and Database using podman-compose or docker-compose.
	# make compose-rm   Stop and remove containers.
	#
	# make up           Upgrade dependencies patch version (Go/Node).
	# make up+          Upgrade dependencies minor version (Go/Node).
	# make up++         Upgrade dependencies major version (Node only).
	#
	# make fmt      Go: Generate code and Format code
	# make test     Go: Check build and Test
	# make cov      Go: Browse test coverage
	# make vet      Go: Lint
	#
	# Before "git push" concerning the backend:
	#
	#     make up+go fmt test vet

.PHONY: all
all: front quid

.PHONY: front
front: ui/dist

ui/dist: ui/node_modules ui/node_modules/* ui/node_modules/*/* $(shell find ui/src -type f)
	yarn    --cwd ui build || \
	yarnpkg --cwd ui build

ui/node_modules/*/*: ui/yarn.lock
ui/node_modules/*:   ui/yarn.lock
ui/node_modules:     ui/yarn.lock
	yarn    install --cwd ui --link-duplicates || \
	yarnpkg install --cwd ui --link-duplicates

ui/yarn.lock: ui/package.json
	yarn    install --cwd ui --link-duplicates || \
	yarnpkg install --cwd ui --link-duplicates

quid: cmd/quid/main.go go.sum */*.go */*/*.go
	# go build -o $@
	CGO_ENABLED=0 GOFLAGS="-trimpath -modcacherw" GOLDFLAGS="-d -s -w -extldflags=-static" go build -a -tags osusergo,netgo -installsuffix netgo -o $@ $^

go.sum: go.mod

go.mod:
	go mod tidy
	go mod verify

.PHONY: run
run: cmd/quid/main.go go.sum config.json
	go run $^ -dev -v

.PHONY: run-prod
run-dev: cmd/quid/main.go go.sum config.json
	go run $^

config.json:
	# Create an empty config.json file and customize it:
	#
	#    ./quid -conf
	#    vim config.json
	#
	# Initialize the PostreSQL database:
	#
	#    ./quid -init
	#

.PHONY: run-front
run-front:
	cd ui && \
	{ yarn    --link-duplicates && yarn    dev; } || \
	{ yarnpkg --link-duplicates && yarnpkg dev; }

define help

Podman and Docker are both supported.
Four options to setup your container engine:

1. Install commands podman (Go) and podman-compose (Python)

  sudo apt install podman python3-pip
  python3 -m pip install --user --upgrade pip
  python3 -m pip install --user --upgrade podman-compose

2. Install commands docker (Go) and docker-compose (Python)

  sudo apt install docker.io python3-pip
  python3 -m pip install --user --upgrade pip
  python3 -m pip install --user --upgrade docker-compose

3. Install command docker (Go) and its "docker compose" (Go) in hybride mode.
   Replace v2.5.1 by the version you want.

  sudo apt install docker.io curl ca-certificates
  mkdir -pv ~/.docker/cli-plugins
  curl -SL https://github.com/docker/compose/releases/download/v2.5.1/docker-compose-linux-x86_64 -o ~/.docker/cli-plugins/docker-compose
  chmod +x ~/.docker/cli-plugins/docker-compose
  docker compose version

3. Install command docker (Go) and its "docker compose" (Go) using docker.com only

  sudo apt purge --purge --autoremove docker docker-engine docker.io containerd runc

  Apply https://docs.docker.com/engine/install/ and finally:

  sudo apt install docker-compose-plugin
  docker compose version

endef

export help

.PHONY: compose-up
compose-up:
	@export DOCKER_BUILDKIT=1          ; \
	export COMPOSE_DOCKER_CLI_BUILD=1  ; \
	{ command -v podman-compose                         && set -x && podman-compose -f compose.yml up --build -d;} || \
	{ command -v docker-compose                         && set -x && docker-compose -f compose.yml up --build -d;} || \
	{ command -v docker && docker help|grep -wq compose && set -x && docker compose -f compose.yml up --build -d;} || \
	{ echo "$$help"; false; }

	# Open browser on localhost:8082 if Quid is running
	@{ command -v podman && set -x && podman ps -qf name=quid || set -x && docker ps -qf name=quid ; } | \
	grep -s . && xdg-open http://localhost:8082

	# Print containers logs. [Ctrl+C] to stop the logs printing.
	@{ command -v podman-compose                         && set -x && docker-compose -f compose.yml logs --follow;} || \
	{  command -v docker-compose                         && set -x && podman-compose -f compose.yml logs --follow;} || \
	{  command -v docker && docker help|grep -wq compose && set -x && docker compose -f compose.yml logs --follow;}

.PHONY: compose-rm
compose-rm:
	{ command -v podman-compose                         && set -x && podman-compose -f compose.yml ps -q  | xargs podman stop && \
	                                                                 podman-compose -f compose.yml ps -aq | xargs podman rm;} || \
	{ command -v docker-compose                         && set -x && docker-compose -f compose.yml down                    ;} || \
	{ command -v docker && docker help|grep -wq compose && set -x && docker compose -f compose.yml down                    ;}

.PHONY: up up+
up: up-ui up-go
up+: up+ui up+go

.PHONY: up-ui
up-ui:
	cd ui && \
	{ yarn    --link-duplicates && yarn    upgrade-interactive --link-duplicates; } || \
	{ yarnpkg --link-duplicates && yarnpkg upgrade-interactive --link-duplicates; }

.PHONY: up+ui
up+ui:
	cd ui && \
	{ yarn --link-duplicates && yarn up+; } || \
	{ yarn --link-duplicates && yarn up+; }

.PHONY: up++
up++: up+go
	cd ui && \
	{ yarn    --link-duplicates && yarn    upgrade-interactive --link-duplicates --latest --tilde; } || \
	{ yarnpkg --link-duplicates && yarnpkg upgrade-interactive --link-duplicates --latest --tilde; }
    # flag --tilde prepends the new version with "~" that limits vanilla upgrade to patch only
    # flag --caret prepends the new version with "^" allowing upgrading the minor number

.PHONY: up-go
up-go: go.sum
	GOPROXY=direct go get -t -u=patch all
	go mod tidy

.PHONY: up+go
up+go: go.sum
	go get -t -u all
	go mod tidy

go.sum: go.mod
	go mod tidy

.PHONY: fmt
fmt:
	go generate ./...
	go run mvdan.cc/gofumpt@latest -w -extra -l -lang 1.19 .

.PHONY: test
test:
	go build ./...
	go test -race -vet all -run=. -bench=. -benchmem -benchtime 10ms -tags=quid -coverprofile=code-coverage.out ./...

code-coverage.out: go.sum */*/*.go */*/*/*.go
	go test -race -vet all -run=. -bench=. -benchmem -benchtime 10ms -tags=quid -coverprofile=code-coverage.out ./...

.PHONY: cov
cov: code-coverage.out
	go tool cover -html code-coverage.out

.PHONY: vet
vet:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --fix || true
	go run ./cmd/quid -dev

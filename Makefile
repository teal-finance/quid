help:
	# Use 'make <target>' where <target> is one of:
	#
	# all          Build both frontend and backend.
	# front        Build the frontend UI.
	# quid         Build the backend.
	#
	# run          Run the backend (serves the frontend static files).
	# run-dev      Run the backend in dev mode (also serves the frontend).
	# run-front    Run the frontend in dev mode.
	#
	# compose-up   Run Quid and Database from docker-compose or podman-compose.
	# compose-rm   Stop and remove containers
	#
	# upg-patch    Upgrade dependencies patch version (Go/Node).
	# upg-minor    Upgrade dependencies minor version (Go/Node).
	# upg-more     Upgrade dependencies major version (Node only).

.PHONY: all
all: front quid

.PHONY: front
front: ui/dist

ui/dist: ui/node_modules ui/node_modules/* ui/node_modules/*/* $(shell find ui/src -type f)
	yarn --cwd ui build

ui/node_modules:     ui/yarn.lock
ui/node_modules/*:   ui/yarn.lock
ui/node_modules/*/*: ui/yarn.lock

ui/yarn.lock: ui/package.json
	yarn --cwd ui --link-duplicates

quid: go.sum main.go quidlib/*.go quidlib/*/*.go quidlib/*/*/*.go
	go build -o $@

go.sum: go.mod

go.mod:
	go mod tidy
	go mod verify

.PHONY: run
run: go.sum config.json
	go run main.go

.PHONY: run-dev
run-dev: go.sum config.json
	go run main.go -dev -v

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
	yarn --cwd ui --link-duplicates
	yarn --cwd ui dev

define help
The compose.yml file supports both docker and podman.

To install docker and docker-compose:

sudo apt install docker.io python3-pip
python3 -m pip install --user --upgrade pip
python3 -m pip install --user --upgrade docker-compose

The same to install podman and podman-compose:

sudo apt install podman python3-pip
python3 -m pip install --user --upgrade pip
python3 -m pip install --user --upgrade podman-compose
endef

export help
.PHONY: compose-up
compose-up:
	@COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 \
	docker-compose -f compose.yml up --build -d || \
	podman-compose -f compose.yml up --build -d || \
	{ echo "$$help"; false; }
	docker-compose -f compose.yml logs --follow || \
	podman-compose -f compose.yml logs --follow

.PHONY: compose-rm
compose-rm:
	docker-compose -f compose.yml down || \
	{ podman-compose -f compose.yml ps -q  | xargs podman stop && \
	  podman-compose -f compose.yml ps -aq | xargs podman rm ; }

.PHONY: upg-patch upg-minor
upg-patch: upg-patch-ui upg-patch-go
upg-minor: upg-minor-ui upg-minor-go

.PHONY: upg-patch-ui
upg-patch-ui:
	yarn --cwd ui --link-duplicates
	yarn upgrade-interactive --cwd ui --link-duplicates

.PHONY: upg-patch-go
upg-patch-go:
	GOPROXY=direct go get -t -u=patch

.PHONY: upg-minor-ui
upg-minor-ui:
	yarn --cwd ui --link-duplicates
	yarn --cwd ui upg-minor

.PHONY: upg-minor-go
upg-minor-go:
	go get -t -u

.PHONY: upg-more
upg-more:
	yarn --cwd ui --link-duplicates
    # flag --tilde prepends the new version with "~" that limits vanilla upgrade to patch only
    # flag --caret prepends the new version with "^" allowing upgrading the minor number
	yarn upgrade-interactive --cwd ui --link-duplicates --latest --tilde

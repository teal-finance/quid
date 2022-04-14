help:
	# Use 'make <target>' where <target> is one of:
	#
	# all          Build both frontend and backend
	# front        Build the frontend UI
	# quid         Build the backend
	#
	# run          Run the backend (serves the frontend static files)
	# run-dev       Run the backend in dev mode (also serves the frontend)
	# run-front     Run the frontend in dev mode
	#
	# up-patch     Upgrade dependencies patch version (Go/Node)
	# up-minor     Upgrade dependencies minor version (Go/Node)
	# up-more      Upgrade dependencies major version (Node only)

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

.PHONY: up-patch
.PHONY: up-minor
up-patch: up-patch-ui up-patch-go
up-minor: up-minor-ui up-minor-go

.PHONY: up-patch-ui
up-patch-ui:
	yarn --cwd ui --link-duplicates
	yarn upgrade-interactive --cwd ui --link-duplicates

.PHONY: up-patch-go
up-patch-go:
	GOPROXY=direct go get -t -u=patch

.PHONY: up-minor-ui
up-minor-ui:
	yarn --cwd ui --link-duplicates
	yarn --cwd ui up-minor

.PHONY: up-minor-go
up-minor-go:
	go get -t -u

.PHONY: up-more
up-more:
	yarn --cwd ui --link-duplicates
    # flag --tilde prepends the new version with "~" that limits vanilla upgrade to patch only
    # flag --caret prepends the new version with "^" allowing upgrading the minor number
	yarn upgrade-interactive --cwd ui --link-duplicates --latest --tilde

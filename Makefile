help:
	# Use 'make <target>' where <target> is one of:
	#
	# all          Build both backend and frontend
	# quid         Build the backend
	# front        Build the frontend UI
	#
	# run          Run the backend
	# rundev       Run the backend in dev mode
	# runfront     Run the frontend in dev mode
	#
	# up-patch     Upgrade dependencies patch version (Go/Node)
	# up-minor     Upgrade dependencies minor version (Go/Node)
	# up-more      Upgrade dependencies major version (Node only)

.PHONY: all
all: front quid

quid: go.sum main.go quidlib/*.go quidlib/*/*.go quidlib/*/*/*.go
	go build -o $@

go.mod:
	go mod tidy
	go mod verify

go.sum: go.mod

.PHONY: up-patch
up-patch:
	GOPROXY=direct go get -t -u=patch
	yarn upgrade-interactive --cwd ui --link-duplicates

.PHONY: up-minor
up-minor:
	go get -t -u
	yarn --cwd ui up-minor

.PHONY: up-more
up-more:
	yarn upgrade-interactive --cwd ui --link-duplicates --latest

.PHONY: front
front: ui/dist

ui/dist: ui/node_modules/*/* $(shell find ui/src -type f)
	yarn --cwd ui build

ui/node_modules/*/*: ui/yarn.lock

ui/yarn.lock: ui/package.json
	yarn --cwd ui --link-duplicates

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

.PHONY: run
run: go.sum config.json
	go run main.go

.PHONY: rundev
rundev:
	go run main.go --dev

.PHONY: runfront
runfront:
	yarn --cwd ui dev

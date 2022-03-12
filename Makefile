help:
	# Please use 'make <target>' where <target> is one of:
	#
	#   build          -- build the backend
	#   buildoldfront  -- build the old frontend ui
	#
	#   run            -- run the backend
	#   rundev         -- run the backend in dev mode
	#   runfront       -- run the frontend in dev mode

build:
	go build
.PHONY: build

buildoldfront:
	cd quidui && yarn && yarn build
.PHONY: buildoldfront

run:
	go run main.go
.PHONY: run

rundev:
	go run main.go --dev
.PHONY: rundev

runfront:
	cd adminui && yarn dev
.PHONY: runfront

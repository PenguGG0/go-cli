## help: print this help message
.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ==================================================================================== #
# BUILD
# ==================================================================================== #

# check the time with 'bin/xxx.exe --build-time'
current_time = $(shell powershell -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ss'")
linker_flags = '-s -w'

## build/wc: build wc/ to bin/wc
.PHONY: build/wc
build/wc:
	@echo "Building wc..."
	go build -ldflags=$(linker_flags) -o=./bin/wc ./wc

## build/todo: build todo/ to bin/todo
.PHONY: build/todo
build/todo:
	@echo "Building todo..."
	go build -ldflags=$(linker_flags) -o=./bin/todo ./todo

## build/mdp: build mdp/ to bin/mdp
.PHONY: build/mdp
build/mdp:
	@echo "Building mdp..."
	go build -ldflags=$(linker_flags) -o=./bin/mdp ./mdp

## build/walk: build walk/ to bin/walk
.PHONY: build/walk
build/walk:
	@echo "Building walk..."
	go build -ldflags=$(linker_flags) -o=./bin/walk ./walk

## build/unarchive: build unarchive/ to bin/unarchive
.PHONY: build/unarchive
build/unarchive:
	@echo "Building unarchive..."
	go build -ldflags=$(linker_flags) -o=./bin/unarchive ./unarchive

## build/colStats: build colStats/ to bin/colStats
.PHONY: build/colStats
build/colStats:
	@echo "Building colStats..."
	go build -ldflags=$(linker_flags) -o=./bin/colStats ./colStats

## build/goci: build goci/ to bin/goci
.PHONY: build/goci
build/goci:
	@echo "Building goci..."
	go build -ldflags=$(linker_flags) -o=./bin/goci ./goci

## build/pScan: build pScan/ to bin/pScan
.PHONY: build/pScan
build/pScan:
	@echo "Building pScan..."
	go build -ldflags=$(linker_flags) -o=./bin/pScan./pScan

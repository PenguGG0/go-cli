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

## build/wc: build wc/ to bin/wc.exe
.PHONY: build/wc
build/wc:
	@echo "Building wc..."
	go build -ldflags=$(linker_flags) -o=./bin/wc.exe ./wc

## build/todo: build todo/ to bin/todo.exe
.PHONY: build/todo
build/todo:
	@echo "Building todo..."
	go build -ldflags=$(linker_flags) -o=./bin/todo.exe ./todo

## build/mdp: build mdp/ to bin/mdp.exe
.PHONY: build/mdp
build/mdp:
	@echo "Building mdp..."
	go build -ldflags=$(linker_flags) -o=./bin/mdp.exe ./mdp

## build/walk: build walk/ to bin/walk.exe
.PHONY: build/walk
build/walk:
	@echo "Building walk..."
	go build -ldflags=$(linker_flags) -o=./bin/walk.exe ./walk

## build/unarchive: build unarchive/ to bin/unarchive.exe
.PHONY: build/unarchive
build/unarchive:
	@echo "Building unarchive..."
	go build -ldflags=$(linker_flags) -o=./bin/unarchive.exe ./unarchive

## build/colStats: build colStats/ to bin/colStats.exe
.PHONY: build/colStats
build/colStats:
	@echo "Building colStats..."
	go build -ldflags=$(linker_flags) -o=./bin/colStats.exe ./colStats

## build/goci: build goci/ to bin/goci.exe
.PHONY: build/goci
build/goci:
	@echo "Building goci..."
	go build -ldflags=$(linker_flags) -o=./bin/goci.exe ./goci

## build/pScan: build pScan/ to bin/pScan.exe
.PHONY: build/pScan
build/pScan:
	@echo "Building pScan..."
	go build -ldflags=$(linker_flags) -o=./bin/pScan.exe ./pScan
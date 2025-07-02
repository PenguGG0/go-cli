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

## build/wc: build cmd/wc to bin/wc.exe
.PHONY: build/wc
build/wc:
	@echo "Building cmd/wc..."
	go build -ldflags=$(linker_flags) -o=./bin/wc.exe ./cmd/wc

## build/todo: build cmd/todo to bin/todo.exe
.PHONY: build/todo
build/todo:
	@echo "Building cmd/todo..."
	go build -ldflags=$(linker_flags) -o=./bin/todo.exe ./cmd/todo

## build/mdp: build cmd/mdp to bin/mdp.exe
.PHONY: build/mdp
build/mdp:
	@echo "Building cmd/mdp..."
	go build -ldflags=$(linker_flags) -o=./bin/mdp.exe ./cmd/mdp

## build/walk: build cmd/walk to bin/walk.exe
.PHONY: build/walk
build/walk:
	@echo "Building cmd/mdp..."
	go build -ldflags=$(linker_flags) -o=./bin/walk.exe ./cmd/walk
.PHONY: help build/wc

## help: print this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ==================================================================================== #
# BUILD
# ==================================================================================== #

current_time = $(shell powershell -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ss'")
linker_flags = '-s -X main.buildTime=$(current_time)'

## build/wc: build cmd/wc to bin/wc.exe
build/wc:
	@echo "Building cmd/wc..."
	go build -ldflags=$(linker_flags) -o=./bin/wc.exe ./cmd/wc

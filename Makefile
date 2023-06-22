# Pattern #1 example: "example : description = Description for example target"
# Pattern #2 example: "### Example separator text
help: HELP_SCRIPT = \
	if (/^([a-zA-Z0-9-\.\/]+).*?: description\s*=\s*(.+)/) { \
		printf "\033[34m%-40s\033[0m %s\n", $$1, $$2 \
	} elsif(/^\#\#\#\s*(.+)/) { \
		printf "\033[33m>> %s\033[0m\n", $$1 \
	}

.PHONY: help
help:
	@perl -ne '$(HELP_SCRIPT)' $(MAKEFILE_LIST)

.PHONY: build
build: description = Build everything
build:
	tinygo build -o build/match.wasm -no-debug -scheduler=none -target wasi src/match/main.go
	tinygo build -o build/transform.wasm -no-debug -scheduler=none -target wasi src/transform/main.go


.PHONY: setup/darwin
setup/darwin: description = Install toolkit for building on macOS M1
setup/darwin:
	brew tap tinygo-org/tools && brew install tinygo wasmtime


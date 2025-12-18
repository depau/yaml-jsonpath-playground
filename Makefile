DIST := dist
WASM_EXEC := $(shell go env GOROOT)/lib/wasm/wasm_exec.js

.PHONY: all build copy-assets clean

all: build

build: $(DIST) $(DIST)/main.wasm copy-assets

$(DIST):
	mkdir -p $(DIST)

$(DIST)/main.wasm:
	GOOS=js GOARCH=wasm go build -o $@ .

copy-assets:
	cp index.html $(DIST)/
	cp $(WASM_EXEC) $(DIST)/

clean:
	rm -rf $(DIST)

GO_CMD=go
BUILD_PATH=bin/wombat

# default
.PHONY: all
all: build

.PHONY: build
build:
	$(GO_CMD) build -o $(BUILD_PATH) main.go

.PHONY: run
run:
	$(BUILD_PATH)

s.PHONY: clean
clean:
	rm -rf bin
	rm -rf *.db
	rm -rf *.db-journal
	
.PHONY: rebuild
rebuild: clean build

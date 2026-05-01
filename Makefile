BIN_NAME := rdw



## ---------- UTILS
default: help

.PHONY: help
help: ## Show this menu
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_][a-zA-Z0-9._-]*:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## Clean all temp files
	@sudo rm -f coverage.* .run .build



## ---------- SETUP
.PHONY: install
install: ## install all requirements
	@yay -S --needed libayatana-appindicator-glib gtk3



## ---------- FORMAT & LINT
.PHONY: fmt
fmt: ## format the code
	@go fmt ./...

.PHONY: vet
vet: ## run static analysis
	@go vet ./...



## ---------- MAIN
.PHONY: build
build: ## Build the code
	@go build -o .build/$(BIN_NAME) ./...

.PHONY: run
run: ## run the app
	@go run ./...

.PHONY: deploy
deploy: build ## deploy the app
	@sudo mv .build/$(BIN_NAME) /usr/local/bin/$(BIN_NAME)
	@rmdir .build
	@sudo chmod +x /usr/local/bin/$(BIN_NAME)
	@systemctl --user restart $(BIN_NAME)

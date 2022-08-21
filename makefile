GO               = go
GOBIN			 ?= $(PWD)/bin
MIGRATIONS_DIR	 = ./db/migrations/
STORAGE_DSN      = postgres://wallet_user:psw@localhost:5432/wallet_db?sslmode=disable

run: stop up

build: build-service build-walletctl

build-service: ## Build binary
	$(info $(M) building service...)
	@GOARCH=$(GOARCH) CGO_ENABLED=0 GOOS=linux $(GO) build -a -installsuffix cgo -o $(GOBIN)/wallet ./cmd/service/*.go

build-walletctl: ## Build walletctl
	$(info $(M) building walletctl...)
	@GOARCH=$(GOARCH) CGO_ENABLED=0 GOOS=linux $(GO) build $(GCFLAGS)  $(LDFLAGS) -o $(GOBIN)/walletctl ./cmd/walletctl/*.go

mod:
	GO111MODULE=on go mod tidy

up:
	docker-compose -f docker-compose.yml up

stop:
	docker-compose -f docker-compose.yml stop

down:
	docker-compose -f docker-compose.yml down

test:
	docker-compose -f docker-compose.test.yml up --build

test-cleanup:
	docker-compose -f docker-compose.test.yml down --volumes

watch: install-tools ; ## Run binaries that rebuild themselves on changes
	$(info $(M) run...)
	@$(GOBIN)/air -c .air.toml

.PHONY: install-tools
install-tools: $(GOBIN) ## Install tools needed for development
	$(info $(GOBIN) install tools needed for development...)
	@GOBIN=$(GOBIN) $(GO) install -tags 'postgres' \
		github.com/cosmtrek/air \
		github.com/golang-migrate/migrate/v4/cmd/migrate \
		github.com/golang/mock/mockgen \
		github.com/deepmap/oapi-codegen/cmd/oapi-codegen

db-migrate: ## Run migrate command
	$(info $(M) running DB migrations...)
	$(info $(GOBIN) running DB migrations...)
	@$(GOBIN)/migrate -path "$(MIGRATIONS_DIR)" -database "$(STORAGE_DSN)" up

db-migrate-down: ## Run migrate command
	$(info $(M) running DB migrations...)
	@$(GOBIN)/migrate -path "$(MIGRATIONS_DIR)" -database "$(STORAGE_DSN)" down

db-load-fixtures: ## Load fixtures into database
	$(info $(M) loading fixtures into DB (example data)...)
	@docker-compose exec server ./bin/walletctl load-fixtures

install-linter: $(GOBIN)
	@GOBIN=$(GOBIN) $(GO) install -mod=readonly \
	github.com/golangci/golangci-lint/cmd/golangci-lint

lint: install-linter ; $(info $(M) running linters...)
	@$(GOBIN)/golangci-lint run --timeout 5m0s ./...

.PHONY: $(GOBIN)
$(GOBIN):
	@mkdir -p $(GOBIN)

COVERAGE_FILE := coverage.out

.PHONY: test test-unix coverage

# For Windows (PowerShell or CMD with findstr)
test:
	@gotestsum -f dots -- -failfast -covermode=count -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -func=$(COVERAGE_FILE) | findstr "total"
	@echo [make test] Done

# For Linux/macOS (Bash, Zsh, etc)
test-unix:
	@gotestsum -f dots -- -failfast -covermode=count -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -func=$(COVERAGE_FILE) | grep 'total' | sed -e 's/\t\+/ /g'
	@echo [make test-unix] Done

coverage:
	@go tool cover -html=$(COVERAGE_FILE)

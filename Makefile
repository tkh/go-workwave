NAME := workwave
BIN_DIR := $(GOPATH)/bin

.PHONY: staticcheck test

STATICCHECK := $(BIN_DIR)/staticcheck
$(STATICCHECK):
	@echo "+ $@"
	@go get honnef.co/go/tools/cmd/staticcheck

staticcheck: $(STATICCHECK)
	@echo "+ $@"
	@$(STATICCHECK) ./...

test:
	@echo "+ $@"
	@go test -v -cover ./...

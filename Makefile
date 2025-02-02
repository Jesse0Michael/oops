.PHONY: test build

# Source a local .env
ifneq (,$(wildcard ./.env))
    include .env
    export
endif
	
#################################################################################
# TEST COMMANDS
#################################################################################
test:
	go test -cover ./... 

lint:
	golangci-lint run ./...

cover:
	go test -coverpkg ./... -coverprofile coverage.out ./... && go tool cover -html=coverage.out

vuln: dependencies
	govulncheck -test ./...


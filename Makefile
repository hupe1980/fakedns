PROJECTNAME=$(shell basename "$(PWD)")

# Go related variables.
# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: help
## help: Prints this help message
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: setup
## setup: Setup installes dependencies
setup:
	@go mod tidy

.PHONY: run
## run: Run fileserve
run: 
	@go run $$(ls -1 cmd/*.go | grep -v _test.go) -h 

.PHONY: test
## test: Runs go test with default values
test: 
	@go test -v -race -count=1  ./...
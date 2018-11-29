COVERAGE_FILE=coverage.txt

all: test
clean:
	go clean
	rm -f $(COVERAGE_FILE)
cover:
	go test -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic -v
	go tool cover -html=$(COVERAGE_FILE)
get:
	go get -d -t ./...
test:
	go test -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic -v

.PHONY: all clean cover get test
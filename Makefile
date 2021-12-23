COVERAGE_FILE=coverage.txt

all: test
clean:
	go clean
	rm -f $(COVERAGE_FILE)
cover: test
	go tool cover -html=$(COVERAGE_FILE)
doc:
	lsof -nti:6060 | xargs kill -9
	godoc -http=:6060 &
	sleep 3
	open http://localhost:6060/pkg/github.com/weathersource/
get:
	go get -d -t ./...
test:
	clear
	go test -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic -v

.PHONY: all clean cover get test
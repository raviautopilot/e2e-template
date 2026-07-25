.PHONY: all deps test-api test-ui test-all clean

all: deps test-all

deps:
	go mod download
	go mod tidy

test-api:
	go test -v ./tests/api/...

test-ui:
	go test -v ./tests/ui/...

test-all:
	go test -v ./tests/...

clean:
	rm -rf requests/ reports/ screenshots/

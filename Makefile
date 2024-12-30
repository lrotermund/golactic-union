API_MAIN=cmd/api/main.go
API_BINARY_NAME=bin/api.out
CLI_MAIN=cmd/cli/main.go
CLI_BINARY_NAME=bin/cli.out

build-api:
	go build -o ${API_BINARY_NAME} ${API_MAIN}

build-cli:
	go build -o ${CLI_BINARY_NAME} ${CLI_MAIN}

test:
	go test -v ./...

run-api:
	go build -o ${API_BINARY_NAME} ${API_MAIN}
	./${API_BINARY_NAME}

run-cli:
	go build -o ${CLI_BINARY_NAME} ${CLI_MAIN}
	./${CLI_BINARY_NAME}

clean:
	go clean
	rm ${API_BINARY_NAME}
	rm ${CLI_BINARY_NAME}

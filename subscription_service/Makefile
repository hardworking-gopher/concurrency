BINARY_NAME=myapp
PG_CONN_STR="host=localhost port=5432 user=postgres password=pwd dbname=subscription sslmode=disable timezone=UTC connect_timeout=5"
REDIS_ADDRESS="127.0.0.1:6379"

## build: Build binary
build:
	@echo "Building..."
	env CGO_ENABLED=0  go build -ldflags="-s -w" -o ${BINARY_NAME} ./cmd
	@echo "Built!"

## run: builds and runs the application
run: build
	@echo "Starting..."
	@env PG_CONN_STR=${PG_CONN_STR} REDIS_ADDRESS=${REDIS_ADDRESS} ./${BINARY_NAME}
	@echo "Started!"

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME}
	@echo "Cleaned!"

## start: an alias to run
start: run

## stop: stops the running application
stop:
	@echo "Stopping..."
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Stopped!"

## restart: stops and starts the application
restart: stop start

## test: runs all tests
test:
	go test -v ./...
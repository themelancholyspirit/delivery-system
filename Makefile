build:
	@go build -o ./api/bin
run: build
	@./api/bin
test:
	@go test ./...
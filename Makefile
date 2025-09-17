run:
	@echo "Running app..."
	go fmt ./...
	go run cmd/web/main.go -m

format:
	@echo "Formatting code..."
	go fmt ./...
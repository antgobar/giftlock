run:
	@echo "Building frontend..."
	cd frontend && npm install && npm run build
	@echo "Running app..."
	go fmt ./...
	go run cmd/web/main.go -m

run-frontend:
	@echo "Running frontend..."
	cd frontend && npm run dev

format:
	@echo "Formatting code..."
	go fmt ./...
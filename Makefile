run:
# 	@echo "Building frontend..."
# 	cd frontend && npm install && npm run build
	@echo "Running app..."
	go fmt ./...
	go run cmd/web/main.go -m

run-frontend:
	@echo "Running frontend..."
	cd frontend && nvm use 22; npm install && npm run build && npm run dev

build-frontend:
	@echo "Building frontend..."
	cd frontend && npm install && npm run build

format:
	@echo "Formatting code..."
	go fmt ./...
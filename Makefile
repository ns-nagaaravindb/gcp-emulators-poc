.PHONY: start stop demo gcs pubsub clean logs

# Start emulators
start:
	@echo "Starting emulators..."
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 10

# Stop emulators  
stop:
	@echo "Stopping emulators..."
	docker-compose down

# Run full demo
demo: start
	@echo "Running full demo..."
	@./run-demo.sh

# Run only GCS demo
gcs: 
	@echo "Running GCS demo..."
	@export STORAGE_EMULATOR_HOST=localhost:4443 && go run gcs.go

# Run only Pub/Sub demo
pubsub:
	@echo "Running Pub/Sub demo..."
	@export PUBSUB_EMULATOR_HOST=localhost:8085 && timeout 15s go run pubsubemulator.go

# View logs
logs:
	docker-compose logs -f

# Clean up everything
clean: stop
	docker-compose down -v
	docker system prune -f

# Install dependencies
deps:
	go mod tidy
	go mod download


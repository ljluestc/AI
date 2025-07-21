.PHONY: build run test clean docker-build docker-run install-deps frontend-dev backend-dev

# Default Go binary name
BINARY_NAME=teathis-server
# Default Docker image name
DOCKER_IMAGE=teathis

# Build the Go application
build:
	go build -o $(BINARY_NAME) .

# Run the application locally
run: build
	./$(BINARY_NAME)

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf ./web/build

# Install dependencies
install-deps:
	go mod download
	go get github.com/go-git/go-git/v5
	go get github.com/neo4j/neo4j-go-driver/v4
	go get gonum.org/v1/gonum
	go mod tidy
	go get github.com/go-git/go-git/v5
	go get github.com/neo4j/neo4j-go-driver/v4/neo4j
	go get github.com/gorilla/mux
	go get github.com/stretchr/testify
	go get gonum.org/v1/gonum
	cd web && npm install
	
	# Fetch all dependencies and update go.sum
	fetch-deps:
	chmod +x scripts/fetch_dependencies.sh
	./scripts/fetch_dependencies.sh

# Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# Run Docker container
docker-run: docker-build
	docker run -p 8080:8080 $(DOCKER_IMAGE)

# Start development environment with Docker Compose
dev-up:
	docker-compose up -d

# Stop development environment
dev-down:
	docker-compose down

# Run frontend development server
frontend-dev:
	cd web && npm start

# Run backend in development mode
backend-dev:
	go run main.go

# Build frontend
frontend-build:
	cd web && npm run build

# Deploy to production (placeholder - customize as needed)
deploy:
	@echo "Deploying to production..."
	# Add your deployment commands here

# GCS and Pub/Sub Emulator Development Environment

A complete local development environment for testing Google Cloud Storage (GCS) and Pub/Sub functionality using Docker-based emulators. This project provides a fully working proof-of-concept with Go applications demonstrating real-world cloud operations locally.

## 🎯 **TL;DR - Quick Demo**

```bash
# Prerequisites: Docker Desktop running + Go installed
./run-demo.sh

# Or just verify code compiles (no Docker needed):
./verify-code.sh
```
**That's it!** One command runs the complete GCS and Pub/Sub demonstration.

> **⚠️ If you encounter issues**: Check the [Troubleshooting](#-troubleshooting) section below for common problems and solutions.



---

## 🚀 Features

- **Local GCS Emulation**: Test bucket operations, file uploads/downloads, and object listing
- **Local Pub/Sub Emulation**: Test topic creation, message publishing, and subscription handling
- **Docker-based Setup**: Isolated, reproducible environment with no external dependencies
- **Go Demonstration Code**: Complete working examples for both services
- **Automated Testing**: Scripts and Makefile for easy setup and testing

## 📁 Project Structure

```
gcs-emulator/
├── docker-compose.yml       # Docker services configuration for emulators
├── gcs.go                   # GCS operations demo (bucket, file operations)
├── pubsubemulator.go        # Pub/Sub operations demo (publish/subscribe)
├── go.mod                   # Go module dependencies
├── go.sum                   # Go dependency checksums
├── Makefile                 # Build and run commands
├── run-demo.sh              # Automated demo script
├── verify-code.sh           # Code compilation verification
└── README.md                # This documentation
```



## 🚀 Quick Start

###  One-Command Demo (Simplest Way)
```bash
# 1. Start Docker Desktop application first
# 2. Navigate to project directory
cd gcs-emulator

# 3. Run the complete demo with one command
./run-demo.sh
```



### Alternative Options

#### Using Makefile Commands
```bash
make demo       # Same as run-demo.sh but through Makefile
make start      # Start emulators only
make gcs        # Run GCS demo only  
make pubsub     # Run Pub/Sub demo only
make stop       # Stop emulators
make clean      # Stop and remove all containers/volumes
```

```bash
🚀 Starting GCS and Pub/Sub Emulator Demo
========================================
📦 Starting emulators with Docker Compose...
⏳ Waiting for services to start...
🔍 Checking service status...
✅ Fake GCS Server is running on http://localhost:4443
✅ Pub/Sub Emulator is running on http://localhost:8085
📥 Downloading Go dependencies...

🎯 Running GCS Demo...
=====================
Created bucket 'test-bucket'
Successfully wrote file 'test-file.txt' to bucket 'test-bucket'
Listing files in bucket:
- test-file.txt (size: 44 bytes, created: 2025-09-30T02:02:42Z)
Reading file 'test-file.txt' content:
Content: Hello, this is a test file for GCS emulator!

🎯 Running Pub/Sub Demo...
==========================
Created topic: test-topic
Created subscription: test-subscription
Starting message subscriber...
Published: Message 1
Received message: Message 1: Hello from Pub/Sub emulator!
  Attributes: map[timestamp:2025-09-30T07:32:45+05:30]
...

✨ Demo completed!
```



## Manual Step-by-Step

#### Step 1: Start Emulators
```bash
# Start Docker containers
docker-compose up -d

# Verify containers are running
docker-compose ps
```

#### Step 2: Verify Services
```bash
# Wait for services to initialize (10-15 seconds)
sleep 10

# Test GCS emulator
curl http://localhost:4443
# Should return: "ok"

# Test Pub/Sub emulator  
curl http://localhost:8085
# Should return: "Ok"
```

#### Step 3: Install Dependencies
```bash
# Download Go modules
go mod tidy

# Verify code compiles correctly
./verify-code.sh
```

#### Step 4: Set Environment Variables
```bash
export STORAGE_EMULATOR_HOST=localhost:4443
export PUBSUB_EMULATOR_HOST=localhost:8085
```

#### Step 5: Run Demonstrations
```bash
# Run GCS demo
go run gcs.go

# Run Pub/Sub demo (in separate terminal or with timeout)
timeout 15s go run pubsubemulator.go
```

#### Step 6: Cleanup
```bash
# Stop emulators
docker-compose down

# Remove volumes (optional)
docker-compose down -v
```

## 📋 What the Programs Demonstrate

### GCS Operations (`gcs.go`)
The GCS demo showcases complete cloud storage workflows:

**Operations Performed:**
1. **Bucket Management**
   - Creates `test-bucket` if it doesn't exist
   - Handles bucket existence checks gracefully

2. **File Operations**
   - Writes a test file with sample content
   - Demonstrates proper error handling
   - Shows file metadata (size, creation time)

3. **Object Listing**
   - Lists all objects in the bucket
   - Displays object properties and timestamps

4. **File Reading**
   - Reads file content back from storage
   - Demonstrates complete round-trip operation

**Sample Output:**
```
Created bucket 'test-bucket'
Successfully wrote file 'test-file.txt' to bucket 'test-bucket'

Listing files in bucket:
- test-file.txt (size: 44 bytes, created: 2025-09-30T02:02:42Z)

Reading file 'test-file.txt' content:
Content: Hello, this is a test file for GCS emulator!
```

### Pub/Sub Operations (`pubsubemulator.go`)
The Pub/Sub demo demonstrates real-time messaging:

**Operations Performed:**
1. **Infrastructure Setup**
   - Creates topic `test-topic`
   - Creates subscription `test-subscription`
   - Handles existing resources gracefully

2. **Message Publishing**
   - Publishes multiple test messages
   - Adds timestamps as message attributes
   - Demonstrates error handling

3. **Message Subscription**
   - Starts background subscriber
   - Processes messages in real-time
   - Shows message acknowledgment

4. **Concurrent Operations**
   - Publisher and subscriber run simultaneously
   - Demonstrates typical pub/sub patterns

**Sample Output:**
```
Created topic: test-topic
Created subscription: test-subscription
Starting message subscriber...
Published: Message 1
Received message: Message 1: Hello from Pub/Sub emulator!
  Attributes: map[timestamp:2025-09-30T07:32:45+05:30]
Published: Message 2
Received message: Message 2: This is message number 2
  Attributes: map[timestamp:2025-09-30T07:32:46+05:30]
```

## 🌐 Service Endpoints

| Service | Local URL | Purpose |
|---------|-----------|---------|
| **Fake GCS Server** | http://localhost:4443 | Storage operations |
| **Pub/Sub Emulator** | http://localhost:8085 | Messaging operations |

## 🔧 Configuration Details

### Docker Services

#### Fake GCS Server
- **Image**: `fsouza/fake-gcs-server`
- **Port**: 4443 (HTTP)
- **Features**: Complete GCS API compatibility
- **Storage**: Persistent Docker volume

#### Pub/Sub Emulator  
- **Image**: `messagebird/gcloud-pubsub-emulator`
- **Port**: 8085 (mapped from internal 8681)
- **Features**: Full Pub/Sub API with Java runtime
- **Project ID**: `test-project`

### Environment Variables
```bash
# Required for GCS operations
STORAGE_EMULATOR_HOST=localhost:4443

# Required for Pub/Sub operations  
PUBSUB_EMULATOR_HOST=localhost:8085
```

### Go Dependencies
- `cloud.google.com/go/storage` - GCS client library
- `cloud.google.com/go/pubsub` - Pub/Sub client library
- `google.golang.org/api` - Google API client

## 🛠 Available Commands

### Makefile Commands
```bash
make start      # Start emulators
make stop       # Stop emulators
make demo       # Full automated demo
make gcs        # Run GCS demo only
make pubsub     # Run Pub/Sub demo only
make logs       # View emulator logs
make clean      # Stop and remove everything
make deps       # Install Go dependencies
```

### Docker Commands
```bash
# Container management
docker-compose up -d           # Start services
docker-compose down            # Stop services
docker-compose down -v         # Stop and remove volumes
docker-compose logs -f         # Follow logs
docker-compose ps              # List containers
docker-compose restart [service]  # Restart specific service

# Individual services
docker-compose up fake-gcs-server
docker-compose up pubsub-emulator
```

### Go Commands
```bash
# Run programs
go run gcs.go                  # GCS demo
go run pubsubemulator.go       # Pub/Sub demo

# Build executables
go build -o gcs gcs.go
go build -o pubsub pubsubemulator.go

# Dependency management
go mod tidy                    # Update dependencies
go mod download                # Download modules
```


### Log Analysis
```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs fake-gcs-server
docker-compose logs pubsub-emulator

# Follow logs with timestamps
docker-compose logs -f -t
```

## 🧪 Testing & Verification

### Code Verification
```bash
# Verify Go code compiles without running emulators
./verify-code.sh

# Expected output:
# 🔍 Verifying Go code compilation...
# Checking gcs.go...
# ✅ gcs.go compiles successfully
# Checking pubsubemulator.go...
# ✅ pubsubemulator.go compiles successfully
# ✨ Code verification complete!
```

## 🔧 Troubleshooting

### Common Issues and Solutions

#### 1. Go Version Mismatch Error
```
compile: version "go1.25.0" does not match go tool version "go1.25.1"
```

**Solution**: Ensure Go compiler and tools are the same version
```bash
# Check current Go version
go version
go env GOROOT

# If using goenv, update to latest version
goenv install 1.25.1
goenv global 1.25.1

# Reload shell and clear cache
source ~/.zshrc
go clean -cache
go clean -modcache
```

#### 2. Port Already Allocated Error
```
Bind for 0.0.0.0:4443 failed: port is already allocated
```

**Solution**: Stop existing containers and check for port conflicts
```bash
# Stop any running containers
docker-compose down

# Check what's using the port
lsof -i :4443
lsof -i :8085

# If needed, change ports in docker-compose.yml
```

#### 3. Docker Platform Warning (Apple Silicon Macs)
```
The requested image's platform (linux/amd64) does not match the detected host platform (linux/arm64/v8)
```

**Solution**: This is a warning only - emulators will still work. To suppress:
```bash
# Add platform specification to docker-compose.yml
services:
  pubsub-emulator:
    platform: linux/amd64  # Add this line
    image: messagebird/gcloud-pubsub-emulator
```

#### 4. Docker Not Running
```
❌ Docker is not running. Please start Docker first.
```

**Solution**: Start Docker Desktop application before running the demo

#### 5. Connection Refused Errors
**Solution**: Wait longer for services to start, or check container logs
```bash
# Check container status
docker-compose ps

# View logs for troubleshooting
docker-compose logs fake-gcs-server
docker-compose logs pubsub-emulator
```

### Environment Variables Check
```bash
# Verify environment variables are set correctly
echo $STORAGE_EMULATOR_HOST     # Should be: localhost:4443
echo $PUBSUB_EMULATOR_HOST      # Should be: localhost:8085

# Set manually if needed
export STORAGE_EMULATOR_HOST=localhost:4443
export PUBSUB_EMULATOR_HOST=localhost:8085
```

## 📚 Additional Resources

### GCP Go Client Documentation
- [Google Cloud Storage Go Client](https://pkg.go.dev/cloud.google.com/go/storage)
- [Google Cloud Pub/Sub Go Client](https://pkg.go.dev/cloud.google.com/go/pubsub)

### Emulator Documentation  
- [Fake GCS Server](https://github.com/fsouza/fake-gcs-server)
- [Cloud SDK Pub/Sub Emulator](https://cloud.google.com/pubsub/docs/emulator)



# Ekansh Thakur - VapusData Assignment

This repository contains a gRPC-based Go application that utilizes the Gin framework for APIs, JWT for authentication, and In-Memory storage. The project integrates various third-party libraries and tools to create a robust backend system.

## Features
- **gRPC Server**: High-performance communication between services.
- **REST API**: Built with the Gin framework for seamless interaction.
- **JWT Authentication**: Secure user authentication and authorization.
- **In-Memory Storage**: Fast data access and management.
- **UUID Support**: Unique identifiers for entities.
- **Dotenv Integration**: Environment variable management using `godotenv`.
- **Protobuf Integration**: Protocol buffers for defining gRPC services.

## Prerequisites
- [Go (1.22.0)](https://go.dev/)
- Protocol Buffers Compiler (`protoc`)
- `git` for version control

## Setup Instructions

### 1. Install Protocol Buffers Compiler
#### For Ubuntu/Debian:
```bash
apt install -y protobuf-compiler
```

#### For macOS:
```bash
brew install protobuf
```

#### Install Go protobuf plugins:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 2. Clone the Repository
```bash
git clone https://github.com/ekanshthakur15/VapusData_Assignment.git
cd VapusData_Assignment
```

### 3. Install Dependencies
```bash
go mod download
```

### 4. Configure Environment Variables
Create a `.env` file in the `./server/` directory and populate it with the required configuration:
```env
JWT_SECRET=<your-secret-key>
```

### 5. Generate Protobuf Files
Use `protoc` to generate the gRPC service definitions:
```bash
protoc --go_out=. --go-grpc_out=. ./proto/*.proto
```

### 6. Run the Application
Start the gRPC server:
```bash
go run ./server/server.go
```

## Testing
1. Start the RESTful client in a new terminal:
```bash
go run ./client/client.go
```

The client runs on `localhost:8000` and provides the following API endpoints:

### Authentication Endpoints
- `POST /signup`: Register a new user
- `POST /login`: Authenticate user and receive JWT token

### Book Management Endpoints
All book endpoints require Bearer token authentication in the header.

#### GET Endpoints
- `GET /getBook?id=<book_id>`: Retrieve a specific book
- `GET /listBooks`: Get all books

#### POST Endpoint
- `POST /createBook`: Create a new book
```json
{
  "name": "The Alchemist",
  "author": "Paulo Coelho",
  "published_year": 2001,
  "price": 299.99
}
```

#### PUT/DELETE Endpoints
- `PUT /updateBook?id=<book_id>`: Update a book
- `DELETE /deleteBook?id=<book_id>`: Delete a book

## Project Structure
```plaintext
├── proto/
│   └── *.proto         # Protocol Buffer Definitions
├── server/
│   └── .env            # Environment Variables
│   └── server.go       # gRPC Server Implementation
├── go.mod              # Module Dependencies
├── go.sum              # Dependency Checksums
├── README.md           # Documentation
```

## Dependencies
- **Gin Framework**: `github.com/gin-gonic/gin`
- **JWT**: `github.com/golang-jwt/jwt/v5`
- **UUID**: `github.com/google/uuid`
- **Environment Variables**: `github.com/joho/godotenv`
- **Protobuf Compiler**: `google.golang.org/protobuf`

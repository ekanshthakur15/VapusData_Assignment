# VapusData
VapusData is a gRPC-based Go application that utilizes the Gin framework for APIs, JWT for authentication, and PostgreSQL as the database. The project integrates various third-party libraries and tools to create a robust backend system.

## Features
- **gRPC Server**: High-performance communication between services.
- **REST API**: Built with the Gin framework for seamless interaction.
- **JWT Authentication**: Secure user authentication and authorization.
- **PostgreSQL Database**: Robust data storage and querying.
- **UUID Support**: Unique identifiers for entities.
- **Dotenv Integration**: Environment variable management using `godotenv`.
- **Protobuf Integration**: Protocol buffers for defining gRPC services.

## Prerequisites
- [Go (1.22.0)](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)
- Protocol Buffers Compiler (`protoc`)
- `git` for version control

## Setup Instructions

### 1. Clone the Repository
```bash
git clone https://github.com/ekanshthakur15/vapusdata.git
cd vapusdata
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Configure Environment Variables
Create a `.env` file in the root directory and populate it with the required configuration:
```env
DATABASE_URL=postgres://<username>:<password>@<host>:<port>/<dbname>?sslmode=disable
JWT_SECRET=<your-secret-key>
```

### 4. Generate Protobuf Files
Use `protoc` to generate the gRPC service definitions:
```bash
protoc --go_out=. --go-grpc_out=. ./proto/*.proto
```

### 5. Run the Application
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
│   └── *.proto          # Protocol Buffer Definitions
├── server/
│   └── server.go        # gRPC Server Implementation
├── .env                 # Environment Variables
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

## Contributing
Feel free to submit issues and pull requests. For major changes, please open an issue first to discuss what you would like to change.

## License
This project is licensed under the MIT License.
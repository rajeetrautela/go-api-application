🛡️ Go JWT API Server
A secure and modular RESTful API server built with Go, featuring JWT-based authentication, role-based access control, background workers, file uploads, and gRPC integration.

🚀 Features
  🔐 JWT Authentication with Access & Refresh Tokens
  👥 Role-based Authorization (admin, user)
  📦 CRUD APIs for Items
  🧵 Background Worker System
  📅 Cron Job Scheduler
  📁 File Upload with gRPC Integration
  🧰 Middleware for Logging, Recovery, and Timeout
  🗃️ PostgreSQL Database with Auto Migration

📁 Project Structure
go-jwt-api/
├── auth/                 # JWT generation, validation, and refresh token logic
├── config/               # Database configuration and connection
├── docker-compose.yml    # Docker Compose setup
├── Dockerfile            # Dockerfile for containerizing the app
├── curls.txt             # Sample curl commands for testing APIs
├── grpc_client.go        # gRPC client for file upload
├── grpc_server.go        # gRPC server implementation
├── internal/             # Internal utilities and helpers
├── main.go               # Application entry point
├── middleware/           # Custom middlewares (logging, recovery, timeout, JWT)
├── migrations/           # Database schema migrations
├── model/                # Data models (User, Item, etc.)
├── proto/                # Protocol buffer definitions for gRPC
├── README.md             # Project documentation
├── readme.txt            # Additional notes or legacy documentation
├── repository/           # Data access layer (optional abstraction)
├── scheduler/            # Cron job scheduler
├── uploads/              # Temporary file storage for uploads
├── worker/               # Background job worker logic
├── go.mod                # Go module definition
├── go.sum                # Go module checksums

🧪 API Endpoints

  🔐 Authentication
  POST /login – Login and receive JWT tokens
  POST /refresh – Refresh access token
  POST /logout – Invalidate refresh token

  📦 Items (Protected)
  GET /items – List all items (admin)
  GET /items/{id} – Get item by ID (admin, user)
  POST /items – Create new item (admin)
  PUT /items/{id} – Update item (admin)
  DELETE /items/{id} – Delete item (admin)

🧵 Workers
POST /startworker – Trigger background workers (admin, user)

📁 File Upload
GET / – Upload form
POST /upload – Upload file and send to gRPC server

🛠️ Setup & Run
    git clone https://github.com/rajeetrautela/go-api-application.git
    cd go-jwt-api
    go mod init go-api-application
    go mod tidy


  Run the Server
    go run *.go 

Server will start at http://localhost:8001

🧪 Sample Users

[
  { "username": "admin", "password": "admin123", "role": "admin" },
  { "username": "user", "password": "user123", "role": "user" }
]


📌 Notes
Ensure your PostgreSQL database is running and accessible.
gRPC server must be implemented and running for file uploads to work.
Cron jobs are started automatically in the background.
Also make sure to update SMTP config to make the emailing funtionality works.

📃 License
MIT License

# Chat-Go

A real-time chat application built with Go, featuring MongoDB storage, WebSocket support for live messaging, REST API endpoints, and a web interface.

## Overview

Chat-Go is a full-stack chat application that enables users to register, login, and engage in real-time conversations. The application uses MongoDB for persistent storage, WebSocket connections for instant messaging, and provides both a REST API and a web-based interface for interaction.

## Features

- **User Authentication**: Secure registration and login with JWT-based authentication
- **Real-Time Messaging**: WebSocket-powered instant chat between users
- **Conversation Management**: Automatic conversation creation between users
- **REST API**: Complete API for user management and chat operations
- **Web Interface**: HTML-based frontend for login, registration, and chat dashboard
- **Password Security**: BCrypt hashing with pepper configuration for secure password storage

## Installation

### Prerequisites

- Go 1.25 or higher
- MongoDB (local or remote)

### Clone and Install Dependencies

```bash
# Clone the repository
git clone <repository-url>
cd chat-go

# Download dependencies
go mod download

# Install dependencies
go mod tidy
```

### Environment Setup

Create a `.env` file in the root directory:

```env
# MongoDB Connection
MONGO_URI=mongodb://localhost:27017

# Application Port
PORT=8080

# JWT Secret Key (generate a secure random key)
SUPER_SECRET_KEY=your-super-secret-jwt-key-here

# Password Pepper (for additional password security)
PASSWORD_PEPPER=your-password-pepper-string
```

### Running the Application

**Using Makefile (Recommended):**

```bash
# Run without building
make run

# Build the binary
make build

# Build and run
make start
```

**Using Go Commands:**

```bash
# Run directly
go run ./cmd/main.go

# Build and run binary
go build -o ./bin/chat ./cmd/main.go
./bin/chat
```

## Usage

### Web Interface

Access the web interface at `http://localhost:8080`:

- **Login Page**: `/api/v1/login` (via web route)
- **Registration Page**: `/api/v1/register` (via web route)
- **Dashboard**: `/dashboard` (after authentication)

### REST API Endpoints

#### Authentication

**Register a new user:**

```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

**Login:**

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

Response includes a JWT token for authenticated requests.

#### User Management

**Get all users (Authentication Required):**

```bash
curl -X GET http://localhost:8080/api/v1/allUsers \
  -H "Authorization: Bearer <your-jwt-token>"
```

### WebSocket Connection

Connect to the WebSocket server for real-time chat:

```
ws://localhost:8080/ws/chat?user_id=<recipient-user-id>&token=<jwt-token>
```

**WebSocket JavaScript Example:**

```javascript
const token = 'your-jwt-token';
const recipientUserId = 'recipient-user-id';
const ws = new WebSocket(
  `ws://localhost:8080/ws/chat?user_id=${recipientUserId}&token=${token}`
);

ws.onopen = () => {
  console.log('Connected to chat server');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message);
};

ws.onclose = () => {
  console.log('Disconnected from chat server');
};
```

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `MONGO_URI` | MongoDB connection string | Yes | `mongodb://localhost:27017` |
| `PORT` | Application listening port | No | `8081` |
| `SUPER_SECRET_KEY` | JWT secret key for signing tokens | Yes | - |
| `PASSWORD_PEPPER` | Pepper string for password hashing | Yes | - |

## Project Structure

```
chat-go/
├── cmd/                    # Application entry points
│   ├── main.go            # Main application entry point
│   ├── db/                # Database initialization
│   └── server/            # HTTP server configuration
├── internal/              # Internal packages (application logic)
│   ├── controllers/       # HTTP request handlers
│   │   ├── authController.go
│   │   ├── userController.go
│   │   └── chatController/    # WebSocket chat handling
│   ├── middlewares/       # HTTP middleware (authentication, logging)
│   ├── models/            # Data models and structs
│   ├── repositories/      # Database operations
│   │   ├── user/
│   │   ├── conversations/
│   │   └── messages/
│   ├── routes/            # API route definitions
│   ├── services/          # Business logic services (JWT, password, WebSocket)
│   └── templates/         # HTML templates for web interface
├── static/                # Static assets
│   ├── css/               # Stylesheets
│   ├── images/            # Image files
│   └── script/            # JavaScript files
├── .env                   # Environment variables (not committed)
├── .gitignore
├── go.mod
├── go.sum
└── Makefile
```

### Directory Descriptions

- **cmd/**: Contains the main entry points and server initialization code
- **internal/controllers/**: Handles HTTP requests and business logic
- **internal/middlewares/**: Provides authentication and logging middleware
- **internal/models/**: Defines data structures for users, conversations, and messages
- **internal/repositories/**: Implements database CRUD operations
- **internal/routes/**: Configures API and WebSocket routes
- **internal/services/**: Contains reusable business logic (JWT, password hashing)
- **internal/templates/**: HTML templates for the web interface
- **static/**: Serves CSS, JavaScript, and image assets

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request


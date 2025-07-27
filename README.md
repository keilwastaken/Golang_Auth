# GoAuth - JWT Authentication Service in Go

A secure, production-ready authentication service built with Go, featuring JWT token management, MongoDB integration, and RESTful API design. This project demonstrates modern authentication patterns and Go best practices.

## 🚀 Features

- **JWT Authentication**: Secure access and refresh token implementation
- **User Management**: Registration, login, and logout functionality
- **Password Security**: BCrypt hashing with strong password validation
- **Token Refresh**: Automatic token renewal with refresh token rotation
- **MongoDB Integration**: Persistent storage for users and refresh tokens
- **RESTful API**: Clean, organized endpoint structure
- **Middleware Protection**: Route-level authentication enforcement
- **Input Validation**: Email format and password strength requirements
- **Clean Architecture**: Separation of concerns with controllers, services, and repositories

## 🛠️ Tech Stack

- **Go 1.22**: Modern Go with generics support
- **Gin Framework**: High-performance HTTP web framework
- **MongoDB**: NoSQL database for flexible data storage
- **JWT (golang-jwt/jwt/v5)**: Industry-standard token authentication
- **BCrypt**: Secure password hashing
- **Environment Variables**: Secure configuration management

## 📁 Project Structure

```
GoAuth/
├── Config/              # Application configuration
├── Controller/          # HTTP request handlers
├── Helpers/             # Utility functions (JWT handling)
├── Interfaces/          # Contract definitions
├── MiddleWare/          # Authentication middleware
├── Models/              # Data structures
├── Repository/          # Database access layer
├── Routes/              # API endpoint definitions
├── Server/              # HTTP and MongoDB server setup
├── Services/            # Business logic layer
└── main.go              # Application entry point
```

## 🔐 API Endpoints

### Authentication Routes (`/auth`)

- `POST /auth/register` - Create new user account
- `POST /auth/login` - Authenticate user credentials
- `POST /auth/refresh` - Refresh access token
- `POST /auth/logout` - Invalidate refresh token
- `GET /auth/validate` - Validate current token (protected)

### Task Routes (`/tasks`) - Protected

- `POST /tasks` - Create new task
- `GET /tasks` - List user tasks
- `PUT /tasks/:id` - Update task
- `DELETE /tasks/:id` - Delete task

## 🏗️ Architecture Highlights

### Clean Architecture Pattern
- **Controllers**: Handle HTTP requests/responses
- **Services**: Implement business logic
- **Repositories**: Manage data persistence
- **Interfaces**: Define contracts for dependency injection

### Security Features
- Password requirements: 8+ characters, uppercase, lowercase, number, special character
- Token expiration and rotation
- Secure token storage in MongoDB
- Protected routes with middleware

## 🚦 Getting Started

### Prerequisites
- Go 1.22 or higher
- MongoDB instance
- Environment variables configuration

### Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/goauth.git
cd goauth
```

2. Install dependencies
```bash
go mod download
```

3. Set up environment variables
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the application
```bash
go run main.go
```

## 📝 Example Usage

### Register a new user
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!"
  }'
```

## 💡 Learning Outcomes

This project helped me understand:
- Go's interface-based design patterns
- JWT token lifecycle management
- Clean architecture principles in Go
- MongoDB integration with Go drivers
- Middleware implementation in Gin
- Secure password handling and validation
- RESTful API design best practices

## 🔮 Future Enhancements

- [ ] OAuth2 integration (Google, GitHub)
- [ ] Rate limiting
- [ ] Email verification
- [ ] Password reset functionality
- [ ] User profile management
- [ ] API documentation with Swagger
- [ ] Docker containerization
- [ ] Unit and integration tests

## 📄 License

This project is open source and available under the MIT License.
# Taskify

A robust task management REST API built with Go, featuring comprehensive task organization capabilities.

## Features

- RESTful API endpoints for task management
- Swagger documentation for API reference
- Middleware for error handling
- Environment-based configuration
- Database integration
- Input validation

## Tech Stack

- Go (Golang)
- Gin Web Framework
- GORM (Object Relational Mapper)
- Swagger for API documentation
- Air for live reload during development

## Getting Started

1. Clone the repository
2. Copy `.env.development` to `.env` and configure your environment variables
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Run the application:
   ```bash
   go run main.go
   ```
   
For development with live reload:
```bash
air
```

## API Documentation

Once the server is running, you can access the Swagger documentation at:
```
http://localhost:8080/swagger/index.html
```

## Project Structure

```
taskify/
├── config/         # Configuration setup
├── controllers/    # Request handlers
├── docs/          # Swagger documentation
├── errors/        # Custom error definitions
├── middleware/    # HTTP middleware
├── models/        # Database models
├── routes/        # Route definitions
└── utils/         # Utility functions
```

## Developer

Built with ❤️ by [Eng Abdirahman Baabale](https://github.com/baabale)

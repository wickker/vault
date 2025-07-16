# Vault

A secure key-value data manager with encrypted storage, built with Go and PostgreSQL.

## Overview

Vault is a REST API service that allows users to securely store and manage sensitive information like passwords, API keys, and other confidential data. The application features:

- **Secure Encryption**: All sensitive data is encrypted using AES-256 encryption before storage
- **User Authentication**: Integration with Clerk for secure user authentication and authorization
- **Categorized Storage**: Organize your sensitive data into categories for better management
- **OpenAPI Specification**: Fully documented REST API with OpenAPI 3.0 specification
- **Database-First Approach**: Uses SQLC for type-safe database operations

## Features

- 🔐 **Encrypted Data Storage**: All sensitive values are encrypted before being stored in the database
- 👤 **User Authentication**: Secure authentication using Clerk
- 📁 **Categories**: Organize items into custom categories with colors
- 📝 **Items & Records**: Store items with multiple key-value records
- 🔍 **Search & Filter**: Search items by name and filter by category
- 🚀 **RESTful API**: Clean REST API with comprehensive OpenAPI documentation
- 📊 **Database Migrations**: Version-controlled database schema with PostgreSQL
- 🛡️ **Security Middleware**: CORS, request ID tracking, and authentication middleware

## Tech Stack

- **Backend**: Go 1.23.4 with Gin web framework
- **Database**: PostgreSQL with SQLC for type-safe queries
- **Authentication**: Clerk SDK for user management
- **Encryption**: AES-256 encryption for sensitive data
- **API Documentation**: OpenAPI 3.0 with code generation
- **Logging**: Structured logging with Zerolog

## Prerequisites

- Go 1.23.4 or later
- PostgreSQL database
- Clerk account for authentication

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
ENV=dev
CLERK_SECRET_KEY=your_clerk_secret_key
DATABASE_URL=your_db_url
FRONTEND_ORIGINS=http://localhost:5173
ENCRYPTION_KEY=your_encryption_key
```

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd vault
   ```

2. Install dependencies:
   ```bash
   go mod install
   ```

3. Set up your environment variables in `.env` file

4. Set up your PostgreSQL database and run migrations:
   ```bash
   # Apply the database schemas
   psql $DATABASE_URL -f db/schemas/v1.sql
   psql $DATABASE_URL -f db/schemas/v2.sql
   ```

## Development

1. `make server` to generate Gin server methods 
2. `make models` to generate controller (service) models
3. `sqlc generate` to generate db methods and models

After making changes to the OpenAPI specification (`openapi/openapi.yaml`) or database queries, run the appropriate command above to regenerate the code.

## Running the Application

1. Start the server:
   ```bash
   go run main.go
   ```

2. The API will be available at `http://localhost:9000`

3. Health check endpoint: `GET http://localhost:9000/`

## API Documentation

The API is documented using OpenAPI 3.0 specification. You can find the specification in `openapi/openapi.yaml`.

### Main Endpoints

- **Categories**: `/protected/categories`
  - Create, read, update, delete categories
  - Each category has a name and color

- **Items**: `/protected/items`
  - Create, read, update, delete items within categories
  - Support for search and filtering

- **Records**: `/protected/records`
  - Create, read, update, delete key-value records within items
  - All values are automatically encrypted

### Authentication

All protected endpoints require a valid Clerk JWT token passed in the `Authorization` header:

```
Authorization: Bearer <your-clerk-jwt-token>
```

## Project Structure

```
vault/
├── config/           # Configuration and environment handling
├── db/              # Database schemas, queries, and generated models
│   ├── queries/     # SQL queries for SQLC
│   ├── schemas/     # Database migration files
│   └── sqlc/        # Generated database models and methods
├── middleware/      # HTTP middleware (auth, CORS, request ID)
├── openapi/         # OpenAPI specification and generated code
├── services/        # Business logic and API handlers
├── utils/           # Utility functions (encryption, etc.)
├── main.go          # Application entry point
├── Makefile         # Build automation
└── sqlc.yaml        # SQLC configuration
```

## Security

- All sensitive data (record values) are encrypted using AES-256 encryption
- User authentication is handled by Clerk with JWT verification
- CORS middleware configured for secure cross-origin requests
- Request ID tracking for better debugging and monitoring

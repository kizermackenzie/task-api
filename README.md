# Task API

A production-ready REST API for task management with JWT authentication, built with Go and deployed on Google Cloud Platform using Terraform.

## 🚀 Features

- **JWT Authentication** - Secure user registration and login
- **Task Management** - CRUD operations with priorities and due dates
- **PostgreSQL Database** - Reliable data persistence with GORM
- **Input Validation** - Comprehensive request validation
- **Secure Passwords** - bcrypt hashing
- **Cloud Deployment** - Google Cloud Run with Cloud SQL
- **Infrastructure as Code** - Terraform configuration

## 📁 Project Structure

```
task-api/
├── main.go                 # Application entry point
├── go.mod                  # Go dependencies
├── Dockerfile              # Container configuration
├── docker-compose.yml      # Local development setup
│
├── database/               # Database connection and configuration
├── handlers/               # HTTP request handlers
├── middleware/             # Authentication middleware
├── models/                 # Database models (User, Task)
├── repositories/           # Data access layer
├── services/               # Business logic layer
├── utils/                  # Utility functions (JWT, password)
│
└── terraform/              # Infrastructure as Code
    ├── main.tf             # Provider and API configuration
    ├── database.tf         # Cloud SQL PostgreSQL
    ├── cloud_run.tf        # Cloud Run service
    ├── iam.tf             # Service accounts and permissions
    ├── secrets.tf         # Secret Manager configuration
    ├── variables.tf       # Input variables
    ├── outputs.tf         # Output values
    ├── terraform.tfvars   # Project configuration
    └── README.md          # Terraform documentation
```

## 🛠️ Development

### Local Development

1. **Start PostgreSQL and API**:
   ```bash
   docker-compose up
   ```

2. **Test endpoints**:
   ```bash
   # Register user
   curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"Password123","first_name":"Test","last_name":"User"}'
   
   # Create task
   curl -X POST http://localhost:8080/api/v1/tasks \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"title":"Complete project","priority":"high"}'
   ```

### Build Docker Image

```bash
docker build -t task-api .
```

## ☁️ Production Deployment

### Using Terraform (Recommended)

1. **Navigate to terraform directory**:
   ```bash
   cd terraform
   ```

2. **Initialize and deploy**:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

3. **Access your deployed API**:
   - Service URL will be provided in terraform output
   - Use the generated service account key for authentication

See `terraform/README.md` for detailed deployment instructions.

## 📚 API Documentation

### Authentication Endpoints

- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh JWT token
- `GET /api/v1/auth/profile` - Get user profile (authenticated)

### Task Endpoints

- `POST /api/v1/tasks` - Create task (authenticated)
- `GET /api/v1/tasks` - List user tasks (authenticated)
- `GET /api/v1/tasks/:id` - Get specific task (authenticated)
- `PUT /api/v1/tasks/:id` - Update task (authenticated)
- `DELETE /api/v1/tasks/:id` - Delete task (authenticated)
- `POST /api/v1/tasks/:id/complete` - Mark task complete (authenticated)

### Other Endpoints

- `GET /health` - Health check
- `GET /api/v1/admin/tasks` - Admin: List all tasks (authenticated)

## 🔐 Security Features

- JWT token authentication
- Password hashing with bcrypt
- Input validation and sanitization
- SQL injection protection with GORM
- Environment-based configuration
- Service account-based GCP authentication

## 🏗️ Architecture

- **Go 1.24** - Modern Go with latest features
- **Gin Framework** - Fast HTTP router
- **GORM** - Go ORM with PostgreSQL driver
- **Google Cloud Run** - Serverless container platform
- **Cloud SQL** - Managed PostgreSQL database
- **Secret Manager** - Secure credential storage
- **Terraform** - Infrastructure as Code

## 📄 License

This project is for educational and demonstration purposes.
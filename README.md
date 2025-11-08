# Task Management Backend

## Setup

### 1. Clone the repo

```bash
git clone https://github.com/hutamy/task-management-backend
cd task-management-backend
```

### 2. Set up .env

```
cp .env.example .env
# fill in DATABASE_URL, JWT_SECRET
```

### 3, Run

```bash
make run
```

## 💡 Project Structure

```
├── cmd/                  # Main application entrypoint
├── config/               # Configuration files and helpers
├── internal/             # Internal packages
├── middleware/           # Middleware packages
├── pkg/                  # External packages
├── .env.example
├── .gitignore
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## API Documentation

### Authentication

#### Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "username",
    "password": "password"
  }'
```

### Tasks

#### Get All Tasks
```bash
curl -X GET http://localhost:8080/api/tasks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Create Task
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Complete project documentation",
    "description": "Write comprehensive API documentation"
  }'
```

#### Update Task
```bash
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Updated task title",
    "description": "Updated task description",
    "status": "in progress"
  }'
```

#### Delete Task
```bash
curl -X DELETE http://localhost:8080/api/tasks/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

# Markdown Overview System

A distributed file processing system that generates AI-powered summaries of uploaded markdown files. The system uses a microservices architecture with Go backend, Python worker, Next.js frontend, and AWS services (S3, SQS) via LocalStack.

## Architecture Overview

The system follows an event-driven architecture with the following flow:

1. **Frontend (Next.js)** - User uploads a markdown file
2. **Backend (Go)** - Receives file, stores in S3, sends task to SQS queue
3. **Worker (Python)** - Polls SQS, processes file with OpenRouter LLM, stores summary in S3, sends completion message to response queue
4. **Backend (Go)** - Polls response queue, retrieves summary from S3, broadcasts via Server-Sent Events (SSE)
5. **Frontend (Next.js)** - Receives SSE notification and displays the summary

```
┌──────────────┐     upload      ┌──────────────┐    task msg    ┌──────────────┐
│   Frontend   │ ──────────────> │  Backend-Go  │ ────────────>  │  Task Queue  │
│   (Next.js)  │                 │   (Gin API)  │                │    (SQS)     │
└──────────────┘                 └──────────────┘                └──────────────┘
       ↑                                ↓                                 ↓
       │                          ┌──────────┐                    ┌──────────────┐
       │        SSE               │    S3    │                    │   Worker     │
       │      updates             │  Bucket  │ <────process────── │   (Python)   │
       │                          └──────────┘                    └──────────────┘
       │                                ↑                                 ↓
       │                                │        summary                  │
       │                          store summary                  completion msg
       │                                │                                 ↓
       │                          ┌──────────────┐                ┌──────────────┐
       └────────────────────────  │  Backend-Go  │ <────────────  │  Response    │
              display summary     │ (SSE Worker) │                │    Queue     │
                                  └──────────────┘                │    (SQS)     │
                                                                  └──────────────┘
```

## Technology Stack

- **Frontend**: Next.js 15, React 19, TailwindCSS, TypeScript
- **Backend**: Go 1.23, Gin framework, AWS SDK v2
- **Worker**: Python 3, Boto3, OpenRouter API
- **Database**: PostgreSQL (NeonDB)
- **Infrastructure**: LocalStack (S3, SQS), Docker
- **Authentication**: Session-based with cookies

## Prerequisites

- **Docker & Docker Compose** (for LocalStack)
- **Go 1.23+** (for backend)
- **Python 3.8+** (for worker)
- **Node.js 20+** and **npm** (for frontend)
- **PostgreSQL** database (NeonDB or local)
- **OpenRouter API Key** (for LLM summaries)

## Environment Setup

1. **Copy the example environment file:**
   ```bash
   cp .env.example .env
   ```

2. **Configure environment variables in `.env`:**
   ```env
   # Database (replace with your PostgreSQL connection string)
   DATABASE_URL=postgresql://username:password@host:5432/dbname?sslmode=require
   
   # LocalStack endpoints
   LOCALSTACK_ENDPOINT=http://localhost:4566
   AWS_DEFAULT_REGION=eu-central-1
   
   # S3 and SQS configuration
   S3_BUCKET_NAME=file-overview-system-bucket
   TASK_QUEUE_NAME=task-queue
   RESPONSE_QUEUE_NAME=response-queue
   
   # Worker - OpenRouter API
   OPENROUTER_API_KEY=your_openrouter_api_key_here
   OPENROUTER_URL=https://openrouter.ai/api/v1/chat/completions
   
   # Frontend
   NEXT_PUBLIC_API_BASE=http://localhost:8080
   ALLOWED_ORIGINS=http://localhost:3000
   
   # Backend port
   PORT=8080
   ```

3. **Backend-specific configuration:**
   
   The backend also uses `backend-go/backend-go.env`. Ensure it contains:
   ```env
   DATABASE_URL=postgresql://username:password@host:5432/dbname?sslmode=require
   LOCALSTACK_ENDPOINT=http://localhost:4566
   AWS_DEFAULT_REGION=eu-central-1
   S3_BUCKET_NAME=file-overview-system-bucket
   TASK_QUEUE_NAME=task-queue
   RESPONSE_QUEUE_NAME=response-queue
   ALLOWED_ORIGINS=http://localhost:3000
   PORT=8080
   ```

## Getting Started

### 1. Start LocalStack (Docker)

LocalStack provides local AWS services (S3, SQS) for development.

```bash
docker-compose up -d localstack
```

Verify LocalStack is running:
```bash
docker ps | grep localstack
```

LocalStack will be available at `http://localhost:4566`.

### 2. Run Backend (Go)

The backend handles API requests, file uploads, and SSE connections.

```bash
cd backend-go

# Install Go dependencies
go mod download

# Run the server
go run cmd/server/main.go
```

The backend will:
- Start on `http://localhost:8080`
- Connect to your PostgreSQL database
- Initialize S3 bucket and SQS queues in LocalStack
- Start the response queue worker for SSE broadcasting

**Key endpoints:**
- `POST /register` - User registration
- `POST /login` - User authentication
- `POST /upload` - File upload (authenticated)
- `GET /events` - SSE endpoint for real-time updates
- `POST /files` - List user's uploaded files

### 3. Run Worker (Python)

The worker processes files using OpenRouter's LLM API.

```bash
cd worker-python

# Install Python dependencies
pip install -r requirements.txt

# Run the worker
python app/main.py
```

The worker will:
- Poll the `task-queue` for new file processing tasks
- Download files from S3
- Send content to OpenRouter API for summarization
- Upload summaries back to S3
- Send completion messages to `response-queue`

**Environment variables required:**
- `OPENROUTER_API_KEY` - Your OpenRouter API key
- `LOCALSTACK_ENDPOINT` - LocalStack URL
- `TASK_QUEUE_URL` - SQS task queue URL
- `RESPONSE_QUEUE_URL` - SQS response queue URL

### 4. Run Frontend (Next.js)

The frontend provides the user interface for file uploads and summary viewing.

```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev
```

The frontend will:
- Start on `http://localhost:3000`
- Connect to backend API at `http://localhost:8080`
- Listen for SSE updates from the backend

**Available pages:**
- `/` - Landing page
- `/register` - User registration
- `/login` - User login
- `/dashboard` - Main dashboard for file upload and summary viewing

## System Flow Detailed

### File Upload Flow

1. User uploads a markdown file via the frontend dashboard
2. Frontend sends multipart form data to `POST /upload`
3. Backend (Go):
   - Authenticates the user via session middleware
   - Saves file to S3 at `users/{userId}/{filename}`
   - Creates a task message with `{bucket, key, userId}`
   - Sends task message to SQS `task-queue`
   - Returns success response to frontend

### Processing Flow

4. Worker (Python):
   - Continuously polls `task-queue` with long polling (10s wait)
   - Receives task message
   - Downloads file from S3
   - Prepends prompt: "Summarize following text in two sentences:"
   - Sends to OpenRouter API using `x-ai/grok-4-fast:free` model
   - Extracts summary from response
   - Uploads summary to S3 at `users/{userId}/{filename}_overview.txt`
   - Sends completion message to `response-queue`
   - Deletes processed message from `task-queue`

### Notification Flow

5. Backend Response Worker (Go):
   - Continuously polls `response-queue`
   - Receives completion message
   - Downloads summary from S3
   - Broadcasts summary via SSE to all connected clients
   - Includes `userId` so frontend can filter relevant updates

6. Frontend:
   - Maintains persistent SSE connection to `GET /events`
   - Receives summary updates
   - Displays summary in the dashboard UI

## Development Tips

### Running All Services

You can run all services simultaneously in separate terminals:

```bash
# Terminal 1 - LocalStack
docker-compose up localstack

# Terminal 2 - Backend
cd backend-go && go run cmd/server/main.go

# Terminal 3 - Worker
cd worker-python && python app/main.py

# Terminal 4 - Frontend
cd frontend && npm run dev
```

### Testing the System

1. Navigate to `http://localhost:3000`
2. Register a new account
3. Login with your credentials
4. Upload a markdown file (`.md` or `.txt`)
5. Watch as the summary appears in real-time via SSE

### Viewing LocalStack Resources

Check S3 buckets:
```bash
aws --endpoint-url=http://localhost:4566 s3 ls
aws --endpoint-url=http://localhost:4566 s3 ls s3://file-overview-system-bucket
```

Check SQS queues:
```bash
aws --endpoint-url=http://localhost:4566 sqs list-queues
```

## Troubleshooting

### Backend won't start

- **Database connection error**: Verify `DATABASE_URL` is correct and database is accessible
- **Port already in use**: Change `PORT` in environment variables
- **LocalStack not accessible**: Ensure Docker is running and LocalStack container is up

### Worker not processing files

- **No OpenRouter API key**: Set `OPENROUTER_API_KEY` in environment
- **Can't connect to LocalStack**: Ensure `LOCALSTACK_ENDPOINT` is correct
- **Queue URLs incorrect**: Verify `TASK_QUEUE_URL` and `RESPONSE_QUEUE_URL` match LocalStack format

### Frontend not receiving updates

- **SSE connection failed**: Check CORS settings in backend (`ALLOWED_ORIGINS`)
- **API base URL wrong**: Verify `NEXT_PUBLIC_API_BASE` points to backend
- **Not logged in**: Ensure you're authenticated before uploading files

### LocalStack issues

- **Services not initialized**: Wait a few seconds after starting LocalStack
- **Bucket/Queue errors**: Backend automatically creates resources on startup
- **Port conflicts**: Ensure port 4566 is not in use by another service

## Database Migrations

The backend uses SQL migrations in `backend-go/migrations/`:

- `000001_create_users_table` - Creates users table
- `000002_create_user_sessions` - Creates sessions table

Migrations should be run before starting the backend. The system uses `sqlc` for type-safe SQL queries.

## Project Structure

```
.
├── backend-go/              # Go backend service
│   ├── cmd/server/          # Main application entry point
│   ├── internal/
│   │   ├── clients/         # AWS SDK clients (S3, SQS)
│   │   ├── db/              # Database connection and SQLC queries
│   │   ├── events/          # SSE broadcaster implementation
│   │   ├── handlers/        # HTTP route handlers
│   │   ├── middleware/      # Auth and session middleware
│   │   ├── router/          # Route configuration
│   │   └── worker/          # Response queue worker
│   └── migrations/          # Database migrations
├── worker-python/           # Python processing worker
│   └── app/main.py          # Worker logic
├── frontend/                # Next.js frontend
│   └── app/                 # App router pages
├── infra/docker/            # Docker configurations
└── docker-compose.yml       # Docker Compose for LocalStack
```

## API Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | Health check | No |
| POST | `/register` | User registration | No |
| POST | `/login` | User login | No |
| GET | `/events` | SSE event stream | No |
| POST | `/upload` | Upload file | Yes |
| POST | `/files` | List user files | Yes |

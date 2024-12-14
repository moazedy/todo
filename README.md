# Todo Service

## Overview
This project implements a Todo Service using PostgreSQL, S3, and SQS, following Hexagonal Architecture principles. The service supports creating and managing TodoItems, uploading files to S3, and sending messages to an SQS queue.

## Features
1. **TodoItem Management**
   - Create TodoItems with fields `id`, `description`, `dueDate`, and an optional `fileName` (referencing an uploaded file in S3).
   - Store TodoItems in a PostgreSQL database.
   - Send TodoItem data to an SQS queue.
   - Full CRUD for TodoItem.

2. **File Upload**
   - REST API endpoint `/file` to upload files to S3.
   - Validate file types and sizes.
   - Return a unique `fileName` with its extension as file name representing the file in S3.
   - Download the uploaded file using `/file/:fileName` endpoint and GET method. 

3. **Hexagonal Architecture**
   - Separation of domain logic from infrastructure.
   - Dependency injection for easy mocking of external services (S3, SQS) in tests.

4. **Dockerized Environment**
   - Docker Compose setup for PostgreSQL and LocalStack (mocking S3 and SQS).

5. **Testing and Benchmarking**
   - Unit tests with mocks for S3 and SQS.
   - Benchmarks for key operations (e.g., database insertion, file uploads, and SQS messaging).

6. **Generic repository**
   - Generic repository for increasing development speed.

7. **Service layer Transaction**
   - Completely abstracted core service transaction to be fully separated from database technology.

8. **Custom error wrapping**
   - Custom error generation and handling so the proper error content and http code be delivered to the end user.

9. **Monitoring**
   - Serving prometheous standard metrics for project.





---

## Prerequisites
- Docker
- Docker Compose
- Go 1.23 or newer
- Make

---

## Setup

### 1. Clone the Repository
```bash
git clone <http://github.com/moazedy/todo.git>
cd <todo>
```

### 2. Build and Run the Services
Run the project using Docker Compose:
```bash
make run-build
```
This starts the PostgreSQL database and LocalStack services for S3 and SQS.
The run-build command is for the first time and in the next times make run will do enough.

---

## Database Migrations
Migrations is fully taken by the gorm and its automatically applied when ever the project runs.

---

## API Endpoints

### 1. File Upload
**Endpoint:** `/file`  
**Method:** `POST`

**Description:** Upload a file to S3. Returns a `fileName` created from file id and its extension.

#### Request Example:
```bash
curl -X POST -F "file=@path/to/file" http://localhost:4853/file
```

#### Response Example:
```json
{
  "fileName": "uuid.ext"
}
```

### 2. Create TodoItem
**Endpoint:** `/todo/item`  
**Method:** `POST`

**Description:** Create a new TodoItem.

#### Request Example:
```json
{
  "description": "Buy groceries",
  "dueDate": "2024-12-31T23:59:59Z",
  "fileName": "fe32a790-0232-4ac0-9525-7439d763c155.pdf"
}
```

#### Response Example:
```json
{
  "id": "unique-id",
  "description": "Buy groceries",
  "dueDate": "2024-12-31T23:59:59Z",
  "fileName": "fe32a790-0232-4ac0-9525-7439d763c155.pdf"
}
```

---

## Testing

### Run Unit Tests
```bash
make test
```
This will execute tests for:
- File uploads to S3.
- TodoItem creation and persistence in PostgreSQL.
- SQS messaging.

### Run Benchmarks
```bash
make benchmark
```
This will benchmark:
- Inserting a TodoItem into PostgreSQL.
- Uploading files to S3.
- Sending messages to SQS.

---

## Project Structure
```
.
├── cmd            # Main application entry points
├── config         # Configuration files
├── internal
│   ├── domain     # Domain logic
│   ├── adapter    # Infrastructure implementations (e.g., S3, SQS, PostgreSQL)
│   └── port       # Interfaces for dependency injection
├── pkg			   # Used packages
├── docker-compose.yml
├── Makefile
├── DockerFile
└── README.md
```

---

## Notes
- **Mocking AWS Services:** LocalStack is used to mock S3 and SQS for local development and testing.
- **Dependency Injection:** All external services are abstracted behind interfaces for testability.

---

## Commands Reference
| Command           | Description                                         |
|-------------------|-----------------------------------------------------|
| `make run-build`  | Builds and starts the project using Docker Compose. |
| `make run`        | Starts the project using Docker Compose.            |
| `make test`       | Runs unit tests.                                    |
| `make benchmark`  | Runs benchmarks.                                    |

---

## License
This project is licensed under the MIT License.



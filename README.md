# REST API for user management made with Clean Architecture

Scalable and maintainable RESTful API built with Go and PostgreSQL following Clean Architecture principles. This project provides a solid foundation for managing user data with full **CRUD (Create, Read, Update, Delete)** operations, integrated testing, and containerization using Docker and Docker Compose.

---

- [Features](#features)
- [Architecture](#architecture)
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Running the Application](#running-the-application)
- [Usage](#usage)
  - [API Endpoints](#api-endpoints)
- [Testing](#testing)

---

## Features

- **Clean Architecture**: Promotes separation of concerns, making the codebase easy to maintain and scale.
- **Full CRUD Operations**: Create, retrieve, update, and delete users seamlessly.
- **Dockerized**: Easily deployable with Docker and Docker Compose.
- **Integrated Testing**: Unit tests ensure code reliability.
- **Extensibility**: Easily adaptable for other projects or additional features.
- **Simplicity**: Straightforward setup and intuitive API endpoints.

---

## Architecture

The project follows the principles of **Clean Architecture**, ensuring that the core business logic is independent of external frameworks, databases, and UI. This approach offers several advantages:

- **Maintainability**: Easier to update and refactor code.
- **Testability**: Core logic can be tested without external dependencies.
- **Flexibility**: Swap out technologies (e.g., databases) with minimal impact.

**Layers:**

1. **Entities**: Core business models (e.g., User).
2. **Use Cases**: Application-specific business rules.
3. **Interfaces (Controllers)**: Frameworks and drivers (e.g., HTTP handlers).
4. **Infrastructure**: External agents like databases and web frameworks.

---

## Technologies Used

- **Go**: The main programming language.
- **Gin**: HTTP web framework.
- **GORM**: ORM library for Go.
- **PostgreSQL**: Relational database.
- **Docker & Docker Compose**: Containerization tools.
- **Testify**: Testing toolkit for Go.

---

### Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/interimme/userapi.git
   cd userapi
   ```

2. **Update Environment Variables (Optional)**

   - The default environment variables are set in `docker-compose.yml`.
   - Modify them if necessary.

### Running the Application

**Build and Start the Containers**

```bash
docker-compose up --build
```

This command will:

- Build the Docker image for the API.
- Pull and start the PostgreSQL database container.
- Run migrations to set up the database schema.
- Start the API service on port `8080`.

**Verify the Application is Running**

Access `http://localhost:8080` in your browser or use `curl`:

```bash
curl http://localhost:8080
```

You should receive a `404 Not Found` response, indicating that the server is running.

---

## Usage

### API Endpoints

#### 1. Create a User

- **Method:** `POST`
- **URL:** `http://localhost:8080/users`

**Request Body:**

```json
{
  "firstname": "Alice",
  "lastname": "Smith",
  "email": "alice.smith@example.com",
  "age": 28
}
```

**Example using `curl`:**

```bash
curl -X POST http://localhost:8080/users \
-H 'Content-Type: application/json' \
-d '{
  "firstname": "Alice",
  "lastname": "Smith",
  "email": "alice.smith@example.com",
  "age": 28
}'
```

**Expected Response:**

- **Status Code:** `201 Created`
- **Body:** JSON representation of the created user.

#### 2. Get a User

- **Method:** `GET`
- **URL:** `http://localhost:8080/user/{id}`

**Example using `curl`:**

```bash
curl http://localhost:8080/user/{user-id}
```

**Expected Response:**

- **Status Code:** `200 OK`
- **Body:** JSON representation of the user.

#### 3. Update a User

- **Method:** `PATCH`
- **URL:** `http://localhost:8080/user/{id}`

**Request Body:**

```json
{
  "firstname": "Alice",
  "lastname": "Johnson",
  "email": "alice.johnson@example.com",
  "age": 29
}
```

**Example using `curl`:**

```bash
curl -X PATCH http://localhost:8080/user/{user-id} \
-H 'Content-Type: application/json' \
-d '{
  "firstname": "Alice",
  "lastname": "Johnson",
  "email": "alice.johnson@example.com",
  "age": 29
}'
```

**Expected Response:**

- **Status Code:** `200 OK`
- **Body:** JSON representation of the updated user.

#### 4. Delete a User

- **Method:** `DELETE`
- **URL:** `http://localhost:8080/user/{id}`

**Example using `curl`:**

```bash
curl -X DELETE http://localhost:8080/user/{user-id}
```

**Expected Response:**

- **Status Code:** `200 OK`
- **Body:**

  ```json
  {
    "message": "User deleted successfully"
  }
  ```

---

## Testing

**Run Unit Tests**

```bash
go test ./...
```

This command runs all the unit tests in the project, ensuring that the business logic and controllers function as expected.

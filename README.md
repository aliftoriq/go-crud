# go-crud RESTful API Documentation

This documentation provides an overview and instructions for running a RESTful API built with Golang. The API can be executed locally using a local PostgreSQL database server or by building a container image with Docker using the provided Dockerfile.

The go-crud API incorporates the following features and technologies:

- **JWT Refresh Tokens**: Authentication in the API is handled using JWT (JSON Web Tokens) and includes support for refresh tokens to maintain user sessions securely.

- **Middleware Cookie Authorization**: Middleware is used for cookie-based authorization, enhancing security and user authentication.

- **GORM PostgreSQL**: The API utilizes GORM, a popular Object-Relational Mapping (ORM) library for Go, to interact with a PostgreSQL database, ensuring efficient and reliable data storage and retrieval.

- **Minio for Storage**: Minio is employed for storage management, allowing users to store and retrieve files and data efficiently.

- **Redis for Cache Management**: Redis is utilized for caching to optimize API response times and enhance overall performance.

## Prerequisites

Before you begin, ensure that you have the following prerequisites installed and set up:

- Golang installed on your local machine.
- PostgreSQL database server running locally (if running locally).
- Docker installed (if running with Docker).
- Minio (if using Minio for image storage locally).
- Redis (if using Redis for caching locally).
- Docker Compose. Before using Docker Compose, ensure you have [Docker](https://docs.docker.com/get-docker/) installed on your system.
- `.env` file containing necessary environment variables (see Setup Instructions).

## Setup Instructions

1. Clone the repository from the [go-crud Repo](https://github.com/alifotoriq/go-crud) main branch.

   ```
   git clone -b dev github.com/aliftoriq/go-crud
   ```

2. Ensure that the PostgreSQL database, minio and redis server is running locally. If not, please install and set it up accordingly.

3. Install the required dependencies by running the following command in the project directory:

   ```bash
   go mod download
   ```

4. Create a `.env` file in the project directory and provide the necessary environment variables :

   ```
    PORT=4001

    # jwt secret key
    SECRET=your_jwt_secret_key

    DB_USER=your_database_user
    DB_PASSWORD=your_database_password
    DB_NAME=your_database_name
    DATABASE_PORT=5432

    # Database Url
    DATABASE_URL=postgres://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable

    # Minio ccredentials
    ACCESKEY=your_acces_key
    SECRETKEY=your_Secret_key
    BUCKETNAME=your_bucket_name

    # Redis
    REDIS_PASSWORD=your_redis_password
    REDIS_DB=your_redis_db

    # docker
    DATABASE_HOST=host.docker.internal
    DB="host=host.docker.internal user=dev_crud password=dev_password dbname=dev_go_crud port=5432 sslmode=disable"
    REDIS_ADDR=host.docker.internal:6379
    ENDPOINT=host.docker.internal:9000
   ```

5. Run the API locally using the following command:

   - using go - run locally

   ```
   go run main.go
   ```

   - docker-compose build

   ```
   docker-compose up --build -d web
   ```

   The API should now be running on `http://localhost:4001`.

## Routes

This is an overview of the available routes and endpoints for the RESTful API. The API requires authentication using a Firebase token, which is obtained from cookies. Users must log in before accessing any other route, except for the signup and login routes.

**Base Url:** `https://baseurl/` + endpoint

## User Routes

These routes are responsible for managing user-related operations such as user registration, login, get and delete user.

### Sign Up

- **Route**: `POST /signup`
- **Description**: Create a new user.
- **Headers**: none
- **JSON Request**:
  ```json
  {
    "Name": "User name",
    "Email": "user@gmail.com",
    "Password": "user_password"
  }
  ```
- **JSON Response**:
  - succes | **HTTP Status Code** : `200`
    ```json
    {
      "message": "Create User Succesfuly"
    }
    ```
  - email already exist | **HTTP Status Code** : `409`
    ```json
    {
      "error": "User with this email already exists"
    }
    ```

### Login

- **Route**: `POST /login`
- **Description**: login with email and password.
- **Headers**: none
- **JSON Request**:
  ```json
  {
    "Email": "user@gmail.com",
    "Password": "user_password"
  }
  ```
- **JSON Response**:
  - succes | **HTTP Status Code** : `200`
    ```json
    {
      "data": {
        "ID": 1,
        "name": "User Name",
        "email": "user@gmail.com",
        "password": "userpassword"
      },
      "message": "Logged in",
      "token": "jwt token"
    }
    ```
  - invalid email or password | **HTTP Status Code** : `400`
    ```json
    {
      "error": "Invalid Email or Password"
    }
    ```

### Get User

- **Route**: `GET /user`
- **Description**: Get user data.
- **Headers**: Required (JWT token obtained from login set cookies).
- **JSON Request**:
  ```json
  none
  ```
- **JSON Response**:
  - succes | **HTTP Status Code** : `200`
    ```json
    {
        "CreatedAt": "2023-09-25T03:34:19.917335Z",
        "UpdatedAt": "2023-09-25T03:34:19.917335Z",
        "ID": user id,
        "name": "user name",
        "email": "user@gmail.com",
        "password": "user password"
    }
    ```

### Delete User

- **Route**: `DELETE /user`
- **Description**: Delete user data.
- **Headers**: Required (JWT token obtained from login set cookies).
- **JSON Response**:
  ```json
  {
    "message": "User deleted successfully"
  }
  ```
  | **HTTP Status Code** : `200`

## Article Routes

These routes are responsible for managing article-related operations such as creating, retrieving, updating, and deleting articles.

### Create Article

- **Route**: `POST /articles`
- **Description**: Create a new article.
- **Headers**: Required (JWT token obtained from login set cookies).
- **JSON Request**:
  ```json
  {
    "email": "aliftoriq52@gmail.com",
    "title": "Postmant",
    "content": "Postman merupakan tool untuk menguji API"
  }
  ```
- **JSON Response**:
  ```json
  {
    "message": "Article Created Succesfuly"
  }
  ```
  | **HTTP Status Code** : `200`

### Get All Article

- **Route**: `GET /articles`
- **Description**: Get all article data from database/from cache.
- **Headers**: Required (JWT token obtained from login set cookies).
- **JSON Response**:
  - database | **HTTP Status Code** : `200`
    ```json
    {
      "data": [
        {
          "CreatedAt": "2023-09-25T03:34:40.53053Z",
          "UpdatedAt": "2023-09-25T03:34:40.53053Z",
          "DeletedAt": null,
          "ID": 1,
          "email": "user1@example.com",
          "title": "Sample Article 1",
          "content": "This is the content of sample article 1.",
          "created_at": "2023-09-25T03:34:40.53053Z",
          "updated_at": "2023-09-25T03:34:40.53053Z",
          "deleted_at": null
        },
        {
          "CreatedAt": "2023-09-25T04:33:47.241906Z",
          "UpdatedAt": "2023-09-25T04:33:47.241906Z",
          "DeletedAt": null,
          "ID": 2,
          "email": "user2@example.com",
          "title": "Sample Article 2",
          "content": "This is the content of sample article 2.",
          "created_at": "2023-09-25T04:33:47.241906Z",
          "updated_at": "2023-09-25T04:33:47.241906Z",
          "deleted_at": null
        },
        {
          "CreatedAt": "2023-09-25T04:33:48.144175Z",
          "UpdatedAt": "2023-09-25T04:33:48.144175Z",
          "DeletedAt": null,
          "ID": 3,
          "email": "user3@example.com",
          "title": "Sample Article 3",
          "content": "This is the content of sample article 3.",
          "created_at": "2023-09-25T04:33:48.144175Z",
          "updated_at": "2023-09-25T04:33:48.144175Z",
          "deleted_at": null
        }
      ],
      "message": "Get Articles Successfully (from database)"
    }
    ```
    - redis cache | **HTTP Status Code** : `200`
    ```json
    {
    "data": [
        {
            ...
        }
    ],
    "message": "Get Articles Successfully (from cache)"
    }
    ```

### Get Article by ID

- **Route**: `GET /articles/:id`
- **Description**: Get article data by ID from database/from cache.
- **Headers**: Required (JWT token obtained from login set cookies).
- **JSON Response**:
  - database | **HTTP Status Code** : `200`
    ```json
    {
      "data": {
        "CreatedAt": "2023-09-25T03:34:40.53053Z",
        "UpdatedAt": "2023-09-25T03:34:40.53053Z",
        "DeletedAt": null,
        "ID": 1,
        "email": "user1@example.com",
        "title": "Sample Article 1",
        "content": "This is the content of sample article 1.",
        "created_at": "2023-09-25T03:34:40.53053Z",
        "updated_at": "2023-09-25T03:34:40.53053Z",
        "deleted_at": null
      },
      "message": "Get Article by ID Successfully (from database)"
    }
    ```
  - redis cache | **HTTP Status Code** : `200`
    ```json
    {
    "data":{
            ...
    },
    "message": "Get Article by ID Successfully (from cache)"
    }
    ```

### Update Article

- **Route**: `PUT /articles/:id`
- **Description**: Update article datda by ID.
- **Headers**: Required (JWT token obtained from login set cookies).
- **JSON Request**:
  ```json
  {
    "email": "user@gmail.com",
    "title": "update example",
    "content": "This is the content of sample update"
  }
  ```
- **JSON Response**:
  - succes | **HTTP Status Code** : `200`
    ```json
    {
      "message": "Article updated successfully"
    }
    ```

### Delete Article

- **Route**: `DELETE /articles/:id`
- **Description**: Delete article datda by ID.
- **Headers**: Required (JWT token obtained from login set cookies).
- **JSON Response**:
  - succes | **HTTP Status Code** : `200`
    ```json
    {
      "message": "Article deleted successfully"
    }
    ```

## Bucket Routes Documentation

These routes are responsible for managing operations related to object storage (bucket).

### Upload Image to Bucket

- **Route**: `POST /upload`
- **Description**: Upload an image to the bucket.
- **Headers**: Required (JWT token obtained from login set cookies).
- **Request Multipart Form**:
  - Field `image` (File): The image to be uploaded.
- **JSON Response**:
  - Success
    ```json
    {
      "message": "Image uploaded successfully",
      "fileName": "object_name.jpg"
    }
    ```
  - Failure
    ```json
    {
      "error": "Error message"
    }
    ```

### Get Image from Bucket

- **Route**: `GET /image/:id`
- **Description**: Download an image from the bucket based on its ID.
- **Headers**: Required (JWT token obtained from login set cookies).
- **Response**: The image in the format corresponding to the original content type (JPEG in this example).

### Delete Image from Bucket

- **Route**: `DELETE /image/:id`
- **Description**: Delete an image from the bucket based on its ID.
- **Headers**: none
- **JSON Response**:

  - Success
    ```json
    {
      "message": "Image deleted successfully"
    }
    ```
  - Failure

    ```json
    {
      "error": "Error message"
    }
    ```

    <br> <br/>

# Unit Testing Documentation

This documentation provides an overview of unit testing strategies and tools used in your project. Unit testing is essential to ensure that individual components of your application work correctly in isolation. In this project, we use Postman for integration testing to verify the entire architecture's functionality, and we utilize the testify and mock libraries to test controllers without establishing connections to the database, Redis, or Minio.

## Testing Approach

### Unit Testing

Unit testing focuses on testing individual components or units of your code in isolation. In your project, unit tests are primarily applied to the controllers. These tests ensure that your controllers handle requests and produce responses correctly without interacting with external services like the database, Redis, or Minio. Instead, we use mocking to simulate these external services.

### Integration Testing (Postman)

Integration testing verifies that the different parts of your application work together as expected. In this project, we use Postman for integration testing. Postman allows us to send HTTP requests to your API endpoints and check the responses. These tests ensure that your entire architecture functions correctly.

## Testing Tools

### Testify

[Testify](https://pkg.go.dev/github.com/stretchr/testify) is a popular testing framework for the Go programming language. It provides various assertion functions and tools for writing clean and efficient unit tests. We use Testify to write and run unit tests for the controllers.

### Mock

[Mock](https://pkg.go.dev/github.com/golang/mock) is a mocking framework for Go. It allows us to create mock implementations of interfaces, which helps us simulate the behavior of external dependencies (such as the database, Redis, or Minio) during unit tests. By using Mock, we can isolate the controllers from the actual services.

### Postman

[Postman](https://www.postman.com/) is a powerful tool for testing APIs. It allows us to create collections of API requests and run them as part of our integration testing process. Postman provides features for defining test scripts, setting up environments, and generating reports to ensure the correctness of your API endpoints.

## Writing Unit Tests

In your project, unit tests are written for controllers using Testify and Mock. These tests focus on verifying that the controllers handle requests and produce responses correctly. They do not interact with the actual database, Redis, or Minio services. Instead, we use mocks to simulate the behavior of these services.

## Integration Testing with Postman

Integration tests are performed using Postman. You can create collections of API requests in Postman to test the entire architecture, including interactions with databases, Redis, and Minio. These tests ensure that your API endpoints function correctly in a real-world scenario.

## Running Tests

To run unit tests, use the following command:

```bash
go test ./controllers -cover
```

# Code Editor Application
This README provides an overview of the Code Editor application, its functionalities, and APIs.

**Overview**
    The Code Editor application allows users to run code submissions, manage submissions, and perform user-related operations. It utilizes WebSocket connections for real-time communication and integrates with Redis and message queues for asynchronous processing of code submissions. Additionally

**Functionality**
  /run API: Establishes WebSocket connection between client and server to handle code execution requests.
  
  /submission API: Allows users to submit code for execution. Requests are processed asynchronously using Redis and message queues.
  
  /user API: Provides CRUD operations for managing user accounts.

**Submission API Example**
Request

    curl --location 'http://localhost:8080/submission' \
    --header 'Content-Type: application/json' \
    --data '{
        "input_file": "print(\'Hello, World!\')",
        "created_by": "123e4567-e89b-12d3-a456-426614174000",
        "lang": "python"
    }'
    Path: /submission
    Request Body: JSON
    Method: POST
    Response Type: JSON

**User API Example**
Sample Request
    {
        "user_name": "example_user",
        "email": "user@example.com"
    }
    Path: /user
    Request Body: JSON
    Method: POST
    Response Type: JSON

**Database Schema**
Submission Table
CREATE TABLE IF NOT EXISTS submissions (
    id              UUID PRIMARY KEY,
    input_file_s3_key  VARCHAR(255),
    created_by       UUID,
    created_at       TIMESTAMP,
    run_time         VARCHAR(50),
    lang            VARCHAR(50),
    status          VARCHAR(50),
    output_file      VARCHAR(255)
);

User Table
    CREATE TABLE IF NOT EXISTS users (
        id       UUID PRIMARY KEY,
        user_name VARCHAR(50),
        email    VARCHAR(100)
    );

Language Management
    Currently, the application uses LanguageEnum for language management. In the future, a separate table (language) can be utilized for better organization and management of languages.

Language Table
    CREATE TABLE IF NOT EXISTS languages (
        id             SERIAL PRIMARY KEY,
        language_name  VARCHAR(50),
        docker_image   VARCHAR(255),
        run_command    VARCHAR(255)
    );
    
*Additional Information*
*GORM is used for interacting with the database tables.
*S3 is utilized for persistently storing input files, and the S3 reference   key is stored in the database.

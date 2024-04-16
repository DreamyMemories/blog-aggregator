# Blog Aggregator

## Introduction
This project is a blog aggregator that fetches and consolidates blog posts from various sources, allowing users to follow different feeds and stay updated with the latest blog posts in their areas of interest.

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Features](#features)
- [Dependencies](#dependencies)
- [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Contributors](#contributors)
- [License](#license)

## Installation
To set up the project, follow these steps:
1. Clone the repository.
2. Install required dependencies using `go get`.
3. Make sure to have goose installed to manage database migrations.
4. Set up the PostgreSQL database using the provided SQL files.
5. Configure environment variables for database connection and other settings following `.env.example`.

## Usage
Run the project by executing:
```bash
go build; ./blog-aggregator.exe
```

## Features
- User authentication
- RSS feed parsing and storing
- Blog scraping functionality
- RESTful API for fetching blog posts

## Dependencies
- PostgreSQL
- Go standard library
- Goose for database migrations and `sqlc` for generating database code

## Configuration
The project uses environment variables for configuration. The following variables are required:
- `PORT`: The port on which the server will run.
- `CONNECTION_STRING`: The connection string for the PostgreSQL database

## API Documentation

### User Management
- **Endpoint**: `/v1/users`
- **Method**: POST, GET
- **Description**: Manages user creation and retrieval.
- **Requires Authentication**: Yes

### Feed Management
- **Endpoint**: `/v1/feeds`
- **Method**: POST
- **Description**: Manages creation of feeds.
- **Requires Authentication**: Yes

### Retrieve All Feeds
- **Endpoint**: `/v1/allfeeds`
- **Method**: GET
- **Description**: Retrieves all feeds.
- **Requires Authentication**: No

### Feed Follows Management
- **Endpoint**: `/v1/feed_follows`
- **Method**: POST
- **Description**: Manages creating feed follows.
- **Requires Authentication**: Yes

### Delete Feed Follow
- **Endpoint**: `/v1/feed_follows/{feedFollowID}`
- **Method**: DELETE
- **Description**: Deletes a specific feed follow.
- **Requires Authentication**: Yes

### Get Posts by User with Limit
- **Endpoint**: `/v1/posts/{limit}`
- **Method**: GET
- **Description**: Retrieves posts by a user with an optional limit.
- **Requires Authentication**: Yes

## Notes
- All endpoints that modify data require authentication.
- Data responses are in JSON format.

## Usage
When sending a request to the API, make sure to include the `Authorization` header with the value `Authorization <token>` where `<token>` is the token received after creation of user.

## Database Schema
The database schema consists of the following tables:
- `users`: Stores user information.
- `feeds`: Stores feed information.
- `feed_follows`: Stores feed follow information.
- `posts`: Stores post information.

## Examples
1. Create a user:
```bash
curl -X POST http://localhost:8080/v1/users -d '{"name": "test"}'
```

2. Create a feed (Note: A follow is automatically created based on the user that created the feed):
```bash
curl -X POST http://localhost:8080/v1/feeds -d '{"name": "Example blog", "url": "https://example.com/feed"}' -H "Authorization <token>"
```

3. Follow a feed:
```bash
curl -X POST http://localhost:8080/v1/feed_follows -d '{"feed_id": 1}' -H "Authorization <token>"
```

4. Retrieve posts:
```bash
curl -X GET http://localhost:8080/v1/posts/10 -H "Authorization <token>"
```
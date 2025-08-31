# gologhttpjson

A lightweight HTTP server for logging JSON requests with configurable features.

## Description

This application accepts POST requests with JSON payloads and logs them in a structured format. It's designed to be simple and focused on logging JSON data with optional header and metadata logging.

## Features

- **JSON Body Logging**: Captures and logs the JSON payload in the `body` field
- **Path Logging**: Records the HTTP path in the `path` field
- **Optional Header Logging**: Enable with `LOG_HEADERS` environment variable
- **Metadata Logging**: Log environment variables prefixed with `METADATA_` under the `metadata` field
- **Request Validation**: Only accepts POST requests with valid JSON payloads

## Environment Variables

- `PORT`: Server port (default: 8080)
- `LOG_HEADERS`: Set to any non-empty value to enable HTTP header logging
- `METADATA_*`: Any environment variable prefixed with `METADATA_` will be logged under metadata

## Usage

### Basic Usage
```bash
go run main.go
```

### With Header Logging
```bash
LOG_HEADERS=true go run main.go
```

### With Metadata
```bash
METADATA_VERSION=1.0.0 METADATA_ENVIRONMENT=production LOG_HEADERS=true go run main.go
```

## API

### POST /any-path
Accepts JSON payload and logs it.

**Request:**
- Method: POST
- Content-Type: application/json
- Body: Any valid JSON

**Response:**
- Status: 200 OK
- Body: "OK\n"

**Error Responses:**
- 400 Bad Request: For non-POST requests or invalid JSON
- 500 Internal Server Error: For server-side errors

## Log Format

The application logs entries in JSON format with the following structure:

```json
{
  "body": {...},           // The original JSON payload
  "headers": {...},        // HTTP headers (if LOG_HEADERS is set)
  "metadata": {...},       // Environment variables with METADATA_ prefix
  "path": "/example/path"  // The requested HTTP path
}
```

## Examples

### Simple JSON Request
```bash
curl -X POST http://localhost:8080/test \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, World!"}'
```

**Logged output:**
```json
{"body":{"message":"Hello, World!"},"path":"/test"}
```

### With Headers and Metadata
```bash
# Start server with headers and metadata
METADATA_SERVICE=api METADATA_VERSION=1.0 LOG_HEADERS=true go run main.go

# Make request
curl -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token123" \
  -d '{"event": "user.created", "data": {"id": 123}}'
```

**Logged output:**
```json
{
  "body": {
    "event": "user.created",
    "data": {"id": 123}
  },
  "headers": {
    "Authorization": "Bearer token123",
    "Content-Type": "application/json",
    "User-Agent": "curl/7.68.0"
  },
  "metadata": {
    "SERVICE": "api",
    "VERSION": "1.0"
  },
  "path": "/api/webhook"
}
```

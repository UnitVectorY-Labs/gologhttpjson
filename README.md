[![GitHub release](https://img.shields.io/github/release/UnitVectorY-Labs/gologhttpjson.svg)](https://github.com/UnitVectorY-Labs/gologhttpjson/releases/latest) [![License](https://img.shields.io/badge/license-MIT-blue)](https://opensource.org/licenses/MIT) [![Active](https://img.shields.io/badge/Status-Active-green)](https://guide.unitvectorylabs.com/bestpractices/status/#active) [![Go Report Card](https://goreportcard.com/badge/github.com/UnitVectorY-Labs/gologhttpjson)](https://goreportcard.com/report/github.com/UnitVectorY-Labs/gologhttpjson)

# gologhttpjson

A lightweight HTTP server that logs HTTP requests containing JSON payloads, with optional header logging and environment-based metadata.

## Purpose

This application captures incoming HTTP POST requests and logs the JSON body, request path, optional headers, and environment variable metadata. It is useful for debugging, testing, and local development when you want to inspect JSON payloads and related request context.

**Why use this?** This tool provides a simple way to log JSON requests for debugging, testing, and inspection. It allows for opt-in logging of request headers and provides a mechanism to inject metadata from environment variables.

**Should I run this in production?** No. This could expose sensitive information from headers and payloads. It is intended strictly for debugging, development, and testing.

## Usage

The latest gologhttpjson Docker image is available for deployment from GitHub Packages at [gologhttpjson on GitHub Packages](https://github.com/UnitVectorY-Labs/gologhttpjson/pkgs/container/gologhttpjson).

You can deploy this application locally with Docker:

```bash
docker run -p 8080:8080 ghcr.io/unitvectory-labs/gologhttpjson:latest
```

## Example Log Output

All responses return an HTTP 200 status code with a body of `OK`. This application is designed to log request payloads for inspection.

The log output is structured as JSON with the following attributes:

- `body` – the original JSON payload
- `path` – the HTTP request path
- `headers` – optional, logged if `LOG_HEADERS` is set
- `metadata` – optional, populated from environment variables prefixed with `METADATA_`

Example (pretty-printed for readability):

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

## Configuration

This application runs as a Docker container or via Go directly. It supports the following environment variables:

- `PORT` – the port the application listens on (default `8080`)
- `LOG_HEADERS` – set to any non-empty value to enable logging of HTTP headers
- `METADATA_*` – any environment variable prefixed with `METADATA_` will be logged under the `metadata` field

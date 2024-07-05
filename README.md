# Go Proxy Server

This repository contains a Go proxy server designed to forward HTTP requests and handle Cross-Origin Resource Sharing (CORS) errors. It listens for incoming requests and forwards them to the specified target URL, making it a simple solution for handling CORS issues.

## Features

- Handles CORS preflight requests (OPTIONS method).
- Forwards HTTP GET, POST, PUT, and DELETE requests to the target server.
- Copies headers from the original request and response.
- Adds CORS headers to the response.
- Logs requests and responses for debugging purposes.

## Getting Started

### Prerequisites

- Go 1.16 or later installed on your system. You can download Go from [here](https://golang.org/dl/).

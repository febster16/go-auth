# go-auth
Simple authentication / identity provider using Go

## Getting Started

1. After installing Go, navigate to the project directory in your terminal.
2. To start the service, run the following command:

   ```shell
   go mod download
   go run cmd/main.go
   ```
3. Or, docker:
   ```shell
   docker compose build
   docker compose up
   ```
The service should now be running locally on `http://localhost:{YOUR_PORT}`.
To access the Swagger docs, go to `/docs/index.html`

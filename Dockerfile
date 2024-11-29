# Build stage
FROM golang:1.23.3-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 

# Build the binary with the necessary environment variables
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o invoice-service ./cmd/api

# Final stage
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/invoice-service .

ENV GO_ENV=development

# Expose port 80
EXPOSE 80

# Run the binary
CMD ["./app/invoice-service"]

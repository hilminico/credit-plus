# Build stage
FROM golang:1.23 as builder

WORKDIR /app

# Enable Go modules
ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/main ./cmd

# Final stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary and locale files from builder
COPY --from=builder /app/bin/main .
COPY --from=builder /app/locale ./locale

# Expose port (adjust as needed)
EXPOSE 8080

CMD ["./main"]
# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o royal-api cmd/royal-api/*.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/royal-api .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

ENTRYPOINT ["./royal-api"]
CMD ["-command", "serve"]

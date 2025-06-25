# ---- Build Stage ----
FROM golang:1.24.2 AS builder

WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

# Copy the rest of the source code
COPY . .

# Create uploads directory in builder stage (optional, so it's copied later)
RUN mkdir -p /app/uploads

# Build the Go binary for Linux, static, no CGO
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o go-jwt-api .

# ---- Final Stage ----
FROM alpine:latest

# Create a non-root user and group
RUN addgroup -g 65532 appgroup && adduser -D -u 65532 -G appgroup appuser

WORKDIR /app

# Copy the Go binary from builder
COPY --from=builder /app/go-jwt-api .

# Copy uploads directory from builder
COPY --from=builder /app/uploads ./uploads

# Copy the .env file (if you want to bake it in)
COPY .env.docker .env

# Fix permissions so appuser can write to uploads
RUN chown -R 65532:65532 /app/uploads

# Switch to non-root user
USER 65532

# Run the app
CMD ["./go-jwt-api"]

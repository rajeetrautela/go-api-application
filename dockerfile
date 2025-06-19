# Stage 1: Build
FROM golang:1.24.2 AS builder

WORKDIR /go-jwt-api

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-jwt-api .

# Stage 2: Run
FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY --from=builder /go-jwt-api .

CMD ["/go-jwt-api"]

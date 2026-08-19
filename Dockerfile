# Stage 1: build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Download dependencies first (cached layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: minimal runtime image
FROM alpine:3.20

WORKDIR /app

# ca-certificates needed for TLS (MySQL over Railway private network)
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]

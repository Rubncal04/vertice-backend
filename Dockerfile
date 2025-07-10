# Build stage
FROM golang:1.23.5 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o vertice-backend ./cmd/main.go

# Run stage
FROM alpine:3.21.0
WORKDIR /app
COPY --from=builder /app/vertice-backend .
COPY .env .
EXPOSE 8080
CMD ["./vertice-backend"]
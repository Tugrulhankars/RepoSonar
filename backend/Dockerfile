# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app .


# Runtime stage
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8082

CMD ["./app"]






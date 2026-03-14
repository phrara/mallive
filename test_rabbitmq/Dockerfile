FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o test-rabbitmq .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/test-rabbitmq .

CMD ["./test-rabbitmq"]

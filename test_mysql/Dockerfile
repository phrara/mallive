FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY main.go .

RUN go build -o test main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/test .
CMD ["./test"]

# 构建阶段
FROM golang:1.25-alpine AS builder

# 安装构建工具
RUN apk add --no-cache git make

WORKDIR /build

# 复制依赖文件
COPY internal/common/go.mod go.mod
COPY internal/common/go.sum go.sum
COPY internal/payment/go.mod go.mod
COPY internal/payment/go.sum go.sum

# 下载依赖
RUN go mod download

# 复制源代码
COPY internal/common internal/common
COPY internal/payment internal/payment

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payment ./internal/payment/main.go

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/payment .

# 暴露端口
EXPOSE 8084


# 启动命令
CMD ["./payment"]

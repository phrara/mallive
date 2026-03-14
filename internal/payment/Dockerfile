# 构建阶段
FROM golang:1.25-alpine AS builder

# 安装构建工具
RUN apk add --no-cache git make

WORKDIR /build

# 复制项目结构
COPY internal/common internal/common
COPY internal/payment internal/payment

# 进入 payment 目录构建
WORKDIR /build/internal/payment

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payment .

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/internal/payment/payment .

# 暴露端口
EXPOSE 8084

# 启动命令
CMD ["./payment"]

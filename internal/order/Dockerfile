# 构建阶段
FROM golang:1.25-alpine AS builder

# 安装构建工具
RUN apk add --no-cache git make

WORKDIR /build

# 复制项目结构
COPY internal/common internal/common
COPY internal/order internal/order

# 进入 order 目录构建（因为 go.mod 在那里）
WORKDIR /build/internal/order

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order .

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/internal/order/order .

# 暴露端口
EXPOSE 8082 5002

# 启动命令
CMD ["./order"]

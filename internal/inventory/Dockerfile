# 构建阶段
FROM golang:1.25-alpine AS builder

# 安装构建工具
RUN apk add --no-cache git make

WORKDIR /build

# 复制项目结构
COPY internal/common internal/common
COPY internal/inventory internal/inventory

# 进入 inventory 目录构建
WORKDIR /build/internal/inventory

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o inventory .

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/internal/inventory/inventory .

# 暴露端口
EXPOSE 5003

# 启动命令
CMD ["./inventory"]

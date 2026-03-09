# 构建阶段
FROM golang:1.25-alpine AS builder

# 安装构建工具
RUN apk add --no-cache git make

WORKDIR /build

# 复制依赖文件
COPY internal/common/go.mod go.mod
COPY internal/common/go.sum go.sum
COPY internal/order/go.mod go.mod
COPY internal/order/go.sum go.sum

# 下载依赖
RUN go mod download

# 复制源代码
COPY internal/common internal/common
COPY internal/order internal/order

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order ./internal/order/main.go

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/order .

# 暴露端口
EXPOSE 8082 5002


# 启动命令
CMD ["./order"]

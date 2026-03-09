# ============================================
# Docker 镜像构建脚本
# ============================================
#
# 使用方法:
#   # 构建所有服务镜像
#   make build-images
#
#   # 单独构建某个服务
#   make build-order
#   make build-inventory
#   make build-payment
#   make build-kitchen
#
# ============================================

# 设置变量
REGISTRY ?= mallive
TAG ?= latest

.PHONY: build-images build-order build-inventory build-payment build-kitchen

# 构建所有服务镜像
build-images: build-order build-inventory build-payment build-kitchen

# 构建 Order 服务镜像
build-order:
	@echo "Building order image..."
	cd internal/order && \
	docker build -t $(REGISTRY)/order:$(TAG) .

# 构建 Inventory 服务镜像
build-inventory:
	@echo "Building inventory image..."
	cd internal/inventory && \
	docker build -t $(REGISTRY)/inventory:$(TAG) .

# 构建 Payment 服务镜像
build-payment:
	@echo "Building payment image..."
	cd internal/payment && \
	docker build -t $(REGISTRY)/payment:$(TAG) .

# 构建 Kitchen 服务镜像
build-kitchen:
	@echo "Building kitchen image..."
	cd internal/kitchen && \
	docker build -t $(REGISTRY)/kitchen:$(TAG) .

# 推送到镜像仓库
push: push-order push-inventory push-payment push-kitchen

push-order:
	docker push $(REGISTRY)/order:$(TAG)

push-inventory:
	docker push $(REGISTRY)/inventory:$(TAG)

push-payment:
	docker push $(REGISTRY)/payment:$(TAG)

push-kitchen:
	docker push $(REGISTRY)/kitchen:$(TAG)

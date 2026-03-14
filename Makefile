.PHONY: gen
gen: genproto genopenapi


.PHONY: genproto
genproto:
	@./scripts/genproto.sh

.PHONY: genopenapi
genopenapi:
	@./scripts/genopenapi.sh

.PHONY: root
root:
	cd ~/go-proj/mallive



# 设置变量
REGISTRY ?= mallive
TAG ?= latest

.PHONY: build-images build-order build-inventory build-payment build-kitchen

# 构建所有服务镜像
build-images: build-order build-inventory build-payment build-kitchen

# 构建 Order 服务镜像
build-order:
	@echo "Building order image..."
	docker build -t $(REGISTRY)/order:$(TAG) -f internal/order/Dockerfile .

# 构建 Inventory 服务镜像
build-inventory:
	@echo "Building inventory image..."
	cd internal/inventory && \
	docker build -t $(REGISTRY)/inventory:$(TAG) -f internal/inventory/Dockerfile .

# 构建 Payment 服务镜像
build-payment:
	@echo "Building payment image..."
	cd internal/payment && \
	docker build -t $(REGISTRY)/payment:$(TAG) -f internal/payment/Dockerfile .

# 构建 Kitchen 服务镜像
build-kitchen:
	@echo "Building kitchen image..."
	cd internal/kitchen && \
	docker build -t $(REGISTRY)/kitchen:$(TAG) .
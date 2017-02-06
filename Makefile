PREFIX=$(shell pwd)
COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date +%m%d%H%M)
export GOPATH=${PREFIX}
TAG=${COMMIT}-${DATE}

default: build
build:
	@echo "begin to build"
	@cd ./src/weixinGate && govendor sync -v
	@go build -v -ldflags ""  -o ${PREFIX}/bin/weixin-gate ./src/weixinGate
	@echo "build success"
clean:
	@echo "begin clean"	
	@rm bin/* || true
	@echo "clean success"
run:
	@echo "run"
	@cd ./bin && ./weixin-gate

pub-image: build-image publish-image
build-image:
	@echo "begin build images"
	@sudo docker build -t chengjt/weixin_gate:${TAG} .
	@echo "build success  chengjt/weixin_gate:${TAG}"
publish-image:
	@sudo docker login --username=1016890794@qq.com registry.cn-hangzhou.aliyuncs.com
	@sudo docker tag chengjt/weixin_gate:${TAG} registry.cn-hangzhou.aliyuncs.com/chengjt/weixin_gate:${TAG}
	@sudo docker push registry.cn-hangzhou.aliyuncs.com/chengjt/weixin_gate:${TAG}
	@sudo docker tag chengjt/weixin_gate:${TAG} registry.cn-hangzhou.aliyuncs.com/chengjt/weixin_gate:latest
	@sudo docker push registry.cn-hangzhou.aliyuncs.com/chengjt/weixin_gate:latest

push2registry: build build-image _push2registry
_push2registry:
	@sudo docker tag chengjt/weixin_gate:${TAG} ${REGISTRY}/chengjt/weixin_gate:${TAG}
	@sudo docker tag chengjt/weixin_gate:${TAG} ${REGISTRY}/chengjt/weixin_gate:latest
	@sudo docker push ${REGISTRY}/chengjt/weixin_gate:latest



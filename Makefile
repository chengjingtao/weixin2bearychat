PREFIX=$(shell pwd)
COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date +%m%d%H%M)
BUILDDATE=$(shell date +%Y%m%d%H%M)
export GOPATH=${PREFIX}
TAG=${COMMIT}-${DATE}

default: build
build:
	@echo "begin to build"
	@cd ./src/weixinGate && govendor sync -v
	@go build -v -ldflags '-X main.version=${COMMIT} -X main.buildDate=${BUILDDATE} -extldflags "-static"'  -o ${PREFIX}/bin/weixin-gate ./src/weixinGate
	@cp -r ./src/weixinGate/tmpl ./bin
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

push2registry: build build-image _push2registry
_push2registry:
	@sudo docker tag chengjt/weixin_gate:${TAG} ${REGISTRY}/chengjt/weixin_gate:${TAG}
	@sudo docker tag chengjt/weixin_gate:${TAG} ${REGISTRY}/chengjt/weixin_gate:latest
	@sudo docker push ${REGISTRY}/chengjt/weixin_gate:latest




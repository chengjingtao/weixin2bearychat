PREFIX=$(shell pwd)
COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date +%m%d%H%M)
BUILDDATE=$(shell date +%Y%m%d%H%M)
export GOPATH=${PREFIX}
TAG=${COMMIT}-${DATE}

default: build
build:
	@echo "begin to build"
	@cd ./src/weixin2bearychat && govendor sync -v
	@go build -v -ldflags '-X main.version=${COMMIT} -X main.buildDate=${BUILDDATE} -extldflags "-static"'  -o ${PREFIX}/bin/weixin2bearychat ./src/weixin2bearychat
	@cp -r ./src/weixin2bearychat/tmpl ./bin
	@echo "build success"
clean:
	@echo "begin clean"	
	@rm bin/* || true
	@echo "clean success"
run:
	@echo "run"
	@cd ./bin && ./weixin2bearychat

pub-image: build-image publish-image
build-image:
	@echo "begin build images"
	@sudo docker build -t chengjt/weixin2bearychat:${TAG} .
	@echo "build success  chengjt/weixin2bearychat:${TAG}"

push2registry: build build-image _push2registry
_push2registry:
	@sudo docker tag chengjt/weixin2bearychat:${TAG} ${REGISTRY}/chengjt/weixin2bearychat:${TAG}
	@sudo docker tag chengjt/weixin2bearychat:${TAG} ${REGISTRY}/chengjt/weixin2bearychat:latest
	@sudo docker push ${REGISTRY}/chengjt/weixin2bearychat:latest




PREFIX=$(shell pwd)
export GOPATH=${PREFIX}

default: build
build:
	@echo "begin to build"
	@cd ./src/weixinGate && govendor sync -v
	@go build  -o ${PREFIX}/bin/weixin-gate ./src/weixinGate
	@echo "build success"
clean:
	@echo "begin clean"	
	@rm bin/* || true
	@echo "clean success"
run:
	@echo "run"
	@cd ./bin && ./weixin-gate




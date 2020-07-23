#! /usr/bin/make
#
# Makefile for themis
#
# Targets:
# - clean     delete all generated files
# - generate  (re)generate all goagen-generated files.
# - build     compile executable
# - ae-build  build appengine
# - ae-dev    deploy to local (dev) appengine
# - ae-deploy deploy to appengine
#
# Meta targets:
# - all is the default target, it runs all the targets in the order above.
#
DEPEND=	bitbucket.org/pkg/inflect \
	github.com/goadesign/goa \
	github.com/goadesign/goa/goagen \
	github.com/goadesign/goa/logging/logrus \
	github.com/sirupsen/logrus \
	gopkg.in/yaml.v2 \
	golang.org/x/tools/cmd/goimports

CURRENT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

Version:=v1.0.2

all: depend clean generate build docker docker-release push-release

depend:
	@go get $(DEPEND)

clean:
	@rm -rf gen
	@rm -f themis
	@docker rm -f $(shell docker ps -a --filter status=exited -q);docker rmi -f $(shell docker images --filter dangling=true -q);docker rmi -f $(shell docker images --filter='reference=*/vdo/*' -q)

generate:
	@rm -rf svc/{gen,tmp}
	@goa gen github.com/JermineHu/themis/design/apis -o svc
	@goa example github.com/JermineHu/themis/design/apis -o svc
	@goa gen github.com/JermineHu/themis/design/apis -o svc/tmp
	@goa example github.com/JermineHu/themis/design/apis -o svc/tmp
	@mv svc/gen/http/openapi.* docs/
docker:
	@docker build --build-arg SHOW_SWAGGER=true -t  registry.cn-hangzhou.aliyuncs.com/vdo/themis:${Version} .
push:
	@docker push registry.cn-hangzhou.aliyuncs.com/vdo/themis:${Version}

docker-release:
	@docker build --build-arg SHOW_SWAGGER=fasle -t registry.cn-hangzhou.aliyuncs.com/vdo/themis:release-${Version} .

img-release:
	@img build --build-arg SHOW_SWAGGER=fasle -t registry.cn-hangzhou.aliyuncs.com/vdo/themis:release-${Version} .

push-release:
	@docker push registry.cn-hangzhou.aliyuncs.com/vdo/themis:release-${Version}

push-img-release:
	@img push registry.cn-hangzhou.aliyuncs.com/vdo/themis:release-${Version}

docker-push-release:docker-release push-release
img-push-release:img-release push-img-release
build:
	@CGO_ENABLED=0 GOOS=linux go build -a -v -installsuffix cgo -o themis github.com/JermineHu/themis/svc/cmd/themissvr
	@GITHUB_TOKEN=96ad1436db3c0302181aac88bf90eb0c178c938c golicense themis
	@strip --strip-unneeded themis
	@upx -9 -o themis

ae-build:
	@if [ ! -d $(HOME)/themis ]; then \
		mkdir $(HOME)/themis; \
		ln -s $(CURRENT_DIR)/appengine.go $(HOME)/themis/appengine.go; \
		ln -s $(CURRENT_DIR)/app.yaml     $(HOME)/themis/app.yaml; \
	fi

ae-deploy: ae-build
	cd $(HOME)/themis
	gcloud app deploy --project themis
FROM  golang:alpine as go
MAINTAINER Jermine.hu@qq.com
ARG SHOW_SWAGGER=false
ENV TIME_ZONE Asia/Shanghai
ENV APP_PACKAGE github.com/JermineHu/themis
ENV APP_SVC ${APP_PACKAGE}/svc
ENV APP_HOME /go/src/${APP_PACKAGE}
ENV APP_NAME themis
ENV SHOW_SWAGGER=$SHOW_SWAGGER
WORKDIR $APP_HOME
RUN apk update && apk add upx git binutils ca-certificates
#echo -e "http://dl-cdn.alpinelinux.org/alpine/edge/main\nhttp://dl-cdn.alpinelinux.org/alpine/edge/community" > /etc/apk/repositories && \
#   sed -i 's@http://dl-cdn.alpinelinux.org@https://mirrors.ustc.edu.cn@g' /etc/apk/repositories && apk update && apk add upx git binutils ca-certificates
RUN apk add --no-cache tzdata make protoc && \
    echo ${TIME_ZONE} > /etc/timezone && \
    ln -sf /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime
COPY . $APP_HOME
COPY .git $APP_HOME
RUN  flags="-X '${APP_SVC}.GoVersion=$(go version)' -X '${APP_SVC}.BuildTime=$(date '+%Y-%m-%d %H:%M:%S')' -X '${APP_SVC}.GitHash=$(git describe --tags --dirty --abbrev=8 --long)' -X '${APP_SVC}.GitLog=$(git log --pretty=oneline -n 1)' -X '${APP_SVC}.GitStatus=$(git status -s)'" && \
     go env -w GOPROXY=https://goproxy.cn,direct && go version \
     && go get -u goa.design/goa/v3/...@v3.2.0  && go get -u github.com/golang/protobuf/protoc-gen-go \
     && go get -u github.com/golang/protobuf/protoc-gen-go && \
     # make generate && \
     CGO_ENABLED=0 GOOS=linux go build -ldflags "$flags"  -a -installsuffix cgo -o $APP_NAME github.com/JermineHu/themis/svc/cmd/themissvr
# strip and compress the binary
RUN strip --strip-unneeded $APP_NAME
RUN upx -9 $APP_NAME
RUN if [ "$SHOW_SWAGGER" != "true" ] ; then rm -rf $APP_HOME/svc/gen/http/* ; echo "swagger 文档删除成功！"; fi #根据环境变量决定是否删除swagger文档
# 最终镜像设置
FROM busybox
MAINTAINER Jermine.hu@qq.com
ENV TIME_ZONE Asia/Shanghai
ENV BASE_API="https://jermine.vdo.pub"
ENV APP_HOME /go/src/github.com/JermineHu/themis
WORKDIR /bin
COPY --from=go $APP_HOME/themis /bin/
#COPY --from=go $APP_HOME/.gitignore $APP_HOME/svc/gen/http/* /bin/
COPY --from=go /etc/timezone /etc/timezone
COPY --from=go /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime
COPY --from=go /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY jwtkey /bin/jwtkey
COPY ui /bin/ui
ENV APP_PORT 8081
ENV ENV development
EXPOSE  $APP_PORT
CMD sed -i "s@BASE_API@$BASE_API@g" `grep BASE_API -rl .` && ./themis

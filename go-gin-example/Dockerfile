FROM registry.cn-beijing.aliyuncs.com/baxiang/golang:1.14 as build
ENV GO111MODULE=off
ENV GOBUILDPATH=code-ai.cubigdata.cn/kg/luban-backend
RUN mkdir -p /go/src/${GOBUILDPATH}&&mkdir -p /go/bin/
COPY ./ /go/src/${GOBUILDPATH}
RUN cd /go/src/${GOBUILDPATH}/cmd && \
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /go/bin/luban luban.go

FROM  centos:8 as dev
WORKDIR /luban
ARG ENV=release
ENV RUN_MODE=$ENV
ENV TIME_ZONE Asia/Shanghai
RUN ln -fs /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime \
    && echo ${TIME_ZONE} > /etc/timezone  \
    && mkdir logs \
    && mkdir keys \
    && mkdir configs \
    && mkdir dict \
    && mkdir download
COPY --from=build /go/bin/luban /usr/local/bin/
COPY ./configs/*.yaml  ./
COPY ./keys ./keys
COPY ./dict ./dict
COPY ./scripts/docker-entrypoint.sh ./
ENTRYPOINT ["./docker-entrypoint.sh"]
EXPOSE 8880
CMD ["luban"]


FROM registry.cn-beijing.aliyuncs.com/baxiang/golang:1.14 as build
ENV GOPROXY https://goproxy.cn,direct
ENV GOBUILDPATH=github.com/baxiang/soldiers-sortie/blog
RUN mkdir -p /go/src/${GOBUILDPATH}&&mkdir -p /go/bin/
COPY ./ /go/src/${GOBUILDPATH}
RUN cd /go/src/${GOBUILDPATH}/cmd && \
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/blog main.go

FROM alpine as dev
WORKDIR /blog
ENV TIME_ZONE Asia/Shanghai
RUN ln -fs /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime \
    && echo ${TIME_ZONE} > /etc/timezone  \
    && mkdir logs \
    && mkdir conf
COPY --from=build /go/bin/blog ./
COPY ./conf/app.ini  ./
EXPOSE 8004
CMD ["./blog"]


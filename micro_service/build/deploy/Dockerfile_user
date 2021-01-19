FROM alpine:3.9.3

COPY ./configs/ /data/configs/
COPY ./userservice /data/

RUN ls -alh /data/configs/
RUN ls -alh /data/

WORKDIR /data/
CMD ./userservice


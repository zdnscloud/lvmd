FROM golang:1.13.7-alpine3.11 AS build
ENV GOPROXY=https://goproxy.cn
  
RUN mkdir -p /go/src/github.com/zdnscloud/lvmd
COPY . /go/src/github.com/zdnscloud/lvmd

WORKDIR /go/src/github.com/zdnscloud/lvmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/src/github.com/zdnscloud/lvmd/lvmd

FROM alpine:3.10.0

LABEL maintainers="Zdns Authors"
LABEL description="K8S Lvmd"
RUN apk update && apk add udev blkid file util-linux e2fsprogs lvm2 udev sgdisk device-mapper e2fsprogs e2fsprogs-extra cfdisk thin-provisioning-tools
COPY --from=build /go/src/github.com/zdnscloud/lvmd/lvmd /lvmd
ENTRYPOINT ["/bin/sh"]

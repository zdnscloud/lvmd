FROM golang:1.12.5-alpine3.9 AS build
  
RUN mkdir -p /go/src/github.com/zdnscloud/lvmd-server
COPY . /go/src/github.com/zdnscloud/lvmd-server

WORKDIR /go/src/github.com/zdnscloud/lvmd-server
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/src/github.com/zdnscloud/lvmd-server/lvmd

FROM alpine:3.9.4

LABEL maintainers="Zdns Authors"
LABEL description="K8S Lvmd"
RUN apk update && apk add lvm2
ADD shell/lvmd.sh /lvmd.sh
ADD shell/getblocks /getblocks
COPY --from=build /go/src/github.com/zdnscloud/lvmd-server/lvmd /lvmd
RUN chmod 755 /getblocks
RUN chmod 755 /lvmd.sh
RUN chmod 755 /lvmd
ENTRYPOINT ["/bin/sh"]

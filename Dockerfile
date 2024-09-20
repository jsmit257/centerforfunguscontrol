FROM golang:bookworm AS build
ENV gopath=/dev
COPY . /go/src/github.com/jsmit257/huautla
WORKDIR /go/src/github.com/jsmit257/huautla
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o /huautla -a -installsuffix cgo -cover ./ingress/http/...

FROM alpine:edge AS deploy
ENV GOCOVERDIR=/tmp
RUN apk update
RUN apk add jq curl
COPY --from=build /huautla /huautla
COPY ./www /www
ENTRYPOINT [ "/huautla" ]

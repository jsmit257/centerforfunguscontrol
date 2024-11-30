FROM golang:1.23.1 AS build
ENV gopath=/dev
COPY . /go/src/github.com/jsmit257/cffc
WORKDIR /go/src/github.com/jsmit257/cffc
RUN grep -v -e '^replace' go.mod > new.mod && mv new.mod go.mod
RUN go mod tidy
RUN go mod vendor
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go get net/http
RUN go build -o /cffc -a -installsuffix cgo -cover ./ingress/http/...

FROM alpine:edge AS deploy
ENV GOCOVERDIR=/tmp
RUN apk update
RUN apk add jq curl
COPY --from=build /cffc /cffc
ENTRYPOINT [ "/cffc" ]

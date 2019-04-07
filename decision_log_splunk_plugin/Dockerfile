FROM golang:1.11-alpine
RUN mkdir /src
WORKDIR /src
RUN apk add --update git make curl bash gcc
RUN apk add libc-dev
RUN go get github.com/open-policy-agent/opa/plugins
COPY main.go /src
RUN go build -buildmode=plugin -o=plugin.so main.go



# Compile stage
FROM golang:1.14.3 AS build-env

# Build Delve
# RUN go env -w GO111MODULE=on
# RUN go env -w GOPROXY=https://goproxy.cn
# RUN go get github.com/go-delve/delve/cmd/dlv

ADD . /dockerdev
WORKDIR /dockerdev

RUN go build -o /server
# RUN go build -gcflags="all=-N -l" -o /server

# Final stage
FROM debian:buster

EXPOSE 1833

WORKDIR /
# COPY --from=build-env /go/bin/dlv /
COPY --from=build-env /server /

CMD ["/server"]
# CMD ["/dlv", "--listen=:7777", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/server"]
FROM golang:1.18-alpine3.17 as build
LABEL author="yongjie.zhuang"

WORKDIR /go/src/build/

# for golang env
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# dependencies
COPY go.mod .
COPY go.sum .

RUN go mod download

# build executable
COPY . .
RUN go build -o main


FROM alpine:3.17
WORKDIR /usr/src/
COPY --from=build /go/src/build/main ./main
COPY --from=build /go/src/build/app-conf-dev.yml ./app-conf-dev.yml

CMD ["./main"]
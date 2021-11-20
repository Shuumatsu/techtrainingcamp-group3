FROM --platform=linux/x86_64 golang:1.17.2

RUN sed -i s@/deb.debian.org/@/mirrors.aliyun.com/@g /etc/apt/sources.list
RUN sed -i s@/security.debian.org/@/mirrors.aliyun.com/@g /etc/apt/sources.list
RUN apt-get update && apt-get -y install apt-utils git unzip build-essential autoconf libtool protobuf-compiler

ADD . /workspace
WORKDIR /workspace

RUN go env -w GO111MODULE="on" && go env -w GOPROXY="https://goproxy.cn,direct" 
RUN go get github.com/twitchtv/retool
RUN go get golang.org/x/tools/cmd/goimports

RUN go mod vendor && go mod tidy
RUN make all

ENV DB_PASSWD=Group3123456 \
    DB_HOST=rdsmysqlh91d80c1b0e1f4d55.rds.ivolces.com \
    REDIS_HOST=redis-cn024n27ch1c7nrnx.redis.ivolces.com \
    KAFKA_HOST=kafka-ylsiccmljqg9rb.cn-beijing.kafka-internal.ivolces.com


CMD ["./bin/http"]

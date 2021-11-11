FROM golang:1.17.2
ADD . /workspace
WORKDIR /workspace
RUN go env -w GO111MODULE="on" && go env -w GOPROXY="https://goproxy.cn,direct"
RUN go mod vendor && go mod tidy

ENV DB_PASSWD=Group3123456 \
    DB_HOST=rdsmysqlh91d80c1b0e1f4d55.rds.ivolces.com \
    REDIS_HOST=redis-cn024n27ch1c7nrnx.redis.ivolces.com \
    KAFKA_HOST=kafka-ylsiccmljqg9rb.cn-beijing.kafka-internal.ivolces.com

CMD ["go", "run", "main.go"]
EXPOSE 8080

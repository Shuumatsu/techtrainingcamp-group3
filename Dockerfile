FROM golang:1.17.2
ADD . /workspace
WORKDIR /workspace
RUN go env -w GO111MODULE="on" \
&& go env -w GOPROXY="https://goproxy.cn,direct"
RUN go mod tidy
CMD ["go", "run", "main.go"]
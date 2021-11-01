FROM ubuntu
# 添加项目到工作目录
ADD ./ /workspace/
# 设置国内镜像
RUN	sed -i s@/archive.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list \
&& 	sed -i s@/security.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list \
&&	apt-get clean \
# 更新apt-get
&& apt-get update \
&& apt-get upgrade -y
# 安装vim
RUN apt-get install -y vim
# 安装wget
RUN apt-get install -y wget
# 下载golang1.17.2并解压
RUN wget -c https://dl.google.com/go/go1.17.2.linux-amd64.tar.gz -O - | tar -xz -C /
# 设置golang环境变量
ENV PATH=$PATH:/go/bin
RUN go env -w GOPATH="/go/gopath" \
&& go env -w GO111MODULE="on" \
&& go env -w GOPROXY="https://goproxy.cn,direct"

# 安装mysql
RUN apt-get install -y mysql-server mysql-client
# 启动mysql服务
RUN service mysql start

EXPOSE 8080

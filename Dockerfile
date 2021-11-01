# 使用请运行: "docker build -t [imageName] ."
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
# 安装lsof
RUN apt-get install -y lsof
# 安装wget
RUN apt-get install -y wget \
# 下载golang1.17.2并解压
&& wget -c https://dl.google.com/go/go1.17.2.linux-amd64.tar.gz -O - | tar -xz -C /
# 设置golang环境变量
ENV PATH=$PATH:/go/bin
RUN go env -w GOPATH="/go/gopath" \
&& go env -w GO111MODULE="on" \
&& go env -w GOPROXY="https://goproxy.cn,direct"

# 安装mysql
RUN apt-get install -y mysql-server mysql-client
# 启动mysql服务
RUN service mysql start \
# 创建mysql用户
&& echo "create user 'group3'@'localhost' identified by '123456';grant all privileges on *.* to 'group3'@'localhost' with grant option;" > /tmp/createuser.sql \
&& mysql -u root < /tmp/createuser.sql \
&& rm /tmp/createuser.sql \
# 运行建库脚本
&& mysql -u root < /workspace/db/database_builder.sql

# 下载golang依赖
RUN cd /workspace && go mod tidy && go mod vendor

# 启动mysql服务
ENTRYPOINT ["service", "mysql", "start"]
# 运行项目
CMD cd /workspace;go run main.go
# 暴露8080端口
EXPOSE 8080

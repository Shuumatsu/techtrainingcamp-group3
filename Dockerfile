FROM golang:1.17.2
ADD . /workspace
WORKDIR /workspace
RUN go env -w GO111MODULE="on" \
&& go env -w GOPROXY="https://goproxy.cn,direct"
<<<<<<< HEAD

# 安装mysql
RUN apt-get install -y mysql-server mysql-client
# 启动mysql服务
RUN usermod -d /var/lib/mysql/ mysql \
&& service mysql start \
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
CMD cd /workspace; go run main.go

# 声明：该项目使用8080端口
EXPOSE 8080
=======
RUN go mod tidy
CMD ["go", "run", "main.go"]
>>>>>>> dockerCompose

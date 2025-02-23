# 第一阶段：构建应用程序
FROM golang:1.24 AS builder
# 定义构建参数，在构建镜像的时候需要传入 SERVICE 来表示构建的模块
ARG SERVICE

# 设置环境变量
ENV TZ=Asia/Shanghai
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct
ENV GOOS=linux
ENV GOARCH=amd64
# 创建工作目录c
RUN mkdir -p /app
RUN mkdir -p /app/output

WORKDIR /app

# 复制所有文件到工作目录
COPY . .

# 下载依赖
RUN go mod tidy

# 编译应用程序
RUN bash ./docker/script/build.sh ${SERVICE}

# 第二阶段：创建最终运行环境
FROM alpine

ARG SERVICE

# 环境变量
ENV TZ=Asia/Shanghai
ENV SERVICE=${SERVICE}
# ETCD_ADDR 默认值在 entrypoint 中设置，不在 Dockerfile 中提供

# 换源，更新软件依赖
RUN sed -i 's#https\?://dl-cdn.alpinelinux.org/alpine#https://mirrors.tuna.tsinghua.edu.cn/alpine#g' /etc/apk/repositories
RUN apk update --no-cache && apk --no-cache add ca-certificates tzdata bash

# 创建工作目录
WORKDIR /app

# 复制 Binary 及启动脚本，调整权限设置
COPY --from=builder /app/output /app/output
COPY --from=builder /app/docker/script/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# 需要注意的是，这个命令是在程序执行时运行，不能直接在这里填入 SERVICE 参，需要绕个弯设置为环境变量
CMD ["sh", "-c", "./entrypoint.sh"]



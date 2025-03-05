# Deploy

这份指南将从零开始部署本项目。

在部署过程中，我们将使用 Makefile 中的命令来启动依赖的环境（包括 MySQL、etcd、Redis 等数据库和缓存服务），然后编译并运行特定的服务。

如果对具体的构建和启动流程感兴趣，请查阅 [Build](./build.zh.md) 页面。

## 先决条件

需要安装的软件包：

- Docker
    - 请参考 [Docker 官方文档](https://docs.docker.com/get-docker/) 安装 Docker。
    - 用于启动 MySQL、etcd、Redis 等服务。
- tmux
    - 便于在终端中管理多个会话，提高运维效率。
    - Makefile 中的命令会默认在 tmux 中执行。

## 本地部署

### 预备操作

修改 `config/config.yaml` 的配置，将数据库等配置的 ip 修改为 `localhost`（如果没有请新增这个文件）。  
我们使用了 oss 来存储图像数据，所以你还需要在 `config/config.yaml` 中配置 oss（upyun） 的相关信息。  
出于安全性考虑，你需要修改 `config/config.yaml` 中的密钥部分，将其设置为你自己的密钥与公钥。  
配置示例请参考 `config.example.yaml`。

### 启动环境

#### 清理本地环境（可选）

```shell
make clean-all
```

#### 启动环境基础容器（数据库等）

```shell
make env-up
```

### 启动特定服务

```shell
make <target>
```

`<target>` 为服务名称，即 `cmd` 目录下的文件夹名称。

完成后的目录结构应该与下面的结构类似：

```shell
.
├── docker
│   ├── script
│   │   └── etcd-monitor.sh
│   ├── env
│   │   ├── redis.env
│   │   ├── mysql.env
│   │   └── etcd.env
│   ├── docker-compose.yaml
├── config
│   ├── sql
│   │   └── init.sql
│   └── config.yaml
└── hack
    │── image-refresh.sh
    └── docker-run.sh
```

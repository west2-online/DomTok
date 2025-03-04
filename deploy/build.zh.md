# Build

该指南将主要介绍项目的构建和启动流程，以及相关的脚本。

后续的 `<target>` 指代某个服务，例如 `api`，具体可以通过 `make help` 获取可构建服务列表。

## 大致流程

1. 使用命令 `make env-up` 启动环境（MySQL、etcd、Redis 等）。
2. `make <target>` 编译并运行具体的服务。
3. 服务从 etcd 中获取 `config.yaml`。
4. 读取 `config.yaml` 中的配置，将 Env 映射到对应的 **结构体** 中。
5. 从 `config.yaml` 中获取可用地址。
6. 初始化服务，将服务注册到 `etcd` 中。
7. 启动服务。

```mermaid
sequenceDiagram
    participant User as 用户
    participant Makefile as Makefile
    participant Services as 服务
    participant Env as 环境 (MySQL, etcd, Redis)

    User->>Makefile: 运行 `make env-up`
    Makefile->>Env: 启动 MySQL、Redis 和 etcd

    User->>Makefile: 运行 `make <target>`
    Makefile->>Services: 编译并启动指定服务

    Services->>Env: 从 etcd 获取 `config.yaml`
    Env-->>Services: 返回 `config.yaml`

    Services->>Services: 读取 `config.yaml` 并映射到结构体, 并获取可用 IP 地址

    Services->>Env: 将服务注册到 etcd
    Services->>Services: 初始化完成, 启动服务
```

## 构建和启动

### 目录结构

项目的关键目录如下：

- `cmd/`：包含各服务模块的启动入口。
- `output/`：构建产物的输出目录。

### 构建流程

此处阐述当我们敲下 `make <target>` 时具体的工作流程，我们省略了 tmux 环境的相关内容。

构建过程主要通过 [build.sh](../docker/script/build.sh) 脚本完成，用于编译指定服务模块的二进制文件或进行系统测试：

1. 进入到 cmd 中的对应服务文件夹下。
2. 执行 `go build` 以编译该服务的二进制文件，并存放到 `output` 文件夹内。

```mermaid
flowchart TD
    A[启动脚本] --> B{检查是否输入了参数, 即 target}
    B --> |为空| C[输出错误信息并退出]
    B --> |不为空| D[设置ROOT_DIR为当前工作目录]

    D --> E[进入指定模块目录 ./cmd/RUN_NAME]
    E --> F[创建文件夹 output/RUN_NAME 并设置权限]

    F --> G{判断是否为测试环境}
    G --> |是测试环境| H[执行测试构建]
    G --> |非测试环境| I[执行构建]

    H --> J[生成测试二进制文件 output/RUN_NAME/fzuhelper-RUN_NAME]
    I --> K[生成构建二进制文件 output/RUN_NAME/fzuhelper-RUN_NAME]
```

### output 目录结构

```text
 output
 └── target
        └── binary
```

### 启动流程

当我们敲下 `make <target>` 而没有设置仅构建标志（`BUILD_ONLY`）时，会自动启动，这里介绍本地调试启动的过程。

> Docker 容器的启动过程是类似的，只是将其移动至容器内

启动过程主要通过 [entrypoint.sh](/docker/script/entrypoint.sh) 脚本完成

1. 通过`export`设置 etcd 的地址的环境变量，为后续程序在 **运行时** 能够获取到 etcd 的地址并获取 `config.yaml`。
2. `cd` 到构建阶段生成的 output 目录，执行对应服务的二进制文件。

```mermaid
flowchart TD
    A[启动 entrypoint.sh] --> B{检查 ETCD_ADDR 是否已设置}
    B --> |未设置| C[设置默认 ETCD_ADDR=localhost:2379]
    B --> |已设置| D[保持现有 ETCD_ADDR]

    C --> E[导出 ETCD_ADDR 环境变量]
    D --> E

    E --> F[启动服务]
```

## 使用方式

两份脚本都由 `Makefile` 中的命令接管，可以通过以下命令调用：

```shell
make <target> [option]   # option = BUILD_ONLY
```

以下是 `make <target>` 的大致流程图：

```mermaid
flowchart TD
    A[启动 make <target> 命令] --> B{检查是否传入 BUILD_ONLY 设置}

    B -- 未设置 --> C[构建并运行]
    B -- 已设置 --> D[只进行构建]

    D --> E[创建 output 目录]
    E --> F[运行 build.sh 脚本进行编译]

    C --> G[创建 output 目录]
    G --> H[运行 build.sh 脚本进行编译]
    H --> I[运行 entrypoint.sh 启动服务]
    I --> J[服务启动完成]
```

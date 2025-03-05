<div align="center">
  <h1 style="display: inline-block; vertical-align: middle;">DomTok</h1>
</div>

## 概述
DomTok 是一个基于 HTTP 和 RPC 协议的简单抖音电商后端项目，采用了整洁架构和分布式架构。它使用了 Kitex、Hertz、Eino、Mysql、Redis、Etcd、Kafka、Elasticsearch、Kibana、Filebeat、Jaeger、Rocketmq、Otel - Collector、多个导出器、Prometheus、VictoriaMetrics、Cadvisor 和 Grafana 等技术。

## 特性
- 云原生：采用原生 Go 语言分布式架构设计，基于字节跳动的最佳实践。
- 高性能：支持异步 RPC、非阻塞 I/O、消息队列和即时（JIT）编译。
- 可扩展性：基于整洁架构进行模块化和分层结构设计，代码清晰易读，降低了开发难度。
- 可观测性：基于 OpenTelemetry 进行分布式追踪，使用 Prometheus 进行监控，利用 Elasticsearch 进行日志收集，通过 Grafana 进行可视化展示。
- 代码质量：基于 GitHub Actions 实现 CI/CD 流程，拥有丰富的单元测试，代码质量高且安全性强。
- AI 功能：基于字节跳动的 Eino 框架和大语言模型（LLM），通过**函数调用**实现文本输入调用接口。
- 开发运维：丰富的脚本和工具减少了不必要的手动操作，简化了使用和部署流程。

## 架构
![架构图](./img/Architecture.png)

### 编码架构
我们基于整洁架构对项目进行了分层设计，如下图所示：
![编码架构图](./img/Coding-architecture.png)

## 项目结构

### 整体结构
```text
.
├── LICENSE
├── Makefile                # 一些 make 命令
├── README.md     
├── app                     # 各种微服务的实现
├── cmd                     # 各种微服务的启动入口
├── config                  # 配置文件
├── deploy                  # 部署文件
├── docker                  # 与 Docker 相关
├── go.mod
├── go.sum
├── hack                    # 用于自动化开发、构建和部署任务的工具
├── idl                     # 接口定义
├── kitex_gen               # Kitex 生成的代码
└── pkg
    ├── base                # 通用基础服务
    │   ├── client    	    # 相应组件的客户端（如 redis、mysql）
    │   └── context         # 用于在服务之间传递数据的自定义上下文
    ├── constants           # 存储常量
    ├── errno               # 自定义错误
    ├── kafka               # Kafka 功能的一些封装
    ├── logger              # 日志系统
    ├── middleware          # 中间件
    ├── upyun               # 又拍云的一些封装
    └── utils               # 一些实用函数
```

### 网关/API 模块
```text
./app/gateway
├── handler                 # 处理请求的处理器
├── model                   # hz 生成的模型
├── mw                      # 中间件
├── pack                    # 封装请求和响应
├── router                  # 路由
└── rpc                     # 发送 RPC 请求
```

### 微服务（订单模块）
```text
./app/order
├── controllers       # rpcService 接口的实现层，负责转换请求和响应
├── domain            # 整洁架构中的领域层
│   ├── model         # 定义模块内使用的结构体
│   ├── repository    # 定义模块内使用的接口
│   └── service       # 实现可复用的核心业务逻辑
├── infrastructure    # 整洁架构中的接口层，命名为 infrastructure 以避免歧义
│   ├── locker        # 领域仓库中 locker 接口的具体实现
│   ├── mq            # 领域仓库中 mq 接口的具体实现
│   ├── mysql         # 领域仓库中 db 接口的具体实现
│   ├── redis         # 领域仓库中缓存接口的具体实现
│   └── rpc           # 领域仓库中 rpc 接口的具体实现
└── usecase
```

## LLM 应用
基于字节跳动的 Eino 框架和大语言模型（LLM），本项目实现了一个简单的 AI 功能。

### 节点拓扑
最开始我们打算用 Eino 官方提供的 Goland 插件 **Eino Dev** 来完成 LLM 的节点编排操作：
![使用 Eino Dev 插件完成节点编排](./img/llm/llm-node-compile-1.png)

后来，在官方文档中，我们发现一个很相似的功能——[React Agent](https://www.cloudwego.io/zh/docs/eino/core_modules/flow_integration_components/react_agent_manual/)

**React Agent** 的节点拓扑和数据流图与我们最初的设计非常相似：
![React Agent 节点拓扑](./img/llm/llm-node-compile-2.png)

同时，**React Agent** 的使用方式更加简单，所以我们最终选择了 **React Agent** 来完成 LLM 应用的实现（图片来自[React Agent使用手册](https://www.cloudwego.io/zh/docs/eino/core_modules/flow_integration_components/react_agent_manual/)）：
![React Agent 流程](./img/llm/llm-node-compile-3.png)

## 测试
- 单元测试：本项目使用 `GitHub/bytedance/mockey` 和 `GitHub.com/smartystreets/goconvey/convey` 编写了丰富的单元测试。你可以使用 `make test` 来运行这些测试。
- 带环境的单元测试：除了需要模拟的单元测试外，我们还使用环境变量来控制测试环境，使我们的一些单元测试可以在真实环境中运行。你可以使用 `make with-env-test` 来启动环境并运行这些测试。
- API 接口测试：我们使用 **Apifox** 对接口进行全自动测试以确保接口的正确性。我们将测试报告导出为 **html** 格式，你可以[点击此处](./resource/domtok-apifox-reports.html)查看我们的测试用例。
- GitHub Workflow：我们使用 GitHub Actions 进行 CI/CD 流程，确保代码质量和安全性。

## 快速启动和部署
本项目通过脚本极大地简化了流程。你可以参考[部署文档](deploy.zh.md)来快速启动和部署项目。


## 可视化示例

### 链路追踪
我们将 span 埋点在中间件中，使用 `Jaeger` 进行链路追踪，如下图所示：
#### Jaeger
![Jaeger 监控图](./img/metrics/jaeger.png)
![jaeger-tracing.png](img/jaeger-tracing.png)

### 日志体系
我们利用 `go.uber.org/zap` 将日志序列化为 json 格式并输出到特定的文件，让 `filebeat` 对日志文件进行收集 并发送到 `elasticsearch`，最终在 `kibana` 展示。

我们使用 **"efk"** 而不是 **"elk"** 体系，原因与分析在[这里](./efk-log.zh.md)。

#### 可视化日志查询
我们使用 `kibana` 的创建了简洁有效的数据视图，并且可以使用 
[KQL](https://learn.microsoft.com/en-us/kusto/query/?view=microsoft-fabric) 进行可视化查询
样例如下：
![kibana-data-graph.png](img/kibana-data-graph.png)
![kql-demo.png](img/kql-demo.png)

#### Kibana-dev-tools 查询
除了利用 **kql** 与 **kibana-dataGraph** 进行可视化查询之外，当然也可以进行常规的查询。如下：
![kibana-dev-tool.png](img/kibana-dev-tool.png)


### 监控
我们将展示通过 `Prometheus`、`Grafana`、`VictoriaMetrics`、`Otel - Collector` 等诸多工具实现了监控的可视化（由于示例较多，仅展示部分）。

#### Docker
![docker 监控图](./img/metrics/docker.png)

#### Go 程序（总计）
![Go 程序监控图](./img/metrics/go.png)

#### Mysql
![Mysql 监控图](./img/metrics/mysql.png)

#### Redis
![Redis 监控图](./img/metrics/redis.png)

#### 系统
![系统监控图](./img/metrics/system.png)

## 贡献者

<a href="https://github.com/west2-online/DomTok/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=west2-online/DomTok" />
</a>

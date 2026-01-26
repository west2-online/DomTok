# Project Context

## Purpose
DomTok is a high-performance, cloud-native distributed e-commerce backend system designed to mimic Douyin E-commerce. It leverages a microservices architecture to provide core e-commerce functionalities such as user management, product browsing (commodity), shopping cart operations, order processing, and payments. Additionally, it integrates an AI assistant powered by ByteDance's Eino framework to enable natural language interactions for tasks like order querying and cart management.

## Tech Stack
- **Language**: Go (Golang) v1.23+
- **Web/Gateway Framework**: Hertz
- **RPC Framework**: Kitex (with Thrift IDL)
- **AI/LLM Framework**: Eino (ByteDance's LLM framework)
- **Database**: 
  - MySQL (Primary relational DB, accessed via Gorm)
  - Redis (Caching and distributed locks)
  - Elasticsearch (Search engine)
- **Message Queue**: 
  - RocketMQ
  - Kafka
- **Service Discovery**: Etcd
- **Observability**: 
  - **Tracing**: OpenTelemetry, Jaeger
  - **Metrics**: Prometheus, Grafana, VictoriaMetrics, Cadvisor
  - **Logging**: Zap (Logger), Filebeat, Kibana (EFK stack)
- **Infrastructure**: Docker, Docker Compose
- **Other Libraries**: 
  - `gorm.io/gorm` (ORM)
  - `github.com/spf13/viper` (Configuration)
  - `github.com/bytedance/mockey`, `github.com/smartystreets/goconvey` (Testing)

## Project Conventions

### Code Style
- Follows standard Go formatting and idioms (`gofmt`).
- Uses `golangci-lint` for static analysis.
- **Logging**: Structured JSON logging using `uber-go/zap`.
- **Error Handling**: Custom error codes and stack traces via `pkg/errno`.

### Architecture Patterns
- **Microservices**: The system is divided into distinct services:
  - `gateway`: HTTP entry point / API Gateway.
  - `user`, `cart`, `commodity`, `order`, `payment`: Core domain services.
  - `assistant`: AI assistant service.
- **Clean Architecture**: Each microservice strictly follows Clean Architecture principles:
  - **Controller/Adapter Layer** (`controllers`): Handles incoming RPC/HTTP requests and converts them to domain models. Isolated from business logic.
  - **Use Case Layer** (`usecase`): Orchestrates business logic by coordinating domain services and repositories. Pure business logic, framework-agnostic.
  - **Domain Layer** (`domain`):
    - `model`: Business entities/structs.
    - `repository`: Interfaces for data access (DB, Cache, RPC).
    - `service`: Core domain logic and rules.
  - **Infrastructure Layer** (`infrastructure`): Concrete implementations of repository interfaces (MySQL, Redis, RPC clients, MQ producers).

### Testing Strategy
- **Unit Tests**: Extensive unit testing using `Mockey` for mocking and `GoConvey` for assertions.
- **Integration Tests**: Environment-aware tests (`make with-env-test`) that run against real dependencies (controlled via env vars).
- **API Tests**: Automated API testing using Apifox.
- **CI/CD**: GitHub Actions workflows for build, test, and linting (`.github/workflows`).

### Git Workflow
- Pull Request (PR) based workflow using `PULL_REQUEST_TEMPLATE.md`.
- `Makefile` provided for common development tasks (`make build`, `make test`, etc.).

## Domain Context
- **E-commerce Model**:
  - **Commodity**: Hierarchical structure with Categories, SPUs (Standard Product Units), and SKUs (Stock Keeping Units). Supports coupons.
  - **Order**: Lifecycle management, inventory reservation, and distributed transaction handling.
  - **Cart**: User-specific shopping cart management.
  - **Payment**: Integration with payment providers (mocked or real).
- **AI Integration**: The `assistant` service uses `Eino` to map natural language inputs to backend "Tools" (Function Calls) like `cart_show`, `order_list`, etc., allowing users to perform actions via chat.

## Important Constraints
- **Distributed System**: Must handle distributed transactions (likely eventual consistency via MQ), distributed locking (Redis), and service discovery.
- **Performance**: High concurrency requirements; usage of asynchronous RPC and caching is critical.
- **Protocol**: 
  - External: HTTP (Hertz)
  - Internal: Thrift RPC (Kitex)

## External Dependencies
- **UpYun**: Used for object storage (images, static assets).
- **VolcEngine**: ByteDance cloud services SDK usage.

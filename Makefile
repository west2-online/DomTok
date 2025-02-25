# 辅助工具安装列表
# 执行 go install github.com/cloudwego/hertz/cmd/hz@latest
# 执行 go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
# 执行 go install golang.org/x/tools/cmd/goimports@latest
# 执行 go install golang.org/x/vuln/cmd/govulncheck@latest
# 执行 go install mvdan.cc/gofumpt@latest
# 访问 https://golangci-lint.run/welcome/install/ 以查看安装 golangci-lint 的方法



# 默认输出帮助信息
.DEFAULT_GOAL := help
# 检查 tmux 是否存在
TMUX_EXISTS := $(shell command -v tmux)
# 远程仓库
REMOTE_REPOSITORY = registry.cn-hangzhou.aliyuncs.com/west2-online/DomTok
# 项目 MODULE 名
MODULE = github.com/west2-online/DomTok
# 当前架构
ARCH := $(shell uname -m)
PREFIX = "[Makefile]"
# 目录相关
DIR = $(shell pwd)
CMD = $(DIR)/cmd
CONFIG_PATH = $(DIR)/config
IDL_PATH = $(DIR)/idl
OUTPUT_PATH = $(DIR)/output
API_PATH= $(DIR)/cmd/api

#es
ES_ANALYSIS = domtok-elasticsearch

# 服务名
SERVICES := gateway user commodity orders cart payment assistant
service = $(word 1, $@)

EnvironmentStartEnv = DOMTOK_ENVIRONMENT_STARTED
EnvironmentStartFlag = true

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  {service name}    : Build a specific service (e.g., make api). use BUILD_ONLY=1 to avoid auto bootstrap."
	@echo "                      Available service list: [${SERVICES}]"
	@echo "  env-up            : Start the docker-compose environment."
	@echo "  env-down          : Stop the docker-compose environment."
	@echo "  kitex-gen-%       : Generate Kitex service code for a specific service (e.g., make kitex-gen-user)."
	@echo "  kitex-update-%    : Update Kitex generated code for a specific service (e.g., make kitex-update-user)."
	@echo "  hertz-gen-api     : Generate Hertz scaffold based on the API IDL."
	@echo "  test              : Run unit tests for the project."
	@echo "  clean             : Remove the 'output' directories and related binaries."
	@echo "  clean-all         : Stop docker-compose services if running and remove 'output' directories and docker data."
	@echo "  push-%            : Push a specific service to the remote repository (e.g., make push-api)."
	@echo "  fmt               : Format the codebase using gofumpt."
	@echo "  import            : Optimize import order and structure."
	@echo "  vet               : Check for possible errors with go vet."
	@echo "  lint              : Run golangci-lint on the codebase."
	@echo "  verify            : Format, optimize imports, and run linters and vet on the codebase."
	@echo "  license           : Check and add license to go file and shell script."

## --------------------------------------
## 构建与调试
## --------------------------------------

# 启动必要的环境，比如 etcd、mysql
.PHONY: env-up
env-up:
	@ docker compose -f ./docker/docker-compose.yml up -d

# 关闭必要的环境，但不清理 data（位于 docker/data 目录中）
.PHONY: env-down
env-down:
	@ cd ./docker && docker compose down

# 基于 idl 生成相关的 go 语言描述文件
.PHONY: kitex-gen-%
kitex-gen-%:
	@ kitex -module "${MODULE}" \
		-thrift no_default_serdes \
		${IDL_PATH}/$*.thrift
	@ go mod tidy

# 生成基于 Hertz 的脚手架
.PHONY: hz-%
hz-%:
	hz update -idl ${IDL_PATH}/api/$*.thrift

# 单元测试
# -gcflags="all=-l -N": -l 表示禁用内联优化，-N 表示禁用优化
# -parallel=16: 可以并行运行的测试数量，这里设置为 16
# -p=16: 指定并行构建的最大数量，这里设置为 16
# -covermode=atomic: 设置覆模式为原子模式
# -race: 启用竞态检测，检查并发代码中的数据竞争问题
# 我们通过`go list`来列出所有的包，然后通过`grep`来过滤掉不需要测试的包
.PHONY: test
test:
	go test -v -gcflags="all=-l -N" -coverprofile=coverage.txt -parallel=16 -p=16 -covermode=atomic -race -coverpkg=./... \
		`go list ./... | grep -E -v "kitex_gen|.github|idl|docs|config|deploy|docker"`

# 构建指定对象，构建后在没有给 BUILD_ONLY 参的情况下会自动运行，需要熟悉 tmux 环境
# 用于本地调试
.PHONY: $(SERVICES)
$(SERVICES):
	@if [ -z "$(TMUX_EXISTS)" ]; then \
		echo "$(PREFIX) tmux is not installed. Please install tmux first."; \
		exit 1; \
	fi
	@if [ -z "$$TMUX" ]; then \
		echo "$(PREFIX) you are not in tmux, press ENTER to start tmux environment."; \
		read -r; \
		if tmux has-session -t fzuhelp 2>/dev/null; then \
			echo "$(PREFIX) Tmux session 'fzuhelp' already exists. Attaching to session and running command."; \
			tmux attach-session -t fzuhelp; \
			tmux send-keys -t fzuhelp "make $(service)" C-m; \
		else \
			echo "$(PREFIX) No tmux session found. Creating a new session."; \
			tmux new-session -s fzuhelp "make $(service); $$SHELL"; \
		fi; \
	else \
		echo "$(PREFIX) Build $(service) target..."; \
		mkdir -p output; \
		bash $(DIR)/docker/script/build.sh $(service); \
		echo "$(PREFIX) Build $(service) target completed"; \
	fi
ifndef BUILD_ONLY
	@echo "$(PREFIX) Automatic run server"
	@if tmux list-windows -F '#{window_name}' | grep -q "^domtok-$(service)$$"; then \
		echo "$(PREFIX) Window 'domtok-$(service)' already exists. Reusing the window."; \
		tmux select-window -t "domtok-$(service)"; \
	else \
		echo "$(PREFIX) Window 'domtok-$(service)' does not exist. Creating a new window."; \
		tmux new-window -n "domtok-$(service)"; \
		tmux split-window -h ; \
		tmux select-layout -t "domtok-$(service)" even-horizontal; \
	fi
	@echo "$(PREFIX) Running $(service) service in tmux..."
	@tmux send-keys -t domtok-$(service).0 'export SERVICE=$(service) && bash ./docker/script/entrypoint.sh' C-m
	@tmux select-pane -t domtok-$(service).1
endif

# 推送到镜像服务中，需要提前 docker login，否则会推送失败
# 不设置同时推送全部服务，这是一个非常危险的操作
.PHONY: push-%
push-%:
	@read -p "Confirm service name to push (type '$*' to confirm): " CONFIRM_SERVICE; \
	if [ "$$CONFIRM_SERVICE" != "$*" ]; then \
		echo "Confirmation failed. Expected '$*', but got '$$CONFIRM_SERVICE'."; \
		exit 1; \
	fi; \
	if echo "$(SERVICES)" | grep -wq "$*"; then \
		if [ "$(ARCH)" = "x86_64" ] || [ "$(ARCH)" = "amd64" ]; then \
			echo "Building and pushing $* for amd64 architecture..."; \
			docker build --build-arg SERVICE=$* -t $(REMOTE_REPOSITORY):$* -f docker/Dockerfile .; \
			docker push $(REMOTE_REPOSITORY):$*; \
		else \
			echo "Building and pushing $* using buildx for amd64 architecture..."; \
			docker buildx build --platform linux/amd64 --build-arg SERVICE=$* -t $(REMOTE_REPOSITORY):$* -f docker/Dockerfile --push .; \
		fi; \
	else \
		echo "Service '$*' is not a valid service. Available: [$(SERVICES)]"; \
		exit 1; \
	fi
## --------------------------------------
## 清理与校验
## --------------------------------------

# 清除所有的构建产物
.PHONY: clean
clean:
	@find . -type d -name "output" -exec rm -rf {} + -print

# 清除所有构建产物、compose 环境和它的数据
.PHONY: clean-all
clean-all: clean
	@echo "$(PREFIX) Checking if docker-compose services are running..."
	@docker-compose -f ./docker/docker-compose.yml ps -q | grep '.' && docker-compose -f ./docker/docker-compose.yml down || echo "$(PREFIX) No services are running."
	@echo "$(PREFIX) Removing docker data..."
	rm -rf ./docker/data

# 格式化代码，我们使用 gofumpt，是 fmt 的严格超集
.PHONY: fmt
fmt:
	gofumpt -l -w .

# 优化 import 顺序结构
.PHONY: import
import:
	goimports -w -local github.com/west2-online .

# 检查可能的错误
.PHONY: vet
vet:
	go vet ./...

# 代码格式校验
.PHONY: lint
lint:
	golangci-lint run --config=./.golangci.yml --tests --allow-parallel-runners --sort-results --show-stats --print-resources-usage

# 检查依赖漏洞
.PHONY: vulncheck
vulncheck:
	govulncheck ./...

.PHONY: tidy
tidy:
	go mod tidy

# 一键修正规范并执行代码检查，同时运行 license 检查
.PHONY: verify
verify: license vet fmt import lint vulncheck tidy

# 补齐 license
.PHONY: license
license:
	sh ./hack/add-license.sh

# 为集成测试做准备
.PHONY: with-env-test
with-env-test:
	@ make env-up ;\
		echo "waiting for env up" ;\
		sleep 3 ;\
		export $(EnvironmentStartEnv)=$(EnvironmentStartFlag) ;\
		make test

.PHONY: with-env-test-nowait
with-env-test-nowait:
	@ make env-up ;\
    	export $(EnvironmentStartEnv)=$(EnvironmentStartFlag) ;\
    	make test

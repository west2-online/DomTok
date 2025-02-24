
# 该脚本适用于 Docker 及本地调试，作为执行程序前的 presetting（即 entrypoint）
# 请不要直接执行这个脚本，这个脚本应当由 Makefile/Dockerfile 接管

#! /usr/bin/env bash
CURDIR=$(pwd)

# 此处只涉及 Kitex，但是 Hertz 使用这个没有影响，保留即可
export KITEX_RUNTIME_ROOT=$CURDIR
export KITEX_LOG_DIR="$CURDIR/log"

if [ ! -d "$KITEX_LOG_DIR/app" ]; then
    mkdir -p "$KITEX_LOG_DIR/app"
fi

if [ ! -d "$KITEX_LOG_DIR/rpc" ]; then
    mkdir -p "$KITEX_LOG_DIR/rpc"
fi

# 参数替换，检查 ETCD_ADDR 是否已经设置，没有将会设置默认值
: ${ETCD_ADDR:="localhost:2379"}
export ETCD_ADDR

# 这个 SERVICE 环境变量会自动地由 Dockerfile/Makefile 设置
exec "$CURDIR/output/$SERVICE/domtok-$SERVICE"

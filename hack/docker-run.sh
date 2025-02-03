# Copyright 2024 The west2-online Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/bin/bash

# 这个脚本用于管理一组服务的Docker容器。
# 它允许使用单个Docker镜像启动特定服务或所有服务。

CONFIG_PATH="../config/config.yaml" # related to project folder

get_port() {
    local server_name="$1"

    port=$(awk -v name="$server_name" '
        $0 ~ "name: " name {found=1}
        found && $0 ~ "port:" {print $2; exit}
        $0 ~ "services:" {found=0}
    ' "$CONFIG_PATH")


    echo "$port"
}

# Docker容器的镜像名称
IMAGE_NAME="registry.cn-hangzhou.aliyuncs.com/west2-online/fzuhelper-server"

# 要启动的服务，默认为 "all" 如果没有提供参数
SERVICE_TO_START=${1:-all}

# 脚本所在的目录
DIR=$(cd "$(dirname "$0")" && pwd)

# 可用服务的列表。在真实场景中，这应该是服务名称的数组。
# 例如：SERVICES=("service1" "service2" "service3")
SERVICES=(api user classroom course launch_screen paper academic version common)
# 1. 编译镜像时将多个服务打包为单一镜像（当然不建议这么做，不过镜像小的且不需要频繁更新这么做很方便）
# 2. 启动时根据SERVICE_TO_START在SERVICES中查找对应服务名
# 3. 如果查找的到，则启动容器，查找不到则抛出错误

# 删除Docker容器的函数
remove_container() {
    local container_name="$1"
    local container_status=$(docker inspect -f '{{.State.Status}}' "$container_name" 2>/dev/null)
    echo "remove container $container_name"

    docker stop "$container_name"

    docker rm "$container_name"

}

# 启动新Docker容器的函数
start_container() {
    # 启动容器前先删除旧容器
    remove_container "$1"
    # 之后更新 image
    sh image-refresh.sh "$1"

    local service_name="$1"
    echo "serverice_name is ${service_name}"
    local server_port=$(get_port "$service_name")
    local image="$IMAGE_NAME:$service_name"

    echo "port is $server_port"
    echo "start container $service_name"
    docker run -d --name $service_name \
    --network fzu-helper \
    -p $server_port:$server_port \
    -e ETCD_ADDR="fzu-helper-etcd:2379" \
    --restart always \
    $image
}

# 启动新容器
if [ "$SERVICE_TO_START" == "all" ]; then
    for service in "${SERVICES[@]}"; do
        start_container "$service"
    done
else
    start_container "$SERVICE_TO_START"
fi

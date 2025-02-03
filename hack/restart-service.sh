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

# This script file will automatically pull the latest image (if available), exit, delete the original image and container, and finally launch a new container.
# Usually, we only need to modify the initial few configurations, and the rest of the content does not need to be changed.

IMAGE_NAME="registry.cn-hangzhou.aliyuncs.com/west2-online/fzuhelper-server:latest"
CONTAINER_NAME_PREFIX="fzuhelper"
NET_MODE="host"

DIR=$(cd $(dirname $0); pwd)
CONFIG_PATH=$DIR/config
CONTAINER_CONFIG_PATH=/app/config

SERVICE_TO_START=${1:-all} # default start all

SERVICES=(api classroom user)

remove_container() {
    container_status=$(docker inspect -f '{{.State.Status}}' "$1")
    if [ "$container_status" == "running" ]; then
        echo "Stopping container $1..."
        docker stop "$1"
    elif [ "$container_status" == "paused" ]; then
        echo "Unpausing and then stopping container $1..."
        docker unpause "$1"
        docker stop "$1"
    fi
    echo "Remove container $1..."
    docker rm "$1"
}

start_container() {
    echo "Starting container for $1..."
    docker run -d --name "$CONTAINER_NAME_PREFIX-$1" \
    -e service=$1 \
    --net=$NET_MODE \
    -v $CONFIG_PATH:$CONTAINER_CONFIG_PATH \
    "$IMAGE_NAME"
}

containers_to_stop=$(docker ps -aq --filter "ancestor=$IMAGE_NAME")
if [ "$SERVICE_TO_START" == "all" ]; then
    for container_id in $containers_to_stop; do
        remove_container $container_id
    done
else
    for container_id in $containers_to_stop; do
        container_id=$(docker inspect -f '{{.Name}}' "$container_id")
        if [ "$container_id" != "/$CONTAINER_NAME_PREFIX-$SERVICE_TO_START" ]; then
            continue
        fi
        remove_container $container_id
    done
fi

echo "Pulling the latest image..."
docker pull "$IMAGE_NAME"

if [ "$SERVICE_TO_START" == "all" ]; then
    for service in "${SERVICES[@]}"; do
        start_container $service
    done
else
    start_container $SERVICE_TO_START
fi

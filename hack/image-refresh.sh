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
# 使用脚本前请保证docker login

# 镜像名称和标签，请根据实际情况进行替换
IMAGE_NAME="registry.cn-hangzhou.aliyuncs.com/west2-online/fzuhelper-server"
IMAGE_TAG="$1"
FULL_IMAGE_NAME="${IMAGE_NAME}:${IMAGE_TAG}"

# 检查jq是否安装
if ! type jq > /dev/null; then
    echo "错误：需要jq工具来解析JSON。请安装jq。 e.g: yum install jq"
    exit 1
fi

# 获取远程镜像的摘要信息
get_remote_image_digest() {
    docker manifest inspect "$FULL_IMAGE_NAME" 2>/dev/null | jq -r '.config.digest'
}

# 获取本地镜像的摘要信息
get_local_image_digest() {
    docker image inspect --format='{{.Id}}' "$FULL_IMAGE_NAME" 2>/dev/null
}

# 拉取最新的Docker镜像
pull_new_image() {
    echo "正在拉取最新镜像: $FULL_IMAGE_NAME..."
    if docker pull "$FULL_IMAGE_NAME"; then
        echo "镜像更新成功。"
    else
        echo "错误：无法拉取镜像。"
        exit 1
    fi
}

# 清理tag为<none>的镜像
function clean_images() {
    # 找到所有 tag 为 <none> 的镜像
    docker images | grep "<none>" | awk '{print $3}' | while read image_id; do
        # 删除找到的镜像
        docker rmi -f $image_id
    done

    echo "成功清理旧镜像"
}

# 主流程
main() {
    local remote_digest=$(get_remote_image_digest)
    if [ -z "$remote_digest" ]; then
        echo "错误：无法获取远程镜像的摘要信息。"
        exit 1
    fi

    local local_digest=$(get_local_image_digest)
    if [ -z "$local_digest" ]; then
        echo "警告：未找到本地镜像 $FULL_IMAGE_NAME 或无法获取摘要。"
        echo "尝试直接拉取远程镜像..."
        pull_new_image
        exit 0
    fi

    echo "本地镜像摘要: $local_digest"
    echo "远程镜像摘要: $remote_digest"

    if [ "$local_digest" != "$remote_digest" ]; then
        echo "本地镜像与远程镜像不一致，需要更新。"
        pull_new_image
    else
        echo "本地镜像已是最新。"
    fi
    echo "开始清理旧镜像"
    clean_images
}

# 执行主流程
main


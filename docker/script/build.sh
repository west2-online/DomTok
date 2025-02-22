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

# 该脚本负责构建二进制文件
# 请不要直接执行这个脚本，这个脚本应当由 Makefile/Dockerfile 接管

#!/usr/bin/env bash
# Usage: ./build.sh {SERVICE}

RUN_NAME="$1"
ROOT_DIR=$(pwd) # 二进制文件将会编译至执行脚本时的目录

if [ -z "$RUN_NAME" ]; then
    echo "Error: Service name is required."
    exit 1
fi

# 进入模块列表
cd "./cmd/${RUN_NAME}" || exit

# 创建产物文件夹并提供权限
mkdir -p ${ROOT_DIR}/output/${RUN_NAME}

# 基于环境变量决策是构建还是测试
if [ "$IS_SYSTEM_TEST_ENV" != "1" ]; then
    go build -o ${ROOT_DIR}/output/${RUN_NAME}/domtok-${RUN_NAME}
else
    go test -c -covermode=set -o ${ROOT_DIR}/output/${RUN_NAME}/domtok-${RUN_NAME} -coverpkg=./...
fi

# 构造结果
# output
# └── {SERVICE}
#     ├── bin
#     │   └── domtok-{SERVICE}
#     └── entrypoint.sh

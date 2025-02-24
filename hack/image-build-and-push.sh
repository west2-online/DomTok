
#!/bin/bash
# 该脚本将在 github/workflow 中自动将项目进行 image 构建和推送

service="$1"
echo "Building and pushing image for service $service..."
echo "$service" | make push-$service


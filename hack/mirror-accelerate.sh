
#!/bin/bash

# This script is used to modify the Docker image source of the current machine to Aliyun, accelerating image pulling.

sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<-'EOF'
{
  "registry-mirrors": ["https://o3nc7upe.mirror.aliyuncs.com"]
}
EOF
sudo systemctl daemon-reload
sudo systemctl restart docker


#!/bin/bash

# This script is used to monitor local files.
# If there are any changes in the local files, they will be synchronized and updated to etcd.
# Then etcd will update the configuration on other servers.

# wait etcd complete
while true; do
  if etcdctl endpoint health; then
    echo "etcd is ready."
    break
  fi
  echo "waiting for etcd to start..."
  sleep 1
done

# upload config
etcdctl put /config -- < /config/config.yaml

# create backup
cp /config/config.yaml /config/config.yaml.bak


# continuous listen
previous_hash=$(sha256sum /config/config.yaml | awk '{print $1}')

while true; do
  current_hash=$(sha256sum /config/config.yaml | awk '{print $1}')

  if [ "$current_hash" != "$previous_hash" ]; then
    etcdctl put /config -- < /config/config.yaml
    echo "spot update, config updated in etcd. $(date +'%Y-%m-%d %H:%M:%S')"

    # diff
    echo "Detected changes in configuration:"
    diff /config/config.yaml.bak /config/config.yaml
    cp /config/config.yaml /config/config.yaml.bak
    previous_hash="$current_hash"
  fi

  sleep 60
done

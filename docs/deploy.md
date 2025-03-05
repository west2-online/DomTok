# Deploy

This guide will deploy this project from scratch.

During the deployment process, we will use the commands in the Makefile to start the dependent environments (including databases and cache services such as MySQL, etcd, and Redis), and then compile and run specific services.

If you are interested in the specific build and startup process, please refer to the [Build](./build.md) page.

## Prerequisites

Packages that need to be installed:

- Docker
    - Please refer to the [Docker official documentation](https://docs.docker.com/get-docker/) to install Docker.
    - It is used to start services such as MySQL, etcd, and Redis.
- tmux
    - It is convenient to manage multiple sessions in the terminal and improve the operation and maintenance efficiency.
    - The commands in the Makefile will be executed in tmux by default.

## Local Deployment

### Preliminary Operations

Modify the configuration in `config/config.yaml` and change the IP of database configurations to `localhost` (if the file doesn't exist, please create it).  
We use OSS to store image data, so you also need to configure the relevant information of OSS (Upyun) in `config/config.yaml`.  
For security reasons, you need to modify the key part in `config/config.yaml` and set it to your own private key and public key.  
Please refer to `config.example.yaml` for the configuration example.

### Start the Environment

#### Clean the local environment (optional)

```shell
make clean-all
```

#### Start the basic containers of the environment (databases, etc.)

```shell
make env-up
```

### Start Specific Services

```shell
make <target>
```

`<target>` is the service name, which is the folder name under the `cmd` directory.

The directory structure after completion should be similar to the following:

```shell
.
├── docker
│   ├── script
│   │   └── etcd-monitor.sh
│   ├── env
│   │   ├── redis.env
│   │   ├── mysql.env
│   │   └── etcd.env
│   ├── docker-compose.yaml
├── config
│   ├── sql
│   │   └── init.sql
│   └── config.yaml
└── hack
    │── image-refresh.sh
    └── docker-run.sh
```

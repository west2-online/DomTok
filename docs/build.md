# Build

This guide will mainly introduce the project's building and startup processes, as well as related scripts.

Subsequent `<target>` refers to a certain service, such as `api`. You can specifically obtain the list of buildable services through `make help`.

## General Process

1. Use the command `make env-up` to start the environment (MySQL, etcd, Redis, etc.).
2. `make <target>` compiles and runs the specific service.
3. The service retrieves `config.yaml` from etcd.
4. Reads the configuration in `config.yaml` and maps the Env to the corresponding **struct**.
5. Obtains the available address from `config.yaml`.
6. Initializes the service and registers the service in `etcd`.
7. Starts the service.

```mermaid
sequenceDiagram
    participant User as User
    participant Makefile as Makefile
    participant Services as Service
    participant Env as Environment (MySQL, etcd, Redis)

    User->>Makefile: Run `make env-up`
    Makefile->>Env: Start MySQL, Redis, and etcd

    User->>Makefile: Run `make <target>`
    Makefile->>Services: Compile and start the specified service

    Services->>Env: Get `config.yaml` from etcd
    Env-->>Services: Return `config.yaml`

    Services->>Services: Read `config.yaml` and map it to the struct, and obtain the available IP address

    Services->>Env: Register the service in etcd
    Services->>Services: Initialization completed, start the service
```

## Building and Starting

### Directory Structure

The key directories of the project are as follows:

- `cmd/`: Contains the startup entry points of each service module.
- `output/`: The output directory for build products.

### Building Process

Here we explain the specific workflow when we type `make <target>`. We omit the content related to the tmux environment.

The building process is mainly completed through the [build.sh](../docker/script/build.sh) script, which is used to compile the binary files of the specified service module or conduct system tests:

1. Enter the corresponding service folder in `cmd`.
2. Execute `go build` to compile the binary file of this service and store it in the `output` folder.

```mermaid
flowchart TD
    A[Start script] --> B{Check if a parameter, i.e., target, is input}
    B --> |Empty| C[Output an error message and exit]
    B --> |Not empty| D[Set ROOT_DIR as the current working directory]

    D --> E[Enter the specified module directory ./cmd/RUN_NAME]
    E --> F[Create the folder output/RUN_NAME and set permissions]

    F --> G{Determine if it is a test environment}
    G --> |It is a test environment| H[Execute the test build]
    G --> |Not a test environment| I[Execute the build]

    H --> J[Generate the test binary file output/RUN_NAME/fzuhelper-RUN_NAME]
    I --> K[Generate the build binary file output/RUN_NAME/fzuhelper-RUN_NAME]
```

### Output Directory Structure

```text
 output
 └── target
        └── binary
```

### Startup Process

When we type `make <target>` without setting the build-only flag (`BUILD_ONLY`), it will start automatically. Here we introduce the process of local debugging startup.

> The startup process of Docker containers is similar, except that it is moved into the container.

The startup process is mainly completed through the [entrypoint.sh](/docker/script/entrypoint.sh) script.

1. Use `export` to set the environment variable of the etcd address, so that the subsequent program can obtain the etcd address during **runtime** and get `config.yaml`.
2. `cd` to the `output` directory generated during the build stage and execute the binary file of the corresponding service.

```mermaid
flowchart TD
    A[Start entrypoint.sh] --> B{Check if ETCD_ADDR is set}
    B --> |Not set| C[Set the default ETCD_ADDR=localhost:2379]
    B --> |Set| D[Keep the existing ETCD_ADDR]

    C --> E[Export the ETCD_ADDR environment variable]
    D --> E

    E --> F[Start the service]
```

## Usage

Both scripts are managed by the commands in `Makefile` and can be called through the following commands:

```shell
make <target> [option]   # option = BUILD_ONLY
```

The following is a rough flowchart of `make <target>`:

```mermaid
flowchart TD
    A[Start the make <target> command] --> B{Check if the BUILD_ONLY setting is passed in}

    B -- Not set --> C[Build and run]
    B -- Set --> D[Only build]

    D --> E[Create the output directory]
    E --> F[Run the build.sh script for compilation]

    C --> G[Create the output directory]
    G --> H[Run the build.sh script for compilation]
    H --> I[Run the entrypoint.sh to start the service]
    I --> J[The service is successfully started]
```

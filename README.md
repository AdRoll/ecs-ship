# :ship::it: ecs-ship

Yet another implementation of the ecs-deploy tool. `ecs-ship` will allow you to
deploy changes to your ECS services by patching its current configuration.

Here are the steps that `ecs-ship` does for you:

1. Grab the configuration from the `cluster` and `service` you specify.
2. Apply the patches you can specify from the command line.
3. Create a task definition out of the patched configuration.
4. Deploy the service with the new task definition.
5. Wait for stability.

## Usage

```
NAME:
   ecs-ship - Deploy your aws ecs services.

USAGE:
   ecs-deploy [options] <cluster> <service>

VERSION:
   2.0.0

GLOBAL OPTIONS:
   --updates FILE, -u FILE          Use an input FILE to describe service updates (default: stdin)
   --timeout DURATION, -t DURATION  Wait this DURATION for the service to be correctly updated (default: 5m0s)
   --no-color, -n                   Disable colored output (default: false)
   --no-wait, -w                    Disable waiting for updates to be completed. (default: false)
   --dry, -d                        Don't deploy just show what would change in the remote service (default: false)
   --help, -h                       show help
   --version, -v                    print the version
```

For the input file you can use this yaml schema:

```yml
cpu: "95"
memory: "34"
containerDefinitions:
  someContainer:
    cpu: 95
    environment:
      NAME: value
    image: "some-image:someTag"
    memory: 34
    memoryReservation: 34
```

**Notice** that every part of the input is optional, so the idea is that you
just pass in the values that you need.

## Getting `ecs-ship`

You can grab ssh ship from [DockerHub][docker-hub] from the repository
[nextroll/ecs-ship][docker-repo] or from the console by:

```bash
docker pull nextroll/ecs-ship
docker run --rm nextroll/ecs-ship ecs-ship --help
```

You can also grab a Linux staticly linked binary from the [releases
page][releases] and drop it in your environment.

Finally you can grab the source code either from the [AdRoll/ecs-ship][repo]
repo or from the releases page and build yourself a fresh version using go.

```bash
go build .
```

## Examples

Here's an example for updating a container image to the latest version:

```bash
cat << EOF | ecs-ship cluster service
containerDefinitions:
  someContainer:
    image: "some-image:latestTag"
EOF
```

Here's an example where you just change an environment variable of a sevice:

```bash
cat << EOF | ecs-ship cluster service
containerDefinitions:
  someContainer:
    environment:
      SECRETS_BUCKET: "s3://updated-secrets-bucket"
EOF
```

Here's yet ahother example where you just want to lower the cpu requirements of
your service to lower your costs:

```bash
cat << EOF | ecs-ship cluster service
cpu: "50"
containerDefinitions:
  someContainer:
    cpu: 50
```

## Development

We provide some make commands for your development convenience.

```shell
make mod      # load dependencies
make lint     # check code
make test     # run unit tests
make build    # build app (see artifacts/ecs-deploy)
make mockgen  # rebuild mocked entities
```

[docker-hub]: https://hub.docker.com/r/nextroll/ecs-ship
[docker-repo]: https://hub.docker.com/r/nextroll/ecs-ship
[releases]: https://github.com/AdRoll/ecs-ship/releases
[repo]: https://github.com/AdRoll/ecs-ship

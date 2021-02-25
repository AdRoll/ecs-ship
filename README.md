# :ship::it: ecs-ship

Yet another implementation of the ecs-deploy tool. `ecs-ship` will allow you
to deploy changes to your ECS services by patching its current configuration.

Here are the steps that `ecs-ship` does for you:

1. Grab the configuration from the `cluster` and `service` you specify.
2. Apply the patches you can specify from the command line.
3. Create a task definition out of the patched configuration.
4. Deploy the service with the new task definition, wait for stability and if
   it's not reached roll back.

## Usage

```
NAME:
   ecs-ship - Deploy your aws ecs services.

USAGE:
   ecs-deploy [options] <cluster> <service>

VERSION:
   0.1.0

GLOBAL OPTIONS:
   --updates FILE, -u FILE          Use an input FILE to describe service updates (default: stdin)
   --timeout DURATION, -t DURATION  Wait this DURATION for the service to be correctly updated (default: 5m0s)
   --no-color, -n                   Disable colored output (default: false)
   --help, -h                       show help (default: false)
   --version, -v                    print the version (default: false)
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

Here's yet ahother example where you just want to lower the cpu requirements
of your service to lower your costs:

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

# Velocity

## Overview

Velocity is a simple test running service. It is designed to be self hosted, easy to understand, and easy to extend.

### Units

The main units of the system are:

-   Routines: A collection of jobs
-   Jobs: A collection of tests and images to run those tests on
-   Images: A description of a docker image
-   Tests: A collection of steps to run or auto-generated steps with a given language/library/framework

A minimal example of a complete configuration would be:

```yaml
tests:
    - name: Self test
      language: golang

images:
    - name: Latest supported golang
      image: golang:1.23-rc-bookworm
    - name: Oldest supported golang
      image: golang:1.22

jobs:
    - name: Self test
      tests:
          - Self test
      images:
          - Latest supported golang
          - Oldest supported golang

routines:
    - name: Full coverage
      jobs:
          - Self test
```

This configuration would expose a routine called "Full coverage" that would run `go test ./...` on the latest and oldest supported golang images.

## Development

### Running the project

The project is ran using [docker compose](https://docs.docker.com/compose/). To start the project in detached mode, run the following command:

```bash
make dev
```

All docker commands are wrapped in the `make` command. For example:

```bash
# To stop the containers
make dev-stop
# To start the containers
make dev-start
# Deletes the containers
make dev-down
```

To create and start the containers, run:

```bash
make dev
```

When you start the services in dev mode, some services might take a little to fully start up (like kafka). If you see errors, try waiting 1-2 minutes to make sure that all the services are stable.

These commands are also available for the prod environment (which pulls from the .env file).

### Developing on the project

Once you have it running, you can access the API at `http://localhost:8080`, the CLI at `bin/velocity`, and the mongo database at `mongodb://localhost:27017/?directConnection=true`.

For example, you could start a routine like:

```bash
bin/velocity routine run <routine-name>
```

Which would kick off a routine and eventually the agent will start the tests for the routine.

# Velocity

Welcome to your own CI/CD system! This is suppose to be a basic example in which you can run a series of tests defined in [velocity.yml](velocity.yml) either locally or pushed to a database to be ran on an agent. The main components of this system are:

-   [CLI](internal/operations/operations.go): User-facing tool that talks to the API

    -   Use cases:

        -   `velocity run-local`: Runs tests locally
        -   `velocity run`: Posts tests to the API
        -   `velocity validate`: Validates the velocity.yml file

    -   Build:

        -   `make build-cli` creates the cli at $GOPATH/bin/velocity

    -   Architectural Decision Record (ADR):

        -   The CLI reads off of the velocity.yml the CWD holds (overriden by flags)
        -   When posting to the API, attempts to authentication with credentials at ~/.velocityrc

-   [Client](client) (TBA): A UI that displays the status of the tests

    -   TBA

-   [Agent](internal/agent/agent.go): Runs tests and reports back to the API

    -   Build:

        -   `make agent` runs the agent with a velocity provider
        -   `make agent-mongodb` runs the agent with a mongodb provider

    -   ADR:
        -   The agent can have different [providers](internal/jobs/provider.go), like:
            -   `InMemoryJobProvider`: Runs jobs that are in memory of the system. This is a finite number of jobs- great for local runs
            -   `VelocityJobProvider`: Queries your API for jobs and posts them to the API. Great for distributed runs
            -   `MongoDBJobProvider`: Directly queries MongoDB for jobs and posts them to MongoDB. Previews flexibility of providers, can be used in proof of concept

-   [API v1](internal/api/v1/v1.go): Handles requests from the CLI, Agent, and Client

    -   Build:

        -   `make server` starts the server

    -   ADR:

        -   The API is built with versioning and has a contract with these [types](src/clients/v1types)

-   [Config](src/config/types.go): The [velocity.yml](velocity.yml) is opinionated in how you should architect your tests

    -   `config`:

        -   `project`: The project id of the repository
        -   `registry`: The docker registry to pull your images from
        -   `server`: The velocity server
        -   `ui`: The ui server

    -   `tests`:

        -   [test-name] (using language + framework):

            -   `language`: [Supported](internal/jobs/defaults.go)
            -   `framework`: [Supported](internal/jobs/defaults.go)

        -   [test-name] (using bash script)

            -   `run`: A shell command
            -   `image`: The docker image to run with

        -   [test-name] (general options)

            -   `directory`: The directory to run the test in

    -   `images`:

        -   [image-name]:

            -   `image`: The docker image to build

    -   `builds`:

        -   [build-name]:

            -   `input`:

                -   `image`: The image (specified above) to build
                -   `directory`: The directory to run the build in
                -   `run`: A shell command to run that creates the build
                -   `path`: The path to the build artifact

            -   `output` (to local or remote):

                -   `path`: A local path, if wanted
                -   `url`: The url to post the build artifact to
                -   `method`: The method to post the build artifact with
                -   `headers`: The headers to post the build artifact with

    -   `workflows`:

        -   [workflow-name]:

            -   `tests`: An array of the tests to run

---
sidebar_position: 1
---

# Tests

The building blocks of your configuration. They are a series of commands that are executed in a specific order. If any command fails, the test fails.

Tests are a wrapper around a list of commands, there are different tyypes of tests usable in a test:

-   [Shell command](#shell-command)
    -   A shell command
-   [Prebuilt command](#prebuilt-command)
    -   A command that is built in to the platform
    -   A list of prebuilt commands can be found [here](#prebuilt-commands)
-   [Operation](#operation)
    -   A user-defined command that can be used in multiple tests
    -   This is used to create reusable components that can be used in multiple tests

Tests are defined in the `test` section. Every test has the following structure (an example of the commands as well):

```yaml title="velocity.yml"
tests:
    - name: Test name
      commands:
          - prebuilt: git.clone
          - command: npm run test
            working_dir: ./backend #optional
            env: #optional
                - NODE_ENV=test
                - PORT=3000
          - operation: my-operation
```

Here might be the operation definition:

```yaml title="velocity.yml"
tests:
    - name: Test name
      commands:
          - prebuilt: git.clone
          - command: npm run test
            working_dir: ./backend #optional
            env: #optional
                - NODE_ENV=test
                - PORT=3000
          - operation: my-operation
```

Here might be the operation definition:
operations: - name: my-operation
commands: - command: npm run lint - command: npm run prechecks

## All commands

You can include `env`, `working_dir` in every command, they are optional.

## Shell command

A shell command dependant on the image you provide. Images are provided where you use the tests, view [images](./images) for more information. For example, if the image the test should run on does not include git, then the test will fail to find the command. Prebuilt commands are recommended over shell commands when possible.

## Prebuilt command

A command that is built in to the platform. It performs the task as best as possible, which may not be possible via a shell command. For example, the `git.clone` command will use a separate image to clone the repository and then copy the files to the image the test is running on. This is then cached and used for other tests that use the same git.clone command.

## Operation

A user-defined command that is a list of other commands, these commands can be of any type.

## Prebuilt commands

Prebuilt commands are separated in to packages.

### git

-   #### git.clone
    -   Description: Clones a git repository
    -   Params:
        -   repository
            -   Description: The git repository to clone
            -   Type: string
            -   Default: Current repository
            -   Format: `username/repository` for github or `https://....`
        -   branch
            -   Description: The branch to clone
            -   Type: string
            -   Default: main branch
        -   directory
            -   Description: The directory to clone the repository to
            -   Type: string
            -   Default: Current directory

### request

-   #### request.[method]
    -   Description: Makes a [method] (get, post, put, delete) request to a URL
    -   Params:
        -   url
            -   Description: The URL to make the request to
            -   Type: string
            -   Default: None
        -   body
            -   Description: The body of the request
            -   Type: string
            -   Default: None
        -   headers
            -   Description: The headers of the request
            -   Type: array
            -   Default: None
            -   Format: Array of objects with the following keys
                -   key
                    -   Description: The key of the header
                    -   Type: string
                    -   Required
                -   value
                    -   Description: The value of the header
                    -   Type: string
                    -   Required
        -   timeout
            -   Description: The timeout of the request
            -   Type: number
            -   Default: 5000
        -   response
            -   Description: The response of the request to validate against, if needed
            -   Type: object
            -   Default: None
            -   Format:
                -   status
                    -   Description: The status of the response
                    -   Type: number
                    -   Default: None
                -   body
                    -   Description: The body of the response
                    -   Type: object
                    -   Default: None
                    -   Format:
                        -   type
                            -   Description: The type of the body
                            -   Type: string
                            -   Default: application/json
                            -   Format: application/json, text/html, etc.
                        -   values
                            -   Description: The value of the body
                            -   Type: object
                            -   Default: None
                            -   Format: Object of expected values

### build

-   #### build.[build-name]
    -   Description: Builds [build-name] and drops the output in the cwd
    -   Params:
        -   directory
            -   Description: The directory to place the output in
            -   Type: string
            -   Default: Current directory
        -   env
            -   Description: The environment variables to use
            -   Type: array
            -   Default: None
            -   Format: Array of strings in the format `key=value`

### run

-   #### run.[build-name]
    -   Description: Runs [build-name] in a separate container. This has to be ran after the build.[build-name] command
    -   Params:
        -   host:
            -   Description: A label used to route requests to
            -   Type: string
            -   Default: localhost
            -   Example: `host: app` would route requests from `http://app` to this app
        -   env
            -   Description: The environment variables to use
            -   Type: array
            -   Default: None
            -   Format: Array of strings in the format `key=value`

### build-run

-   #### build-run.[build-name]
    -   Description: Builds and runs [build-name]
    -   Params:
        -   directory
            -   Description: The directory to place the output in
            -   Type: string
            -   Default: Current directory
        -   build-env
            -   Description: The environment variables to use when building
            -   Type: array
            -   Default: None
            -   Format: Array of strings in the format `key=value`
        -   run-env
            -   Description: The environment variables to use when running
            -   Type: array
            -   Default: None
            -   Format: Array of strings in the format `key=value`

### command

-   #### command.tar
    -   Description: Tars a directory
    -   Params:
        -   directory
            -   Description: The directory to tar
            -   Type: string
            -   Default: None
        -   output
            -   Description: The output file
            -   Type: string
            -   Default: None
-   #### command.gzip
    -   Description: Gzips a file
    -   Params:
        -   file
            -   Description: The file to gzip
            -   Type: string
            -   Default: None
        -   output
            -   Description: The output file
            -   Type: string
            -   Default: None
-   #### command.targzip
    -   Description: Tar and gzip a directory
    -   Params:
        -   directory
            -   Description: The directory to tar and then gzip
            -   Type: string
            -   Default: None
        -   output
            -   Description: The output file
            -   Type: string
            -   Default: None

```

```

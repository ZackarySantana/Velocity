# Velocity

Velocity is a self-hosted testing platform.

## Table of Contents

-   [Overview](#overview)
    -   [Concepts](#concepts)
    -   [Images](#images)
        -   [Example](#example)
    -   [Tests](#tests)
        -   [Example](#example-1)
    -   [Builds](#builds)
        -   [Example](#example-2)
    -   [Deployments](#deployments)
        -   [Example](#example-3)
    -   [Workflows](#workflows)
        -   [Example](#example-4)
    -   [Config](#config)
        -   [Example](#example-5)
    -   [Prebuilt commands](#prebuilt-commands)
-   [Installation](#installation)
-   [Components](#components)
    -   [CLI](#cli)
    -   [Client](#client)
    -   [API](#api)
    -   [Agent](#agent)

## Overview

[Velocity](https://velocity-ci.com) is a self-contained testing platform that's designed to be plug-and-play. After understanding the main concepts, you should be able to deploy your own version, integrate it in to your developer workflow, CI pipeline, and deployment process.

### Concepts

Velocity starts in the 'velocity.yml' file. This file contains:

-   [Images](#images)
    -   Docker images that are used in tests and deployments
-   [Tests](#tests)
    -   Units of process that can be used in a workflow
-   [Deployments](#deployments)
    -   The process of building your app
-   [Workflows](#workflows)
    -   A series of tests that are ran
-   [Config](#config)
    -   Options to connect your local environment to the platform

### Images

Images are the building blocks of your tests. They are Docker image aliases with additional options that you can use throughout your tests and deployments.

#### Example

```yaml
images:
    - name: node
      image: node:latest
      env:
          - NODE_ENV=development
          - PORT=3000
```

### Tests

Tests are the building blocks of your workflows. They are a series of commands that can be ran. If any of the commands fail, the test fails.

There are a couple different types of commands that can be used in a test.

-   [Shell command](#shell-command)
    -   A shell command that is ran in the container
    -   This can only run commands that are available in the used image in the workflow
-   [Prebuilt command](#prebuilt-command)
    -   A command that is built in to the platform
    -   A list of prebuilt commands can be found [here](#prebuilt-commands)

#### Example

##### Shell command

```yaml
tests:
    - name: test
      commands:
          - command: npm run test
            env:
                - NODE_ENV=test
                - PORT=3000
```

##### Prebuilt command

```yaml
tests:
    - name: test
      commands:
          - prebuilt: git.clone
```

### Builds

#### Example

```yaml
builds:
    - name: build
      build_image: node
      runtime_image: node
      output: dist
      run: node index.js
      commands:
          - command: npm run build
            env:
                - NODE_ENV=production
                - PORT=3000
```

### Deployments

#### Example

```yaml
builds:
    - name: app
      build_image: node
      runtime_image: node
      output: dist
      run: node index.js
      commands:
          - command: npm run build
            env:
                - NODE_ENV=production
                - PORT=3000

deployments:
    - name: deploy
      image: node
      commands:
          - prebuilt: build.app
          - prebuilt: command.tar
            params:
                directory: dist
                output: output.tar
          - prebuilt: command.gzip
            params:
                input: output.tar
                output: output.tar.gz
          - prebuilt: request.post
            params:
                url: https://api.example.com/deploy
                body: output.tar.gz
                headers:
                    - key: Content-Type
                      value: application/gzip
                timeout: 10000
```

### Workflows

#### Example

```yaml
TBA
```

### Config

#### Example

```yaml
TBA
```

### Prebuilt commands

#### git

-   git.clone
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

#### request

-   request.[method] (get, post, put, delete)
    -   Description: Makes a [method] request to a URL
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

#### build

-   build.[build-name]
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

#### run

-   run.[build-name]
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

#### build-run

-   build-run.[build-name]
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

### Examples

#### Test some endpoints

```yaml
builds:
    - name: app
      build_image: node
      runtime_image: node
      output: dist
      run: node index.js
      commands:
          - command: npm run build
            env:
                - NODE_ENV=production
                - PORT=3000

tests:
    - name: Health endpoint
      commands:
          - prebuilt: build-run.app
          - prebuilt: request.get
            params:
                url: http://localhost:3000/health
                response:
                    status: 200
                    body:
                        values:
                            status: ok
    - name: 404 page
      commands:
          - prebuilt: build-run.app
          - prebuilt: request.get
            params:
                url: http://localhost:3000/404
                response:
                    status: 404
```

## Installation

## Components

### CLI

### Client

### API

### Agent

````

```

```
````

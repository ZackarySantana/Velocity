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
    -   [Workflows](#workflows)
        -   [Example](#example-3)
    -   [Config](#config)
        -   [Example](#example-4)
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

-   Images
    -   Docker images that are used in tests and builds
-   Tests
    -   Units of process that can be used in a workflow
-   Builds
    -   The process of building your app
-   Workflows
    -   A series of tests that are ran
-   Config
    -   Options to connect your local environment to the platform

### Images

Images are the building blocks of your tests. They are Docker image aliases with additional options that you can use throughout your tests and builds.

#### Example

```yaml
images:
    - name: node
        image: node:latest
        commands:
        - npm install
        - npm run build
```

### Tests

#### Example

```yaml
TBA
```

### Builds

#### Example

```yaml
TBA
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

## Installation

## Components

### CLI

### Client

### API

### Agent

```

```

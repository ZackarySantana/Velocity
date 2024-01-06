---
sidebar_position: 4
---

# Deployments

A unit of process that handles a deployment of some kind. This is similar to [workflows](./workflows) but is given different permissions as well as only running a list of commands.

Deployments can list workflows they depend on passing before running. This can help prevent deploying broken code. This is done by listing the workflow names in the `workflows` field. The workflow name is lowercased and spaces are replaced with dashes (e.g. 'Workflow test NAME' -> 'workflow-test-name').

The following assumes there is some build called 'app' that outputs a directory called 'dist'.

```yaml title="velocity.yml"
deployments:
    - name: Deploy to staging
      workflows:
          - node-test
      commands:
          - prebuilt: build.app
          - prebuilt: command.targzip
            directory: dist
            output: output.tar.gz
          - prebuilt: request.post
            params:
                url: https://example.com/staging
                body: output.tar.gz
                headers:
                    - key: Content-Type
                      value: application/gzip
                    - key: Authorization
                      value: Bearer $TOKEN_STAGING # TODO: add support for env vars
                timeout: 10000
    - name: Deploy to production
      commands:
          - prebuilt: build.app
          - prebuilt: command.targzip
            directory: dist
            output: output.tar.gz
          - prebuilt: request.post
            params:
                url: https://example.com/prod
                body: output.tar.gz
                headers:
                    - key: Content-Type
                      value: application/gzip
                    - key: Authorization
                      value: Bearer $TOKEN_PROD # TODO: add support for env vars
                timeout: 10000
```

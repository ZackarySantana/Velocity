---
sidebar_position: 3
---

# Runtimes

Runtimes are where your instructions (tests, builds, deployments) will run. They are used in [workflows](./workflows), [builds](./builds), and [deployments](./deployments). Typically, these are docker iamges pulled from a registry. Depending on your deployment options, you can do bare metal runtimes or different types of runtimes as well. They are intended to be extensible, so consult your organization on what runtimes might be available to you.

The most commonm runtimes are docker images. You can also tie environmental variables to your runtimes, which will be available to your tests, builds, and deployments.

You can optionally change the docker registry that is pulled from, by default it is `docker.io`. View the [configuration sectiion](./config) for more information.

```yaml title="velocity.yml"
runtimes:
    - name: node
      image: node:latest
      env: # optional
          - NODE_ENV=test
          - PORT=3000
    - name: node-bare-metal
      machine: mac # Only available if your organization supports it. Check with your administrator(s).
        env: # optional
            - NODE_ENV=test
            - PORT=3000
```

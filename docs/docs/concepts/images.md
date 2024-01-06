---
sidebar_position: 2
---

# Images

Images are self-labeled docker images that are used in [workflows](./workflows), [builds](./builds), and [deployments](./deployments). They are used to define the environment that the tests will run in. For example, if you are testing a Node.js application, you would use the `node` image. If you are testing a Python application, you would use the `python` image.

You can optionally change the docker registry that is pulled from, by default it is `docker.io`. View the [configuration sectiion](./config) for more information.

```yaml title="velocity.yml"
images:
    - name: node
      image: node:latest
      env: # optional
          - NODE_ENV=test
          - PORT=3000
```

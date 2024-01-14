---
sidebar_position: 7
---

# Config

The config section helps the CLI know where to find your project.

```yaml title="velocity.yml"
config:
    project: velocity
    registry: docker.io/library
    server: https://velocity-ci.com/api
    ui: https://velocity-ci.com
```

As well, the registry, server, and UI will default to the settings displayed above. If you are hosting your own instance of velocity-ci, you can change these settings to point to your own instance.

---
sidebar_position: 2
---

# Operations

Operations are user-defined commands that are a list of other commands. They are used to create reusable components that can be used in multiple tests. They can be recursive and include 'env' or a 'working_dir'.

```yaml title="velocity.yml"
operations:
    - name: lint current directory
      commands:
          - prebuilt: git.clone
          - command: npm run lint
```

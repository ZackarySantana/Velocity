---
sidebar_position: 5
---

# Workflows

Workflows are sets of images and tests that are ran in any order. They are used in CI pipelines to run tests.

```yaml
workflows:
    - name: Backwards compatibility tests for Nodejs
      groups:
          - name: Node
            images:
                - node-16
                - node-17
                - node-18
                - node-19
            tests:
                - test-node
                - test-node-lint
    - name: Python tests
      groups:
          - name: Python
            images:
                - python-3.9
            tests:
                - test-python
```

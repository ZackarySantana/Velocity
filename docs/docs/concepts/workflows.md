---
sidebar_position: 6
---

# Workflows

Workflows are sets of runtimes and tests that are ran in any order. They are used in CI pipelines to run tests.

```yaml title="velocity.yml"
workflows:
    - name: Backwards compatibility tests for Nodejs
      groups:
          - name: Node
            runtimes:
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
            runtimes:
                - python-3.9
            tests:
                - test-python
```

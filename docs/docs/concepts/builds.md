---
sidebar_position: 3
---

# Builds

A unit of process that compiles an artifact. This artifact can be a binary, a tarball, json, etc. It can be used in tests and deployments.

## Usage

Builds can be used with the [build](./tests#build) prebuilt command package.

```yaml
builds:
    - name: app
      build_image: node
      output: dist # specifies what directory or file is the build output
      runtime_image: node # optional. used in build-run
      runtime_cmd: node index.js # optional. used in build-run
      commands:
          - command: npm install
          - command: npm run build
            env:
                - NODE_ENV=production
    - name: json
      build_image: node
      output: output.csv
      commands:
          - command: npm install
          - command: npm run generate-json
```

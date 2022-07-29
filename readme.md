```yaml
services:
  - hostname: golang184a
    type: go@1
    ports:
      - port: 1999
    buildFromGit: https://github.com/tikinang/zps@gitlab-runners-v2
    minContainers: 1
    maxContainers: 1
```

```yaml
services:
  - hostname: db
    type: keydb@6
    mode: NON_HA
  - hostname: stressor
    type: go@1
    ports:
      - port: 8080
    buildFromGit: https://github.com/tikinang/zps@stressor
    enableSubdomainAccess: true
    minContainers: 1
```

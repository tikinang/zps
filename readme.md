```yaml
services:
  - hostname: dummy
    type: go@1
    ports:
      - port: 1999
    buildFromGit: https://github.com/tikinang/zps@dummy
    enableSubdomainAccess: true
    minContainers: 1
    maxContainers: 1
```

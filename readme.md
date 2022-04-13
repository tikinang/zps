```yaml
services:
  - hostname: zps
    type: go@1
    mode: NON_HA
    ports:
      - port: 8080
    buildFromGit: https://github.com/tikinang/zps@logging
```

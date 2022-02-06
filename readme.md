```yaml
services:
  - hostname: zps
    type: go@1
    mode: NON_HA
    ports:
      - port: 7100
      - port: 22
        protocol: tcp
        httpSupport: false
    buildFromGit: https://github.com/tikinang/zps
    enableSubdomainAccess: true
```
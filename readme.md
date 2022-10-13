```yaml
services:
  - hostname: dummy
    type: ubuntu@22.04
    ports:
      - port: 1999
    minContainers: 1
    maxContainers: 1
  - hostname: old
    type: go@1
    ports:
      - port: 1999
    minContainers: 1
    maxContainers: 1
```

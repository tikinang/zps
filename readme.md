```yaml
services:
  - hostname: elasticsearch
    type: go@1
    ports:
      - port: 9200
    envVariables:
      ES_HOST: ${ZEROPS_Hostname}
      ES_HOST_LIST: ${ZEROPS_Hostnames|pipeToComma}
    buildFromGit: https://github.com/tikinang/zps@elasticsearch
    enableSubdomainAccess: true
    verticalAutoscaling:
      minRam: 6
    minContainers: 3
    maxContainers: 3
```

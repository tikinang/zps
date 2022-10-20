```yaml
services:

  - hostname: ubuntu
    type: ubuntu@22.04
    ports:
      - port: 1999
    minContainers: 1
    maxContainers: 1

  - hostname: go
    type: go@1
    ports:
      - port: 1999
    minContainers: 1
    maxContainers: 1

  - hostname: dotnet
    type: dotnet@6
    ports:
      - port: 1999
    envVariables:
      MY_X_VAR: "1"
    minContainers: 1
    maxContainers: 1

  - hostname: elasticsearch
    type: ubuntu@22.04
    ports:
      - port: 9200
    envVariables:
      ES_HOST: ${ZEROPS_Hostname}
      ES_HOST_LIST: ${ZEROPS_Hostnames|pipeToComma}
    verticalAutoscaling:
      minRam: 6
    minContainers: 3
    maxContainers: 3
```

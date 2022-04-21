```yaml
services:
  - hostname: db
    type: keydb@6
    mode: NON_HA
  - hostname: zps
    type: go@1
    mode: NON_HA
    ports:
      - port: 8080
    envVariables:
      JSON_ENV: |-
        {
            "intro": "👋👋👋",
            "dbHostname": "$db_hostname",
            "subdomain": "${zeropsSubdomain}",
            "below": "${JSON_ENV_BELOW|stringify}"
        }
      JSON_ENV_BELOW: |-
        {
          "bar": "$db_hostname",
          "baz": "👋👋👋"
        }
    buildFromGit: https://github.com/tikinang/zps@env
    enableSubdomainAccess: true
```

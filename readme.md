# import.yml

```yaml
  - hostname: stelabook
    type: go@1
    mode: NON_HA
    ports:
      - port: 2922
    envVariables:
      HTTP_BASIC_AUTH_USERNAME: fill_me_out_please
      HTTP_BASIC_AUTH_PASSWORD: fill_me_out_please
    buildFromGit: https://github.com/tikinang/zps@stelabook
    enableSubdomainAccess: true
```

# build client
TODO
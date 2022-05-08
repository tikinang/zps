# import.yml

```yaml
  - hostname: stelabook
    type: go@1
    mode: NON_HA
    ports:
      - port: 2922
    envVariables:
      HTTP_BASIC_AUTH_USERNAME: <your_basic_auth_username>
      HTTP_BASIC_AUTH_PASSWORD: <your_basic_auth_password>
    buildFromGit: https://github.com/tikinang/zps@stelabook
    enableSubdomainAccess: true
```

# build client
`GOOS=linux GOARCH=arm go build -ldflags="-X 'main.Username=<your_basic_auth_username>' -X 'main.Password=<your_basic_auth_password>' -X 'main.ServerAddress=<server_address>'" -o /media/tikinang/PB616W/stelabook_client cmd/client/main.go`

contents of `/media/<user>/PB616W/applications/stelabook.app`
```shell
#!/bin/sh

chmod +x /mnt/ext1/stelabook_client
/mnt/ext1/stelabook_client >> /mnt/ext1/stelabook.log 2>&1
```
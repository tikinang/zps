```yaml
services:
  - hostname: redis0
    type: keydb@6
    mode: NON_HA
  - hostname: zps
    type: go@1
    mode: NON_HA
    ports:
      - port: 7100
      - port: 22
        protocol: tcp
        httpSupport: false
    envVariables:
      - key: SSH_PUBLIC_KEY
        content: ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+PDlUXMuxWR/zHKTYCOsItQ7x84WLwyBVPRpo5i95ZRXLObtsBuPQHwNSwW8NF2Jw4oEkjP9GDA3mM9yDrq3qLL6vPLg2/yM77wkhlID0FovUMfoc3cvo8yix7RCOMygyQjk7+S/WtwF5fL8GvdLgDXQDtAU86vdnfHxtpzjH58d8sqoW1zqY/+UBFF2xHleFNNdovSRI41k5GJOuLDQhwzIjl8GXYO9V3u/0gMwAn1zm1qqe2r+TvzjPPPMpCVAz4CIxmSYaDzL+waSJift7goadypolTvkMLltJmFZlpmBwri5cjEkK0+SVz/uZohswXtxUKpS0FUOUhY/w7KBj tikinang@gmail.com
      - key: ZPS_LISTEN_PORT
        content: 7100
    buildFromGit: https://github.com/tikinang/zps
    enableSubdomainAccess: true
```

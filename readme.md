```yaml
services:
  - hostname: db
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
      SSH_PUBLIC_KEY: |
        ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+PDlUXMuxWR/zHKTYCOsItQ7x84WLwyBVPRpo5i95ZRXLObtsBuPQHwNSwW8NF2Jw4oEkjP9GDA3mM9yDrq3qLL6vPLg2/yM77wkhlID0FovUMfoc3cvo8yix7RCOMygyQjk7+S/WtwF5fL8GvdLgDXQDtAU86vdnfHxtpzjH58d8sqoW1zqY/+UBFF2xHleFNNdovSRI41k5GJOuLDQhwzIjl8GXYO9V3u/0gMwAn1zm1qqe2r+TvzjPPPMpCVAz4CIxmSYaDzL+waSJift7goadypolTvkMLltJmFZlpmBwri5cjEkK0+SVz/uZohswXtxUKpS0FUOUhY/w7KBj tikinang@gmail.com
        ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCrgObBsxG4PGoAZtLDRVqkTbfLG3vwwHI1jeP7BzD0k+t3/6qXC1WCakG8J01oiF2JSfHfIiAMC82oJ8PdMsMaod7N3ck1mAfLyQn5r+ElgepSfX8iDLJ9B0aEXgHaYBYMToPRyIO+L7X8aU3ckxlBeCsuxPtkriwKgJf9eA44omzuTj8L65XlKw6E1hLkVOqXm4AUCz/rkq2iq9VasXvxEPmv2Wk7WsVKOrViFp48qmfra9btFR1mig0oia9okXZNqPCmjMPJCmd9U54ln99W6qeV2BKBhsYQrwiaX8a8/Ja3Nezl0B74cNTaIozeJPTAHFV6Bhf/J/hklSEkZfG/yoqEyZA2dsjs2VBkDlIHSM7IQISGvelcbhhpV9R0CKWUR8Ae4M0B2GKKNBU6jv0zNrnDzI3Ho6fVmHTm3GhqgoaPujCXicRn6W5ob7XMxtStdCccHocXv6g6D2iz0jHx7jDIZR2XFyuX6cWYexmEapwKRTXLNS9sfg0iKcb6ayU= tikinang@matej
      ZPS_LISTEN_PORT: 7100
    buildFromGit: https://github.com/tikinang/zps
    enableSubdomainAccess: true
```

# Telegram 赌场机器人

# 获取代码
```
go get -u github.com/kardianos/govendor
go get -u github.com/zhangpanyi/botcasino
govendor fetch +external
```

# 生成密钥
```shell
cd certs

openssl genrsa -out server.key 2048

openssl ecparam -genkey -name secp384r1 -out server.key

openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650 -subj "/C=US/ST=New York/L=Brooklyn/O=Brooklyn Company/CN=YOURDOMAIN"

```

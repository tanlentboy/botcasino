# Telegram 红包机器人
botcasino 是基于 [Telegram](https://telegram.org/) 机器人的红包机器人。自从中国禁止加密货币与法币的交易后，中国越来越多的加密货币玩家开始使用Telegram聊天。但是Telegram一直没有发红包功能，对于习惯使用QQ和微信的朋友来说，确实是一个遗憾。botcasino 正是为了解决这个问题而被创造出来。由于加密货币的价格不够稳定，在支付或交易的过程中会有很多不方便的地方。所有 botcasino 选择了[比特股](https://bitshares.org/)系统内的智能货币 [BitCNY](https://coinmarketcap.com/currencies/bitcny/)/[BitUSD](https://coinmarketcap.com/currencies/bitusd/) 作为红包的基础货币。

请在 Telegram 中 @luck_money_bot，或者打开 [http://telegram.me/luck_money_bot](http://telegram.me/luck_money_bot) 体验吧。

![](http://i796.photobucket.com/albums/yy247/zhangpanyi/1_zpsuxxjuzgp.png)

# 获取代码
```
git clone https://github.com/zhangpanyi/botcasino.git
glide installs
```

# 解析域名
botcasino 服务使用了 [Webhook](https://core.telegram.org/bots/api#setwebhook) 的方式接收机器人消息更新。所有必须准备一个域名，并解析到运行 botcasino 的服务器上。

# 生成密钥
```shell
cd certs

openssl genrsa -out server.key 2048

openssl ecparam -genkey -name secp384r1 -out server.key

openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650 -subj "/C=US/ST=New York/L=Brooklyn/O=Brooklyn Company/CN=YOURDOMAIN"

```
> 注意： 请将 YOURDOMAIN 修改为你自己的域名。


# 数据库
botcasino 服务使用 MySQL 数据库存储操作记录，请预先准备 MySQL 数据库。
1. 创建一个新的数据库。数据库名为 `casino_logs`，字符集为 `utf8mb4 -- UTF-8 Unicode`。
2. 将 `mysql/casino_logs.sql` 导入到这个数据库中。

# 钱包服务
钱包服务的作用是监控机器人的比特股钱包。如果有人转账到机器人的比特股钱包它就会通知 botcasino，用户的提现申请 botcasino 也是通过钱包服务去处理的，它们之间通过 gRPC 进行交流。
```
git clone https://github.com/zhangpanyi/btsmonitor.git
```
请按照说明文档启动钱包服务。

# 配置文件
1. `dynamic.yml` 是动态配置文件，可在服务运行期间修改生效，使用默认配置就可以了。
2. `master.yml` 是服务的基本配置文件，启动服务前必须将 `domain`、`token`、`mysql`字段改为自己的配置。`domain` 字段请使用 `www.google.com` 格式，不要使用 `https://www.google.com/` 格式。

# 启动服务
```
go build
./botcasino
```

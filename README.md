# pixiv-tg-bot

🤖 方便 Pixiv 上看小说的 Telegram Bot

<div >
    <img style="" src="./screenshots/1.png"/>
    <img style="" src="./screenshots/2.png"/>
</div>

## 部署

### 机器人设置

```txt
# 设置机器人命令
/setcommands

start - 快速开始
help - 查看帮助信息
subnovels - 订阅小说
showsubnovels - 查看已经订阅的小说
checknovelupdate - 查看订阅的小说是否更新
removesubnovels - 移除订阅的小说
```

### 运行

```shell
pixiv-tg-bot -t [机器人token]

# 通过代理运行
pixiv-tg-bot -t [机器人token] -p [host:port]

# 指定数据库
pixiv-tg-bot -t [机器人token] -p [host:port] -d [数据路径]

# 后台运行
nohup pixiv-tg-bot -t [机器人token] -p [host:port] -d ./database/pixiv.db > out.log &
```

## 使用

```shell
# pixiv-tg-bot --help
NAME:
   pixiv telegram bot - A new cli application

USAGE:
   pixiv telegram bot [global options] command [command options] [arguments...]

VERSION:
   1.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --token value, -t value   机器人的 Token
   --proxy value, -p value   代理地址, 比如(127.0.0.1:10808)
   --dbPath value, -d value  Sqlite3的数据库路径(默认为: './database/pixiv.db')
   --help, -h                show help (default: false)
   --version, -v             print the version (default: false)
```

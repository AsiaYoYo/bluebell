# bluebell论坛

# 使用指南
## 下载
```
git clone https://github.com/AsiaYoYo/bluebell.git
```

## 配置MySQL
1. 在你的数据库中执行以下命令，创建本项目所用的数据库：
```
CREATE DATABASE bluebell DEFAULT CHARSET=utf8mb4;
```
2. 在bluebell/conf/config.yaml文件中按如下提示配置数据库连接信息。
```
name: "bluebell"
# 生产环境可配置为"release"
mode: "dev"
start_time: "2021-03-01"
machine_id: 1
port: 8081
version: "v0.1.4"

log:
  level: "debug"
  filename: "bluebell.log"
  max_size: 200
  max_age: 30
  max_backups: 7

mysql:
  host: "你的数据库host地址"
  port: 你的数据库端口
  user: "你的数据库用户名"
  password: "你的数据库密码"
  dbname: "bluebell"
  max_open_conns: 20
  max_idle_conns: 10

redis:
  host: "你的redis host地址"
  port: 你的redis端口
  password: ""
  db: 0
  pool_size: 100
```

## 编译
```
go build
```

## 执行
Mac/Unix：
```
./bluebell conf/config.yaml
```
Windows:
```
bluebell.exe conf/config.yaml
```

# 使用docker部署
## 构建docker image
```
docker build . -t bluebell_app
```

## 启动容器
```
docker run -itd -v /YOUR_DIR/config.yaml:/conf/config.yaml -p 8081:8081 --name bluebell bluebell_app

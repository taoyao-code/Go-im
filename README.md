# Go-im实现

## 简单介绍
一个即时通讯服务器，基于WebSocket协议,使用Golang语言完成,

## 实现功能
- 登录注册
- 验证码
- 上传文件
- 单聊/群聊
- 文字、表情、图片、语音等消息推送
- 添加好友
- 创建群/加入群
- 七牛云对象存储
- 消息持久化
- 使用redis存储数据，减少数据库io操作
- JWT-Token认证模式
- 日志记录
- 文字敏感信息过滤


## 使用说明
- 1、安装数据库
- 2、创建数据库
- 3、在config目录下，创建config.yaml文件进行数据配置
```
MySQL:
  Username: xxx
  Password: xxxx
  Address: xxx
  Port: xxx
  Database: xxxx
QNY: #七牛云配置
  QINIU_DOMAIN: xxx
  QINIU_ACCESS_KEY: xxxx
  QINIU_SECRET_KEY: xxxx
  QINIU_TEST_BUCKET: xxxx
Redis: # redis
  Address: 127.0.0.1:6379
  Port: 6379
  Password: foobared
token_key:
```
- 4、根目录下运行：go build
- 5、运行生成的文件
- 6、apidoc.json为apidoc文档生成配置文件

## api文档
```
http://chat.bo5.xyz/apidoc/
```
## 说明
[前端文件](https://github.com/ltsj404/chat-im.git)
```shell script
https://github.com/ltsj404/chat-im.git
```
```
个人项目，知识有限，欢迎 issue
```
后续请自行扩展

原项目参考：
```
https://github.com/winlion/chat
```

# Go-im实现

## 简单介绍
一个即时通讯服务器，基于WebSocket协议,使用Golang完成

## 实现功能
- 登录注册
- 验证码
- 上传文件
- 单聊/群聊
- 文字、表情、图片、语音等消息推送
- 添加好友
- 创建群/加入群
- 七牛云对象存储


## 使用说明
- 1、安装数据库
- 2、创建数据库
- 3、在跟目录下创建config目录，创建config.yaml文件进行数据配置
```
mysql:
  username: xxxx
  password: xxxx
  host: xxxx
  port: xxxx
  dbname: xxx
qiniu:
  QINIU_DOMAIN: xxx
  QINIU_ACCESS_KEY: xxxx
  QINIU_SECRET_KEY: xxxxx
  QINIU_TEST_BUCKET: xxxxx
```
- 4、根目录下运行：go build
- 5、运行生成的文件

## 说明
前端文件
```shell script
https://github.com/ltsj404/chat-im.git
```
后续请自行扩展

原项目参考：
```
https://github.com/winlion/chat
```
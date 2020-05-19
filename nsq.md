下载二进制安装包

进入bin目录
-  打开一个终端，启动`nsqlookupd`
```sh
nohup ./nsqlookupd &
```
- 打开另一个终端，启动`nsqd`
```sh
nohup ./nsqd --lookupd-tcp-address=127.0.0.1:4160 &
```
- 打开另一个终端，启动`nsqadmin`
```sh
nohup ./nsqadmin --lookupd-http-address=127.0.0.1:4161 &
```

打开另外一个终端，启动nsq_to_file，将消息写入/tmp文件的日志文件，文件名默认由主题topic+主机+日期时间戳组成
```sh
nohup ./nsq_to_file --topic=test --output-dir=/tmp --lookupd-http-address=127.0.0.1:4161 &

```

使用curl命令，发布一条消息
```sh
curl -d 'hello world' 'http://127.0.0.1:4151/pub?topic=test'
```

3.相关概念

- nsqlookupd：管理nsqd节点拓扑信息并提供最终一致性的发现服务的守护进程
- nsqd：负责接收、排队、转发消息到客户端的守护进程，并且定时向nsqlookupd服务发送心跳
- nsqadmin：nsq的web统计界面，可实时查看集群的统计数据和执行一些管理任务
- utilities：常见基础功能、数据流处理工具，如nsq_stat、nsq_tail、nsq_to_file、nsq_to_http、nsq_to_nsq、to_nsq



### 开启和关闭nsq shell脚本

进入bin目录，创建data目录

```sh
mkdir data
```

开启nsq

```sh
sudo chmod -R u+x nsq_start.sh
./nsq_start.sh
```

 nsq_start.sh 

```shell

#!/bin/sh
#服务启动
#lookupd:151 152
#更改 --data-path 所指定的数据存放路径，否则无法运行
echo '删除日志文件'
rm -f nsqlookupd.log
rm -f nsqd1.log
rm -f nsqd2.log
rm -f nsqadmin.log

echo '启动nsqlookupd服务'
nohup ./nsqlookupd >nsqlookupd.log 2>&1 &

echo '启动nsqd服务'
#nohup ./nsqd --lookupd-tcp-address=0.0.0.0:4160 -tcp-address="0.0.0.0:4153" --data-path=./data1  >nsqd1.log 2>&1 &
nohup ./nsqd --lookupd-tcp-address=0.0.0.0:4160 -tcp-address="0.0.0.0:4154" -http-address="0.0.0.0:4155" --data-path=./data >nsqd2.log 2>&1 &
echo '启动nsqdadmin服务'
nohup ./nsqadmin --lookupd-http-address=127.0.0.1:4161 >nsqadmin.log 2>&1 &


```

 关闭nsq 

```shell
sudo chmod -R u+x nsq_shutdown.sh
./nsq_shutdown.sh
```

 nsq_shutdown.sh 

```shell

#!/bin/sh
#服务停止
ps -ef | grep nsq| grep -v grep | awk '{print $2}' | xargs kill -2
```


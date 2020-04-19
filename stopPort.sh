# 1、根据名称称查找并关闭：
pgrep -f reptile-go | xargs kill -9

# 2、根据端口称查找并关闭：
#lsof -i:端口 | grep LISTEN|awk '{print $2}' | xargs kill -9

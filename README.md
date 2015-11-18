# agent -- golang代理服务器

# 安装
```go
go get github.com/guhao022/agent
```

# 使用
比如代理localhost:80
```go
agent run localhost:80  //代理localhost:80
```

比如访问localhost:80?id=1&name=aa
访问代理服务器 localhost:9900?id=1&name=aa

访问可以使用post和get方法

控制台可以输出想要的信息

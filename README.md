# 功能介绍
## 说明
将微信公众号上接收到的消息转发到配置的bearychat 讨论组中。  
目前支持
- 文字型消息
- 图片类型消息
# 使用
## 参数说明

```
NAME:
   weixin2bearychat server - start weixin gate server

USAGE:
   weixin2bearychat server [command options] [arguments...]

OPTIONS:
   --host value, -H value    server bind host to (default: "0.0.0.0")
   --port value, -p value    bind port (default: 80)
   --tmplpath value          msg template path (default: "/etc/weixin2bearychat/tmpl/")
   --target value, -t value  which url that post msg to

```
## 启动转发服务
```
make build
cd ./bin
./weixin2bearychat server --target https://bearychaturl
```
登陆服务器，将./bin/tmpl 拷贝至/etc/weixin2bearychat/tmpl   
启动服务，bearychaturl 配置为 对应讨论组添加的incoming 机器人的消息地址  
注意将端口绑定到公网的ip和80端口上。（微信公众号只支持80端口）
## 绑定微信公众号
打开微信公众号配置，将转发服务监听的ip地址 例如 http://x.x.x.x 配置公众号server地址，并启用。
# 开发
go  > v1.6.3  
govendor  > v1.0.8  
编译 `make build`

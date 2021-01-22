### Token中控服务器
微信Access Token中控服务器，用来统一管理各个小程序/公众号的access_token，提供统一的接口进行获取和自动刷新Access Token。

### 简介
微信access_token每日有一个次数限制，所以客户服务器不能每次都去请求一个新的access_token，每次请求之后，access_token都有一个过期时间。因此微信平台建议你使用一个中控服务器来定时刷新token，取得之后存起来不用再去请求token，因为access_token请求有次数限制。这样处理只有有两个好处：

1. 保证access_token每日都不会超出访问限制，保证服务的正常。
2. 提高服务的性能，不用每次发送业务请求之前都先发送一次获取access_token请求。（将access_token保存在内存中，直到过期的时候再去请求一个新的来替代）

>> 微信建议开发者使用中控服务器统一获取和刷新Access_token，其他业务逻辑服务器所使用的access_token均来自于该中控服务器，不应该各自去刷新，否则容易造成冲突，导致access_token覆盖而影响业务

### 特点
* Basic Auth的认证方式，需要通过HTTP Basic认证才能访问
* 使用BuntDB作为access_token的内存缓存数据库
* Echo框架编写的REST API服务

### 获取令牌
http://localhost:8080/token?appid=appid

### 制作镜像
1. CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
2. DOCKER_HOST=tcp://ip:port docker build -t 项目名.镜像名:版本号 .

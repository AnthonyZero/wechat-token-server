#源镜像
FROM alpine:latest
#作者
MAINTAINER anthonyzero "736252868@qq.com"
#设置工作目录
WORKDIR /app
ADD main /app
ADD wechat.json /app
#暴露端口
EXPOSE 8080
ENTRYPOINT  ["./main"]
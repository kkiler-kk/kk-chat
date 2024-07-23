#源镜像
FROM golang:alpine3.13

#将二进制文件拷贝进容器的GOPATH目录中
ADD main   /go/src/main

#为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.cn,direct"

#暴露端口
EXPOSE 3000

#工作目录
WORKDIR /go/src

#最终运行docker的命令
CMD ["./main"]
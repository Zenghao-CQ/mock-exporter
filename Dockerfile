#源镜像
FROM golang:1.16
WORKDIR /run
COPY ./mock-exporter .
#go构建可执行文件
EXPOSE 8000
#最终运行docker的命令
ENTRYPOINT  ["./mock-exporter"]
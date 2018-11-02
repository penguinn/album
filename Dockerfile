# 编译源代码
FROM golang:1.11.0 AS compile
COPY . /go/src/github.com/penguinn/album/
WORKDIR /go/src
RUN go build -o album github.com/penguinn/album/

# 把编译好的二进制放入到运行时容器
FROM golang:1.11.0
COPY --from=compile /go/src/album /album
CMD /album
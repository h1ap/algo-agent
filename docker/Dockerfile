FROM docker.mybacc.com/golang:1.24 AS builder

COPY . /src
WORKDIR /src

# 构建 dlv 调试工具
RUN GOPROXY=https://goproxy.cn make build

FROM docker.mybacc.com/debian:stable-slim

RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
# 先安装ca-certificates
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

# 更新并安装必要的包
RUN apt-get update --allow-insecure-repositories && \
    apt-get install -y --no-install-recommends \
        netbase \
    && rm -rf /var/lib/apt/lists/ \
    && apt-get autoremove -y \
    && apt-get autoclean -y

COPY --from=builder /src/bin /app

WORKDIR /app

EXPOSE 30011
VOLUME /data/conf

CMD ["./server", "-conf", "/data/conf"]

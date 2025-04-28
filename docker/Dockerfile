FROM docker.mybacc.com/golang:1.24 AS builder

COPY . /src
WORKDIR /src

# 构建 dlv 调试工具
RUN GOPROXY=https://goproxy.cn make build

FROM docker.mybacc.com/debian:stable-slim

# 使用中科大镜像源
RUN sed -i 's/deb.debian.org/mirrors.huaweicloud.com/g' /etc/apt/sources.list && \
    sed -i 's|security.debian.org/debian-security|mirrors.huaweicloud.com/debian-security|g' /etc/apt/sources.list

# 更新并安装必要的包
RUN apt-get update --allow-insecure-repositories && \
    apt-get install -y --no-install-recommends ca-certificates && \
    apt-get install -y netbase \
    && rm -rf /var/lib/apt/lists/ \
    && apt-get autoremove -y \
    && apt-get autoclean -y

COPY --from=builder /src/bin /app

WORKDIR /app

EXPOSE 30011
VOLUME /data/conf

CMD ["./server", "-conf", "/data/conf"]

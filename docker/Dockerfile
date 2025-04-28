FROM docker.mybacc.com/golang:1.24 AS builder

COPY . /src
WORKDIR /src

# 构建 dlv 调试工具
RUN GOPROXY=https://goproxy.cn make build

FROM docker.mybacc.com/debian:stable-slim

# 先安装ca-certificates
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# 使用中科大镜像源
RUN echo "deb https://mirrors.ustc.edu.cn/debian/ stable main contrib non-free non-free-firmware" > /etc/apt/sources.list && \
    echo "deb https://mirrors.ustc.edu.cn/debian/ stable-updates main contrib non-free non-free-firmware" >> /etc/apt/sources.list && \
    echo "deb https://mirrors.ustc.edu.cn/debian/ stable-backports main contrib non-free non-free-firmware" >> /etc/apt/sources.list && \
    echo "deb https://mirrors.ustc.edu.cn/debian-security stable-security main contrib non-free non-free-firmware" >> /etc/apt/sources.list

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

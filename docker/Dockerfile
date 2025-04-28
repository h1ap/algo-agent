FROM docker.mybacc.com/golang:1.24 AS builder

COPY . /src
WORKDIR /src

# 构建 dlv 调试工具
RUN GOPROXY=https://goproxy.cn make build

FROM docker.mybacc.com/debian:stable-slim

# 创建并配置华为云镜像源
RUN echo "deb http://mirrors.huaweicloud.com/debian/ stable main contrib non-free non-free-firmware" > /etc/apt/sources.list && \
    echo "deb http://mirrors.huaweicloud.com/debian/ stable-updates main contrib non-free non-free-firmware" >> /etc/apt/sources.list && \
    echo "deb http://mirrors.huaweicloud.com/debian/ stable-backports main contrib non-free non-free-firmware" >> /etc/apt/sources.list && \
    echo "deb http://mirrors.huaweicloud.com/debian-security stable-security main contrib non-free non-free-firmware" >> /etc/apt/sources.list

# 更新并安装必要的包
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    apt-get install -y netbase \
    && rm -rf /var/lib/apt/lists/ \
    && apt-get autoremove -y \
    && apt-get autoclean -y

COPY --from=builder /src/bin /app
COPY --from=builder /src/configs /app

WORKDIR /app
RUN ls -la

EXPOSE 30011
VOLUME /data/conf

ENTRYPOINT ["/docker-entrypoint.sh"]

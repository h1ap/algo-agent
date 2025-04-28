#!/bin/bash

# 检查是否传入镜像名称
if [ -z "$1" ]; then
    echo "请提供镜像名称作为参数。"
    echo "用法: $0 <镜像名称>"
    exit 1
fi

# 从参数中获取镜像名称
IMAGE_NAME="$1"

# 检查镜像是否存在
IMAGE_ID=$(docker images -q $IMAGE_NAME)

if [ -n "$IMAGE_ID" ];then
    echo "找到镜像: $IMAGE_NAME (ID: $IMAGE_ID)"

    # 检查是否有相关的容器在运行
    CONTAINER_IDS=$(docker ps -a -q --filter "ancestor=$IMAGE_NAME")

    if [ -n "$CONTAINER_IDS" ]; then
        echo "找到相关容器，停止并删除容器..."
        # 停止容器
        docker stop $CONTAINER_IDS
        # 删除容器
        docker rm $CONTAINER_IDS
    else
        echo "没有相关的容器。"
    fi

    # 删除镜像
    echo "删除镜像..."
    docker rmi $IMAGE_ID --force
else
    echo "镜像 $IMAGE_NAME 不存在。"
fi


# 设置脚本为可执行
# chmod +x remove_docker_image.sh

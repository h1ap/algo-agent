#!/bin/sh

# 如果没有指定配置文件路径，bootstrap.yml
CONFIG_LOCATION=${CONFIG_LOCATION:-"/app/config.yaml"}

exec /app/algo-agent -conf ${CONFIG_LOCATION}
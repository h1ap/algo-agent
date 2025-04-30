package mq

import (
	"errors"
	"fmt"
)

// MqMessageTypeEnum 消息枚举
type MqMessageTypeEnum struct {
	code int
	name string
}

// 定义枚举值
var (
	UNKNOWN          = &MqMessageTypeEnum{0, "未知"}
	SYSTEM_METRICS   = &MqMessageTypeEnum{1, "系统指标"}
	DOCKER_LOG       = &MqMessageTypeEnum{2, "docker日志"}
	TASK_EVALUATE    = &MqMessageTypeEnum{3, "评估任务"}
	TRAIN_TASK       = &MqMessageTypeEnum{4, "训练任务"}
	TRAIN_PUBLISH    = &MqMessageTypeEnum{5, "训练发布"}
	TRAIN_POOL_CLOSE = &MqMessageTypeEnum{6, "训练池关闭"}
	TASK_DEPLOY      = &MqMessageTypeEnum{7, "部署任务"}
)

// Code 获取编码
func (e *MqMessageTypeEnum) Code() int {
	return e.code
}

// Name 获取枚举名
func (e *MqMessageTypeEnum) Name() string {
	return e.name
}

// codeToString 将编码转换为字符串
func (e *MqMessageTypeEnum) CodeToString() string {
	return fmt.Sprintf("%d", e.code)
}

// ENUM_MAPS 存储枚举值和对应的枚举实例
var ENUM_MAPS = make(map[int]*MqMessageTypeEnum)

// 初始化ENUM_MAPS
func init() {
	ENUM_MAPS[SYSTEM_METRICS.code] = SYSTEM_METRICS
	ENUM_MAPS[DOCKER_LOG.code] = DOCKER_LOG
	ENUM_MAPS[TASK_EVALUATE.code] = TASK_EVALUATE
	ENUM_MAPS[TRAIN_TASK.code] = TRAIN_TASK
	ENUM_MAPS[TRAIN_PUBLISH.code] = TRAIN_PUBLISH
	ENUM_MAPS[TRAIN_POOL_CLOSE.code] = TRAIN_POOL_CLOSE
	ENUM_MAPS[TASK_DEPLOY.code] = TASK_DEPLOY
}

// GetNameByCode 根据编码获取枚举名
func GetNameByCode(code int) (string, error) {
	if enum, ok := ENUM_MAPS[code]; ok {
		return enum.name, nil
	}
	return "", errors.New("枚举未找到")
}

// GetEnumByCode 根据编码获取枚举实例
func GetEnumByCode(code int) (*MqMessageTypeEnum, error) {
	if enum, ok := ENUM_MAPS[code]; ok {
		return enum, nil
	}
	return nil, errors.New("枚举未找到")
}

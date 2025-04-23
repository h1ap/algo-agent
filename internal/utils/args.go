package utils

import (
	taskArgs "algo-agent/internal/cons/task"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// 添加标签
func AddLabels(datasetLabel string, args *[]string) {
	// 解析JSON字符串为map
	stringObjectMap, err := ParseToMap(datasetLabel)
	if err != nil {
		// 如果解析失败，直接返回
		return
	}

	if stringObjectMap != nil && len(stringObjectMap) > 0 {
		// 创建一个可排序的切片
		type KeyValue struct {
			Key   int
			Value interface{}
		}

		sortedList := make([]KeyValue, 0, len(stringObjectMap))

		// 将map转换为切片
		for k, v := range stringObjectMap {
			key, err := strconv.Atoi(k)
			if err != nil {
				continue
			}
			sortedList = append(sortedList, KeyValue{Key: key, Value: v})
		}

		// 按key排序
		sort.Slice(sortedList, func(i, j int) bool {
			return sortedList[i].Key < sortedList[j].Key
		})

		// 创建一个字符串切片存储值
		values := make([]string, 0, len(sortedList))
		for _, entry := range sortedList {
			values = append(values, fmt.Sprintf("%v", entry.Value))
		}

		// 添加类名参数
		*args = append(*args, taskArgs.ArgClassNames)
		*args = append(*args, strings.Join(values, ","))
	}
}

package file

import (
	"path/filepath"
)

// 系统文件分隔符
var SEPARATOR = string(filepath.Separator)

// 文件路径常量
const (
	// Unix系统的文件分割线
	UNIX_SEPARATOR = "/"

	// 训练任务 根目录
	TRAIN = "train"

	// 提取任务
	EXTRACT = "extract"

	// 评估任务 根目录
	EVAL = "eval"

	// 部署任务 根目录
	DEPLOY = "deploy"

	// 镜像目录
	IMAGE = "image"

	// 数据集目录
	DATASET = "dataset"

	// 数据集数据目录
	DATASET_DATA = "data"

	// 数据集标注目录
	DATASET_ANNOTATION = "annotation"

	// 算法脚本目录
	SCRIPT = "script"

	// 模型目录
	MODEL = "model"

	// 检查点目录
	CHECKPOINT = "checkpoint"

	// 容器工作路径
	WORKSPACE = "workspace"
)

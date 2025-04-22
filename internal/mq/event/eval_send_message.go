package event

// EvalSendMessage 评估作业发送消息
type EvalSendMessage struct {
	// 评估作业主键
	TaskId string `json:"taskId"`

	// 操作类型：0启动 1终止
	Op int32 `json:"op"`

	// 算法名称
	AlgorithmName string `json:"algorithmName"`

	// 算法脚本所在minio存储桶
	AlgorithmScriptMinioBucket string `json:"algorithmScriptMinioBucket"`

	// 算法脚本所在minio的存储地址
	AlgorithmScriptZipMinioURL string `json:"algorithmScriptZipMinioUrl"`

	// 算法训练用docker image名称
	AlgorithmTrainImageName string `json:"algorithmTrainImageName"`

	// 算法训练用docker image所在minio存储桶
	AlgorithmTrainImageMinioBucket string `json:"algorithmTrainImageMinioBucket"`

	// 算法训练用docker image在minio的存储地址
	AlgorithmTrainImageMinioURL string `json:"algorithmTrainImageMinioUrl"`

	// 训练数据所在Bucket
	DatasetBucket string `json:"datasetBucket"`

	// 训练数据所在的数据目录集合
	DatasetFolders []string `json:"datasetFolders"`

	// 数据文件是否为zip
	DataZip bool `json:"dataZip"`

	// 标注文件据所在的数据目录集合
	AnnotationFolders []string `json:"annotationFolders"`

	// 标注文件是否zip格式
	AnnotationZip bool `json:"annotationZip"`

	// 评估类型：1 模型 2 检查点
	EvalType string `json:"evalType"`

	// 模型或检查点的 存储桶
	ModelOrCheckpointBucket string `json:"modelOrCheckpointBucket"`

	// 模型 或检查点的文件地址
	ModelOrCheckpointFileURL string `json:"modelOrCheckpointFileUrl"`

	// 标注标签
	DatasetLabel string `json:"datasetLabel"`

	// 训练参数
	Args []string `json:"args"`
}

package event

// TrainTaskReqMessage 训练任务请求消息
type TrainTaskReqMessage struct {
	// 训练任务ID
	TaskId string `json:"taskId"`

	// 操作类型：0-启动，1-终止
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

	// 预训练权重文件所在bucket
	PreModelBucket string `json:"preModelBucket"`

	// 预训练权重文件在minio的存储地址
	PreModelFileURL string `json:"preModelFileUrl"`

	// 标注标签
	DatasetLabel string `json:"datasetLabel"`

	// 训练参数
	Args []string `json:"args"`
}

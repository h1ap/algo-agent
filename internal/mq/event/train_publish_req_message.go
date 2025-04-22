package event

// TrainPublishReqMessage 训练发布请求消息体
type TrainPublishReqMessage struct {
	// 训练详情id
	TrainDetailId string `json:"trainDetailId"`

	// 模型版本Id
	ModelVersionId int64 `json:"modelVersionId"`

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

	// 模型或检查点的 存储桶
	CheckpointBucket string `json:"checkpointBucket"`

	// 模型 或检查点的文件地址
	CheckpointFileURL string `json:"checkpointFileUrl"`

	// 训练参数
	Args []string `json:"args"`
}

// 实现ReqMessage接口
func (m *TrainPublishReqMessage) IsReqMessage() {}

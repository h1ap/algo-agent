package train

// TrainJobStatus 训练作业状态
type TrainJobStatus struct {
	Code int32
	Name string
}

var (
	// TrainJobStatusUnstart 未启动
	TrainJobStatusUnstart = TrainJobStatus{
		Code: 0,
		Name: "未启动",
	}

	// TrainJobStatusPreparing 准备中
	TrainJobStatusPreparing = TrainJobStatus{
		Code: 1,
		Name: "准备中",
	}

	// TrainJobStatusStarting 启动中
	TrainJobStatusStarting = TrainJobStatus{
		Code: 2,
		Name: "启动中",
	}

	// TrainJobStatusDownloadTrainImage 下载训练镜像
	TrainJobStatusDownloadTrainImage = TrainJobStatus{
		Code: 3,
		Name: "下载训练镜像",
	}

	// TrainJobStatusDownloadAlgoScripts 下载算法镜像
	TrainJobStatusDownloadAlgoScripts = TrainJobStatus{
		Code: 4,
		Name: "下载算法镜像",
	}

	// TrainJobStatusDownloadTrainData 下载训练数据
	TrainJobStatusDownloadTrainData = TrainJobStatus{
		Code: 5,
		Name: "下载训练数据",
	}

	// TrainJobStatusDownloadPreWeights 下载初始权重
	TrainJobStatusDownloadPreWeights = TrainJobStatus{
		Code: 6,
		Name: "下载初始权重",
	}

	// TrainJobStatusDownloadCheckpointFile 下载检查点文件
	TrainJobStatusDownloadCheckpointFile = TrainJobStatus{
		Code: 61,
		Name: "下载检查点文件",
	}

	// TrainJobStatusRunning 运行中
	TrainJobStatusRunning = TrainJobStatus{
		Code: 7,
		Name: "运行中",
	}

	// TrainJobStatusStopped 已终止
	TrainJobStatusStopped = TrainJobStatus{
		Code: 8,
		Name: "已终止",
	}

	// TrainJobStatusFail 已失败
	TrainJobStatusFail = TrainJobStatus{
		Code: 9,
		Name: "已失败",
	}

	// TrainJobStatusSucceed 已完成
	TrainJobStatusSucceed = TrainJobStatus{
		Code: 10,
		Name: "已完成",
	}
)

// GetCode 获取状态编码
func (t TrainJobStatus) GetCode() int32 {
	return t.Code
}

// GetName 获取状态名称
func (t TrainJobStatus) GetName() string {
	return t.Name
}

// GetByCode 根据编码获取对应的状态
func GetTrainJobStatusByCode(code int32) *TrainJobStatus {
	allStatuses := []TrainJobStatus{
		TrainJobStatusUnstart,
		TrainJobStatusPreparing,
		TrainJobStatusStarting,
		TrainJobStatusDownloadTrainImage,
		TrainJobStatusDownloadAlgoScripts,
		TrainJobStatusDownloadTrainData,
		TrainJobStatusDownloadPreWeights,
		TrainJobStatusDownloadCheckpointFile,
		TrainJobStatusRunning,
		TrainJobStatusStopped,
		TrainJobStatusFail,
		TrainJobStatusSucceed,
	}

	for _, status := range allStatuses {
		if status.Code == code {
			return &status
		}
	}
	return nil
}

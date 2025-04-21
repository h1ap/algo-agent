package data

import (
	"algo-agent/internal/cons/gpu"
	"algo-agent/internal/mq/event"
	"bufio"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

// NvidiaGpuManager NVIDIA GPU管理器实现
type NvidiaGpuManager struct {
	log *log.Helper
}

// GetGpuInfo 获取GPU信息
func (m *NvidiaGpuManager) GetGpuInfo() []event.GpuInfo {
	return m.getNvidiaGpuMetrics()
}

// getNvidiaGpuMetrics 调用nvidia驱动程序所带命令，获取nvidia显卡运行参数信息
func (m *NvidiaGpuManager) getNvidiaGpuMetrics() []event.GpuInfo {
	gpuList := make([]event.GpuInfo, 0)

	// 执行 nvidia-smi 命令
	cmd := exec.Command("nvidia-smi", "--query-gpu=index,name,utilization.gpu,memory.total,memory.used,gpu_uuid", "--format=csv,noheader,nounits")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		m.log.Infof("Error creating stdout pipe for nvidia-smi command: %v", err)
		return gpuList
	}

	if err := cmd.Start(); err != nil {
		m.log.Infof("Error starting nvidia-smi command: %v", err)
		return gpuList
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		metrics := strings.Split(line, ", ")

		if len(metrics) >= 6 {
			index, err := strconv.Atoi(metrics[0])
			if err != nil {
				continue
			}

			cudaUsageStr := metrics[2]
			cudaUsage, err := strconv.ParseFloat(cudaUsageStr, 64)
			if err != nil {
				continue
			}

			memoryTotal, err := strconv.ParseInt(metrics[3], 10, 64)
			if err != nil {
				continue
			}

			memoryUsed, err := strconv.ParseInt(metrics[4], 10, 64)
			if err != nil {
				continue
			}

			gpuInfo := event.GpuInfo{
				Index:       index,
				Name:        metrics[1],
				CudaUsage:   float64(int(cudaUsage*100)) / 100.0, // 保留两位小数
				MemoryTotal: memoryTotal,
				MemoryUsed:  memoryUsed,
				DeviceId:    metrics[5],
				Vendor:      m.GetIdentify().GetVendor(),
			}

			gpuList = append(gpuList, gpuInfo)
		}
	}

	if err := cmd.Wait(); err != nil {
		m.log.Infof("Error waiting for nvidia-smi command: %v", err)
	}

	return gpuList
}

// GetIdentify 获取GPU唯一标识
func (m *NvidiaGpuManager) GetIdentify() gpu.GpuVendor {
	return gpu.NVIDIA
}

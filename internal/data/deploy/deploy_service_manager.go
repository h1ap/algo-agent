package deploy

import (
	"algo-agent/internal/conf"
	"algo-agent/internal/cons/file"
	"algo-agent/internal/model/deploy"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/log"
)

type DeployServiceManager struct {
	serviceList []deploy.DeployServiceInfo

	filePath        string // 部署服务文件路径
	mappingFilePath string // 部署服务映射文件路径

	addChan    chan addRequest
	removeChan chan removeRequest
	getChan    chan chan []deploy.DeployServiceInfo
	updateChan chan updateRequest
	stopChan   chan struct{}
	writeChan  chan chan struct{}
}

// 添加服务的请求结构
type addRequest struct {
	service deploy.DeployServiceInfo
	result  chan error // 用于返回添加结果
}

// 删除服务的请求结构
type removeRequest struct {
	id     string
	result chan bool // 用于返回删除结果
}

// 更新服务的请求结构
type updateRequest struct {
	service deploy.DeployServiceInfo
	result  chan error // 用于返回更新结果
}

func NewDeployServiceManager(c *conf.Data, logger log.Logger) *DeployServiceManager {
	manager := &DeployServiceManager{
		filePath:        file.DEPLOY + file.SEPARATOR + "deploy.json",
		mappingFilePath: c.MappingFilePath,

		serviceList: make([]deploy.DeployServiceInfo, 0),
		addChan:     make(chan addRequest),
		removeChan:  make(chan removeRequest),
		getChan:     make(chan chan []deploy.DeployServiceInfo),
		updateChan:  make(chan updateRequest),
		stopChan:    make(chan struct{}),
		writeChan:   make(chan chan struct{}),
	}

	// 加载服务列表
	manager.loadServicesFromFile()

	// 启动处理 goroutine
	go manager.process()
	return manager
}

// 加载服务列表
func (m *DeployServiceManager) loadServicesFromFile() {
	filePath := filepath.Join(m.mappingFilePath, m.filePath)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Info("服务文件不存在，初始化空服务列表")
		m.serviceList = make([]deploy.DeployServiceInfo, 0)
		return
	}

	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Error("读取服务文件失败", err)
		return
	}

	// 检查文件是否为空
	if len(data) == 0 {
		log.Info("服务文件为空，初始化空服务列表")
		m.serviceList = make([]deploy.DeployServiceInfo, 0)
		return
	}

	// 解析文件内容
	var services []deploy.DeployServiceInfo
	if err := json.Unmarshal(data, &services); err != nil {
		log.Error("解析服务文件失败", err)
		return
	}

	// 更新服务列表
	m.serviceList = services
	log.Info("成功从文件加载了服务列表")
}

func (m *DeployServiceManager) process() {
	for {
		select {
		case req := <-m.addChan:
			// 检查是否已存在相同 serviceId 的服务
			exists := false
			for _, existingService := range m.serviceList {
				if existingService.ServiceId == req.service.ServiceId {
					exists = true
					break
				}
			}

			if exists {
				// 返回错误
				req.result <- fmt.Errorf("服务ID %s 已存在，无法添加", req.service.ServiceId)
			} else {
				// 添加服务
				m.serviceList = append(m.serviceList, req.service)

				// 同步写入文件
				writeDone := make(chan struct{})
				m.writeChan <- writeDone
				m.writeServicesToFile()
				close(writeDone)

				// 返回成功
				req.result <- nil
			}

		case req := <-m.removeChan:
			found := false
			for i, s := range m.serviceList {
				if s.ServiceId == req.id {
					m.serviceList = append(m.serviceList[:i], m.serviceList[i+1:]...)
					found = true

					// 同步写入文件
					writeDone := make(chan struct{})
					m.writeChan <- writeDone
					m.writeServicesToFile()
					close(writeDone)

					break
				}
			}
			// 返回删除结果
			req.result <- found

		case responseChan := <-m.getChan:
			copyList := make([]deploy.DeployServiceInfo, len(m.serviceList))
			copy(copyList, m.serviceList)
			responseChan <- copyList

		case req := <-m.updateChan:
			// 处理更新请求
			found := false
			index := -1

			// 查找要更新的服务
			for i, s := range m.serviceList {
				if s.ServiceId == req.service.ServiceId {
					found = true
					index = i
					break
				}
			}

			if !found {
				req.result <- fmt.Errorf("服务ID %s 不存在，无法更新", req.service.ServiceId)
			} else {
				// 更新服务
				m.serviceList[index] = req.service

				// 同步写入文件
				writeDone := make(chan struct{})
				m.writeChan <- writeDone
				m.writeServicesToFile()
				close(writeDone)

				// 返回成功
				req.result <- nil
				log.Info(fmt.Sprintf("服务ID %s 已成功更新", req.service.ServiceId))
			}

		case done := <-m.writeChan:
			// 这个 case 只用于同步，实际写入操作在 addChan 和 removeChan 的 case 中
			// 我们可以添加其他保存触发逻辑
			// 不需要做任何事情，只需等待关闭
			<-done

		case <-m.stopChan:
			return
		}
	}
}

func (m *DeployServiceManager) writeServicesToFile() {
	filePath := filepath.Join(m.mappingFilePath, m.filePath)

	// 检查目录是否存在，如果不存在则创建
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("创建目录失败: %v\n", err)
		return
	}

	// 写入文件
	data, err := json.Marshal(m.serviceList)
	if err != nil {
		fmt.Printf("序列化数据失败: %v\n", err)
		return
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
	}
}

// 添加服务，同步等待写入完成
func (m *DeployServiceManager) AddService(service deploy.DeployServiceInfo) error {
	resultChan := make(chan error)
	m.addChan <- addRequest{
		service: service,
		result:  resultChan,
	}
	return <-resultChan // 等待结果
}

// 删除服务，同步等待写入完成
func (m *DeployServiceManager) RemoveService(id string) bool {
	resultChan := make(chan bool)
	m.removeChan <- removeRequest{
		id:     id,
		result: resultChan,
	}
	return <-resultChan // 等待结果
}

// 获取服务列表
func (m *DeployServiceManager) GetServiceList() []deploy.DeployServiceInfo {
	responseChan := make(chan []deploy.DeployServiceInfo)
	m.getChan <- responseChan
	return <-responseChan
}

// 根据ID查找服务
func (m *DeployServiceManager) FindServiceById(serviceId string) *deploy.DeployServiceInfo {
	serviceList := m.GetServiceList()
	for _, service := range serviceList {
		if service.ServiceId == serviceId {
			return &service
		}
	}
	return nil
}

// 更新服务，线程安全
func (m *DeployServiceManager) UpdateService(updatedService deploy.DeployServiceInfo) error {
	resultChan := make(chan error)
	m.updateChan <- updateRequest{
		service: updatedService,
		result:  resultChan,
	}
	return <-resultChan // 等待结果
}

// 手动触发写入文件
func (m *DeployServiceManager) SaveToFile() {
	writeDone := make(chan struct{})
	m.writeChan <- writeDone
	<-writeDone // 等待写入完成
}

// 停止管理器
func (m *DeployServiceManager) Stop() {
	close(m.stopChan)
}

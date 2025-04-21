package data

import (
	"algo-agent/internal/model/deploy"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
)

type DeployServiceManagerRepo struct {
	serviceList []*deploy.DeployServiceInfo
	mu          sync.RWMutex // 添加互斥锁保护服务列表

	filePath        string // 部署服务文件路径
	mappingFilePath string // 部署服务映射文件路径

	log *log.Helper
}

// 创建新的DeployServiceManagerRepo实例
func NewDeployServiceManager(filePath, mappingFilePath string, logger log.Logger) *DeployServiceManagerRepo {
	repo := &DeployServiceManagerRepo{
		serviceList:     make([]*deploy.DeployServiceInfo, 0),
		filePath:        filePath,
		mappingFilePath: mappingFilePath,
		log:             log.NewHelper(logger),
	}

	// 加载服务列表
	repo.loadServicesFromFile(context.Background())
	return repo
}

// 加载服务列表
func (m *DeployServiceManagerRepo) loadServicesFromFile(ctx context.Context) {
	filePath := filepath.Join(m.mappingFilePath, m.filePath)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		m.log.WithContext(ctx).Info("服务文件不存在，初始化空服务列表")
		m.serviceList = make([]*deploy.DeployServiceInfo, 0)
		return
	}

	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		m.log.WithContext(ctx).Errorf("读取服务文件失败: %v", err)
		return
	}

	// 检查文件是否为空
	if len(data) == 0 {
		m.log.WithContext(ctx).Info("服务文件为空，初始化空服务列表")
		m.serviceList = make([]*deploy.DeployServiceInfo, 0)
		return
	}

	// 解析文件内容
	var services []*deploy.DeployServiceInfo
	if err := json.Unmarshal(data, &services); err != nil {
		m.log.WithContext(ctx).Errorf("解析服务文件失败: %v", err)
		return
	}

	// 更新服务列表
	m.mu.Lock()
	defer m.mu.Unlock() // 使用defer确保锁被释放

	m.serviceList = services
	m.log.WithContext(ctx).Info("成功从文件加载了服务列表")
}

func (m *DeployServiceManagerRepo) writeServicesToFile(ctx context.Context) error {
	filePath := filepath.Join(m.mappingFilePath, m.filePath)

	// 检查目录是否存在，如果不存在则创建
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		m.log.WithContext(ctx).Errorf("创建目录失败: %v", err)
		return err
	}

	// 使用读锁保护服务列表读取
	m.mu.RLock()
	data, err := json.Marshal(m.serviceList)
	m.mu.RUnlock() // 立即释放读锁，不用defer，避免在文件IO期间持有锁

	if err != nil {
		m.log.WithContext(ctx).Errorf("序列化数据失败: %v", err)
		return err
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		m.log.WithContext(ctx).Errorf("写入文件失败: %v", err)
		return err
	}

	return nil
}

// 添加服务
func (m *DeployServiceManagerRepo) AddService(ctx context.Context, service *deploy.DeployServiceInfo) error {
	// 第一阶段：检查和添加服务（需要锁）
	serviceAdded := func() bool {
		m.mu.Lock()
		defer m.mu.Unlock() // 使用defer确保锁被释放

		// 检查是否已存在相同 serviceId 的服务
		for _, existingService := range m.serviceList {
			if existingService.ServiceId == service.ServiceId {
				return false
			}
		}

		// 添加服务
		m.serviceList = append(m.serviceList, service)
		return true
	}()

	if !serviceAdded {
		return fmt.Errorf("服务ID %s 已存在，无法添加", service.ServiceId)
	}

	// 第二阶段：写入文件（不需要锁）
	err := m.writeServicesToFile(ctx)
	if err == nil {
		m.log.WithContext(ctx).Infof("服务ID %s 已成功添加", service.ServiceId)
	}

	return err
}

// 删除服务
func (m *DeployServiceManagerRepo) RemoveService(ctx context.Context, id string) bool {
	// 第一阶段：查找和删除服务（需要锁）
	var found bool
	func() {
		m.mu.Lock()
		defer m.mu.Unlock() // 使用defer确保锁被释放

		for i, s := range m.serviceList {
			if s.ServiceId == id {
				m.serviceList = append(m.serviceList[:i], m.serviceList[i+1:]...)
				found = true
				break
			}
		}
	}()

	// 第二阶段：如果找到并删除，则写入文件（不需要锁）
	if found {
		err := m.writeServicesToFile(ctx)
		if err != nil {
			m.log.WithContext(ctx).Errorf("删除服务后写入文件失败: %v", err)
		} else {
			m.log.WithContext(ctx).Infof("服务ID %s 已成功删除", id)
		}
	}

	return found
}

// 获取服务列表
func (m *DeployServiceManagerRepo) GetServiceList(ctx context.Context) []*deploy.DeployServiceInfo {
	m.mu.RLock()
	defer m.mu.RUnlock() // 使用defer确保锁被释放

	copyList := make([]*deploy.DeployServiceInfo, len(m.serviceList))
	copy(copyList, m.serviceList)
	return copyList
}

// 根据ID查找服务
func (m *DeployServiceManagerRepo) FindServiceById(ctx context.Context, serviceId string) *deploy.DeployServiceInfo {
	m.mu.RLock()
	defer m.mu.RUnlock() // 使用defer确保锁被释放

	for _, service := range m.serviceList {
		if service.ServiceId == serviceId {
			// 返回一个副本以避免竞态条件
			copiedService := *service // 假设DeployServiceInfo可以通过值复制
			return &copiedService
		}
	}
	return nil
}

// 更新服务
func (m *DeployServiceManagerRepo) UpdateService(ctx context.Context, updatedService *deploy.DeployServiceInfo) error {
	// 第一阶段：查找和更新服务（需要锁）
	var updated bool
	var err error
	func() {
		m.mu.Lock()
		defer m.mu.Unlock() // 使用defer确保锁被释放

		found := false
		index := -1

		// 查找要更新的服务
		for i, s := range m.serviceList {
			if s.ServiceId == updatedService.ServiceId {
				found = true
				index = i
				break
			}
		}

		if !found {
			err = fmt.Errorf("服务ID %s 不存在，无法更新", updatedService.ServiceId)
			return
		}

		// 更新服务
		m.serviceList[index] = updatedService
		updated = true
	}()

	if err != nil {
		return err
	}

	// 第二阶段：如果更新成功，则写入文件（不需要锁）
	if updated {
		err = m.writeServicesToFile(ctx)
		if err == nil {
			m.log.WithContext(ctx).Infof("服务ID %s 已成功更新", updatedService.ServiceId)
		}
	}

	return err
}

// 手动触发写入文件
func (m *DeployServiceManagerRepo) SaveToFile(ctx context.Context) {
	err := m.writeServicesToFile(ctx)
	if err != nil {
		m.log.WithContext(ctx).Errorf("手动保存文件失败: %v", err)
	}
}

// 停止管理器 - 在基于锁的实现中不需要特别的停止逻辑
func (m *DeployServiceManagerRepo) Stop(ctx context.Context) {
	// 可以在这里添加资源清理代码
	m.log.WithContext(ctx).Info("管理器已停止")
}

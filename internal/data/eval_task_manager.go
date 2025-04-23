package data

import (
	"algo-agent/internal/model/eval"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
)

// EvalTaskManagerRepo 评估任务管理仓库
type EvalTaskManagerRepo struct {
	taskList []*eval.EvalTaskInfo
	mu       sync.RWMutex // 互斥锁保护任务列表

	filePath        string // 评估任务文件路径
	mappingFilePath string // 评估任务映射文件路径

	log *log.Helper
}

// NewEvalTaskManager 创建新的EvalTaskManagerRepo实例
func NewEvalTaskManager(filePath, mappingFilePath string, logger log.Logger) *EvalTaskManagerRepo {
	repo := &EvalTaskManagerRepo{
		taskList:        make([]*eval.EvalTaskInfo, 0),
		filePath:        filePath,
		mappingFilePath: mappingFilePath,
		log:             log.NewHelper(logger),
	}

	// 加载任务列表
	repo.loadTasksFromFile(context.Background())
	return repo
}

// 加载任务列表
func (m *EvalTaskManagerRepo) loadTasksFromFile(ctx context.Context) {
	filePath := filepath.Join(m.mappingFilePath, m.filePath)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		m.log.WithContext(ctx).Info("评估任务文件不存在，初始化空任务列表")
		m.taskList = make([]*eval.EvalTaskInfo, 0)
		return
	}

	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		m.log.WithContext(ctx).Errorf("读取评估任务文件失败: %v", err)
		return
	}

	// 检查文件是否为空
	if len(data) == 0 {
		m.log.WithContext(ctx).Info("评估任务文件为空，初始化空任务列表")
		m.taskList = make([]*eval.EvalTaskInfo, 0)
		return
	}

	// 解析文件内容
	var tasks []*eval.EvalTaskInfo
	if err := json.Unmarshal(data, &tasks); err != nil {
		m.log.WithContext(ctx).Errorf("解析评估任务文件失败: %v", err)
		return
	}

	// 更新任务列表
	m.mu.Lock()
	defer m.mu.Unlock() // 使用defer确保锁被释放

	m.taskList = tasks
	m.log.WithContext(ctx).Info("成功从文件加载了评估任务列表")
}

// writeTasksToFile 写入任务列表到文件
func (m *EvalTaskManagerRepo) writeTasksToFile(ctx context.Context) error {
	filePath := filepath.Join(m.mappingFilePath, m.filePath)

	// 检查目录是否存在，如果不存在则创建
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		m.log.WithContext(ctx).Errorf("创建目录失败: %v", err)
		return err
	}

	// 使用读锁保护任务列表读取
	m.mu.RLock()
	data, err := json.Marshal(m.taskList)
	m.mu.RUnlock() // 立即释放读锁，避免在文件IO期间持有锁

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

// AddTask 添加任务
func (m *EvalTaskManagerRepo) AddTask(ctx context.Context, task *eval.EvalTaskInfo) error {
	// 第一阶段：检查和添加任务（需要锁）
	taskAdded := func() bool {
		m.mu.Lock()
		defer m.mu.Unlock() // 使用defer确保锁被释放

		// 检查是否已存在相同 taskId 的任务
		for _, existingTask := range m.taskList {
			if existingTask.TaskId == task.TaskId {
				return false
			}
		}

		// 添加任务
		m.taskList = append(m.taskList, task)
		return true
	}()

	if !taskAdded {
		return fmt.Errorf("任务ID %s 已存在，无法添加", task.TaskId)
	}

	// 第二阶段：写入文件（不需要锁）
	err := m.writeTasksToFile(ctx)
	if err == nil {
		m.log.WithContext(ctx).Infof("评估任务ID %s 已成功添加", task.TaskId)
	}

	return err
}

// RemoveTask 删除任务
func (m *EvalTaskManagerRepo) RemoveTask(ctx context.Context, id string) bool {
	// 第一阶段：查找和删除任务（需要锁）
	var found bool
	func() {
		m.mu.Lock()
		defer m.mu.Unlock() // 使用defer确保锁被释放

		for i, t := range m.taskList {
			if t.TaskId == id {
				m.taskList = append(m.taskList[:i], m.taskList[i+1:]...)
				found = true
				break
			}
		}
	}()

	// 第二阶段：如果找到并删除，则写入文件（不需要锁）
	if found {
		err := m.writeTasksToFile(ctx)
		if err != nil {
			m.log.WithContext(ctx).Errorf("删除评估任务后写入文件失败: %v", err)
		} else {
			m.log.WithContext(ctx).Infof("评估任务ID %s 已成功删除", id)
		}
	}

	return found
}

// GetTaskList 获取任务列表
func (m *EvalTaskManagerRepo) GetTaskList(ctx context.Context) []*eval.EvalTaskInfo {
	m.mu.RLock()
	defer m.mu.RUnlock() // 使用defer确保锁被释放

	copyList := make([]*eval.EvalTaskInfo, len(m.taskList))
	copy(copyList, m.taskList)
	return copyList
}

// FindTaskById 根据ID查找任务
func (m *EvalTaskManagerRepo) FindTaskById(ctx context.Context, taskId string) *eval.EvalTaskInfo {
	m.mu.RLock()
	defer m.mu.RUnlock() // 使用defer确保锁被释放

	for _, task := range m.taskList {
		if task.TaskId == taskId {
			// 返回一个副本以避免竞态条件
			copiedTask := *task
			return &copiedTask
		}
	}
	return nil
}

// UpdateTask 更新任务
func (m *EvalTaskManagerRepo) UpdateTask(ctx context.Context, updatedTask *eval.EvalTaskInfo) error {
	// 第一阶段：查找和更新任务（需要锁）
	var updated bool
	var err error
	func() {
		m.mu.Lock()
		defer m.mu.Unlock() // 使用defer确保锁被释放

		found := false
		index := -1

		// 查找要更新的任务
		for i, t := range m.taskList {
			if t.TaskId == updatedTask.TaskId {
				found = true
				index = i
				break
			}
		}

		if !found {
			err = fmt.Errorf("评估任务ID %s 不存在，无法更新", updatedTask.TaskId)
			return
		}

		// 更新任务
		m.taskList[index] = updatedTask
		updated = true
	}()

	if err != nil {
		return err
	}

	// 第二阶段：如果更新成功，则写入文件（不需要锁）
	if updated {
		err = m.writeTasksToFile(ctx)
		if err == nil {
			m.log.WithContext(ctx).Infof("评估任务ID %s 已成功更新", updatedTask.TaskId)
		}
	}

	return err
}

// SaveToFile 保存到文件
func (m *EvalTaskManagerRepo) SaveToFile(ctx context.Context) {
	err := m.writeTasksToFile(ctx)
	if err != nil {
		m.log.WithContext(ctx).Errorf("手动保存评估任务文件失败: %v", err)
	}
}

// Stop 停止管理器
func (m *EvalTaskManagerRepo) Stop(ctx context.Context) {
	// 可以在这里添加资源清理代码
	m.log.WithContext(ctx).Info("评估任务管理器已停止")
}

package container

// ContainerInfo 容器信息
type ContainerInfo struct {
	ContainerID   string
	ContainerName string
}

// ContainerState 容器状态
type ContainerState struct {
	ContainerID string
	State       string
	ExitCode    int
}

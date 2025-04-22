package container

// ContainerInfo 容器信息
type ContainerInfo struct {
	ContainerId   string
	ContainerName string
}

// ContainerState 容器状态
type ContainerState struct {
	ContainerId string
	State       string
	ExitCode    int
}

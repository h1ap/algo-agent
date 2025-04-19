package container

// DockerContainerState 表示容器的状态
type DockerContainerState string

const (
	// RUNNING 表示容器正在运行
	RUNNING DockerContainerState = "running"

	// STOPPED 表示容器已停止
	STOPPED DockerContainerState = "stopped"

	// PAUSED 表示容器已暂停
	PAUSED DockerContainerState = "paused"

	// RESTARTING 表示容器正在重启
	RESTARTING DockerContainerState = "restarting"

	// CREATED 表示容器已创建
	CREATED DockerContainerState = "created"

	// DESTROYED 表示容器已销毁
	DESTROYED DockerContainerState = "destroyed"

	// EXITED 表示容器已退出
	EXITED DockerContainerState = "exited"
)

// String 实现Stringer接口，返回容器状态的字符串表示
func (s DockerContainerState) String() string {
	return string(s)
}

package data

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"algo-agent/internal/conf"
	mc "algo-agent/internal/model/container"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-kratos/kratos/v2/log"
)

// DockerRepo 实现Docker客户端，满足biz.DockerService接口
type DockerRepo struct {
	client     *client.Client
	conf       *conf.Data_Docker
	log        *log.Helper
	lock       sync.Mutex
	logStreams map[string]struct {
		close     func() error
		isRunning bool
	}
}

// FindContainerByName 根据容器名查找容器
func (r *DockerRepo) FindContainerByName(ctx context.Context, containerName string) (*mc.ContainerInfo, error) {
	// Docker API的容器名称前面会加"/"
	dockerAPIContainerName := "/" + containerName

	args := filters.NewArgs()
	args.Add("name", containerName)

	containers, err := r.client.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: args,
	})
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to list containers: %v", err)
		return nil, err
	}

	for _, c := range containers {
		for _, name := range c.Names {
			if name == dockerAPIContainerName {
				return &mc.ContainerInfo{
					ContainerID:   c.ID,
					ContainerName: strings.TrimPrefix(name, "/"),
				}, nil
			}
		}
	}

	return nil, nil
}

// RunContainer 通过镜像和自定义参数运行一个容器
func (r *DockerRepo) RunContainer(ctx context.Context, imageName string, customArgs []string) (string, error) {
	hostConfig := &container.HostConfig{}

	resp, err := r.client.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageName,
			Cmd:   customArgs,
			Tty:   true,
		},
		hostConfig,
		nil,
		nil,
		"",
	)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to create container: %v", err)
		return "", err
	}

	err = r.client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to start container: %v", err)
		return "", err
	}

	return resp.ID, nil
}

// RunAndStartContainer 运行并启动容器，返回容器信息
func (r *DockerRepo) RunAndStartContainer(ctx context.Context, imageName, hostPath, containerPath, scriptPath string, args []string) (*mc.ContainerInfo, error) {
	// 配置主机卷挂载
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: hostPath,
				Target: containerPath,
			},
		},
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
		Resources: container.Resources{
			Memory: 16 * 1024 * 1024 * 1024, // 16GB
		},
	}

	// 检查是否需要使用GPU
	useGPU := false
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "--device_type" && strings.EqualFold(args[i+1], "GPU") {
			useGPU = true
			break
		}
	}

	if useGPU {
		r.log.WithContext(ctx).Info("GPU is used")
		hostConfig.DeviceRequests = []container.DeviceRequest{
			{
				Driver:       "nvidia",
				Count:        -1, // -1 表示使用所有可用的 GPU
				Capabilities: [][]string{{"gpu"}},
			},
		}
	}

	// 拼接启动命令
	cmd := []string{"python", scriptPath}
	cmd = append(cmd, args...)

	r.log.WithContext(ctx).Infof("Container creation configuration: image=%s, host_path=%s, container_path=%s",
		imageName, hostPath, containerPath)

	// 创建容器
	resp, err := r.client.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageName,
			Cmd:   cmd,
			Tty:   true,
		},
		hostConfig,
		nil,
		nil,
		"",
	)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to create container: %v", err)
		return nil, err
	}

	// 启动容器
	err = r.client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to start container: %v", err)
		return nil, err
	}

	// 获取容器信息
	info, err := r.client.ContainerInspect(ctx, resp.ID)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to inspect container: %v", err)
		return nil, err
	}

	containerName := info.Name
	if strings.HasPrefix(containerName, "/") {
		containerName = containerName[1:]
	}

	return &mc.ContainerInfo{
		ContainerID:   resp.ID,
		ContainerName: containerName,
	}, nil
}

// RunAndStartContainerWithPort 运行并启动带端口映射的容器
func (r *DockerRepo) RunAndStartContainerWithPort(ctx context.Context, imageName, hostPath, containerPath, command string, args []string, hostPort int) (string, error) {
	// 配置主机卷挂载
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: hostPath,
				Target: containerPath,
			},
		},
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
		Resources: container.Resources{
			Memory: 16 * 1024 * 1024 * 1024, // 16GB
		},
	}

	// 配置端口映射
	portMap := nat.PortMap{
		nat.Port("5000/tcp"): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: fmt.Sprintf("%d", hostPort),
			},
		},
	}
	hostConfig.PortBindings = portMap

	// 拼接启动命令
	cmd := []string{"python", command}
	cmd = append(cmd, args...)

	r.log.WithContext(ctx).Infof("Container creation configuration: image=%s, host_path=%s, container_path=%s, port_mapping=%d:5000",
		imageName, hostPath, containerPath, hostPort)

	// 创建容器
	resp, err := r.client.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageName,
			Cmd:   cmd,
			Tty:   true,
			ExposedPorts: nat.PortSet{
				nat.Port("5000/tcp"): struct{}{},
			},
		},
		hostConfig,
		&network.NetworkingConfig{},
		nil,
		"",
	)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to create container: %v", err)
		return "", err
	}

	// 启动容器
	err = r.client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to start container: %v", err)
		return "", err
	}

	// 获取容器信息
	info, err := r.client.ContainerInspect(ctx, resp.ID)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to inspect container: %v", err)
		return "", err
	}

	containerName := info.Name
	if strings.HasPrefix(containerName, "/") {
		containerName = containerName[1:]
	}

	return containerName, nil
}

// GetContainerLastLogs 获取指定容器的最后若干行日志
func (r *DockerRepo) GetContainerLastLogs(ctx context.Context, containerID string, tail int) (string, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       fmt.Sprintf("%d", tail),
	}

	logs, err := r.client.ContainerLogs(ctx, containerID, options)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to get container logs: %v", err)
		return "", err
	}
	defer logs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, logs)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to read container logs: %v", err)
		return "", err
	}

	return buf.String(), nil
}

// processLogFrame 处理日志帧
func (r *DockerRepo) processLogFrame(log string, isError bool, logConsumer func(string)) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	var formattedLog string
	if isError {
		formattedLog = fmt.Sprintf("[%s] [ERROR] %s", timestamp, log)
	} else {
		formattedLog = fmt.Sprintf("[%s] [INFO] %s", timestamp, log)
	}
	logConsumer(formattedLog)
}

// StartLogStream 开始日志流
func (r *DockerRepo) StartLogStream(ctx context.Context, containerID string, logConsumer func(string)) {
	r.lock.Lock()
	defer r.lock.Unlock()

	// 如果已经有日志流在运行，关闭它
	if stream, exists := r.logStreams[containerID]; exists && stream.isRunning {
		if stream.close != nil {
			_ = stream.close()
		}
	}

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	}

	logs, err := r.client.ContainerLogs(ctx, containerID, options)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to start log stream: %v", err)
		return
	}

	r.logStreams[containerID] = struct {
		close     func() error
		isRunning bool
	}{
		close:     logs.Close,
		isRunning: true,
	}

	r.log.WithContext(ctx).Infof("started log streaming for container: %s", containerID)

	// 在单独的goroutine中处理日志流
	go func() {
		defer logs.Close()
		defer func() {
			r.lock.Lock()
			defer r.lock.Unlock()
			if stream, exists := r.logStreams[containerID]; exists {
				stream.isRunning = false
				r.logStreams[containerID] = stream
			}
		}()

		buf := make([]byte, 8192)
		for {
			// 检查容器是否仍在运行
			r.lock.Lock()
			stream, exists := r.logStreams[containerID]
			isRunning := exists && stream.isRunning
			r.lock.Unlock()

			if !isRunning {
				return
			}

			// 读取日志
			n, err := logs.Read(buf)
			if err != nil {
				if err != io.EOF {
					r.log.Errorf("error reading container logs: %v", err)
				}
				return
			}

			if n > 0 {
				// Docker日志格式：前8个字节为头部，第一个字节表示流类型（1=stdout，2=stderr）
				headerSize := 8
				for i := 0; i < n; i += headerSize {
					headerStart := i
					if headerStart+headerSize > n {
						break
					}

					// 确定有效载荷的大小
					payloadSize := 0
					for j := 4; j < 8; j++ {
						payloadSize = payloadSize*256 + int(buf[headerStart+j])
					}

					// 确定有效载荷的结束位置
					payloadEnd := headerStart + headerSize + payloadSize
					if payloadEnd > n {
						break
					}

					isError := buf[headerStart] == 2 // 2表示stderr
					payload := buf[headerStart+headerSize : payloadEnd]
					r.processLogFrame(string(payload), isError, logConsumer)

					// 移动到下一个消息
					i = payloadEnd - headerSize
				}
			}
		}
	}()
}

// StopLogStream 停止日志流
func (r *DockerRepo) StopLogStream(ctx context.Context, containerName string, logConsumer func(string)) {
	containerInfo, err := r.FindContainerByName(ctx, containerName)
	if err != nil || containerInfo == nil {
		r.log.WithContext(ctx).Errorf("failed to find container: %s, error: %v", containerName, err)
		return
	}

	containerID := containerInfo.ContainerID

	// 获取最后一行日志
	lastLog, err := r.GetContainerLastLogs(ctx, containerID, 1)
	if err == nil && lastLog != "" {
		r.processLogFrame(lastLog, true, logConsumer)
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	// 关闭日志流
	if stream, exists := r.logStreams[containerID]; exists {
		stream.isRunning = false
		if stream.close != nil {
			_ = stream.close()
		}
		delete(r.logStreams, containerID)
	}

	r.log.WithContext(ctx).Infof("stopped log streaming for container: %s", containerID)
}

// StopContainer 停止容器
func (r *DockerRepo) StopContainer(ctx context.Context, containerID string) error {
	r.log.WithContext(ctx).Infof("stopping container: %s", containerID)
	err := r.client.ContainerStop(ctx, containerID, container.StopOptions{})
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to stop container: %v", err)
		return err
	}
	return nil
}

// StopContainerByName 通过名称停止容器
func (r *DockerRepo) StopContainerByName(ctx context.Context, containerName string, remove bool) error {
	if containerName == "" {
		return errors.New("container name is empty")
	}

	// 查找容器
	containerInfo, err := r.FindContainerByName(ctx, containerName)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to find container: %v", err)
		return err
	}

	if containerInfo == nil {
		r.log.WithContext(ctx).Infof("no container found with name: %s", containerName)
		return nil
	}

	// 获取容器详情
	inspect, err := r.client.ContainerInspect(ctx, containerInfo.ContainerID)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to inspect container: %v", err)
		return err
	}

	// 如果容器正在运行，则停止容器
	if inspect.State.Running {
		err = r.client.ContainerStop(ctx, containerInfo.ContainerID, container.StopOptions{})
		if err != nil {
			r.log.WithContext(ctx).Errorf("failed to stop container: %v", err)
			return err
		}
		r.log.WithContext(ctx).Infof("container stopped: %s", containerName)
	} else {
		r.log.WithContext(ctx).Infof("container already stopped: %s", containerName)
	}

	// 如果需要删除容器
	if remove {
		err = r.client.ContainerRemove(ctx, containerInfo.ContainerID, container.RemoveOptions{})
		if err != nil {
			r.log.WithContext(ctx).Errorf("failed to remove container: %v", err)
			return err
		}
		r.log.WithContext(ctx).Infof("container removed: %s", containerName)
	}

	return nil
}

// GetContainerStopTime 获取容器的停止时间戳
func (r *DockerRepo) GetContainerStopTime(ctx context.Context, containerID string) (int64, error) {
	inspect, err := r.client.ContainerInspect(ctx, containerID)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to inspect container: %v", err)
		return 0, err
	}

	if inspect.State.Running {
		return 0, nil
	}

	finishedAt := inspect.State.FinishedAt
	if finishedAt == "" {
		return 0, nil
	}

	// Docker时间格式为: yyyy-MM-dd'T'HH:mm:ss.SSSSSSSSZ
	t, err := time.Parse(time.RFC3339Nano, finishedAt)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to parse container stop time: %v", err)
		return 0, err
	}

	return t.Unix(), nil
}

// FindImageByName 查找指定名称的镜像
func (r *DockerRepo) FindImageByName(ctx context.Context, imageName string) (bool, error) {
	// 检查镜像名称是否包含标签，如果不包含，则添加:latest标签
	if !strings.Contains(imageName, ":") {
		imageName = imageName + ":latest"
	}

	// 列出所有镜像
	images, err := r.client.ImageList(ctx, image.ListOptions{})
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to list images: %v", err)
		return false, err
	}

	// 调试信息：打印所有镜像信息
	r.log.WithContext(ctx).Info("Local mirror list: ")
	for _, img := range images {
		if img.RepoTags != nil {
			for _, tag := range img.RepoTags {
				r.log.WithContext(ctx).Info(tag)
			}
		}
	}

	// 查找匹配的镜像
	for _, img := range images {
		if img.RepoTags != nil {
			for _, tag := range img.RepoTags {
				if tag == imageName {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// ImportAndTagImage 导入并标记镜像
func (r *DockerRepo) ImportAndTagImage(ctx context.Context, tarFilePath, fullImageName string) error {
	if _, err := os.Stat(tarFilePath); err != nil {
		r.log.WithContext(ctx).Errorf("image file does not exist: %s", tarFilePath)
		return fmt.Errorf("image file does not exist: %s", tarFilePath)
	}

	// 分析fullImageName是否带有tag
	var repository, tag string
	parts := strings.Split(fullImageName, ":")
	if len(parts) == 1 {
		repository = fullImageName
		tag = "latest"
	} else {
		repository = parts[0]
		tag = parts[1]
	}

	finalFullImageName := repository + ":" + tag
	r.log.WithContext(ctx).Infof("importing image: %s", finalFullImageName)

	r.lock.Lock()
	defer r.lock.Unlock()

	// 导入前的镜像列表
	beforeImages, err := r.client.ImageList(ctx, image.ListOptions{})
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to list images: %v", err)
		return err
	}

	beforeImageIDs := make([]string, len(beforeImages))
	for i, img := range beforeImages {
		beforeImageIDs[i] = img.ID
	}

	// 打开文件
	file, err := os.Open(tarFilePath)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to read image file: %v", err)
		return err
	}

	// 导入镜像
	resp, err := r.client.ImageLoad(ctx, file)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to import image: %v", err)
		return err
	}
	defer resp.Body.Close()

	// 读取响应
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to read import response: %v", err)
		return err
	}

	// 导入后的镜像列表
	afterImages, err := r.client.ImageList(ctx, image.ListOptions{})
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to list images: %v", err)
		return err
	}

	// 找出新增的镜像
	var newImageID string
	found := false
	for _, img := range afterImages {
		isNew := true
		for _, beforeID := range beforeImageIDs {
			if img.ID == beforeID {
				isNew = false
				break
			}
		}
		if isNew {
			newImageID = img.ID
			found = true
			break
		}
	}

	if !found {
		return errors.New("no new image found after import")
	}

	// 为新镜像打上指定的标签
	err = r.client.ImageTag(ctx, newImageID, finalFullImageName)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to tag image: %v", err)
		return err
	}

	r.log.WithContext(ctx).Infof("image imported and tagged: %s", finalFullImageName)
	return nil
}

// Close 关闭连接
func (r *DockerRepo) Close() {
	r.lock.Lock()
	defer r.lock.Unlock()

	// 关闭所有日志流
	for containerID, stream := range r.logStreams {
		if stream.isRunning && stream.close != nil {
			_ = stream.close()
		}
		delete(r.logStreams, containerID)
	}

	// 关闭Docker客户端
	if r.client != nil {
		_ = r.client.Close()
	}
}

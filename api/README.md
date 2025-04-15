# API 说明文档

## 目录结构

```
api/
  oss/             # OSS服务API目录
    v1/            # 版本
      oss.proto    # OSS服务接口定义
  docker/          # Docker服务API目录
    v1/            # 版本
      docker.proto # Docker服务接口定义
  rabbitmq/        # RabbitMQ服务API目录
    v1/            # 版本
      rabbitmq.proto # RabbitMQ服务接口定义
```

## 生成Proto

在项目根目录下执行以下命令生成API相关文件：

```bash
# 安装依赖
make init

# 生成proto文件
make api
```

## 使用方法

1. 确保已经执行了上述的生成命令
2. 在服务实现中取消`internal/service/service.go`文件中的相关服务函数注释
3. 在服务注册处(如`internal/server/grpc.go`)注册相应服务实现:

```go
// RegisterGRPCServer 注册gRPC服务
func RegisterGRPCServer(server *grpc.Server, 
    ossServer *service.OSSServer,
    dockerServer *service.DockerServer,
    rabbitMQServer *service.RabbitMQServer) {
    
    // 注册OSS服务
    v1.RegisterOSSServiceServer(server, ossServer)
    
    // 注册Docker服务
    dockerv1.RegisterDockerServiceServer(server, dockerServer)
    
    // 注册RabbitMQ服务
    mqv1.RegisterRabbitMQServiceServer(server, rabbitMQServer)
}
```

## 服务说明

### OSS服务
OSS服务提供对象存储功能，包括文件上传、下载、删除等操作。

### Docker服务
Docker服务提供容器管理功能，包括查找容器、运行容器、停止容器等操作。

### RabbitMQ服务
RabbitMQ服务提供消息队列功能，包括发送消息到交换机、队列以及服务等操作。 
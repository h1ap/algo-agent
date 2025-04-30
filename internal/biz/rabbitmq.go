package biz

import (
	v1 "algo-agent/api/deploy/v1"
	"algo-agent/internal/cons/mq"
	"algo-agent/internal/mq/event"
	"context"
	"encoding/json"
	"runtime"

	"github.com/wagslane/go-rabbitmq"

	"github.com/go-kratos/kratos/v2/log"
)

// MqService 消息队列发送接口
type MqService interface {
	// SendMessage 发送字符串消息
	SendMessage(ctx context.Context, exchangeName, routingKey string, message *event.ReqMessage) error

	// SendToQueue 发送消息到队列
	SendToQueue(ctx context.Context, queueName string, message *event.ReqMessage) error

	// SendToService 发送消息到特定服务
	SendToService(ctx context.Context, service string, message *event.ReqMessage) error

	GetOrCreateConsumer(ctx context.Context) (*rabbitmq.Consumer, error)

	// Close 关闭连接
	Close()
}

// RabbitMQUsecase 是RabbitMQ用例
type RabbitMQUsecase struct {
	mq  MqService
	cs  *rabbitmq.Consumer
	log *log.Helper

	trainingUsecase *TrainingTaskUsecase
	evalUsecase     *EvalTaskUsecase
	deployUsecase   *DeployUsecase
	extractUsecase  *ExtractTaskUsecase
}

// SendMessage 发送消息
func (uc *RabbitMQUsecase) SendMessage(ctx context.Context, exchangeName, routingKey string, message *event.ReqMessage) error {
	uc.log.WithContext(ctx).Infof("SendMessage: exchangeName=%s, routingKey=%s, message=%s", exchangeName, routingKey, message)
	return uc.mq.SendMessage(ctx, exchangeName, routingKey, message)
}

// SendToQueue 发送消息到队列
func (uc *RabbitMQUsecase) SendToQueue(ctx context.Context, queueName string, message *event.ReqMessage) error {
	uc.log.WithContext(ctx).Infof("SendToQueue: queueName=%s, message=%s", queueName, message)
	return uc.mq.SendToQueue(ctx, queueName, message)
}

// SendToService 发送消息到特定服务
func (uc *RabbitMQUsecase) SendToService(ctx context.Context, service string, message *event.ReqMessage) error {
	uc.log.WithContext(ctx).Infof("SendToService: service=%s, message=%s", service, message)
	return uc.mq.SendToService(ctx, service, message)
}

func (uc *RabbitMQUsecase) Subscribe(ctx context.Context) error {
	consumer, err := uc.mq.GetOrCreateConsumer(ctx)

	// 设置消费者
	uc.cs = consumer

	if err != nil {
		uc.log.Errorf("创建RabbitMQ消费者失败: %v", err)
		return err
	}
	// 启动消费处理
	uc.log.Info("RabbitMQ消费者已启动")

	// 使用 defer 和 recover 捕获可能的 panic
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 2048)
			runtime.Stack(buf, false)
			uc.log.Errorf("RabbitMQ消费者运行时发生panic: %v\n堆栈信息:\n%s", r, buf)
			if uc.cs != nil {
				uc.cs.Close()
			}
		}
	}()

	err = uc.cs.Run(uc.messageHandler)
	if err != nil {
		uc.log.Errorf("启动RabbitMQ消费者失败: %v", err)
		if uc.cs != nil {
			uc.cs.Close()
		}
		return err
	}

	return nil
}

// messageHandler 处理收到的消息
func (uc *RabbitMQUsecase) messageHandler(d rabbitmq.Delivery) rabbitmq.Action {
	ctx := context.Background()
	message := string(d.Body)
	contentType := d.ContentType

	uc.log.Infof("收到消息: %s, contentType: %s", message, contentType)

	// 检查内容类型
	if contentType == "application/json" {
		// 尝试解析为不同类型的消息并处理
		uc.processJSONMessage(ctx, message)
	} else {
		// 处理纯文本消息
		uc.log.Infof("收到纯文本消息: %s", message)
	}

	// 确认消息已处理
	return rabbitmq.Ack
}

// processJSONMessage 处理JSON格式的消息
func (uc *RabbitMQUsecase) processJSONMessage(ctx context.Context, jsonMessage string) {

	var reqMessage event.ReqMessage
	unmarshalErr := json.Unmarshal([]byte(jsonMessage), &reqMessage)
	if unmarshalErr != nil {
		uc.log.Errorf("解析JSON消息失败: %v", unmarshalErr)
		return
	}

	if reqMessage.Type != mq.UNKNOWN.Code() {
		payloadJson, err := json.Marshal(reqMessage.Payload)
		if err != nil {
			uc.log.Errorf("序列化JSON消息失败: %v", err)
			return
		}

		// 利用switch匹配消息类型
		switch reqMessage.Type {
		case mq.TRAIN_TASK.Code():
			var trainTask event.TrainTaskReqMessage
			err = json.Unmarshal(payloadJson, &trainTask)

			if err != nil || trainTask.TaskId == "" {
				uc.log.Errorf("解析训练任务消息失败: %v", err)
				return
			}

			uc.log.Infof("处理训练任务消息, 任务ID: %s", trainTask.TaskId)
			uc.trainingUsecase.HandleTrainingTask(ctx, &trainTask)
		case mq.TASK_EVALUATE.Code():
			var evalTask event.EvalSendMessage
			err = json.Unmarshal(payloadJson, &evalTask)

			if err != nil || evalTask.TaskId == "" {
				uc.log.Errorf("解析评估任务消息失败: %v", err)
				return
			}

			uc.log.Infof("处理评估任务消息, 任务ID: %s", evalTask.TaskId)
			uc.evalUsecase.HandleEvalTask(ctx, &evalTask)
		case mq.TRAIN_PUBLISH.Code():
			var publishTask event.TrainPublishReqMessage
			err = json.Unmarshal(payloadJson, &publishTask)

			if err != nil || publishTask.TrainDetailId == "" {
				uc.log.Errorf("解析发布任务消息失败: %v", err)
				return
			}

			uc.log.Infof("处理发布任务消息, 训练详情ID: %s", publishTask.TrainDetailId)
			uc.extractUsecase.HandleExtractTask(ctx, &publishTask)
		case mq.TASK_DEPLOY.Code():
			var deployTask v1.DeployRequest
			err = json.Unmarshal(payloadJson, &deployTask)

			if err != nil || deployTask.ServiceId == "" {
				uc.log.Errorf("解析部署任务消息失败: %v", err)
				return
			}

			uc.log.Infof("处理部署任务消息, 服务ID: %s", deployTask.ServiceId)
			uc.deployUsecase.HandleEvent(ctx, &deployTask)
		default:
			// 无法识别的JSON消息
			uc.log.Warnf("无法识别的JSON消息类型: %s", jsonMessage)
		}
	}
}

// Close 关闭消费者
func (uc *RabbitMQUsecase) Unsubscribe() {
	if uc.cs != nil {
		uc.cs.Close()
		uc.log.Info("RabbitMQ消费者已关闭")
	}
}

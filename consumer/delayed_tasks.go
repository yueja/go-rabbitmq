package consumer

import (
	amqp "github.com/rabbitmq/amqp091-go"
	rabbit "github.com/yueja/rabbitmq/mq"
)

type DelayedTaskBody struct {
	TaskKind string      `json:"task_kind"` // 任务类型，自定义字符串
	Data     interface{} `json:"data"`
}

// DelayedTaskConsumer 延时任务
func (c *Consumer) DelayedTaskConsumer() (msgs <-chan amqp.Delivery, err error) {
	// 消费死信消息
	if msgs, err = c.consumer(rabbit.GetRabbitConfig().DeadQueueName); err != nil {
		return
	}
	return
}

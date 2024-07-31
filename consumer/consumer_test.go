package consumer

import (
	"fmt"
	"testing"
	"time"
	"yueja/go-rabbitmq/config"
	"yueja/go-rabbitmq/mq"
)

// Test_QueueConsumer 队列模式
func Test_QueueConsumer(t *testing.T) {
	rabbitmq.InitMQ(config.DefaultConfig)
	defer rabbitmq.CloseMQ()

	option := NewOption().SetConsumerOption(NewConsumerOption().SetConsumer("consumer-queue").SetAutoAck(false))
	consumer, err := NewConsumer(option)
	if err != nil {
		return
	}
	defer consumer.PoolReleaseConnection()

	msg, err := consumer.QueueConsumer("queue-01")
	if err != nil {
		return
	}

	for info := range msg {
		fmt.Println("消费者收到消息：", string(info.Body))
		// 默认为手动提交
		info.Ack(true)
	}
	fmt.Println("消费者退出")
	return
}

// Test_WorkQueueConsumer 多个消费者绑定到同一个队列，即为WorkQueue模式
func Test_WorkQueueConsumer(t *testing.T) {
	rabbitmq.InitMQ(config.DefaultConfig)
	defer rabbitmq.CloseMQ()

	ch := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func(num int) {
			option := NewOption().SetConsumerOption(NewConsumerOption().SetConsumer("consumer-work-queue").SetAutoAck(false))
			consumer, err := NewConsumer(option)
			if err != nil {
				return
			}
			defer consumer.PoolReleaseConnection()

			msg, err := consumer.QueueConsumer("queue-01")
			if err != nil {
				return
			}

			for info := range msg {
				fmt.Println(fmt.Sprintf("消费者%d收到消息：%s", num, string(info.Body)))
				info.Ack(true)
			}
		}(i)
	}

	<-ch
	fmt.Println("消费者退出")
	return
}

// Test_WorkQueueConsumerAutoAck 多个消费者绑定到同一个队列，即为WorkQueue模式，自动提交ACK
func Test_WorkQueueConsumerAutoAck(t *testing.T) {
	rabbitmq.InitMQ(config.DefaultConfig)
	defer rabbitmq.CloseMQ()

	ch := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func(num int) {
			option := NewOption().SetConsumerOption(NewConsumerOption().SetConsumer("consumer-work-queue-autoAck").SetAutoAck(true))

			consumer, err := NewConsumer(option)
			if err != nil {
				return
			}
			defer consumer.PoolReleaseConnection()

			msg, err := consumer.QueueConsumer("queue-01")
			if err != nil {
				return
			}

			for info := range msg {
				fmt.Println(fmt.Sprintf("消费者%d收到消息：%s", num, string(info.Body)))
			}
		}(i)
	}

	<-ch
	fmt.Println("消费者退出")
	return
}

// Test_PublishSubscribeConsumer 队列模式
func Test_PublishSubscribeConsumer(t *testing.T) {
	rabbitmq.InitMQ(config.DefaultConfig)
	defer rabbitmq.CloseMQ()

	option := NewOption().SetConsumerOption(NewConsumerOption().SetConsumer("consumer-publish-subscribe"))
	consumer, err := NewConsumer(option)
	if err != nil {
		return
	}
	defer consumer.PoolReleaseConnection()

	msg, err := consumer.PublishSubscribeConsumer("exchange-publish-subscribe", "queue-publish-subscribe")
	if err != nil {
		return
	}

	for info := range msg {
		fmt.Println("消费者收到消息：", string(info.Body))
		// 默认为手动提交
		info.Ack(true)
	}
	fmt.Println("消费者退出")
	return
}

// Test_RoutingKeyConsumer routing路由模式
func Test_RoutingKeyConsumer(t *testing.T) {
	rabbitmq.InitMQ(config.DefaultConfig)
	defer rabbitmq.CloseMQ()

	option := NewOption().SetConsumerOption(NewConsumerOption().SetConsumer("consumer-routing-key-01"))
	consumer, err := NewConsumer(option)
	if err != nil {
		return
	}
	defer consumer.PoolReleaseConnection()

	msg, err := consumer.RoutingKeyConsumer("exchange-routingKey", "routingKey-yja")
	if err != nil {
		return
	}

	for info := range msg {
		fmt.Println("消费者收到消息：", string(info.Body))
		// 默认为手动提交
		info.Ack(true)
	}

	fmt.Println("消费者退出")
	return
}

// Test_TopicConsumer topic模式
func Test_TopicConsumer(t *testing.T) {
	rabbitmq.InitMQ(config.DefaultConfig)
	defer rabbitmq.CloseMQ()

	option := NewOption().SetConsumerOption(NewConsumerOption().SetConsumer("consumer-topic-01"))
	consumer, err := NewConsumer(option)
	if err != nil {
		return
	}
	defer consumer.PoolReleaseConnection()

	msg, err := consumer.TopicConsumer("exchange-topic", "topic-yja")
	if err != nil {
		return
	}

	for info := range msg {
		fmt.Println("消费者收到消息：", string(info.Body))
		// 默认为手动提交
		info.Ack(true)
	}
	fmt.Println("消费者退出")
	return
}

// Test_DelayedTaskConsumer 延时任务
func Test_DelayedTaskConsumer(t *testing.T) {
	rabbitmq.InitMQ(config.DefaultConfig)
	defer rabbitmq.CloseMQ()

	option := NewOption().SetConsumerOption(NewConsumerOption().SetConsumer("consumer-delayedTask-01"))

	consumer, err := NewConsumer(option)
	if err != nil {
		return
	}
	defer consumer.PoolReleaseConnection()

	msg, err := consumer.DelayedTaskConsumer()
	if err != nil {
		return
	}

	for info := range msg {
		fmt.Println(fmt.Sprintf("消费者收到消息，时间：%+v,消息内容：%+v", time.Now(), string(info.Body)))
		// 默认为手动提交
		info.Ack(true)
	}
	fmt.Println("消费者退出")
	return
}

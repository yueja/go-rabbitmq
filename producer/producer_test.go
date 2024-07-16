package producer

import (
	"fmt"
	"github.com/yueja/rabbitmq/config"
	"github.com/yueja/rabbitmq/mq"
	"testing"
	"time"
)

// 发布完成自动释放连接（默认）
func Test_QueueProducer(t *testing.T) {
	rabbitmq.InitMQ(config.DefaultConfig)
	defer rabbitmq.CloseMQ()

	for i := 0; i < 10; i++ {
		producer, err := NewProducer(nil)
		if err != nil {
			return
		}
		if err = producer.QueueProducer("queue-01", fmt.Sprintf("我是一条消息：%+v", i)); err != nil {
			return
		}
		fmt.Println("我是一条消息：", i)
		time.Sleep(2 * time.Second)
	}
	fmt.Println("消息发送成功")
	return
}

// 发布完成手动释放连接
func Test_QueueProducerReleaseConn(t *testing.T) {
	rabbitmq.InitMQ(config.DefaultConfig)
	defer rabbitmq.CloseMQ()

	option := NewOption().SetPublishOption(NewPublishOption().SetAutoReleaseConn(false))
	producer, err := NewProducer(option)
	if err != nil {
		return
	}
	defer producer.PoolReleaseConnection()

	for i := 0; i < 100; i++ {
		if err = producer.QueueProducer("queue-01", fmt.Sprintf("我是一条消息：%+v", i)); err != nil {
			return
		}
		fmt.Println("我是一条消息：", i)
		time.Sleep(2 * time.Second)
	}

	// 手动释放连接
	fmt.Println("消息发送成功")
	return
}

func Test_PublishSubscribeProducer(t *testing.T) {
	rabbitmq.InitMQ(config.DefaultConfig)
	defer rabbitmq.CloseMQ()

	producer, err := NewProducer(nil)
	if err != nil {
		return
	}

	for i := 0; i < 10000; i++ {
		if err := producer.PublishSubscribeProducer("exchange-publish-subscribe", fmt.Sprintf("我是一条消息：%+v", i)); err != nil {
			return
		}
		fmt.Println("我是一条消息：", i)
		time.Sleep(2 * time.Second)
	}
	fmt.Println("消息发送成功")
	return
}

func Test_RoutingKeyProducer(t *testing.T) {
	rabbitmq.InitMQ(config.DefaultConfig)
	defer rabbitmq.CloseMQ()

	producer, err := NewProducer(nil)
	if err != nil {
		return
	}

	for i := 0; i < 10000; i++ {
		if err := producer.RoutingKeyProducer("exchange-routingKey", "routingKey-yja", fmt.Sprintf("我是一条消息：%+v", i)); err != nil {
			return
		}
		fmt.Println("我是一条消息：", i)
		time.Sleep(2 * time.Second)
	}
	fmt.Println("消息发送成功")
	return
}

func Test_TopicProducer(t *testing.T) {
	rabbitmq.InitMQ(config.DefaultConfig)
	defer rabbitmq.CloseMQ()

	producer, err := NewProducer(nil)
	if err != nil {
		return
	}

	for i := 0; i < 10000; i++ {
		if err := producer.TopicProducer("exchange-topic", "topic-yja", fmt.Sprintf("我是一条消息：%+v", i)); err != nil {
			return
		}
		fmt.Println("我是一条消息：", i)
		time.Sleep(2 * time.Second)
	}
	fmt.Println("消息发送成功")
	return
}

func Test_DelayedTaskPublish(t *testing.T) {
	rabbitmq.InitMQ(config.DefaultConfig)
	defer rabbitmq.CloseMQ()

	for i := 0; i < 10000; i++ {
		producer, err := NewProducer(nil)
		if err != nil {
			return
		}

		if err = producer.DelayedTaskPublish(DelayedTask{
			ExchangeName: "exchange-delayedTask",
			QueueName:    "queue-delayedTask",
			RoutingKey:   "routingKey-yja-delayedTask",
			TaskKind:     "task-kind-01",
		}, 60000, fmt.Sprintf("我是一条消息：%+v", i)); err != nil {
			return
		}
		fmt.Println(fmt.Sprintf("生产者发送消息，时间：%+v,消息内容：我是一条消息%+v", time.Now(), i))
		time.Sleep(2 * time.Second)
	}
	fmt.Println("消息发送成功")
	return
}

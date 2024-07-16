package consumer

import (
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yueja/rabbitmq/exchange"
	rabbit "github.com/yueja/rabbitmq/mq"
	mqQueue "github.com/yueja/rabbitmq/queue"
	"log"
)

type AbstractConsumer interface {
	// PoolReleaseConnection 释放连接
	PoolReleaseConnection() (err error)

	// QueueConsumer 队列模式
	QueueConsumer(queueName string) (msgs <-chan amqp.Delivery, err error)

	// PublishSubscribeConsumer 发布订阅模式
	PublishSubscribeConsumer(exchangeName, queueName string) (msgs <-chan amqp.Delivery, err error)

	// RoutingKeyConsumer routing路由模式
	RoutingKeyConsumer(exchangeName, routingKey string) (msgs <-chan amqp.Delivery, err error)

	// TopicConsumer Topic模式
	TopicConsumer(exchangeName, topic string) (msgs <-chan amqp.Delivery, err error)

	// DelayedTaskConsumer 延时任务
	DelayedTaskConsumer() (msgs <-chan amqp.Delivery, err error)
}

type Consumer struct {
	conn *rabbit.MQRabbitConnection
	Option
}

func NewConsumer(option *Option) (AbstractConsumer, error) {
	if option == nil {
		option = NewOption()
	}
	if option.queueOption == nil {
		option.SetQueueOption(mqQueue.NewQueueOption())
	}
	if option.consumerOption == nil {
		option.SetConsumerOption(NewConsumerOption())
	}
	if option.exchangeOption == nil {
		option.SetExchangeOption(exchange.NewExchangeOption()) // 默认扇形交换机
	}

	conn, err := rabbit.GetRabbitMQConnection()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:   conn,
		Option: *option,
	}, nil
}

// PoolReleaseConnection 释放连接
func (c *Consumer) PoolReleaseConnection() (err error) {
	return rabbit.PoolReleaseConnection(c.conn)
}

func (c *Consumer) QueueConsumer(queueName string) (msgs <-chan amqp.Delivery, err error) {
	ch := c.conn.Ch
	if err = c.checkConsumer(); err != nil {
		return
	}
	// 创建存储队列
	_, err = c.queueOption.QueueDeclare(ch, queueName)
	if err != nil {
		log.Fatalf("无法声明队列:%s", err)
		return
	}

	// 消费信息
	if msgs, err = c.consumer(queueName); err != nil {
		return
	}
	return
}

func (c *Consumer) PublishSubscribeConsumer(exchangeName, queueName string) (msgs <-chan amqp.Delivery, err error) {
	ch := c.conn.Ch
	if err = c.checkConsumer(); err != nil {
		return
	}
	if err = c.exchangeOption.ExchangeDeclare(ch, exchangeName); err != nil {
		return
	}

	if _, err = c.queueOption.QueueDeclare(ch, queueName); err != nil {
		log.Fatalf("无法声明队列:%s", err)
		return
	}

	if err = c.queueOption.QueueBind(ch, exchangeName, queueName, ""); err != nil {
		log.Fatalf("无法绑定队列:%s", err)
		return
	}

	// 消费信息
	if msgs, err = c.consumer(queueName); err != nil {
		return
	}
	return
}

// RoutingKeyConsumer routing路由模式
func (c *Consumer) RoutingKeyConsumer(exchangeName, routingKey string) (msgs <-chan amqp.Delivery, err error) {
	ch := c.conn.Ch
	if err = c.checkConsumer(); err != nil {
		return
	}
	if routingKey == "" {
		err = errors.New("routingKey is empty")
		return
	}
	c.exchangeOption.Kind = amqp.ExchangeDirect
	if err = c.exchangeOption.ExchangeDeclare(ch, exchangeName); err != nil {
		return
	}

	// 队列名称，留空表示由RabbitMQ自动生成,因为定义了key所以队列名可以是随意的，毕竟是依靠key来进行匹配的
	var queue amqp.Queue
	if queue, err = c.queueOption.QueueDeclare(ch, ""); err != nil {
		return
	}

	// 将队列绑定到交换机上，并指定要接收的路由键
	if err = c.queueOption.QueueBind(ch, exchangeName, queue.Name, routingKey); err != nil {
		return
	}

	// 消费信息
	if msgs, err = c.consumer(queue.Name); err != nil {
		return
	}
	return
}

// TopicConsumer routing路由模式
func (c *Consumer) TopicConsumer(exchangeName, topic string) (msgs <-chan amqp.Delivery, err error) {
	ch := c.conn.Ch
	if err = c.checkConsumer(); err != nil {
		return
	}
	if topic == "" {
		err = errors.New("topic is empty")
		return
	}
	c.exchangeOption.Kind = amqp.ExchangeTopic
	if err = c.exchangeOption.ExchangeDeclare(ch, exchangeName); err != nil {
		return
	}

	// 声明一个临时队列，队列名称，留空表示由RabbitMQ自动生成
	var queue amqp.Queue
	if queue, err = c.queueOption.QueueDeclare(ch, ""); err != nil {
		return
	}

	// 将队列绑定到交换机上，并指定要接收的路由键
	if err = c.queueOption.QueueBind(ch, exchangeName, queue.Name, topic); err != nil {
		return
	}

	// 消费信息
	if msgs, err = c.consumer(queue.Name); err != nil {
		return
	}
	return
}

func (c *Consumer) consumer(queueName string) (msgs <-chan amqp.Delivery, err error) {
	return c.conn.Ch.Consume(
		queueName,
		c.consumerOption.consumer,
		c.consumerOption.autoAck,
		c.consumerOption.exclusive,
		c.consumerOption.noLocal,
		c.consumerOption.noWait,
		c.consumerOption.args)
}

func (c *Consumer) checkConsumer() error {
	if c.conn.Conn == nil {
		return errors.New("rabbitmq is nil")
	}
	if c.consumerOption.consumer == "" {
		return errors.New("consumer is empty")
	}
	return nil
}

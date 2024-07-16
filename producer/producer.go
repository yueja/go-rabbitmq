package producer

import (
	"encoding/json"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yueja/rabbitmq/exchange"
	rabbit "github.com/yueja/rabbitmq/mq"
	mqQueue "github.com/yueja/rabbitmq/queue"
	"log"
	"strconv"
)

type AbstractProducer interface {
	// PoolReleaseConnection 释放连接
	PoolReleaseConnection() (err error)

	// QueueProducer 队列模式
	QueueProducer(queueName string, body interface{}) (err error)

	// PublishSubscribeProducer 发布订阅模式
	PublishSubscribeProducer(exchangeName string, body interface{}) (err error)

	// RoutingKeyProducer routing路由模式
	RoutingKeyProducer(exchangeName, routingKey string, body interface{}) (err error)

	// TopicProducer Topic模式
	TopicProducer(exchangeName, topic string, body interface{}) (err error)

	// DelayedTaskPublish 延时任务
	DelayedTaskPublish(task DelayedTask, delayTimeStamp int64, body interface{}) (err error)
}

type Producer struct {
	conn *rabbit.MQRabbitConnection
	Option
}

func NewProducer(option *Option) (AbstractProducer, error) {
	if option == nil {
		option = NewOption()
	}
	if option.queueOption == nil {
		option.SetQueueOption(mqQueue.NewQueueOption())
	}
	if option.publishOption == nil {
		option.SetPublishOption(NewPublishOption())
	}
	if option.exchangeOption == nil {
		option.SetExchangeOption(exchange.NewExchangeOption()) // 默认扇形交换机
	}

	conn, err := rabbit.GetRabbitMQConnection()
	if err != nil {
		return nil, err
	}

	return &Producer{
		conn:   conn,
		Option: *option,
	}, nil
}

// PoolReleaseConnection 释放连接
func (p *Producer) PoolReleaseConnection() (err error) {
	return rabbit.PoolReleaseConnection(p.conn)
}

// QueueProducer 队列模式
func (p *Producer) QueueProducer(queueName string, body interface{}) (err error) {
	if p.conn == nil {
		return errors.New("rabbitmq is nil")
	}

	// 创建存储队列
	_, err = p.queueOption.QueueDeclare(p.conn.Ch, queueName)
	if err != nil {
		log.Fatalf("无法声明队列:%s", err)
		return
	}
	if err = p.publish("", queueName, body); err != nil {
		return err
	}
	return
}

// PublishSubscribeProducer 发布订阅模式
func (p *Producer) PublishSubscribeProducer(exchangeName string, body interface{}) (err error) {
	if p.conn == nil {
		return errors.New("rabbitmq is nil")
	}

	if err = p.exchangeOption.ExchangeDeclare(p.conn.Ch, exchangeName); err != nil {
		return err
	}

	if err = p.publish(exchangeName, "", body); err != nil {
		return err
	}
	return err
}

// RoutingKeyProducer routing路由模式
func (p *Producer) RoutingKeyProducer(exchangeName, routingKey string, body interface{}) (err error) {
	if p.conn == nil {
		return errors.New("rabbitmq is nil")
	}

	if routingKey == "" {
		return errors.New("routingKey is empty")
	}
	p.exchangeOption.Kind = amqp.ExchangeDirect

	if err = p.exchangeOption.ExchangeDeclare(p.conn.Ch, exchangeName); err != nil {
		return err
	}

	if err = p.publish(exchangeName, routingKey, body); err != nil {
		return err
	}
	return err
}

// TopicProducer Topic模式
func (p *Producer) TopicProducer(exchangeName, topic string, body interface{}) (err error) {
	if p.conn == nil {
		return errors.New("rabbitmq is nil")
	}

	if topic == "" {
		return errors.New("topic is empty")
	}
	p.exchangeOption.Kind = amqp.ExchangeTopic

	if err = p.exchangeOption.ExchangeDeclare(p.conn.Ch, exchangeName); err != nil {
		return err
	}

	if err = p.publish(exchangeName, topic, body); err != nil {
		return err
	}
	return err
}

func (p *Producer) publish(exchange, routingKey string, body interface{}) (err error) {
	defer func(conn *rabbit.MQRabbitConnection) {
		// 释放连接
		if p.publishOption.autoReleaseConn {
			if err := rabbit.PoolReleaseConnection(conn); err != nil {
				log.Printf("Producer publish PoolReleaseConnection err:%+v", err)
				return
			}
		}
	}(p.conn)

	var bodyByte []byte
	if bodyByte, err = json.Marshal(body); err != nil {
		return err
	}

	err = p.conn.Ch.Publish(
		exchange,
		routingKey,
		p.publishOption.mandatory, // 必须发送到消息队列
		p.publishOption.immediate, // 不等待服务器确认
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyByte,
			Expiration:  strconv.FormatInt(p.publishOption.expiration, 10), // 消息过期时间（消息级别）,毫秒
		})
	if err != nil {
		log.Fatalf("消息生产失败:%s", err)
		return err
	}
	return nil
}

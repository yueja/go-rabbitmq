package producer

import (
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	rabbit "github.com/yueja/rabbitmq/mq"
)

type DelayedTask struct {
	ExchangeName string // 交换机名称
	QueueName    string // 队列名称
	RoutingKey   string // 路由设置
	TaskKind     string // 任务类型，自定义字符串
}

type DelayedTaskBody struct {
	TaskKind string      `json:"task_kind"` // 任务类型，自定义字符串
	Data     interface{} `json:"data"`
}

// DelayedTaskPublish 延时任务
// delayTimeStamp 延时时间，单位毫秒
func (p *Producer) DelayedTaskPublish(task DelayedTask, delayTimeStamp int64, body interface{}) (err error) {
	if p.conn == nil {
		return errors.New("rabbitmq is nil")
	}
	if delayTimeStamp == 0 {
		return errors.New("expiration must greater than zero")
	}

	// # ========== 2.设置队列（队列、交换机、绑定） ==========
	// 声明队列
	p.queueOption.SetArgs(amqp.Table{
		//"x-message-ttl":             expiration,                 // 消息TTL，毫秒(队列级)
		"x-dead-letter-exchange":    rabbit.GetRabbitConfig().DeadExchangeName, // 设置死信交换机
		"x-dead-letter-routing-key": rabbit.GetRabbitConfig().DeadRoutingKey,   // 指定死信routing-key
	})
	if _, err = p.queueOption.QueueDeclare(p.conn.Ch, task.QueueName); err != nil {
		return err
	}
	// 声明交换机
	p.exchangeOption.SetKind(amqp.ExchangeDirect)
	if err = p.exchangeOption.ExchangeDeclare(p.conn.Ch, task.ExchangeName); err != nil {
		return err
	}
	// 队列绑定（将队列、routing-key、交换机三者绑定到一起）
	if err = p.queueOption.QueueBind(p.conn.Ch, task.ExchangeName, task.QueueName, task.RoutingKey); err != nil {
		return err
	}

	// # ========== 3.设置死信队列（队列、交换机、绑定） ==========
	// 声明死信队列
	// args 为 nil。切记不要给死信队列设置消息过期时间,否则失效的消息进入死信队列后会再次过期。
	p.queueOption.SetArgs(nil)
	if _, err = p.queueOption.QueueDeclare(p.conn.Ch, rabbit.GetRabbitConfig().DeadQueueName); err != nil {
		return err
	}
	// 声明交换机
	if err = p.exchangeOption.ExchangeDeclare(p.conn.Ch, rabbit.GetRabbitConfig().DeadExchangeName); err != nil {
		return err
	}
	// 队列绑定（将队列、routing-key、交换机三者绑定到一起）
	if err = p.queueOption.QueueBind(p.conn.Ch, rabbit.GetRabbitConfig().DeadExchangeName, rabbit.GetRabbitConfig().DeadQueueName, rabbit.GetRabbitConfig().DeadRoutingKey); err != nil {
		return err
	}

	// # ========== 发布消息 ==========
	var data = DelayedTaskBody{
		TaskKind: task.TaskKind,
		Data:     body,
	}
	p.publishOption.SetExpiration(delayTimeStamp)
	if err = p.publish(task.ExchangeName, task.RoutingKey, data); err != nil {
		return err
	}
	return
}

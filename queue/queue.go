package queue

import amqp "github.com/rabbitmq/amqp091-go"

type QueueOption struct {
	durable    bool                   // 持久化设置，可以为true根据需求选择
	autoDelete bool                   // 自动删除，没有用户连接删除queue一般不选用
	exclusive  bool                   // 独占
	noWait     bool                   // 等待服务器确认
	args       map[string]interface{} // 参数
}

func NewQueueOption() *QueueOption {
	return &QueueOption{durable: false}
}

func (q *QueueOption) SetDurable(durable bool) *QueueOption {
	q.durable = durable
	return q
}

func (q *QueueOption) SetAutoDelete(autoDelete bool) *QueueOption {
	q.autoDelete = autoDelete
	return q
}

func (q *QueueOption) SetExclusive(exclusive bool) *QueueOption {
	q.exclusive = exclusive
	return q
}

func (q *QueueOption) SetNoWait(noWait bool) *QueueOption {
	q.noWait = noWait
	return q
}

func (q *QueueOption) SetArgs(args map[string]interface{}) *QueueOption {
	q.args = args
	return q
}

func (q *QueueOption) QueueDeclare(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queueName,
		q.durable,
		q.autoDelete,
		q.exclusive,
		q.noWait,
		q.args)
}

func (q *QueueOption) QueueBind(ch *amqp.Channel, exchangeName, queueName string, routingKey string) error {
	return ch.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		q.noWait,
		q.args)
}

package consumer

import (
	"yueja/go-rabbitmq/exchange"
	mqQueue "yueja/go-rabbitmq/queue"
)

type Option struct {
	exchangeOption *exchange.ExchangeOption
	queueOption    *mqQueue.QueueOption
	consumerOption *ConsumerOption
}

func NewOption() *Option {
	return &Option{}
}

func (o *Option) SetExchangeOption(exchangeOption *exchange.ExchangeOption) *Option {
	o.exchangeOption = exchangeOption
	return o
}

func (o *Option) SetQueueOption(queueOption *mqQueue.QueueOption) *Option {
	o.queueOption = queueOption
	return o
}

func (o *Option) SetConsumerOption(consumerOption *ConsumerOption) *Option {
	o.consumerOption = consumerOption
	return o
}

type ConsumerOption struct {
	consumer  string                 // 消费者名称
	autoAck   bool                   // 是否自动提交ack，默认：false
	exclusive bool                   // 独占
	noLocal   bool                   //
	noWait    bool                   // 等待服务器确认
	args      map[string]interface{} // 参数
}

func NewConsumerOption() *ConsumerOption {
	return &ConsumerOption{autoAck: false}
}

func (c *ConsumerOption) SetConsumer(consumer string) *ConsumerOption {
	c.consumer = consumer
	return c
}

func (c *ConsumerOption) SetAutoAck(autoAck bool) *ConsumerOption {
	c.autoAck = autoAck
	return c
}

func (c *ConsumerOption) SetExclusive(exclusive bool) *ConsumerOption {
	c.exclusive = exclusive
	return c
}
func (c *ConsumerOption) SetNoLocal(noLocal bool) *ConsumerOption {
	c.noLocal = noLocal
	return c
}

func (c *ConsumerOption) SetNoWait(noWait bool) *ConsumerOption {
	c.noWait = noWait
	return c
}

func (c *ConsumerOption) SetArgs(args map[string]interface{}) *ConsumerOption {
	c.args = args
	return c
}

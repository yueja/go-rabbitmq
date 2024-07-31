package producer

import (
	"yueja/go-rabbitmq/exchange"
	mqQueue "yueja/go-rabbitmq/queue"
)

type Option struct {
	exchangeOption *exchange.ExchangeOption
	queueOption    *mqQueue.QueueOption
	publishOption  *PublishOption
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

func (o *Option) SetPublishOption(publishOption *PublishOption) *Option {
	o.publishOption = publishOption
	return o
}

type PublishOption struct {
	mandatory       bool  // 必须发送到消息队列
	immediate       bool  // 不等待服务器确认
	expiration      int64 // 消息过期时间（消息级别）,毫秒
	autoReleaseConn bool  // 发布完成后是否自动释放连接(默认：true)
}

func NewPublishOption() *PublishOption {
	return &PublishOption{autoReleaseConn: true}
}

func (p *PublishOption) SetMandatory(mandatory bool) *PublishOption {
	p.mandatory = mandatory
	return p
}

func (p *PublishOption) SetImmediate(immediate bool) *PublishOption {
	p.immediate = immediate
	return p
}

func (p *PublishOption) SetExpiration(expiration int64) *PublishOption {
	p.expiration = expiration
	return p
}

func (p *PublishOption) SetAutoReleaseConn(autoReleaseConn bool) *PublishOption {
	p.autoReleaseConn = autoReleaseConn
	return p
}

package exchange

import amqp "github.com/rabbitmq/amqp091-go"

type ExchangeOption struct {
	Kind       string                 // 交换机类型
	Durable    bool                   // 持久化设置，可以为true根据需求选择
	AutoDelete bool                   // 自动删除，没有用户连接删除queue一般不选用
	Internal   bool                   // 是否内部使用，true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
	NoWait     bool                   // 等待服务器确认.
	Args       map[string]interface{} // 参数
}

func NewExchangeOption() *ExchangeOption {
	return &ExchangeOption{Kind: amqp.ExchangeFanout, Durable: false}
}

func (e *ExchangeOption) SetKind(kind string) *ExchangeOption {
	e.Kind = kind
	return e
}

func (e *ExchangeOption) SetDurable(durable bool) *ExchangeOption {
	e.Durable = durable
	return e
}

func (e *ExchangeOption) SetAutoDelete(autoDelete bool) *ExchangeOption {
	e.AutoDelete = autoDelete
	return e
}

func (e *ExchangeOption) SetInternal(internal bool) *ExchangeOption {
	e.Internal = internal
	return e
}

func (e *ExchangeOption) SetNoWait(noWait bool) *ExchangeOption {
	e.NoWait = noWait
	return e
}

func (e *ExchangeOption) SetArgs(args map[string]interface{}) *ExchangeOption {
	e.Args = args
	return e
}

func (e *ExchangeOption) ExchangeDeclare(ch *amqp.Channel, exchangeName string) error {
	return ch.ExchangeDeclare(
		exchangeName,
		e.Kind,
		e.Durable,
		e.AutoDelete,
		e.Internal,
		e.NoWait,
		e.Args,
	)
}

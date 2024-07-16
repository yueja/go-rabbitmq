package config

type Config struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	VHost    string `mapstructure:"vhost"`

	DeadExchangeName string `mapstructure:"dead_exchange_name"` // 死信交换机
	DeadQueueName    string `mapstructure:"ead_queue_name"`     // 死信队列
	DeadRoutingKey   string `mapstructure:"ead_routing_key"`    // 死信RoutingKey

	MaxConn int `mapstructure:"max_conn"` // 最大连接数
}

var DefaultConfig = Config{
	Host:             "127.0.0.1",
	Port:             "5674",
	UserName:         "admin",
	Password:         "admin123456",
	VHost:            "my_vhost",
	DeadExchangeName: "dead-exchange-name-01",
	DeadQueueName:    "dead-queue-name-01",
	DeadRoutingKey:   "dead-routingKey-01",
}

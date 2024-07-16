## 基于Rabbit MQ实现消息队列工具包

### 实现内容
1、连接池的实现
2、MQ的各个模式实现
3、基于死信队列的延时任务实现

### 基于本地docker安装调试

#### 安装Rabbit MQ
```shell
docker run -d -p 15673:15672 -p 5674:5672 \
        --restart=always \
        -e RABBITMQ_DEFAULT_VHOST=my_vhost  \
        -e RABBITMQ_DEFAULT_USER=admin \
        -e RABBITMQ_DEFAULT_PASS=admin123456 \
        --hostname myRabbit \
        --name rabbitmq\
        rabbitmq:latest

注意：在映射的端口号的时候不要映射 5671端口，端口5671是 RabbitMQ 的默认AMQP over TLS/SSL端口。AMQP（Advanced Message Queuing Protocol）是一种消息传递协议，用于在应用程序之间进行可靠的消息传递。

参数说明：

-d：表示在后台运行容器；
-p：将主机的端口 15673（Web访问端口号）对应当前rabbitmq容器中的 15672 端口，将主机的5674（应用访问端口）端口映射到rabbitmq中的5672端口；
--restart=alawys：设置开机自启动
-e：指定环境变量：
    RABBITMQ_DEFAULT_VHOST：默认虚拟机名；
    RABBITMQ_DEFAULT_USER：默认的用户名；
    RABBITMQ_DEFAULT_PASS：默认的用户密码；
--hostname：指定主机名（RabbitMQ 的一个重要注意事项是它根据所谓的 节点名称 存储数据，默认为主机名）；
--name rabbitmq-new：设置容器名称；
```

#### 启动web客户端
```shell
docker exec -it rabbitmq rabbitmq-plugins enable rabbitmq_management
```

#### 访问rabbitmq的web客户端
在浏览器上输入 ip+端口(15673) 访问rabbitmq的web客户端
输入上面在初始化Rabbitmq容器时我们自己指定了默认账号和密码：admin/admin123456，如果没有指定的话那么rabbitmq的默认账号密码是:guest/guest

### 消息队列使用示例
#### 生产者
```golang
func main() {
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
```
#### 消费者
```golang
func main() {
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
```

备注：更多功能的使用方法见单元测试示例，此处不再罗列。


### 延时任务示例
#### 生产方
```golang
func main() {
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
		}, 10000, fmt.Sprintf("我是一条消息：%+v", i)); err != nil {
			return
		}
		fmt.Println(fmt.Sprintf("生产者发送消息，时间：%+v,消息内容：我是一条消息%+v", time.Now(), i))
		time.Sleep(2 * time.Second)
	}
	fmt.Println("消息发送成功")
	return
}
```
#### 任务执行方
```golang
func main() {
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
```
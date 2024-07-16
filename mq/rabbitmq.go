package rabbitmq

import (
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yueja/rabbitmq/config"
	"log"
	"sync"
)

var mq *Pool
var mqConfig *config.Config

var once sync.Once

// InitMQ 初始化mq连接池
func InitMQ(conf config.Config) {
	once.Do(func() {
		mqConfig = &conf
		mq = NewPool(conf, func() (conn *amqp.Connection, err error) {
			url := fmt.Sprintf(fmt.Sprintf("amqp://%s:%s@%s:%s/%s", conf.UserName, conf.Password, conf.Host, conf.Port, conf.VHost))
			return amqp.DialConfig(url, amqp.Config{})
		})
	})
	fmt.Println("Connected to RabbitMQ")
}

// CloseMQ 关闭连接池，主要用于优雅退出
func CloseMQ() (err error) {
	var mqPool pool
	if mqPool, err = getRabbitMQPool(); err != nil {
		return
	}
	if err = mqPool.close(); err != nil {
		return err
	}
	log.Printf("Closed RabbitMQ pool")
	return nil
}

type MQRabbitConnection struct {
	Conn *amqp.Connection

	Ch *amqp.Channel
}

// GetRabbitMQConnection 获取连接
func GetRabbitMQConnection() (conn *MQRabbitConnection, err error) {
	var mqPool pool
	if mqPool, err = getRabbitMQPool(); err != nil {
		return
	}
	conn = new(MQRabbitConnection)
	if conn.Conn, err = mqPool.get(); err != nil {
		return nil, err
	}
	// 开启通道
	if conn.Ch, err = conn.Conn.Channel(); err != nil {
		return nil, err
	}
	return
}

// PoolReleaseConnection 连接释放
func PoolReleaseConnection(conn *MQRabbitConnection) (err error) {
	var mqPool pool
	if mqPool, err = getRabbitMQPool(); err != nil {
		return
	}
	return mqPool.release(conn.Conn)
}

func getRabbitMQPool() (pool pool, err error) {
	if mq == nil {
		return nil, errors.New("mq is nil")
	}
	return mq, nil
}

func GetRabbitConfig() *config.Config {
	return mqConfig
}

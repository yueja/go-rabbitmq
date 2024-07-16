package rabbitmq

import (
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yueja/rabbitmq/config"
	"github.com/yueja/rabbitmq/constant"
	"log"
	"sync"
	"time"
)

type pool interface {
	// 获取一个连接
	get() (*amqp.Connection, error)
	// 释放一个连接
	release(conn *amqp.Connection) (err error)
	// 关闭连接池
	close() (err error)
	// 关闭某个连接
	closeConn(conn *MQRabbitConnection) (err error)
}

// Pool RabbitMQ连接池结构
// 暂未实现连接池的动态缩容
type Pool struct {
	factory     func() (*amqp.Connection, error) // 创建新连接的函数
	maxConn     int                              // 池中的最大连接数
	connMap     map[interface{}]struct{}         // 所有有效连接map
	connChanMap map[interface{}]struct{}         // 可用空闲连接chan map
	conns       chan *amqp.Connection            // 可用空闲连接chan
	mu          sync.Mutex                       // 用于保护并发访问
	closed      bool                             // 标记池是否已关闭
}

// NewPool 创建一个新的RabbitMQ连接池
func NewPool(config config.Config, createConnection func() (*amqp.Connection, error)) *Pool {
	if config.MaxConn <= 0 {
		config.MaxConn = constant.MQConnectionPoolMaxCountDefault
	}
	return &Pool{
		factory:     createConnection,
		maxConn:     config.MaxConn,
		connMap:     make(map[interface{}]struct{}, config.MaxConn),
		connChanMap: make(map[interface{}]struct{}, config.MaxConn),
		conns:       make(chan *amqp.Connection, config.MaxConn),
		mu:          sync.Mutex{},
		closed:      false,
	}
}

// Get 从池中获取一个连接
func (p *Pool) get() (*amqp.Connection, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil, errors.New("pool is closed")
	}

	if len(p.connMap) >= p.maxConn {
		// 已经达到了最大连接数
		return p.getConn()
	}

	// 尚未达到最大连接数，优先判断是否有空闲连接
	for {
		select {
		case conn := <-p.conns:
			// 有空闲连接
			delete(p.connChanMap, conn)
			if _, ok := p.connMap[conn]; !ok {
				// 无效连接，该连接可能已经被关闭
				continue
			}
			return conn, nil
		default:
			// 无空闲连接，创建
			return p.create()
		}
	}
}

func (p *Pool) create() (*amqp.Connection, error) {
	conn, err := p.factory()
	if err != nil {
		return conn, err
	}
	p.connMap[conn] = struct{}{}
	return conn, nil
}

func (p *Pool) getConn() (*amqp.Connection, error) {
	// 已经达到了最大连接数
	timer := time.NewTimer(constant.MQConnectionPoolGetMaxWaitTimeDefault * time.Second)
	for {
		select {
		case conn := <-p.conns:
			delete(p.connChanMap, conn)
			if _, ok := p.connMap[conn]; !ok {
				// 无效连接，该连接可能已经被关闭
				continue
			}
			// 阻塞等待
			return conn, nil
		case <-timer.C:
			return nil, errors.New("get mq conn timeout")
		}
	}
}

// Release 尝试将连接放回池中
func (p *Pool) release(conn *amqp.Connection) (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		// 如果池已关闭，则关闭连接并返回错误
		if !conn.IsClosed() {
			if err = conn.Close(); err != nil {
				return err
			}
		}
		return errors.New("pool is closed")
	}
	if conn.IsClosed() {
		// 连接已经被关闭，从连接池中剔除，不再释放到连接池
		delete(p.connMap, conn)
		return nil
	}

	if _, ok := p.connChanMap[conn]; ok {
		// 该链接已经在chan中，无需再次释放
		return
	}
	p.conns <- conn
	p.connChanMap[conn] = struct{}{}
	return nil
}

// closeConn 关闭某个连接
func (p *Pool) closeConn(conn *MQRabbitConnection) (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.connMap, conn)
	conn.Ch.Close()
	conn.Conn.Close()
	return nil
}

// Close 关闭连接池并释放所有资源。
func (p *Pool) close() (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return errors.New("pool is already closed")
	}

	p.closed = true
	close(p.conns)

	// 关闭池中的所有连接
	for {
		conn, ok := <-p.conns
		if !ok {
			break
		}
		if conn.IsClosed() {
			// 连接已经被关闭
			continue
		}
		if err = conn.Close(); err != nil {
			log.Printf("close conn conn err:%v", err)
		}
	}
	return nil
}

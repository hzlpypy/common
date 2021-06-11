package topic

import (
	"context"
	"github.com/streadway/amqp"
)

type Interface interface {
	CreateExchange() error
	QueueDeclareAndBindRoutingKey() error
	Send() error // 废弃
}

func NewTopicReq(ctx context.Context, topicReq *TopicReq) (Interface, error) {
	conn := topicReq.Conn
	// 公平分发
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, err
	}
	topicReq.Ch = ch
	go func() {
		select {
		case <-ctx.Done():
			ch.Close()
			return
		}
	}()
	return topicReq, nil
}

// QueueDeclareAndBindRoutingKey: 创建队列并绑定交换机和路由键
func (t *TopicReq) QueueDeclareAndBindRoutingKey() error {
	queue := t.Queue
	ch := t.Ch
	q, err := ch.QueueDeclare(
		queue.QueueName,
		t.Durable,
		false,
		false,
		false,
		queue.QueueDeclareMap,
	)
	if err != nil {
		return err
	}
	return ch.QueueBind(
		q.Name,         // queue name
		t.RoutingKey,   // routing key
		t.ExchangeName, // exchange name
		false,
		queue.QueueBindMap)
}

// CreateExchange: 创建交换机
func (t *TopicReq) CreateExchange() error {

	return t.Ch.ExchangeDeclare(
		t.ExchangeName, // name
		t.ExchangeType, // type
		t.Durable,      // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
}

// Send：消息发送
func (t *TopicReq) Send() error {
	ch := t.Ch
	// 将消息发送到交换机
	err := ch.Publish(
		t.ExchangeName,
		t.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: t.ContentType,
			Body:        []byte(t.Msg),
		},
	)
	return err
}

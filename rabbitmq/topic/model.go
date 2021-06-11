package topic

import "github.com/streadway/amqp"

type TopicReq struct {
	Conn         *amqp.Connection
	Ch           *amqp.Channel
	ExchangeName string
	ExchangeType string
	Durable      bool
	Msg          string
	ContentType  string
	RoutingKey   string
	Queue        *Queue
}

type Queue struct {
	// 队列名称
	QueueName string
	// 创建队列额外参数，(例如："x-max-length": queue.QueueMax"参数限制了一个队列的消息总数，当消息总数达到限定值时，队列头的消息会被抛弃)
	QueueDeclareMap map[string]interface{}
	// 绑定队列额外参数
	QueueBindMap map[string]interface{}
}

type TopicRes struct {
	Msg byte
}

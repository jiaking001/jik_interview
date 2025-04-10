package rabbmit

import (
	"app/internal/model"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

var Ch *amqp.Channel

// InitRabbitMQ 初始化 RabbitMQ
func InitRabbitMQ() {
	conn, err := amqp.Dial("amqp://jiaking:123456@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	Ch, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
}

// SendMessage 生产者
func SendMessage(msg []model.Question) error {

	// 将结构体序列化为 JSON
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = Ch.Publish(
		"",         // 交换机名称（使用默认交换机）
		"my_queue", // 队列名称
		false,      // 是否强制发送
		false,      // 是否立即发送
		amqp.Publishing{
			ContentType: "application/json", // json 类型的消息格式
			Body:        jsonMsg,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

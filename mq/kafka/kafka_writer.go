package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

type KafkaClient struct {
	writer *kafka.Writer
}

func NewKafkaClient(addr string, topic string, partition int, timeOut int) *KafkaClient {
	if timeOut == 0 {
		timeOut = 5
	}
	writer := kafka.Writer{
		Addr:                   kafka.TCP(addr),
		Topic:                  topic,
		Balancer:               &kafka.Hash{}, // 用于对key进行hash，决定消息发送到哪个分区
		MaxAttempts:            0,
		WriteBackoffMin:        0,
		WriteBackoffMax:        0,
		BatchSize:              0,
		BatchBytes:             0,
		BatchTimeout:           0,
		ReadTimeout:            0,
		WriteTimeout:           time.Duration(timeOut) * time.Millisecond, // kafka有时候可能负载很高，写不进去，那么超时后可以放弃写入，用于可以丢消息的场景
		RequiredAcks:           kafka.RequireNone,                         // 不需要任何节点确认就返回
		Async:                  false,
		Completion:             nil,
		Compression:            0,
		Logger:                 nil,
		ErrorLogger:            nil,
		Transport:              nil,
		AllowAutoTopicCreation: false, // 第一次发消息的时候，如果topic不存在，就自动创建topic，工作中禁止使用
	}

	return &KafkaClient{
		writer: &writer,
	}
}

func (k *KafkaClient) Write(key, value []byte) error {
	return k.writer.WriteMessages(context.Background(), kafka.Message{Key: key, Value: value})
}

func (k *KafkaClient) Close() {
	for {
		if err := k.writer.Close(); err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
}

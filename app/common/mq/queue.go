package mq

import (
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/growerlab/backend/app/model/db"
)

const (
	DefaultKey   = "default"
	DefaultValue = "default"
	DefaultGroup = "defaultGroup"
)

type Payload struct {
	ID    string
	Field string
	Value string
}

type Consumer interface {
	Name() string           // consumer name
	Consume(*Payload) error // 进行消费
	// 下面的功能先注释，未来再加
	// Number() int            // 消费者人数（决定提供多少worker调用Consume()
	// RetryCount() int        // 重试次数
	// RetryInterval() int     // 重试间隔，单位s
}

type MessageQueue struct {
	memDB  *redis.Client
	stream *Stream

	Consumers map[string]Consumer // 消费者
}

func (m *MessageQueue) Register(consumers ...Consumer) error {
	if len(consumers) == 0 {
		return nil
	}
	if m.Consumers == nil {
		m.Consumers = map[string]Consumer{}
	}

	for _, c := range consumers {
		if _, exists := m.Consumers[c.Name()]; exists {
			return fmt.Errorf("consumer exists: %s", c.Name())
		}
		m.Consumers[c.Name()] = c
		err := m.createStream(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MessageQueue) createStream(c Consumer) error {
	streamKey := m.streamKey(c.Name())
	_, err := m.stream.AddMessage(streamKey, DefaultKey, DefaultValue)
	if err != nil {
		return err
	}

	err = m.stream.CreateGroup(DefaultGroup, streamKey)
	if err != nil {
		return err
	}
	return nil
}

func (m *MessageQueue) Run() error {

	return nil
}

func (m *MessageQueue) streamKey(name string) string {
	return db.BaseKeyBuilder(name).String()
}

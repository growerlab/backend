package queue

import (
	"encoding/json"
	"sync/atomic"
	"testing"
	"time"
)

var count int32 = 0

func TestNewQueue(t *testing.T) {
	q := newQueue()
	q.AddJob(&SimulateJob{})

	go q.Start()

	var want int32 = 3
	for i := int32(0); i < want; i++ {
		pushPayload(q)
	}

	time.Sleep(1 * time.Second)

	t.Run("test queue", func(t *testing.T) {
		got := count
		if got != want {
			t.Errorf("got = %v, want = %v", got, want)
		}
	})
}

func pushPayload(q *Queue) {
	load := &payload{
		Hello: "moli",
		World: "xxxx",
	}
	body, _ := json.Marshal(load)
	q.PushPayload("test", body)
}

type SimulateJob struct {
}

func (s *SimulateJob) Name() string {
	return "test"
}

func (s *SimulateJob) Eval(payload []byte) (requeue bool, err error) {
	atomic.AddInt32(&count, 1)
	return false, nil
}

func newQueue() *Queue {
	l := &SimulateList{
		dataCh: make(chan []byte, 10),
	}
	return New(l, 10, 2)
}

type payload struct {
	Hello string `json:"hello,omitempty"`
	World string `json:"world,omitempty"`
}

type SimulateList struct {
	dataCh chan []byte
}

func (s *SimulateList) Pop(key string) ([]byte, error) {
	result := <-s.dataCh
	return result, nil
}

func (s *SimulateList) Push(key string, payload []byte) error {
	s.dataCh <- payload
	return nil
}

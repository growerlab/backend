package events

import (
	"github.com/growerlab/backend/app/common/event/queue"
	"github.com/growerlab/backend/app/common/notify"
	"github.com/growerlab/backend/app/model/db"
)

var queueInstance *queue.Queue

func InitEvents() error {
	if err := initQueue(); err != nil {
		return err
	}
	queueInstance.AddJob(NewEmail())
	return nil
}

func initQueue() error {
	workerCount := 5 // worker count，暂时写死
	jobCount := 1    // 每个worker的待处理容器，多出来的会被阻塞

	queueInstance = queue.New(queue.NewList(db.MemDB), workerCount, jobCount)

	go queueInstance.Start()

	notify.Subscribe(func() {
		queueInstance.Release()
	})
	return nil
}

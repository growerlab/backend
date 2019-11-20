package queue

import (
	"github.com/growerlab/backend/app/common/notify"
	"github.com/growerlab/backend/app/common/queue/queue"
	"github.com/growerlab/backend/app/model/db"
)

var queueInstance *queue.Queue

func InitQueue() error {
	workerCount := 5 // worker count，暂时写死
	jobCount := 2    // 每个worker的待处理容器，多出来的会被阻塞

	queueInstance = queue.New(queue.NewList(db.QueueDB), workerCount, jobCount)
	queueInstance.AddJob(NewEmail())

	go queueInstance.Start()

	notify.Subscribe(func() {
		queueInstance.Release()
	})
	return nil
}

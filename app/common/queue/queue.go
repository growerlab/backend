package queue

import (
	"errors"
	"fmt"
	"time"

	"github.com/growerlab/backend/app/common/notify"
	"github.com/growerlab/backend/app/common/queue/job"
	"github.com/growerlab/backend/app/model/db"
	"github.com/growerlab/backend/app/utils/logger"
	"github.com/ivpusic/grpool"
)

// 通用型简易队列

var (
	ErrExists = errors.New("exists job")
)

type Job interface {
	// 唯一性名称
	Name() string
	// 获取payload并执行
	// 当requeue返回true，则payload将重新入队，下次将继续执行
	Eval(payload []byte) (requeue bool, err error)
	// TODO 超时时间、重试次数 等等
}

type Queuer interface {
	// 排队（入队）
	Push(key string, payload []byte) (err error)
	// 出队
	Pop(key string) (payload []byte, err error)
}

func NewQueue() *Queue {
	q := &Queue{
		jobsSet:   make(map[string]Job),
		done:      make(chan struct{}),
		srcQueuer: NewList(db.QueueDB),
	}

	q.AddJob(job.NewEmail(q.onEnter))

	workerCount := len(q.jobsSet) * 2 // 每个jobs至少分配2个worker
	jobCount := 2                     // 每个worker的待处理容器，多出来的会被阻塞
	q.workerPool = grpool.NewPool(workerCount, jobCount)

	notify.Subscribe(func() {
		q.Release()
	})

	return q
}

type Queue struct {
	// 注册job
	jobsSet map[string]Job

	// 是否结束队列
	done chan struct{}

	// 元队列
	srcQueuer Queuer

	// workers
	workerPool *grpool.Pool
}

func (q *Queue) onEnter(jobName string, payload []byte) (err error) {
	key := q.jobKey(jobName)
	err = q.srcQueuer.Push(key, payload)
	if err != nil {
		logger.Error("has err on push: %v", err)
		return err
	}
	return nil
}

func (q *Queue) AddJob(w Job) error {
	if _, ok := q.jobsSet[w.Name()]; ok {
		return ErrExists
	}
	q.jobsSet[w.Name()] = w
	return nil
}

func (q *Queue) Start() {
	run := func() (idle bool, err error) {
		idle = true
		for jobName, _ := range q.jobsSet {
			var payload []byte
			key := q.jobKey(jobName)
			payload, err = q.srcQueuer.Pop(key)
			if err != nil {
				logger.Error("pop job has err: %v", err)
				continue
			}
			if payload == nil {
				continue
			}
			idle = false // 没有闲着
			q.callEval(jobName, payload)
		}
		return
	}

	for {
		select {
		case <-q.done:
			logger.Info("queue was done.")
			return
		default:
			idle, err := run()
			// 对于长时间处于闲置状态时，避免过度/过快从redis数据库中取数据，降低redis压力
			if idle {
				time.Sleep(500 * time.Millisecond)
			}
			if err != nil {
				logger.Error("has err on running: %v", err)
				continue
			}
		}
	}
}

func (q *Queue) callEval(jobName string, payload []byte) {
	q.workerPool.WaitCount(1)

	q.workerPool.JobQueue <- func() {
		defer q.workerPool.JobDone()

		requeue, err := q.jobsSet[jobName].Eval(payload)
		if err != nil {
			logger.Error("queue has err on running: %v", err)
		}
		if requeue {
			err = q.onEnter(jobName, payload)
			if err != nil {
				logger.Error("has err on requeue: %v payload: %v", err, payload)
				return
			}
		}
	}
}

func (q *Queue) Release() {
	close(q.done)
	q.workerPool.WaitAll()
	q.workerPool.Release()
}

func (q *Queue) jobKey(jobName string) string {
	return fmt.Sprintf("queue:%s", jobName)
}

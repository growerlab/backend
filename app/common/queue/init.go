package queue

var queueInstance *Queue

func InitQueue() error {
	queueInstance = NewQueue()
	go queueInstance.Start()
	return nil
}

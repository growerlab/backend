package common

type PushPayloadFunc func(jobName string, payload []byte) error

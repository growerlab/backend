package mq

import (
	"github.com/go-redis/redis/v7"
	"github.com/growerlab/backend/app/common/errors"
)

var (
	ErrNoSuchKey = errors.New("ERR no such key")
)

type Stream struct {
	memDB *redis.Client
}

func NewStream(c *redis.Client) *Stream {
	return &Stream{memDB: c}
}

func (s *Stream) GroupExists(groupName string) bool {
	_, err := s.GroupInfo(groupName)
	if err != nil && err.Error() == ErrNoSuchKey.Error() {
		return false
	}
	return true
}

func (s *Stream) CreateGroup(groupName, streamKey string) error {
	err := s.memDB.XGroupCreate(streamKey, groupName, "0-0").Err()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (s *Stream) GroupInfo(groupName string) ([]redis.XInfoGroups, error) {
	info, err := s.memDB.XInfoGroups(groupName).Result()
	if err != nil {
		return nil, errors.Trace(err)
	}
	return info, nil
}

func (s *Stream) ReadGroupNew(groupName, consumer, streamKey string, count int64) ([]redis.XMessage, error) {
	msgs, err := s.ReadGroupMessages(groupName, consumer, []string{streamKey, ">"}, count)
	if err != nil {
		return nil, err
	}
	if len(msgs) > 0 {
		return msgs, nil
	}

	// 读取历史数据
	msgs, err = s.ReadGroupMessages(groupName, consumer, []string{streamKey, "0-0"}, count)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (s *Stream) ReadGroupMessages(groupName, consumer string, streams []string, count int64) ([]redis.XMessage, error) {
	xstreams, err := s.memDB.XReadGroup(&redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumer,
		Streams:  streams,
		Count:    count,
		Block:    0,
		NoAck:    true,
	}).Result()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return xstreams[0].Messages, nil
}

func (s *Stream) AddMessage(streamKey, field, value string) (id string, err error) {
	id, err = s.memDB.XAdd(&redis.XAddArgs{
		Stream: streamKey,
		ID:     "*",
		Values: map[string]interface{}{field: value},
	}).Result()

	return id, nil
}

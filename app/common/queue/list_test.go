package queue

import (
	"bytes"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/growerlab/backend/app/utils/conf"
)

var redisPool = func() *redis.Pool {
	if _, err := os.Stat(conf.DefaultConfigPath); os.IsNotExist(err) {
		err = os.Chdir(filepath.Join(os.Getenv("GOPATH"), "src", "github.com/growerlab/backend"))
		if err != nil {
			panic(err)
		}
	}

	err := conf.LoadConfig()
	if err != nil {
		panic(err)
	}

	cfg := conf.GetConf().Redis

	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			dbOpt := redis.DialDatabase(11)
			conn, err := redis.Dial("tcp", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)), dbOpt)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}()

func TestNewList(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		key := "test:list"
		want := []byte(`{"hello":"world"}`)

		list := NewList(redisPool)
		err := list.Push(key, want)
		if err != nil {
			t.Errorf("push() = err: %v", err)
			return
		}
		defer list.Release(key)

		got, err := list.Pop(key)
		if err != nil {
			t.Errorf("pop() = err: %v", err)
			return
		}

		if bytes.Compare(got, want) != 0 {
			t.Errorf("got = %v, want: %v", got, want)
			return
		}
	})

}

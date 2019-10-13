package think

import (
	"errors"

	"github.com/growerlab/backend/app/service/graphql/think/types"
)

var ErrAlreadyExistsInHub = errors.New("already exists in hub")

var (
	hub *ObjectHub
)

func init() {
	hub = new(ObjectHub)
	hub.objDict = make(map[string]types.Object)
}

// 管理graphql的object 容器
//
type ObjectHub struct {
	objDict map[string]types.Object
}

func (o *ObjectHub) Each(fn func(types.Object)) {
	for _, obj := range o.objDict {
		fn(obj)
	}
}

func (o *ObjectHub) Register(obj types.Object) error {
	if _, found := o.objDict[obj.Name()]; found {
		return ErrAlreadyExistsInHub
	}
	o.objDict[obj.Name()] = obj
	return nil
}

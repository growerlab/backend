package errors

import (
	"strings"

	pkgerr "github.com/pkg/errors"
)

// 定义Reason
const (
	// 非法参数
	InvalidParameter = "InvalidParameter"
	// 无法找到
	NotFound = "NotFound"
	// 无法找到属性（字段）
	NotFoundField = "NotFoundField"
	// sql错误
	SqlError = "SQLError"
	// 非法长度
	InvalidLength = "InvalidLength"
	// 不等于
	NotEqual = "NotEqual"
	// 失效，过期
	Expired = "Expired"
	// 已被使用过
	Used = "Used"
	// 已存在
	AlreadyExists = "AlreadyExists"
)

var P = InvalidParameterError

func InvalidParameterError(model, field, reason string) string {
	return mustCode(InvalidParameter, model, field, reason)
}

func NotFoundError(model, field string) string {
	return mustCode(NotFound, model, field)
}

func AlreadyExistsError(model, reason string) string {
	return mustCode(AlreadyExists, model, reason)
}

func mustCode(parts ...string) string {
	if len(parts) == 0 {
		panic("parts is required")
	}
	return strings.Join(parts, ".")
}

// 封装（避免在项目中使用时，引用多个包）
var (
	Wrap     = pkgerr.Wrap
	Wrapf    = pkgerr.Wrapf
	Message  = pkgerr.WithMessage
	Messagef = pkgerr.WithMessagef
	Trace    = pkgerr.WithStack
	Cause    = pkgerr.Cause
	Errorf   = pkgerr.Errorf
	New      = pkgerr.New
)

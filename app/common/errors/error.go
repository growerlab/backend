package errors

import (
	"fmt"
	"strings"

	pkgerr "github.com/pkg/errors"
)

// 定义错误
const (
	// 非法参数
	invalidParameter = "InvalidParameter"
	// 无法找到
	notFound = "NotFound"
	// GraphQLError
	graphQLError = "GraphQLError"
	// 已存在
	alreadyExists = "AlreadyExists"
	// AccessDenied
	accessDenied = "AccessDenied"
	// sql错误
	sqlError = "SQLError"
)

// 定义错误原因
const (
	Invalid = "Invalid"
	// 无法找到属性（字段）
	NotFoundField = "NotFoundField"
	// 非法长度
	InvalidLength = "InvalidLength"
	// 失效，过期
	Expired = "Expired"
	// 已被使用过
	Used = "Used"
)

var P = InvalidParameterError

func InvalidParameterError(model, field, reason string) string {
	return mustCode(invalidParameter, model, field, reason)
}

func NotFoundError(model string) string {
	return mustCode(notFound, model)
}

func AlreadyExistsError(model, reason string) string {
	return mustCode(alreadyExists, model, reason)
}

func SQLError() string {
	return mustCode(sqlError)
}

func GraphQLError() string {
	return mustCode(graphQLError)
}

// 必须调用该方法生成<xxx>字符串，便于前端解析数据
func mustCode(parts ...string) string {
	if len(parts) == 0 {
		panic("parts is required")
	}
	return fmt.Sprintf("<%s>", strings.Join(parts, "."))
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

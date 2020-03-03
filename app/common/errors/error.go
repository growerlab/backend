package errors

import (
	"fmt"
	"strings"

	pkgerr "github.com/pkg/errors"
)

// TODO 目前有一些错误不应该在后端输出错误，例如 NotFound， 此类错误是没有必要的

// 定义错误
const (
	// 非法参数
	invalidParameter = "InvalidParameter"
	// 无法找到
	notFoundError = "NotFoundError"
	// GraphQLError
	graphQLError = "GraphQLError"
	// 已存在
	alreadyExists = "AlreadyExists"
	// AccessDenied
	accessDeniedError = "AccessDeniedError"
	// sql错误
	sqlError = "SQLError"
	// 未登录
	unauthorized = "Unauthorized"
	// PermissionError
	permissionError = "PermissionError"
	// 仓库
	repositoryError = "RepositoryError"
)

// 定义错误原因
const (
	// 非法的
	Invalid = "Invalid"
	// 无法找到属性（字段）
	NotFoundField = "NotFoundField"
	// 非法长度
	InvalidLength = "InvalidLength"
	// 失效，过期
	Expired = "Expired"
	// 已被使用过
	Used = "Used"
	// 不匹配
	NotEqual = "NotEqual"
	// 已存在
	AlreadyExists = "AlreadyExists"
	// 未激活
	NotActivated = "NotActivated"
	// 仓库服务异常
	SvcServerNotReady = "SvcServerNotReady"
	// 无权限
	NoPermission = "NoPermission"
)

var P = InvalidParameterError

func InvalidParameterError(model, field, reason string) string {
	return mustCode(invalidParameter, model, field, reason)
}

func NotFoundError(model string) string {
	return mustCode(notFoundError, model)
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

func Unauthorize() string {
	return mustCode(unauthorized)
}

func AccessDenied(model, reason string) string {
	return mustCode(accessDeniedError, model, reason)
}

func PermissionError(reason string) string {
	return mustCode(permissionError, reason)
}

func RepositoryError(reason string) string {
	return mustCode(repositoryError, reason)
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

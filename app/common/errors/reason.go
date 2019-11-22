package errors

// 定义Reason

const (
	// 非法参数
	InvalidParameter = "InvalidParameter"
	// 无法找到
	NotFound = "NotFound"
	// 无法找到属性（字段）
	NotFoundField = "NotFoundField"
	// sql错误
	SqlError = "SqlError"
	// 非法密码
	InvalidPassword = "InvalidPassword"
	// 不等于
	NotEqual = "NotEqual"
	// 失效
	Expired = "Expired"
	// 已被使用过
	Used = "Used"
)

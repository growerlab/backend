package errors

// 定义暴露出来的方法

var (
	P                     = invalidParameterError
	InvalidParameterError = invalidParameterError
	Trace                 = trace
	Sql                   = sqlError
	NotFoundError         = notFoundError
)

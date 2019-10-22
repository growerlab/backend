package user

type CreateUserPayload struct {
	Email    string
	Password string
	Username string
	Name     string
}

// 创建用户
// 1. 将用户信息添加到数据库中
// 2. 发送验证邮件（这里可以考虑使用KeyDB来建立邮件发送队列，避免重启进程后，发送任务丢失）
// 3.
//
func CreateUser(payload *CreateUserPayload) error {
	return nil
}

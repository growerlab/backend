package common

// 上下文
const (
	ContextRepository = 1001 // 仓库上下文
)

// 用户域
const (
	UserDomainSuperAdmin      = 2000 // 超级管理员
	UserDomainPerson          = 2001 // 个人
	UserDomainRepositoryOwner = 2002 // 仓库创建者
)

// 权限
const (
	PermissionViewRepository = 3101 // 「查看仓库」权限
)

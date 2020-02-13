package common

const (
	NamespaceEveryone = -1 // 特殊命名空间：访客
)

// 上下文
const (
	ContextRepository = 1001 // 仓库上下文
)

// 用户域
const (
	UserDomainSuperAdmin      = 2000 // 超级管理员
	UserDomainPerson          = 2001 // 个人
	UserDomainRepositoryOwner = 2002 // 仓库创建者
	UserDomainEveryone        = 2003 // 每个人（含访客）
)

// 权限
const (
	PermissionViewRepository  = 3101 // 「查看仓库」权限
	PermissionCloneRepository = 3102 // 「clone仓库」权限
	PermissionPushRepository  = 3103 // 「push仓库」权限
)

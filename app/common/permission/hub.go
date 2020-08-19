package permission

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/growerlab/backend/app/common/context"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/common/userdomain"
	"github.com/growerlab/backend/app/model/db"
	permModel "github.com/growerlab/backend/app/model/permission"
	"github.com/growerlab/backend/app/utils/timestamp"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFoundRule = errors.New("not found permission rule")
)

type PermissionsFunc func(src sqlx.Queryer, code int, c *context.Context) ([]*permModel.Permission, error)

type Rule struct {
	// Code 具体的权限
	Code int

	// ConstraintUserDomains「约束」权限允许的用户域（例如个人、组织成员等）
	// - 在添加相关权限到数据库时，需要该参数进行验证
	ConstraintUserDomains []int

	// BuiltInUserDomains 默认的、不可删除的特殊用户域（或者说用户角色），例如：「仓库创建者」等等
	// 这里的默认角色，默认就拥有Code所代表的权限
	// - 在构建权限缓存时，这里的用户域将一起初始化到缓存中
	BuiltInUserDomains []int
}

type Hub struct {
	ruleMap       map[int]*Rule
	userDomainHub map[int]UserDomainDelegate
	contextHub    map[int]ContextDelegate

	// PermissionsByContextFunc 独立出来，灵活实现数据源
	// 必须实现
	PermissionsByContextFunc PermissionsFunc
	// DBCtx 数据库操作对象; 内存数据库操作对象等
	DBCtx *context.DBContext
}

func NewPermissionHub(src sqlx.Queryer, memdb *redis.Client) *Hub {
	return &Hub{
		DBCtx: &context.DBContext{
			Src:   src,
			MemDB: memdb,
		},
		ruleMap:                  make(map[int]*Rule),
		userDomainHub:            make(map[int]UserDomainDelegate),
		contextHub:               make(map[int]ContextDelegate),
		PermissionsByContextFunc: permModel.ListPermissionsByContext,
	}
}

func (p *Hub) RegisterRules(rules []*Rule) error {
	for _, r := range rules {
		if _, exist := p.ruleMap[r.Code]; !exist {
			p.ruleMap[r.Code] = r
		} else {
			return fmt.Errorf("permission rule: %d exist", r.Code)
		}
	}
	return nil
}

func (p *Hub) RegisterUserDomains(userDomains []UserDomainDelegate) error {
	for _, u := range userDomains {
		if _, exist := p.userDomainHub[u.Type()]; !exist {
			p.userDomainHub[u.Type()] = u
		} else {
			return fmt.Errorf("permission userdomain: %s exist", u.TypeLabel())
		}
	}
	return nil
}

func (p *Hub) RegisterContexts(contexts []ContextDelegate) error {
	for _, c := range contexts {
		if _, exist := p.contextHub[c.Type()]; !exist {
			p.contextHub[c.Type()] = c
		} else {
			return fmt.Errorf("permission context: %s exist", c.TypeLabel())
		}
	}
	return nil
}

func (p *Hub) CheckCache(namespaceID int64, c *context.Context, code int, rebuild bool) error {
	var (
		keyUser        = p.keyUser(namespaceID)
		withPermission = p.keyContextWithPermission(code, c)
		keyStamp       = p.keyStamp()
	)

	if rebuild {
		lastUpdateStamp, err := p.DBCtx.MemDB.HGet(keyStamp, withPermission).Int64()
		if err != nil && err != redis.Nil {
			return errors.Trace(err)
		}

		existPermissionStamp, err := p.DBCtx.MemDB.HGet(keyUser, withPermission).Int64()
		if err != nil && err != redis.Nil {
			return errors.Trace(err)
		}

		mustRebuild := existPermissionStamp == 0
		if !mustRebuild {
			mustRebuild = lastUpdateStamp > existPermissionStamp
		}
		if mustRebuild {
			// rebuild
			rule, ok := p.ruleMap[code]
			if !ok {
				return ErrNotFoundRule
			}
			if err := p.buildCache(rule, c); err != nil {
				return err
			}
		}
	}

	if b := p.DBCtx.MemDB.HExists(keyUser, withPermission); !b.Val() {
		return errors.New(errors.PermissionError(errors.NoPermission))
	}
	return nil
}

// buildCache 重新构建缓存
// 这里之所以传rule，因为希望rebuild时，尽量只构建小一些的颗粒度缓存
// - 每天凌晨12点自动过期
func (p *Hub) buildCache(rule *Rule, c *context.Context) error {
	userDomains, err := p.listUserDomainsByContext(rule, c)
	if err != nil {
		return err
	}
	if len(userDomains) == 0 {
		return nil
	}

	userContextSet := make(map[string]string)

	for _, u := range userDomains {
		ud, ok := p.userDomainHub[u.Type]
		if !ok {
			return errors.Errorf("not found userdomain: %d", u.Type)
		}
		uids, err := ud.Eval(NewEvalArgs(c, u, p.DBCtx))
		if err != nil {
			return err
		}
		if len(uids) == 0 {
			continue
		}
		//
		for _, id := range uids {
			key := p.keyUser(id)
			withPermissionKey := p.keyContextWithPermission(rule.Code, c)
			userContextSet[key] = withPermissionKey
		}
	}

	now := time.Now()
	todayEndTime := timestamp.DayEnd(now)

	pipe := p.DBCtx.MemDB.Pipeline()
	for key, keyCtxCode := range userContextSet {
		_ = pipe.HSet(key, keyCtxCode, now.Unix())
		_ = pipe.ExpireAt(key, todayEndTime)
		_ = pipe.HSet(p.keyStamp(), keyCtxCode, 0)
	}
	_, err = pipe.Exec()
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (p *Hub) listUserDomainsByContext(rule *Rule, c *context.Context) ([]*userdomain.UserDomain, error) {

	// 默认增加超级管理员的用户域，即超级管理员
	// 这样超级管理员默认就拥有所有的权限
	// userDomains = append(userDomains, &userdomain.UserDomain{
	// 	Type: userdomain.TypeSuperAdmin,
	// })

	permissions, err := p.PermissionsByContextFunc(p.DBCtx.Src, rule.Code, c)
	if err != nil {
		return nil, err
	}

	udLength := len(rule.BuiltInUserDomains) + len(permissions)

	userDomains := make([]*userdomain.UserDomain, 0, udLength)
	for _, domain := range rule.BuiltInUserDomains {
		userDomains = append(userDomains, &userdomain.UserDomain{
			Type: domain,
		})
	}

	for _, p := range permissions {
		userDomains = append(userDomains, &userdomain.UserDomain{
			Type:  p.UserDomainType,
			Param: p.UserDomainParam,
		})
	}
	return userDomains, nil
}

func (p *Hub) keyUser(uid int64) string {
	return db.BaseKeyBuilder(fmt.Sprintf("permission:%d", uid)).String()
}

func (p *Hub) keyContextWithPermission(code int, c *context.Context) string {
	key := fmt.Sprintf("%d:%d:%d:%d", c.Type, c.Param1, c.Param2, code)
	return db.BaseKeyBuilder(key).String()
}

// keyStamp 当permission表或者相关角色变动后，将更新 keyStamp HSET中的stamp，表示memDBKey需要被更新
func (p *Hub) keyStamp() string {
	return db.BaseKeyBuilder("permission", "stamp").String()
}

func (p *Hub) updateStamp(code int, c *context.Context) error {
	key := p.keyContextWithPermission(code, c)
	err := p.DBCtx.MemDB.HSet(p.keyStamp(), key, time.Now().Unix()).Err()
	return errors.Trace(err)
}

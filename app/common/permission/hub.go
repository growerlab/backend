package permission

import (
	"fmt"
	"strconv"
	"time"

	"github.com/growerlab/backend/app/common/context"
	"github.com/growerlab/backend/app/common/userdomain"

	"github.com/go-redis/redis/v7"
	"github.com/growerlab/backend/app/common/errors"
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
	nsID := strconv.FormatInt(namespaceID, 10)
	key := p.memdbKey(code, c)

	if rebuild {
		lastUpdateStamp, err := p.DBCtx.MemDB.HGet(p.stampKey(), key).Int64()
		if err != nil && err != redis.Nil {
			return errors.Trace(err)
		}

		valueStampMap, err := p.DBCtx.MemDB.HGetAll(key).Result()
		if err != nil && err != redis.Nil {
			return errors.Trace(err)
		}

		mustRebuild := len(valueStampMap) == 0
		if !mustRebuild {
			for _, oldStampRaw := range valueStampMap {
				oldStamp, _ := strconv.ParseInt(oldStampRaw, 10, 64)
				mustRebuild = lastUpdateStamp > oldStamp
				break
			}
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

	if b := p.DBCtx.MemDB.HExists(key, nsID); !b.Val() {
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

	userIDValues := make(map[string]interface{})
	now := time.Now()
	stamp := now.UnixNano()

	for _, u := range userDomains {
		ud, ok := p.userDomainHub[u.Type]
		if !ok {
			return errors.Errorf("not found userdomain: %d", u.Type)
		}
		IDs, err := ud.Eval(NewEvalArgs(c, u, p.DBCtx))
		if err != nil {
			return err
		}
		if len(IDs) == 0 {
			continue
		}
		//
		for _, id := range IDs {
			idStr := strconv.FormatInt(id, 10)
			userIDValues[idStr] = stamp
		}
	}

	todayEndTime := timestamp.DayEnd(now)
	key := p.memdbKey(rule.Code, c)
	pipe := p.DBCtx.MemDB.Pipeline()
	_ = pipe.Del(key)
	_ = pipe.HMSet(key, userIDValues)
	_ = pipe.ExpireAt(key, todayEndTime)
	_ = pipe.HSet(p.stampKey(), key, 0)
	_, err = pipe.Exec()
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (p *Hub) listUserDomainsByContext(rule *Rule, c *context.Context) ([]*userdomain.UserDomain, error) {
	userDomains := make([]*userdomain.UserDomain, len(rule.BuiltInUserDomains))
	for i, domain := range rule.BuiltInUserDomains {
		userDomains[i] = &userdomain.UserDomain{
			Type: domain,
		}
	}

	// 默认增加超级管理员的用户域，即超级管理员
	// 这样超级管理员默认就拥有所有的权限
	userDomains = append(userDomains, &userdomain.UserDomain{
		Type: userdomain.TypeSuperAdmin,
	})

	permissions, err := p.PermissionsByContextFunc(p.DBCtx.Src, rule.Code, c)
	if err != nil {
		return nil, err
	}

	for _, p := range permissions {
		userDomains = append(userDomains, &userdomain.UserDomain{
			Type:  p.UserDomainType,
			Param: p.UserDomainParam,
		})
	}
	return userDomains, nil
}

func (p *Hub) memdbKey(code int, c *context.Context) string {
	return db.BaseKeyBuilder(fmt.Sprintf("permission:%d:context:%d:%d:%d", code, c.Type, c.Param1, c.Param2)).String()
}

// stampKey 当permission表或者相关角色变动后，将更新stampKey HSET中的stamp，表示memdbKey需要被更新
func (p *Hub) stampKey() string {
	return db.BaseKeyBuilder("permission", "stamp").String()
}
func (p *Hub) updateKeyStamp(code int, c *context.Context) error {
	key := p.memdbKey(code, c)
	err := p.DBCtx.MemDB.HSet(p.stampKey(), key, time.Now().UnixNano()).Err()
	return errors.Trace(err)
}

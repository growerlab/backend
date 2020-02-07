package permission

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/growerlab/backend/app/common/ctx"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/common/permission/common"
	permModel "github.com/growerlab/backend/app/model/permission"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFoundRule = errors.New("not found permission rule")
)

type Rule struct {
	// Code 具体的权限
	Code int
	// ConstraintUserDomains「约束」权限允许的用户域（例如个人、组织成员等）
	ConstraintUserDomains []int
	// BuiltInUserDomains 默认的、不可删除的特殊用户域（或者说用户角色），例如：「仓库创建者」等等
	// 这里的默认角色，默认就拥有Code所代表的权限
	BuiltInUserDomains []int
}

type PermissionHub struct {
	ruleMap       map[int]*Rule
	userDomainHub map[int]common.UserDomainDelegate
	contextHub    map[int]common.ContextDelegate

	// DBCtx 数据库操作对象; 内存数据库操作对象等
	DBCtx *ctx.DBContext
}

func NewPermissionHub(src sqlx.Queryer, memdb *redis.Client) *PermissionHub {
	return &PermissionHub{
		DBCtx: &ctx.DBContext{
			Src:   src,
			MemDB: memdb,
		},
		ruleMap: make(map[int]*Rule),
	}
}

func (p *PermissionHub) RegisterRules(rules []*Rule) error {
	for _, r := range rules {
		if _, exist := p.ruleMap[r.Code]; !exist {
			p.ruleMap[r.Code] = r
		} else {
			return fmt.Errorf("permission rule: %d exist", r.Code)
		}
	}
	return nil
}

func (p *PermissionHub) RegisterUserDomains(userDomains []common.UserDomainDelegate) error {
	for _, u := range userDomains {
		if _, exist := p.userDomainHub[u.Type()]; !exist {
			p.userDomainHub[u.Type()] = u
		} else {
			return fmt.Errorf("permission userdomain: %s exist", u.TypeLabel())
		}
	}
	return nil
}

func (p *PermissionHub) RegisterContexts(contexts []common.ContextDelegate) error {
	for _, c := range contexts {
		if _, exist := p.contextHub[c.Type()]; !exist {
			p.contextHub[c.Type()] = c
		} else {
			return fmt.Errorf("permission context: %s exist", c.TypeLabel())
		}
	}
	return nil
}

func (p *PermissionHub) CheckCache(namespaceID int64, c *ctx.Context, code int, rebuild bool) error {
	nsID := strconv.FormatInt(namespaceID, 10)
	key := p.memdbKey(code, c)

	if rebuild {
		lastUpdateStamp, err := p.DBCtx.MemDB.HGet(p.stampKey(), key).Int64()
		if err != nil {
			return errors.Trace(err)
		}

		valueStampMap, err := p.DBCtx.MemDB.HGetAll(key).Result()
		if err != nil {
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
		return errors.New(errors.PermisstionErrror(errors.NoPermission))
	}
	return nil
}

// buildCache 重新构建缓存
// 这里之所以传rule，因为希望rebuild时，尽量只构建小一些的颗粒度
// - 每天凌晨12点自动过期
// -
func (p *PermissionHub) buildCache(rule *Rule, c *ctx.Context) error {
	userdomains, err := p.listUserDomainsByContext(rule, c)
	if err != nil {
		return err
	}
	if len(userdomains) == 0 {
		return nil
	}

	userIDValues := make(map[string]interface{})
	now := time.Now()
	stamp := now.UnixNano()

	for _, u := range userdomains {
		ud, ok := p.userDomainHub[u.Type]
		if !ok {
			return errors.Errorf("not found userdomain: %d", u.Type)
		}
		IDs, err := ud.BatchEval(p.DBCtx, &common.EvalArgs{
			Ctx: c,
		})
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

	todayEndTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
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

func (p *PermissionHub) listUserDomainsByContext(rule *Rule, c *ctx.Context) ([]*ctx.UserDomain, error) {
	userdomains := make([]*ctx.UserDomain, len(rule.BuiltInUserDomains))
	for i := range rule.BuiltInUserDomains {
		userdomains[i] = &ctx.UserDomain{
			Type: rule.BuiltInUserDomains[i],
		}
	}

	// 默认增加超级管理员的用户域，既超级管理员
	userdomains = append(userdomains, &ctx.UserDomain{
		Type: common.UserDomainSuperAdmin,
	})

	permissions, err := permModel.ListPermissionsByContext(p.DBCtx.Src, rule.Code, c)
	if err != nil {
		return nil, err
	}

	for _, p := range permissions {
		userdomains = append(userdomains, &ctx.UserDomain{
			Type:  p.UserDomainType,
			Param: p.UserDomainParam,
		})
	}
	return userdomains, nil
}

func (p *PermissionHub) memdbKey(code int, c *ctx.Context) string {
	return fmt.Sprintf("permission:%d:context:%d:%d:%d", code, c.Type, c.Param1, c.Param2)
}

// stampKey 当permission表或者相关角色变动后，将更新stampKey HSET中的stamp，表示memdbKey需要被更新
func (p *PermissionHub) stampKey() string {
	return fmt.Sprintf("permission:stamp")
}
func (p *PermissionHub) updateKeyStamp(code int, c *ctx.Context) error {
	key := p.memdbKey(code, c)
	err := p.DBCtx.MemDB.HSet(p.stampKey(), key, time.Now().UnixNano()).Err()
	return errors.Trace(err)
}

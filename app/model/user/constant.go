package user

// 不允许用户注册的关键字
//
var InvalidUsernameList = []string{
	"admin",
	"username",
	"udmin",
	"settings",
	"setting",
	"profile",
	"dashboard",
	"organizations",
	"repositories",
	"repository",
	"git",
	"gist",
	"team",
	"pulls",
	"issues",
	"explore",
	"blog",
	"home",
	"new",
	"project",
	"projects",
	"help",
	"signin",
	"signout",
}

var InvalidUsernameSet = make(map[string]struct{})

func init() {
	for i := range InvalidUsernameList {
		InvalidUsernameSet[InvalidUsernameList[i]] = struct{}{}
	}
}

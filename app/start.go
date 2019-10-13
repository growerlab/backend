package app

import "github.com/growerlab/backend/app/service/graphql/think"

// 需要初始化的全局数据放在这里
//	eg. onStart(job.Work)
//
func Init() {
	onStart(think.InitGraphQL)
}

func onStart(fn func() error) {
	if err := fn(); err != nil {
		panic(err)
	}
}

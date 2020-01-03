package repository

import "github.com/growerlab/backend/app/service"

func CreateRepository(reqRepo *service.NewRepository) (bool, error) {
	err := validate(reqRepo)
	if err != nil {
		return false, err
	}
	return false, nil
}

// validate
//	req.Owner 是否是自己
//	req.Name 名称是否合法、是否重名
func validate(reqRepo *service.NewRepository) (err error) {
	return nil
}

package repository

import (
	"context"
	"xiaoweishu/webook/internal/domain"
	"xiaoweishu/webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository struct {
	dao *dao.UserDAO
}

// NewUserRepository 要用的东西都不要内部初始化，让他从外面传入参数后调用new方法来进行初始化
func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{dao: dao}

}

// Create repository已经是到达了数据库层面，所以这里不再会有注册的概念，而是要涉及到数据库的操作，所以这里写create
func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	//在这里操作缓存
}
func (r *UserRepository) FindById(int64) {

}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Email:    u.Email,
		Password: u.Password,
	}, nil

}

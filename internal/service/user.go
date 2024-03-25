package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"xiaoweishu/webook/internal/domain"
	"xiaoweishu/webook/internal/repository"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}

}

// SignUp 不能用上一层的结构体数据，而是要用下一层传上来的数据进行操作
// service层的注册要做的事情,不需要传指针，传指针还需要进行判空
// 加密更应该放在sercive层面，这样后面的层也是加密的，更安全
func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash) //返回的加密密码也是字节切片，要转换成字符串

	return svc.repo.Create(ctx, u)
	//再service层直接去调用下一层，service层面一般都不用去做什么具体操作，有错误则返回，再到上一层去处理错误
}
func (svc *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	//先找用户
	u, err := svc.repo.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	//这是找不到邮箱的情况
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	//通过邮箱找到了这个用户，说明数据库中是有这个人的，说明他已经注册过了，所以现在要开始进行比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		//后面要在这里进行打印日志
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return domain.User{}, nil

}

//

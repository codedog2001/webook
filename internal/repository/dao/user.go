package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	//创建时间和更新时间 毫秒数
	Ctime int64
	Utime int64
}

type UserDAO struct {
	db *gorm.DB
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli() //毫秒数在高并发的场景下更有优势
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	var mysqlErr *mysql.MySQLError
	//先判断错误是不是由mysql引起的，是的话再判断是不是由唯一索引错误引起的
	if errors.As(err, &mysqlErr) {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			//只设置一个unique邮箱，所以发送唯一索引冲突的时候，就一定是邮箱出问题了
			return ErrUserDuplicateEmail
		}
	}
	return err //gorm会去数据库中创建u的实列

}

// NewUserDAO new函数都只是做一个初始化操作
func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}
func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	//取到的数据放在u里面
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	//err:=dao.db.WithContext(ctx).First(&u,"email=?",email).Error
	//两种写法都可以
	//如果err= gorm.ErrRecordNotFound,那么就会自动返回errusernotfound
	return u, err

}

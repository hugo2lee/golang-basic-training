package repository

import (
	"context"

	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	// SELECT * FROM `users` WHERE `email`=?
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) UpdateById(ctx context.Context, user domain.User) error {
	return r.dao.UpdateById(ctx, dao.User{
		Id:           user.Id,
		Nickname:     user.Nickname,
		Birthday:     user.Birthday,
		Introduction: user.Introduction,
	})
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 先从 cache 里面找
	// 再从 dao 里面找
	// 找到了回写 cache
	u, err := r.dao.FindById(ctx, id)
	if err != nil {
		if err == dao.ErrUserNotFound {
			err = ErrUserNotFound
		}
		return domain.User{}, err
	}
	return domain.User{
		Id:           u.Id,
		Email:        u.Email,
		Nickname:     u.Nickname,
		Birthday:     u.Birthday,
		Introduction: u.Introduction,
		Ctime:        u.Ctime,
	}, nil
}

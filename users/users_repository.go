package users

import (
	"gorm.io/gorm"
)

type UsersRepository interface {
	SaveUserRepository(user User) (User, error)
	FindUserByEmailRepository(email string) (User, error)
	FindUserByIdRepository(ID int) (User, error)
	UpdateUserRepository(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) SaveUserRepository(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindUserByEmailRepository(email string) (User, error) {
	var user User

	err := r.db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindUserByIdRepository(ID int) (User, error) {
	var user User

	err := r.db.Where("id = ?", ID).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) UpdateUserRepository(user User) (User, error) {

	err := r.db.Save(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

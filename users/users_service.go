package users

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UsersService interface {
	RegisterUserService(input RegisterUserInput) (User, error)
	LoginUserService(input LoginInput) (User, error)
	CheckEmailAvailabilityService(input CheckEmailInput) (bool, error)
	GetUserByIDService(ID int) (User, error)
	UpdateUserByIDService(ID int, input UpdateUserInput) (User, error)
	UpdateAvatarByIDService(ID int, fileLocation string) (User, error)
}

type service struct {
	repository UsersRepository
}

func NewUsersService(repository UsersRepository) *service {
	return &service{repository}
}

func (s *service) RegisterUserService(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.PasswordHash = input.Password

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	user.PasswordHash = string(passwordHash)

	newUser, err := s.repository.SaveUserRepository(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (s *service) LoginUserService(input LoginInput) (User, error) {
	user, err := s.repository.FindUserByEmailRepository(input.Email)
	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) CheckEmailAvailabilityService(input CheckEmailInput) (bool, error) {
	user, err := s.repository.FindUserByEmailRepository(input.Email)
	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *service) GetUserByIDService(ID int) (User, error) {
	user, err := s.repository.FindUserByIdRepository(ID)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) UpdateUserByIDService(ID int, input UpdateUserInput) (User, error) {
	user, err := s.repository.FindUserByIdRepository(ID)
	if err != nil {
		return user, err
	}

	user.Name = input.Name
	user.Email = input.Email
	fmt.Println(input.Name)
	updateUser, err := s.repository.UpdateUserRepository(user)
	if err != nil {
		return updateUser, err
	}

	return updateUser, nil
}

func (s *service) UpdateAvatarByIDService(ID int, fileLocation string) (User, error) {
	user, err := s.repository.FindUserByIdRepository(ID)
	if err != nil {
		return user, err
	}

	user.AvatarFileName = fileLocation
	updateUser, err := s.repository.UpdateUserRepository(user)
	if err != nil {
		return updateUser, err
	}

	return updateUser, nil
}

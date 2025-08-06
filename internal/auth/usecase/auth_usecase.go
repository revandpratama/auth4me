package usecase

import (
	"errors"
	"fmt"

	"github.com/revandpratama/auth4me/internal/auth/dto"
	"github.com/revandpratama/auth4me/internal/auth/entity"
	"github.com/revandpratama/auth4me/internal/auth/repository"
	"github.com/revandpratama/auth4me/pkg"
)

type AuthUsecase interface {
	Login(email string, password string) (string, error)
	Register(registerRequest *dto.RegisterRequest) error
	GetUserByID(id string) (*entity.User, error)
}

type authUsecase struct {
	repository repository.AuthRepository
}

func NewAuthUsecase(repository repository.AuthRepository) AuthUsecase {
	return &authUsecase{
		repository: repository,
	}
}

func (u *authUsecase) Login(email string, password string) (string, error) {

	user, err := u.repository.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	if err := pkg.ValidatePassword(user.Password, password); err != nil {
		return "", err
	}
	mfaCompleted := false
	if user.MFAEnabled {
		// TODO : Validate MFA
	}

	token, err := pkg.GenerateToken(user, "local", mfaCompleted)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Bearer %s", token), nil
}

func (u *authUsecase) Register(registerRequest *dto.RegisterRequest) error {

	if registerRequest.Password != registerRequest.ConfirmPassword {
		return errors.New("password and confirm password does not match")
	}

	exists, err := u.repository.IsEmailExists(registerRequest.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	hashedPassword, err := pkg.EncryptPassword(registerRequest.Password)
	if err != nil {
		return err
	}
	registerRequest.Password = hashedPassword

	newUser := entity.User{
		Email:     registerRequest.Email,
		Password:  registerRequest.Password,
		FullName:  registerRequest.FullName,
		AvatarPath: registerRequest.AvatarPath,
	}

	if err := u.repository.CreateUser(&newUser); err != nil {
		return err
	}

	return nil
}

func (u *authUsecase) GetUserByID(id string) (*entity.User, error) {
	return u.repository.GetUserByID(id)
}

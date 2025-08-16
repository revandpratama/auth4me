package usecase

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/revandpratama/auth4me/internal/auth/dto"
	"github.com/revandpratama/auth4me/internal/auth/entity"
	"github.com/revandpratama/auth4me/internal/auth/repository"
	"github.com/revandpratama/auth4me/pkg"
)

type AuthUsecase interface {
	Login(email string, password string) (string, string, error)
	Register(registerRequest *dto.RegisterRequest) error
	RefreshToken(refreshToken string, accessToken string) (string, string, error)
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

func (u *authUsecase) Login(email string, password string) (string, string, error) {

	user, err := u.repository.GetUserByEmail(email)
	if err != nil {
		return "", "", err
	}

	if err := pkg.ValidatePassword(user.Password, password); err != nil {
		return "", "", err
	}
	mfaCompleted := false
	if user.MFAEnabled {
		// TODO : Validate MFA
	}
	log.Println("password validated")
	permissions, err := u.repository.GetUserPermissionsByRoleID(user.RoleID)
	if err != nil {
		return "", "", err
	}
	log.Println("permissions fetched")

	token, err := pkg.GenerateToken(user, "local", permissions, mfaCompleted)
	if err != nil {
		return "", "", err
	}
	log.Println("token generated")

	token = fmt.Sprintf("Bearer %s", token)

	refreshToken := uuid.NewString()
	pkg.SaveRefreshToken(refreshToken, pkg.TokenData{
		UserID:       user.ID,
		Email:        user.Email,
		RoleID:       user.RoleID,
		Provider:     "local",
		MFACompleted: mfaCompleted,
		ExpiresAt:    time.Now().Add(time.Hour * 24)},
	)

	log.Println("refresh token saved")

	return refreshToken, token, nil
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
		Email:      registerRequest.Email,
		Password:   registerRequest.Password,
		FullName:   registerRequest.FullName,
		RoleID:     registerRequest.RoleID,
		AvatarPath: registerRequest.AvatarPath,
	}

	if err := u.repository.CreateUser(&newUser); err != nil {
		return err
	}

	return nil
}

func (u *authUsecase) RefreshToken(refreshToken string, accessToken string) (string, string, error) {

	//Validate access token
	claims, err := pkg.ParseExpiredToken(accessToken)
	if err != nil {
		return "", "", err
	}

	//Validate refresh token in redis
	refreshTokenData, exists := pkg.GetRefreshToken(refreshToken)
	if !exists || time.Now().After(refreshTokenData.ExpiresAt) {
		return "", "", errors.New("refresh token expired")
	}

	if refreshTokenData.UserID != claims.UserID {
		return "", "", errors.New("refresh token and access token user id does not match")
	}

	//Generate new access token

	user := entity.User{
		ID:     claims.UserID,
		Email:  claims.Email,
		RoleID: claims.RoleID,
	}

	newAccessToken, err := pkg.GenerateToken(&user, claims.Provider, claims.Permissions, claims.MFACompleted)
	if err != nil {
		return "", "", err
	}

	newAccessToken = fmt.Sprintf("Bearer %s", newAccessToken)

	//Generate new refresh token
	data := pkg.TokenData{
		UserID:       claims.UserID,
		Email:        claims.Email,
		RoleID:       claims.RoleID,
		Provider:     claims.Provider,
		SessionID:    claims.SessionID,
		MFACompleted: claims.MFACompleted,
		ExpiresAt:    time.Now().Add(time.Hour * 1),
	}

	newRefreshToken := uuid.New().String()

	pkg.SaveRefreshToken(newRefreshToken, data)

	pkg.DeleteRefreshToken(refreshToken)

	return newRefreshToken, newAccessToken, nil
}

func (u *authUsecase) GetUserByID(id string) (*entity.User, error) {
	return u.repository.GetUserByID(id)
}

package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/revandpratama/auth4me/internal/auth/dto"
	"github.com/revandpratama/auth4me/internal/auth/entity"
	"github.com/revandpratama/auth4me/internal/auth/repository"
	"github.com/revandpratama/auth4me/pkg"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type OAuthUsecase interface {
	GetOAuthURL() (string, string)
	GoogleOAuthCallback(code string) (string, string, error)
}

type oauthUsecase struct {
	oauthCfg  *oauth2.Config
	authRepo  repository.AuthRepository
	oauthRepo repository.OAuthRepository
}

func NewOAuthUsecase(oauthCfg *oauth2.Config, authRepo repository.AuthRepository, oauthRepo repository.OAuthRepository) OAuthUsecase {
	return &oauthUsecase{
		oauthCfg:  oauthCfg,
		authRepo:  authRepo,
		oauthRepo: oauthRepo,
	}
}

func (u *oauthUsecase) GetOAuthURL() (string, string) {
	state, err := generateRandomState()
	if err != nil {
		return "", ""
	}
	return u.oauthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline), state
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (u *oauthUsecase) GoogleOAuthCallback(code string) (string, string, error) {
	ctx := context.Background()
	token, err := u.oauthCfg.Exchange(ctx, code)
	if err != nil {
		return "", "", fmt.Errorf("exchange failed: %w", err)
	}

	client := u.oauthCfg.Client(ctx, token)
	resp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		return "", "", fmt.Errorf("get user info failed: %w", err)
	}
	defer resp.Body.Close()

	var googleProfile dto.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleProfile); err != nil {
		return "", "", fmt.Errorf("decode user info failed: %w", err)
	}

	var userToTokenize *entity.User

	existingUser, err := u.authRepo.GetUserByEmail(googleProfile.Email)
	if err != nil {
		// not found → create new
		if err == gorm.ErrRecordNotFound {
			newUser := &entity.User{
				Email:         googleProfile.Email,
				FullName:      googleProfile.Name,
				AvatarPath:    googleProfile.Picture,
				EmailVerified: googleProfile.EmailVerified,
			}

			createdUser, err := u.authRepo.CreateUser(newUser)
			if err != nil {
				return "", "", err
			}

			userToTokenize = createdUser

			provider := &entity.OAuthProvider{
				UserID:       createdUser.ID,
				Provider:     "google",
				ProviderID:   googleProfile.Sub,
				AccessToken:  token.AccessToken,
				RefreshToken: token.RefreshToken,
				ExpiresAt:    token.Expiry,
			}
			if err := u.oauthRepo.CreateProvider(provider); err != nil {
				return "", "", err
			}
		} else {
			return "", "", err
		}
	} else {
		userToTokenize = existingUser
		provider, err := u.oauthRepo.GetProvider(userToTokenize.ID, "google")
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				provider = &entity.OAuthProvider{
					UserID:       userToTokenize.ID,
					Provider:     "google",
					ProviderID:   googleProfile.Sub,
					AccessToken:  token.AccessToken,
					RefreshToken: token.RefreshToken,
					ExpiresAt:    token.Expiry,
				}
				if err := u.oauthRepo.CreateProvider(provider); err != nil {
					return "", "", fmt.Errorf("failed to link google provider to existing user: %w", err)
				}
			} else {
				return "", "", err
			}
		} else {
			provider.AccessToken = token.AccessToken
			provider.RefreshToken = token.RefreshToken
			provider.ExpiresAt = token.Expiry
			if err := u.oauthRepo.UpdateProvider(provider); err != nil {
				return "", "", fmt.Errorf("failed to update provider tokens: %w", err)
			}
		}
	}

	if !userToTokenize.EmailVerified && googleProfile.EmailVerified {
		userToTokenize.EmailVerified = googleProfile.EmailVerified
		if err := u.authRepo.UpdateUser(userToTokenize); err != nil {
			return "", "", err
		}
	}

	mfaCompleted := false
	if userToTokenize.MFAEnabled {
		// TODO : Validate MFA
	}

	accessToken, err := pkg.GenerateToken(userToTokenize, "google", userToTokenize.Role.Permissions, mfaCompleted)
	if err != nil {
		return "", "", fmt.Errorf("generate token failed: %w", err)
	}

	// Generate Refresh Token
	refreshToken := uuid.NewString()
	pkg.SaveRefreshToken(refreshToken, pkg.TokenData{
		UserID:       userToTokenize.ID,
		Email:        userToTokenize.Email,
		RoleID:       userToTokenize.RoleID,
		Provider:     "google",
		MFACompleted: mfaCompleted,
		ExpiresAt:    time.Now().Add(time.Hour * 24)},
	)

	return refreshToken, accessToken, nil
}

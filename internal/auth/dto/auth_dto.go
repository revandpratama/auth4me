package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	FullName        string `json:"full_name" validate:"required"`
	AvatarPath       string `json:"avatar_path,omitempty" validate:"omitempty"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,eqfield=Password"`
}

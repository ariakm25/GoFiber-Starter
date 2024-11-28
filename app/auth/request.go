package auth

type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"email,required"`
	Password string `json:"password" form:"password" validate:"min=6,required"`
}

type RegisterRequest struct {
	Name                 string `json:"name" form:"name" validate:"min=3,max=100"`
	Email                string `json:"email" form:"email" validate:"required,email,unique=users.email"`
	Password             string `json:"password" form:"password" validate:"min=6"`
	PasswordConfirmation string `json:"password_confirmation" form:"password_confirmation" validate:"eqfield=Password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

type ValidateResetPasswordTokenRequest struct {
	Email string `json:"email" form:"email" validate:"required,email"`
	Token string `json:"token" form:"token" validate:"required"`
}

type ResetPasswordRequest struct {
	Email                string `json:"email" form:"email" validate:"required,email"`
	Token                string `json:"token" form:"token" validate:"required"`
	Password             string `json:"password" form:"password" validate:"min=6"`
	PasswordConfirmation string `json:"password_confirmation" form:"password_confirmation" validate:"eqfield=Password"`
}

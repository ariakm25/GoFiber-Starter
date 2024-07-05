package auth

type LoginRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"min=6,required"`
}

type RegisterRequest struct {
	Name                 string `json:"name" validate:"min=3,max=100"`
	Email                string `json:"email" validate:"email"`
	Password             string `json:"password" validate:"min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"eqfield=Password"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

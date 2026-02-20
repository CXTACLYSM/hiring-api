package dto

type RegisterDTO struct {
	Username             string `json:"username" validate:"required,min=3,max=50"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

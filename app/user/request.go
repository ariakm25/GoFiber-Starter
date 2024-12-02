package user

type CreateUserRequest struct {
	Name     string `json:"name" form:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" form:"email" validate:"required,email,unique=users.email"`
	Status   string `json:"status" form:"status" validate:"required,oneof=active inactive suspended"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	ID       string `json:"id" form:"id" validate:"required,exist=users.uid"`
	Name     string `json:"name" form:"name" validate:"omitempty,min=3,max=100"`
	Email    string `json:"email" form:"email" validate:"omitempty,email"`
	Status   string `json:"status" form:"status" validate:"omitempty,oneof=active inactive suspended"`
	Password string `json:"password" form:"password" validate:"omitempty,min=6"`
}

package dto

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4"`
}

type LoginRes struct {
	Token string `json:"token"`
}

type RegisterReq struct {
	Name     string `json:"name" validate:"required"`
	Role     string `json:"role"  validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4"`
}

type RegisterUserReq struct {
	Name     string `json:"name" validate:"required"`
	Role     string `json:"role"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4"`
}

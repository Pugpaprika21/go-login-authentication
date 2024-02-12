package dto

type UserLoginBodyRequest struct {
	ID       uint
	Username string `validate:"required"`
	Password string `validate:"required"`
	Email    string
	Token    string
}

type UserRegisterBodyRequest struct {
	ID       uint
	Username string `validate:"required"`
	Password string `validate:"required"`
	Email    string `validate:"required,email"`
	Token    string
}

type UserQueryRow struct {
	ID       uint   `gorm:"id"`
	Username string `gorm:"username"`
	Password string `gorm:"password"`
	Email    string `gorm:"email"`
	Token    string `gorm:"token"`
}

type UserRespone struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	TokenJWT string `json:"tokenJWT"`
}

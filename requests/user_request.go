package requests

type CreateUserRequest struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	Role        int16  `json:"role" binding:"omitempty,oneof=1 2"`
}

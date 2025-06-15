package handler


func ceilDiv(a, b int) int {
	if b == 0 {
		return 0
	}
	return (a + b - 1) / b
}

type LoginRequest struct {
	Login string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	Token string `json:"token" binding:"required"`	
}


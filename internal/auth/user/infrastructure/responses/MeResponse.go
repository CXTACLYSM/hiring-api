package responses

import "github.com/CXTACLYSM/hiring-api/internal/auth/user/domain/entities"

type MeResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewMeResponse(user *entities.User) *MeResponse {
	return &MeResponse{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}
}

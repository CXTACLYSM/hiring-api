package responses

type RegisterResponse struct {
	Token string `json:"token"`
}

func NewRegisterResponse(token string) *RegisterResponse {
	return &RegisterResponse{
		Token: token,
	}
}

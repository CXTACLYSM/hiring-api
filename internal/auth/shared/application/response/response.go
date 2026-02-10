package response

type SuccessResponse struct {
	Ok   bool `json:"ok"`
	Data any  `json:"data"`
}

type ErrorResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

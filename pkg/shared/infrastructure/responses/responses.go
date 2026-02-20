package responses

type InfoResponse struct {
	Version string
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type ValidationErrorsResponse struct {
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewInfoResponse(version string) *InfoResponse {
	return &InfoResponse{
		Version: version,
	}
}

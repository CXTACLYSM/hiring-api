package dto

type FindListDTO struct {
	Name    string `json:"name" validate:"omitempty,min=3,max=100"`
	Content string `json:"content" validate:"omitempty,min=9,max=300"`
	UserId  string `json:"-" validate:"required,uuid"`
}

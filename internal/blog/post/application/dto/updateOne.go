package dto

type UpdateOneDTO struct {
	Id      string `json:"-" validate:"required,uuid"`
	Name    string `json:"name" validate:"required,min=3,max=100"`
	Content string `json:"content" validate:"required,min=9,max=300"`
	UserId  string `json:"-" validate:"required,uuid"`
}

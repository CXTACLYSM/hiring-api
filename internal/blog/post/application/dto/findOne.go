package dto

type FindOneDTO struct {
	Id     string `json:"-" validate:"required,uuid"`
	UserId string `json:"-" validate:"required,uuid"`
}

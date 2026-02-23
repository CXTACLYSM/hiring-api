package deleteOne

type Command struct {
	UserId string
	Id     string
}

type Handler interface {
	Handle(Command) error
}

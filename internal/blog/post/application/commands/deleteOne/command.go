package deleteOne

type Command struct {
	Id     string
	UserId string
}

type Handler interface {
	Handle(Command) error
}

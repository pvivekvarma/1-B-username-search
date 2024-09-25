package command

type Command interface {
	Execute() error
	SetNext(Command)
}

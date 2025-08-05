package contracts

import (
	"io"
)

const (
	CommandHelp = "help"
	CommandList = "list"
	CommandOptions = "options"
	CommandRun = "run"
)

type Printer func(writer io.Writer, object any)

type Runner[T any] interface {
	Help(writer io.Writer)
	Print(writer io.Writer, requests []*T)
	Run(requests []*T) (err error)
	Setup(settings Settings) (err error)
}

package notifier

import "fmt"

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

type TerminalNotifier struct{}

func NewTerminalNotifiter() *TerminalNotifier {
	return &TerminalNotifier{}
}

func (t *TerminalNotifier) Info(msg string) {
	fmt.Println(Blue + msg + Reset)
}

func (t *TerminalNotifier) Success(msg string) {
	fmt.Println(Green + msg + Reset)
}

func (t *TerminalNotifier) Error(msg string) {
	fmt.Println(Red + msg + Reset)
}

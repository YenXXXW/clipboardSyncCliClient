package notifier

import "fmt"

type TerminalNotifier struct{}

func NewTerminalNotifiter() *TerminalNotifier {
	return &TerminalNotifier{}
}

func (t *TerminalNotifier) Print(msg string) {
	fmt.Println(msg)
}

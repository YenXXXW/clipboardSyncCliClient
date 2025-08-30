package formatter

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

type Formatter struct{}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func (f *Formatter) Success(msg string) string {
	return Green + msg + Reset
}

func (f *Formatter) Error(msg string) string {
	return Red + msg + Reset
}

func (f *Formatter) Info(msg string) string {
	return Blue + msg + Reset
}

func (f *Formatter) Warn(msg string) string {
	return Yellow + msg + Reset
}

package utils

const (
	BlueText  = "\033[34m"
	RedText   = "\033[31m"
	GreenText = "\033[32m"
	ResetText = "\033[0m"
)

// TextColor represents a text color.
type TextColor string

// Text colors.
const (
	Blue  TextColor = BlueText
	Red   TextColor = RedText
	Green TextColor = GreenText
	Reset TextColor = ResetText
)

// Colorize returns the text with the specified color.
func Colorize(text string, color TextColor) string {
	return string(color) + text + string(Reset)
}

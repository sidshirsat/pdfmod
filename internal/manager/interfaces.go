package manager

// PDFManagerInterface defines the methods available on PDFManager.
type PDFManagerInterface interface {
	Execute() error
}

// Prompter is an interface for prompting the user for input.
type Prompter interface {
	PromptUser(prompt string) string
}

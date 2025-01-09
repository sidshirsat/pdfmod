package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Prompter interface {
	PromptUser(prompt string) string
}

// ConsolePrompter prompts the user for input via the console.
type ConsolePrompter struct {
	Prompter Prompter
}

func (cp *ConsolePrompter) PromptUser(prompt string) string {
	return cp.Prompter.PromptUser(prompt)
}

type BasePrompter struct{}

func (bp *BasePrompter) PromptUser(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

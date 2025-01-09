package utils_test

import (
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sidshirsat/pdfmod/internal/utils"
	"github.com/sidshirsat/pdfmod/mocks"
)

func TestPromptUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPrompter := mocks.NewMockPrompter(ctrl)

	// Inject the mockPrompter into consolePrompter
	consolePrompter := &utils.ConsolePrompter{
		Prompter: mockPrompter,
	}

	mockPrompter.EXPECT().PromptUser("Enter the number of your choice: ").Return("2").Times(1)

	choice := consolePrompter.PromptUser("Enter the number of your choice: ")
	if choice != "2" {
		t.Fatalf("expected choice 2, got %s", choice)
	}
}

func TestPromptUser_EmptyInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPrompter := mocks.NewMockPrompter(ctrl)

	// Inject the mockPrompter into consolePrompter
	consolePrompter := &utils.ConsolePrompter{
		Prompter: mockPrompter,
	}

	mockPrompter.EXPECT().PromptUser("Enter something: ").Return("").Times(1)

	choice := consolePrompter.PromptUser("Enter something: ")
	if choice != "" {
		t.Fatalf("expected empty choice, got %s", choice)
	}
}

func TestBasePrompter_PromptUser(t *testing.T) {
	// Simulate user input
	input := "Test input\n"
	expected := strings.TrimSpace(input)

	// Replace os.Stdin with a bytes.Buffer containing the simulated input
	oldStdin := os.Stdin                   // Save the original os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore it after the test

	r, w, _ := os.Pipe()
	w.WriteString(input)
	w.Close()
	os.Stdin = r

	// Create an instance of BasePrompter
	prompter := &utils.BasePrompter{}

	// Call the method and check the result
	result := prompter.PromptUser("Enter something: ")
	if result != expected {
		t.Fatalf("expected '%s', got '%s'", expected, result)
	}
}

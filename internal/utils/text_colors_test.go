package utils_test

import (
	"testing"

	"github.com/sidshirsat/pdfmod/internal/utils"
)

func TestColorize(t *testing.T) {
	tests := []struct {
		text     string
		color    utils.TextColor
		expected string
	}{
		{"Hello, World!", utils.Blue, utils.BlueText + "Hello, World!" + utils.ResetText},
		{"Error!", utils.Red, utils.RedText + "Error!" + utils.ResetText},
		{"Success", utils.Green, utils.GreenText + "Success" + utils.ResetText},
		{"Reset", utils.Reset, utils.ResetText + "Reset" + utils.ResetText},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result := utils.Colorize(tt.text, tt.color)
			if result != tt.expected {
				t.Errorf("Colorize(%q, %q) = %q; want %q", tt.text, tt.color, result, tt.expected)
			}
		})
	}
}

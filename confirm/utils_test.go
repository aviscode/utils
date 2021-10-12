package confirm

import (
	"os"
	"testing"
)

func TestAskForConfirmation(t *testing.T) {
	tests := []struct {
		name     string
		question string
		answer   string
		want     bool
	}{
		{"test yes", "test yes:", "yes", true},
		{"test YES", "test YES:", "YES", true},
		{"test y", "test y:", "y", true},
		{"test Y", "test Y:", "Y", true},
		{"test no", "test no:", "no", false},
		{"test NO", "test NO:", "NO", false},
		{"test n", "test n:", "n", false},
		{"test N", "test N:", "N", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := []byte(tt.answer)
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			if _, err = w.Write(input); err != nil {
				t.Error(err)
			}
			if err = w.Close(); err != nil {
				t.Error(err)
			}
			stdin := os.Stdin
			// Restore stdin right after the test.
			defer func() { os.Stdin = stdin }()
			os.Stdin = r

			if got := AskForConfirmation(tt.question); got != tt.want {
				t.Errorf("AskForConfirmation() = %v, want %v", got, tt.want)
			}
		})
	}
}

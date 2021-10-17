package spinner

import (
	spinner "github.com/theckman/yacspin"
	"time"
)

// SpinnerConfig this func will create a full jub spinner with custom msg
// prefix - gets the msg that the spinner is running on like the msg to display to the console
// successfulMsg - gets the msg to print when done successfully
// failMsg - gets the msg to print when a failure.
func SpinnerConfig(prefix, successfulMsg, failMsg string) (*spinner.Spinner, error) {
	cfg := spinner.Config{
		Frequency:         100 * time.Millisecond,
		CharSet:           spinner.CharSets[59],
		Prefix:            prefix,
		Colors:            []string{"fgGreen"},
		SuffixAutoColon:   true,
		StopCharacter:     "\n✓  " + successfulMsg,
		StopColors:        []string{"fgGreen"},
		StopFailCharacter: "\n✗  " + failMsg,
		StopFailColors:    []string{"fgRed"},
	}
	return spinner.New(cfg)
}

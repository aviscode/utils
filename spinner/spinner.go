package spinner

import (
	spinner "github.com/theckman/yacspin"
	"time"
)

func spinnerConfig(prefix, successfulMsg, failMsg string) spinner.Config {
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
	return cfg
}

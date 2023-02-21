package iprogressbar

import (
	"fmt"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

func GetProgressbar(length int, desc string) *progressbar.ProgressBar {
	return progressbar.NewOptions(
		length,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Print("\n")
		}),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}

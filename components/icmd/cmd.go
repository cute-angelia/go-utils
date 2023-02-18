package icmd

import (
	"github.com/go-cmd/cmd"
	"log"
	"time"
)

func Exec(exePath string, param []string, timeout time.Duration) cmd.Status {
	findCmd := cmd.NewCmd(exePath, param...)
	statusChan := findCmd.Start() // non-blocking

	ticker := time.NewTicker(2 * time.Second)

	// Print last line of stdout every 2s
	go func() {
		for range ticker.C {
			status := findCmd.Status()
			n := len(status.Stdout)
			if n > 0 {
				log.Println(status.Stdout[n-1])
			}
		}
	}()

	// Stop command after time
	go func() {
		<-time.After(timeout)
		findCmd.Stop()
	}()
	// Block waiting for command to exit, be stopped, or be killed
	finalStatus := <-statusChan

	defer findCmd.Stop()

	// log.Println(ijson.Pretty(finalStatus))
	return finalStatus
}

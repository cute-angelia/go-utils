

bak code 2021年09月24日 

```
package icmd

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func Exec(exePath string, param []string, timeout time.Duration, logFile string) (*exec.Cmd, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, exePath, param...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// 处理日志
	if len(logFile) == 0 {
		var b bytes.Buffer
		cmd.Stdout = &b
		cmd.Stderr = &b
		//if err := cmd.Start(); err != nil {
		//	log.Printf("Exce cmd error: %s", err.Error())
		//	return nil, err
		//}
		cmd.Run()

		// log.Println("result:",b.String())

		return cmd, errors.New(b.String())
	} else {
		logOut, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		defer logOut.Close()
		// 将标准输出和标准错误都写到log中
		cmd.Stdout = logOut
		cmd.Stderr = logOut
		log.Printf("Exce cmd: %v", param)
		if err := cmd.Start(); err != nil {
			log.Printf("Exce cmd error: %s", err.Error())
			return nil, err
		}
		return cmd, nil
	}
}

func KillCmd(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return errors.New("process not found or already stopped")
	}

	log.Printf("Kill CMD. pid: %d", cmd.Process.Pid)

	err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	if err != nil {
		return err
	}

	// 如果不wait，则会产生僵尸进程
	go cmd.Wait() // 为了不阻塞，另外起一个线程去Wait

	return nil
}

```
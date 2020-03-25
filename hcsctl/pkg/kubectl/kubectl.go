package kubectl

import (
	"context"
	"io"
	"os/exec"
	"time"
)

const (
	CMD_TIMEOUT = 15 * time.Minute
)

func Run(stdout, stderr io.Writer, arg ...string) error {
	ctx, cancel:= context.WithTimeout(context.Background(), CMD_TIMEOUT)
	defer cancel()

	cmd := exec.CommandContext(ctx, "kubectl", arg...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

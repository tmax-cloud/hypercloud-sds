package kubectl

import (
	"context"
	"io"
	"os/exec"
	"time"
)

const (
	cmdTIMEOUT = 15 * time.Minute *2
)

// Run execute kubectl command
func Run(stdout, stderr io.Writer, arg ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), cmdTIMEOUT)
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

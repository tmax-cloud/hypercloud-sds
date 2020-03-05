package kubectl

import (
	"io"
	"os/exec"
)

func Run(stdout, stderr io.Writer, arg ...string) error {
	cmd := exec.Command("kubectl", arg...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

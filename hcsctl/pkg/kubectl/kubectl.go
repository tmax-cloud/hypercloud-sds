package kubectl

import (
	"context"
	"io"
	"os/exec"
	"time"
)

const (
	cmdTIMEOUT = 15 * time.Minute * 2
)

// KubeConfig represents the location of kubeconfig file
var KubeConfig *string

// Run execute kubectl command
func Run(stdout, stderr io.Writer, arg ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), cmdTIMEOUT)
	defer cancel()

	orderedArgs := append([]string{}, "--kubeconfig", *KubeConfig)
	orderedArgs = append(orderedArgs, arg...)

	cmd := exec.CommandContext(ctx, "kubectl", orderedArgs...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

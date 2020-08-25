package rook

import (
	"bytes"
	"errors"
	"hypercloud-sds/hcsctl/pkg/kubectl"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/golang/glog"
)

// Status checks current ceph status
// Usage: Status()
func Status() error {
	glog.Info("Check current ceph status")

	cmdStdOut, err := execInToolbox(os.Stderr, "ceph", "-s")
	if err != nil {
		return err
	}

	glog.Infof("Ceph status is : \n%s\n", cmdStdOut.String())

	return nil
}

// Exec executes ceph commands
// Usage: Exec([]string{"ceph", "osd", "status"})
func Exec(args []string) error {
	glog.Infof("Executing '%s' on rook-ceph-toolbox", strings.Join(args, " "))

	cmdStdOut, err := execInToolbox(os.Stderr, args...)
	if err != nil {
		return err
	}

	glog.Infof("Stdout of '%s' is : \n%s\n", strings.Join(args, " "), cmdStdOut.String())

	return nil
}

// IsAvailableToolbox check whether the rook-ceph-toolbox pod is available
// Usage: IsAvailableToolbox()
func IsAvailableToolbox() (bool, error) {
	_, err := execInToolbox(os.NewFile(uintptr(syscall.Stdin), os.DevNull), "ls")
	if err != nil {
		return false, err
	}

	return true, nil
}

func getRunningToolboxName(errWriter io.Writer) (string, error) {
	var cmdStdOut bytes.Buffer

	// get the name of Running phased toolbox
	err := kubectl.Run(&cmdStdOut, errWriter, "get", "pod", "-n", "rook-ceph", "--selector=app=rook-ceph-tools",
		"--field-selector", "status.phase=Running", "-o", "custom-columns=name:.metadata.name", "--no-headers")

	if cmdStdOut.String() == "" {
		return "", errors.New("there isn't any running phased toolbox pod")
	} else if err != nil {
		return "", err
	}

	// if there are many running toolbox pods, just use the first element
	toolboxName := strings.Split(cmdStdOut.String(), "\n")[0]
	glog.Info("rook-ceph-toolbox pod's name is : " + toolboxName)

	return toolboxName, nil
}

func execInToolbox(errWriter io.Writer, cmd ...string) (bytes.Buffer, error) {
	var cmdStdOut bytes.Buffer

	toolboxName, err := getRunningToolboxName(errWriter)

	if err != nil || toolboxName == "" {
		return cmdStdOut, err
	}

	kubectlCmd := []string{"exec", "-n", "rook-ceph", toolboxName, "--"}
	kubectlCmd = append(kubectlCmd, cmd...)

	err = kubectl.Run(&cmdStdOut, errWriter, kubectlCmd...)

	if err != nil {
		return cmdStdOut, err
	}

	return cmdStdOut, nil
}

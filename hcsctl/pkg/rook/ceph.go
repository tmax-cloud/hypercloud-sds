package rook

import (
	"bytes"
	"hypercloud-storage/hcsctl/pkg/kubectl"
	"os"
	"strings"

	"github.com/golang/glog"
)

var cmdStdOut bytes.Buffer

// Status checks current ceph status
// Usage: Status()
func Status() error {
	glog.Info("Check current ceph status")

	err := execInToolbox("ceph", "-s")
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

	err := execInToolbox(args...)
	if err != nil {
		return err
	}

	glog.Infof("Stdout of '%s' is : \n%s\n", strings.Join(args, " "), cmdStdOut.String())

	return nil
}

func execInToolbox(cmd ...string) error {
	// get toolbox pod name
	var toolboxPodName bytes.Buffer

	err := kubectl.Run(&toolboxPodName, os.Stderr, "get", "pod", "-n", "rook-ceph",
		"--selector=app=rook-ceph-tools", "-o", "jsonpath={.items[0].metadata.name}")
	if err != nil {
		return err
	}

	glog.Info("rook-ceph-toolbox pod name is : " + toolboxPodName.String())

	// Execute cmd in the toolbox container
	kubectlCmd := []string{"exec", "-n", "rook-ceph", toolboxPodName.String(), "--"}
	kubectlCmd = append(kubectlCmd, cmd...)

	err = kubectl.Run(&cmdStdOut, os.Stderr, kubectlCmd...)
	if err != nil {
		return err
	}

	return nil
}

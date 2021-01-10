package rook

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"hypercloud-sds/hcsctl/pkg/kubectl"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

const (
	rookInfoDirName      string = "rookInfo"
	clusterInfoDirName   string = "cluster"
	podInfoDirName       string = "pod"
	cephInfoDirName      string = "ceph"
	pvcInfoDirName       string = "pvc"
	cephMountInfoDirName string = "cephMount"
	kubeCommand          string = "kubeCommand"
	cephCommand          string = "cephCommand"
	commonCommand        string = "commonCommand"
)

// GetIssueTemplate gets current rook, node info
func GetIssueTemplate() error {
	var failMessage, errorMessage string

	glog.Info("Start Get IssueTemplate")
	glog.Info("Create info folder")

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	rookInfoDir, err := createInfoDir(wd, rookInfoDirName)
	if err != nil {
		return err
	}

	glog.Info("Get cluster info")

	errorMessage = getClusterInfo(rookInfoDir)
	if errorMessage != "" {
		failMessage += "[Fail][getClusterInfo]\n" + errorMessage + "\n"

		glog.Error("Failure occurred in getClusterInfo")
	}

	glog.Info("Get pvc info")

	errorMessage = getPvcInfo(rookInfoDir)
	if errorMessage != "" {
		failMessage += "[Fail][getPvcInfo]\n" + errorMessage + "\n"

		glog.Error("Failure occurred in getPvcInfo")
	}

	glog.Info("Get pod info")

	errorMessage = getPodInfo(rookInfoDir)
	if errorMessage != "" {
		failMessage += "[Fail][getPodInfo]\n" + errorMessage + "\n"

		glog.Error("Failure occurred in getPodInfo")
	}

	glog.Info("Get ceph info")

	errorMessage = getCephInfo(rookInfoDir)
	if errorMessage != "" {
		failMessage += "[Fail][getCephInfo]\n" + errorMessage + "\n"

		glog.Error("Failure occurred in getCephInfo")
	}

	glog.Info("Get ceph mount info")

	errorMessage = getCephMountInfo(rookInfoDir)
	if errorMessage != "" {
		failMessage += "[Fail][getCephMountInfo]\n" + errorMessage + "\n"

		glog.Error("Failure occurred in getCephMountInfo")
	}

	err = writeInfoFile(rookInfoDir, "Fail", failMessage)
	if err != nil {
		glog.Error(err.Error())
	}

	glog.Info("Compress " + rookInfoDirName + ".tar.gz")

	err = zip(rookInfoDirName, "")
	if err != nil {
		glog.Error(err.Error())
	}

	return nil
}

func getPvcInfo(rookInfoDir string) string {
	var stdout, stderr bytes.Buffer

	var infoName, failMessage string

	var jsonData map[string]interface{}

	infoDir, err := createInfoDir(rookInfoDir, pvcInfoDirName)
	if err != nil {
		failMessage += err.Error() + "\n"
		return failMessage
	}

	cmdList := [][]string{}
	cmdList = append(cmdList, []string{"pvc_list", "get", "pvc", "-A", "-o", "wide"},
		[]string{"pv_list", "get", "pv", "-o", "wide"})

	err = kubectl.Run(&stdout, &stderr, "get", "pvc", "-A", "-o", "json")
	if err != nil {
		failMessage += "fail to kubectl Run, error: " + err.Error() + ", message: " + stderr.String() + "\n"
		return failMessage
	}

	err = json.Unmarshal(stdout.Bytes(), &jsonData)
	if err != nil {
		failMessage += err.Error() + "\n"
		return failMessage
	}

	pvcList := jsonData["items"].([]interface{})

	for _, value := range pvcList {
		pvcData := value.(map[string]interface{})
		pvcPhase := pvcData["status"].(map[string]interface{})["phase"].(string)

		if pvcPhase == "Bound" {
			continue
		}

		pvcName := pvcData["metadata"].(map[string]interface{})["name"].(string)
		pvcNamespace := pvcData["metadata"].(map[string]interface{})["namespace"].(string)
		infoName = "problem_pvc_describe_" + pvcNamespace + "_" + pvcName
		cmdList = append(cmdList, []string{infoName, "describe", "pvc", pvcName, "-n", pvcNamespace})
	}

	for _, cmd := range cmdList {
		err = writeInfo(infoDir, kubeCommand, cmd[0], cmd[1:])
		if err != nil {
			failMessage += err.Error() + "\n"
		}
	}

	return failMessage
}

func getCephMountInfo(rookInfoDir string) string {
	var stdout, stderr bytes.Buffer

	var failMessage, infoName string

	var jsonData map[string]interface{}

	infoDir, err := createInfoDir(rookInfoDir, cephMountInfoDirName)
	if err != nil {
		failMessage += err.Error() + "\n"
		return failMessage
	}

	cmdList := [][]string{}

	err = kubectl.Run(&stdout, &stderr, "get", "pod", "-n", "rook-ceph", "--selector=app=csi-rbdplugin",
		"-o", "json")
	if err != nil {
		failMessage += "fail to kubectl Run, error: " + err.Error() + ", message: " + stderr.String() + "\n"
		return failMessage
	}

	err = json.Unmarshal(stdout.Bytes(), &jsonData)
	if err != nil {
		failMessage += err.Error() + "\n"
		return failMessage
	}

	podList := jsonData["items"].([]interface{})

	for _, value := range podList {
		podData := value.(map[string]interface{})
		podName := podData["metadata"].(map[string]interface{})["name"].(string)
		nodeName := podData["spec"].(map[string]interface{})["nodeName"].(string)
		infoName = "rbd_showmapped_" + nodeName
		cmdList = append(cmdList, []string{infoName, "exec", podName, "-c", "csi-rbdplugin", "-n", "rook-ceph",
			"--", "rbd", "showmapped"})
	}

	for _, cmd := range cmdList {
		err = writeInfo(infoDir, cephCommand, cmd[0], cmd[1:])
		if err != nil {
			failMessage += err.Error() + "\n"
		}
	}

	return failMessage
}

func getPodInfo(rookInfoDir string) string {
	var stdout, stderr bytes.Buffer

	var infoName, failMessage string

	var jsonData map[string]interface{}

	infoDir, err := createInfoDir(rookInfoDir, podInfoDirName)
	if err != nil {
		failMessage += err.Error() + "\n"
		return failMessage
	}

	cmdList := [][]string{}

	cmdList = append(cmdList, []string{"pod_list_all", "get", "pods", "-A", "-o", "wide"},
		[]string{"pod_list_rook-ceph", "get", "pods", "-n", "rook-ceph", "-o", "wide"})

	err = kubectl.Run(&stdout, &stderr, "get", "pods", "-n", "rook-ceph",
		"-o", "json")
	if err != nil {
		failMessage += "fail to kubectl Run, error: " + err.Error() + ", message: " + stderr.String() + "\n"
		return failMessage
	}

	err = json.Unmarshal(stdout.Bytes(), &jsonData)
	if err != nil {
		failMessage += err.Error() + "\n"
		return failMessage
	}

	podList := jsonData["items"].([]interface{})

	for _, value := range podList {
		podData := value.(map[string]interface{})
		podName := podData["metadata"].(map[string]interface{})["name"].(string)
		nodeName := podData["spec"].(map[string]interface{})["nodeName"].(string)
		containerList := podData["spec"].(map[string]interface{})["containers"].([]interface{})

		infoName = "pod_describe_" + nodeName + "_" + podName
		cmdList = append(cmdList, []string{infoName, "describe", "pod", podName, "-n", "rook-ceph"})

		for _, value := range containerList {
			containerData := value.(map[string]interface{})
			containerName := containerData["name"].(string)
			infoName := "pod_log_" + nodeName + "_" + podName + "_" + containerName
			cmdList = append(cmdList, []string{infoName, "logs", podName, "-c", containerName, "-n", "rook-ceph"})
		}
	}

	for _, cmd := range cmdList {
		err = writeInfo(infoDir, kubeCommand, cmd[0], cmd[1:])
		if err != nil {
			failMessage += err.Error() + "\n"
		}
	}

	return failMessage
}

func getClusterInfo(rookInfoDir string) string {
	var stdout, stderr bytes.Buffer

	var infoName, failMessage string

	var jsonData map[string]interface{}

	infoDir, err := createInfoDir(rookInfoDir, clusterInfoDirName)

	if err != nil {
		failMessage += err.Error() + "\n"
		return failMessage
	}

	execCmd := []string{"kubeadm_info", "kubeadm", "config", "view"}
	err = writeInfo(infoDir, commonCommand, execCmd[0], execCmd[1:])

	if err != nil {
		failMessage += err.Error() + "\n"
	}

	cmdList := [][]string{}
	cmdList = append(cmdList, []string{"cephcluster", "get", "cephcluster", "rook-ceph", "-n", "rook-ceph", "-o", "yaml"},
		[]string{"node_list", "get", "nodes", "-o", "wide"})

	err = kubectl.Run(&stdout, &stderr, "get", "nodes", "-o", "json")

	if err != nil {
		failMessage += "fail to kubectl Run, error: " + err.Error() + ", message: " + stderr.String() + "\n"
		return failMessage
	}

	err = json.Unmarshal(stdout.Bytes(), &jsonData)
	if err != nil {
		failMessage += err.Error() + "\n"
		return failMessage
	}

	nodeList := jsonData["items"].([]interface{})

	for _, value := range nodeList {
		nodeData := value.(map[string]interface{})
		nodeName := nodeData["metadata"].(map[string]interface{})["name"].(string)
		infoName = "node_describe_" + nodeName
		cmdList = append(cmdList, []string{infoName, "describe", "node", nodeName})
	}

	for _, cmd := range cmdList {
		err = writeInfo(infoDir, kubeCommand, cmd[0], cmd[1:])
		if err != nil {
			failMessage += err.Error() + "\n"
		}
	}

	return failMessage
}

func getCephInfo(rookInfoDir string) string {
	var failMessage string

	var cephTimeout string = "10"

	infoDir, err := createInfoDir(rookInfoDir, cephInfoDirName)
	if err != nil {
		failMessage += err.Error() + "\n"
		return failMessage
	}

	cmdList := [][]string{}
	cmdList = append(cmdList, []string{"ceph_status", "ceph", "status"},
		[]string{"ceph_health_detail", "ceph", "health", "detail"},
		[]string{"ceph_df_detail", "ceph", "df", "detail"},
		[]string{"ceph_osd_pool_ls_detail", "ceph", "osd", "pool", "ls", "detail"},
		[]string{"ceph_node_ls_all", "ceph", "node", "ls", "all"},
		[]string{"ceph_osd_df", "ceph", "osd", "df"},
		[]string{"ceph_osd_status", "ceph", "osd", "status"},
		[]string{"ceph_osd_tree", "ceph", "osd", "tree"},
		[]string{"ceph_pg_dump", "ceph", "pg", "dump"},
		[]string{"ceph_pg_dump_pgs_brief", "ceph", "pg", "dump", "pgs_brief"},
		[]string{"ceph_device_ls", "ceph", "device", "ls"},
		[]string{"ceph_osd_blacklist_ls", "ceph", "osd", "blacklist", "ls"},
		[]string{"ceph_tell_mds_client_ls", "ceph", "tell", "mds.0", "client", "ls"},
		[]string{"ceph_fs_subvolume_ls", "ceph", "fs", "subvolume", "ls", "myfs", "csi"},
		[]string{"rbd_ls", "rbd", "ls", "replicapool"})

	for _, cmd := range cmdList {
		if cmd[1] == "ceph" {
			cmd = append(cmd, "--connect-timeout", cephTimeout)
		}

		err = writeInfo(infoDir, cephCommand, cmd[0], cmd[1:])

		if err != nil {
			failMessage += err.Error() + "\n"
		}
	}

	return failMessage
}

func createInfoDir(dirPath, dirName string) (string, error) {
	infoDir := path.Join(dirPath, dirName)
	err := os.MkdirAll(infoDir, 0755)

	if err != nil {
		err = errors.Wrap(err, "fail to create directory "+dirName)
	}

	return infoDir, err
}

func writeInfo(infoDir, command, infoName string, cmd []string) error {
	var stdout, stderr bytes.Buffer

	var err error

	message := "[Run][" + command + "] infoName: " + infoName + ", cmd: [" + strings.Join(cmd, " ") + "]"
	glog.Info(message)

	switch command {
	case kubeCommand:
		err = kubectl.Run(&stdout, &stderr, cmd...)
	case cephCommand:
		stdout, err = execInToolbox(&stderr, cmd...)
	case commonCommand:
		err = executeCommand(&stdout, &stderr, cmd[0], cmd[1:]...)
	default:
		err = errors.New("ERROR: unknown command")
	}

	if err != nil {
		errorMessage := "[Fail][" + command + "] infoName: " + infoName +
			", cmd: [" + strings.Join(cmd, " ") + "], message: " + stderr.String()
		err = errors.Wrap(err, errorMessage)

		return err
	}

	return writeInfoFile(infoDir, infoName, stdout.String())
}

func writeInfoFile(infoDir, infoName, output string) error {
	filePath := path.Join(infoDir, infoName)
	f, err := os.Create(filePath)

	if err != nil {
		err = errors.Wrap(err, "fail to write "+infoName)
		return err
	}

	defer deferCheck(f.Close)

	_, err = f.WriteString(output)

	if err != nil {
		err = errors.Wrap(err, "fail to write "+infoName)
	}

	return err
}

func zip(source, target string) error {
	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.tar.gz", filename))
	writer, err := os.Create(target)

	if err != nil {
		return err
	}

	defer deferCheck(writer.Close)

	archiver := gzip.NewWriter(writer)
	defer deferCheck(archiver.Close)

	tarball := tar.NewWriter(archiver)
	defer deferCheck(tarball.Close)

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err2 := tarball.WriteHeader(header); err2 != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}

			defer deferCheck(file.Close)
			_, err = io.Copy(tarball, file)
			return err
		})
}

func deferCheck(f func() error) {
	if err := f(); err != nil {
		glog.Error(err.Error())
	}
}

func executeCommand(stdout, stderr io.Writer, command string, arg ...string) error {
	const execTimeout = 10
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*execTimeout)

	defer cancel()

	orderedArgs := append([]string{}, arg...)

	cmd := exec.CommandContext(ctx, command, orderedArgs...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

#!/bin/bash

# Exit script when any commands failed
set -eo pipefail

#####################
# define common variables
srcDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/../"
multinodeK8sDir="${srcDir}/hack/k8s-vagrant-multi-node"
clusterConfigDir="$multinodeK8sDir"/.created-cluster-config

#####################
# util functions
function kubectl_wait_avail() {
  namespace=$1
  waitFor=$2
  timeoutSecond=$3

  echo "Waiting for available ${waitFor}..."

  for ((seconds = 0; seconds <= timeoutSecond; seconds = seconds + 1)); do
    if kubectl describe -n "${namespace}" "${waitFor}" &>/dev/null; then break; fi
    sleep 1
  done

  kubectl wait -n "${namespace}" --for=condition=available "${waitFor}" --timeout="${timeoutSecond}"s
}

function kubectl_wait_delete() {
  namespace=$1
  waitFor=$2
  timeoutSecond=$3

  echo "Waiting for deleting ${waitFor}..."

  for ((seconds = 0; seconds <= timeoutSecond; seconds = seconds + 1)); do
    if ! kubectl describe -n "${namespace}" "${waitFor}" &>/dev/null; then break; fi
    sleep 1
  done
}

function print_red() {
  msg=$1

  echo -e "\e[0;31;47m${msg}\e[0m"
}

function sourceConfigFromFile() {
  if [ ! -f "$clusterConfigDir" ]; then
    echo "There isn't any config file in $clusterConfigDir"
    exit 1
  fi

  source "$clusterConfigDir"
}

function checkWhetherClusterIsRunning() {
  if [ -f "$clusterConfigDir" ]; then
    echo "There is running cluster already in $clusterConfigDir"
    exit 1
  fi
}

function waitMinikubeSsh() {
  local tries=100
  while ((tries > 0)); do
    if minikube ssh echo connected &>/dev/null; then
      return 0
    fi
    tries=$((tries - 1))
    sleep 0.1
  done
  echo ERROR: ssh did not come up >&2
  exit 1
}
#####################
# main features
function clusterUp() {
  checkWhetherClusterIsRunning

  local nodeCount=${NODE_COUNT:?"need env NODE_COUNT"}
  local os=${BOX_OS:?"need env BOX_OS"}
  local k8sVersion=${KUBERNETES_VERSION:?"need env KUBERNETES_VERSION"}

  # temporarily save config for deleting later
  echo "NODE_COUNT=$nodeCount BOX_OS=$os KUBERNETES_VERSION=$k8sVersion" >"$clusterConfigDir"

  echo "create $os-$nodeCount-worker-node-cluster with k8s version : $k8sVersion"

  DISK_COUNT=2 DISK_SIZE_GB=5 make --directory "${multinodeK8sDir}" up -j2
  print_red "========================== cluster created =========================="
  echo "However, you may need to wait a few seconds until nodes are ready"
  kubectl get nodes
}

function clusterClean() {
  sourceConfigFromFile

  NODE_COUNT="$NODE_COUNT" make --directory "${multinodeK8sDir}" clean -j"$(nproc)"

  rm "$clusterConfigDir"
  echo "The cluster-config file deleted in $clusterConfigDir"
}

function minikubeUp() {
  # TODO: minikube version
  # TODO: multinode cluster
  minikube start --vm-driver=virtualbox
  waitMinikubeSsh

  # Rook에서 사용할 디렉토리를 마운트
  minikube ssh "sudo mkdir -p /mnt/sda1/${PWD}; sudo mkdir -p $(dirname "$PWD"); sudo ln -s /mnt/sda1/${PWD} $(dirname "$PWD")/"
  minikube ssh "sudo mkdir -p /mnt/sda1/var/lib/rook;sudo ln -s /mnt/sda1/var/lib/rook /var/lib/rook"

  kubectl patch storageclass standard '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"false"}}}'
}

function minikubeClean() {
  minikube delete
}

#####################
# main logic
case "$1" in
minikubeUp)
  minikubeUp
  ;;
minikubeClean)
  minikubeClean
  ;;
clusterUp)
  clusterUp
  ;;
clusterClean)
  clusterClean
  ;;
*)
  echo "usage:" >&2
  echo "  $0 minikubeUp" >&2
  echo "  $0 minikubeClean" >&2
  echo "  BOX_OS=fedora NODE_COUNT=3 KUBERNETES_VERSION=1.15.3 $0 clusterUp" >&2
  echo "  $0 clusterClean" >&2
  ;;
esac

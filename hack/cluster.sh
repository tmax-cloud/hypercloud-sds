#!/bin/bash

# Exit script when any commands failed
set -eo pipefail

# define common variables
srcDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/../"
multinodeK8sDir="${srcDir}/hack/k8s-vagrant-multi-node"
clusterConfigDir="$multinodeK8sDir"/.created-cluster-config

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
  echo "========================== cluster created =========================="
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
  # TODO virtualbox driver does not support memory, cpu options
  minikube start --vm-driver=virtualbox --memory='4096mb' --cpus=2 --disk-size='40000mb' --kubernetes-version="$KUBERNETES_VERSION"
  waitMinikubeSsh

  # wait until nodes ready
  sleep 30
  kubectl get nodes

  # TODO 여기서 이 명령을 내리는 게 맞는지 ?
  kubectl patch storageclass standard -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"false"}}}'
}

function minikubeClean() {
  minikube delete
}

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

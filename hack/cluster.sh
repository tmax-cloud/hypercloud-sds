#!/bin/bash

# shellcheck source=common.sh
. "$(dirname "$0")/common.sh"

multinodeK8sDir="${srcDir}/hack/k8s-vagrant-multi-node"

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

function minikubeUp() {
  # TODO: minikube version
  # TODO: multinode cluster
  # TODO: minikube profile

  minikube start
  waitMinikubeSsh

  # Rook에서 사용할 디렉토리를 마운트
  minikube ssh "sudo mkdir -p /mnt/sda1/${PWD}; sudo mkdir -p $(dirname "$PWD"); sudo ln -s /mnt/sda1/${PWD} $(dirname "$PWD")/"
  minikube ssh "sudo mkdir -p /mnt/sda1/var/lib/rook;sudo ln -s /mnt/sda1/var/lib/rook /var/lib/rook"
}

function minikubeClean() {
  minikube delete
}

function clusterUp() {
  # TODO vagrant global-status --prune check or ps -ef | vagrant check
  DISK_COUNT=2 DISK_SIZE_GB=5 NODE_COUNT=3 make --directory "${multinodeK8sDir}" up -j"$(nproc)"
  print_red "========================== cluster created =========================="
  echo "However, you may need to wait some seconds until nodes are ready"
}

function clusterClean() {
  DISK_COUNT=2 DISK_SIZE_GB=5 NODE_COUNT=3 make --directory "${multinodeK8sDir}" clean -j"$(nproc)"
}

function main() {
  case "${1:-}" in
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
    echo "  $0 clusterUp" >&2
    echo "  $0 clusterClean" >&2
    ;;
  esac
}

main "$1"

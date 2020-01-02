#!/bin/bash

scriptdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="${REPO_DIR:-${scriptdir}/../.cache/k8s-vagrant-multi-node/}"

function get_k8s_vagrant() {
  if [ ! -d "${REPO_DIR}" ]; then
    echo "k8s-vagrant-multi-node not found in .cache dir. Cloning.."
    mkdir -p "${REPO_DIR}"
    git clone https://github.com/galexrt/k8s-vagrant-multi-node.git "${REPO_DIR}"
  else
    git -C "${REPO_DIR}" fetch origin
  fi
}

function main() {
  cd "${REPO_DIR}" || {
    echo "failed to access k8s-vagrant-multi-node dir ${REPO_DIR}. exiting."
    exit 1
  }

  case "${1:-}" in
  up)
    NODE_COUNT=3 DISK_COUNT=2 DISK_SIZE_GB=5 make up -j4
    ;;
  clean)
    NODE_COUNT=3 DISK_COUNT=2 DISK_SIZE_GB=5 make clean -j4
    NODE_COUNT=3 DISK_COUNT=2 DISK_SIZE_GB=5 make clean-data
    ;;
  status)
    make status
    ;;
  *)
    echo "usage:" >&2
    echo "  $0 up" >&2
    echo "  $0 clean" >&2
    echo "  $0 status" >&2
    echo "  $0 help" >&2
    ;;
  esac
}

get_k8s_vagrant
main $1

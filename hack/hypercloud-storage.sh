#!/bin/bash

# Exit script when any commands failed
set -e
set -o pipefail

# Print all commands
set -x

srcDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/../"
deployDir="${srcDir}/deploy/"

function install() {
  # CDI
  kubectl create -f "$deployDir/cdi/cdi-operator.yaml"
  kubectl create -f "$deployDir/cdi/cdi-cr.yaml"
  kubectl_wait cdi deployment/cdi-apiserver 60
  kubectl_wait cdi deployment/cdi-deployment 60
  kubectl_wait cdi deployment/cdi-operator 60
  kubectl_wait cdi deployment/cdi-uploadproxy 60
  kubectl_wait cdi cdi.cdi.kubevirt.io/cdi 60
}

function kubectl_wait() {
  namespace=$1
  waitFor=$2
  timeoutSecond=$3

  for ((seconds = 0; seconds <= timeoutSecond; seconds = seconds + 1)); do
    if kubectl describe -n ${namespace} ${waitFor} &>/dev/null; then break; fi
    sleep 1
  done

  kubectl wait -n ${namespace} --for=condition=available ${waitFor} --timeout=${timeoutSecond}s
}

function uninstall() {
  # CDI
  kubectl delete --wait=true --ignore-not-found=true -f "$deployDir/cdi/cdi-cr.yaml"
  kubectl delete --ignore-not-found=true -f "$deployDir/cdi/cdi-operator.yaml"
}

function clusterUp() {
  # TODO: minikube version
  # TODO: multinode cluster
  # TODO: minikube profile

  minikube start
}

function clusterClean() {
  minikube delete
}

function e2e() {
  echo "e2e test ok"
}

function main() {
  case "${1:-}" in
  install)
    install
    ;;
  uninstall)
    uninstall
    ;;
  clusterUp)
    clusterUp
    ;;
  clusterClean)
    clusterClean
    ;;
  e2e)
    e2e
    ;;
  *)
    set +x
    echo "usage:" >&2
    echo "  $0 install" >&2
    echo "  $0 uninstall" >&2
    echo "  $0 clusterUp" >&2
    echo "  $0 clusterClean" >&2
    echo "  $0 e2e" >&2
    echo "  $0 help" >&2
    ;;
  esac
}

main $1

#!/bin/bash

# Exit script when any commands failed
set -eo pipefail

srcDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/../"
manifestsDir="${srcDir}/manifests"
deployDir="${srcDir}/.cache/deploy"

# include
. ${srcDir}/cluster_config

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

function print_red() {
  msg=$1

  echo -e "\e[0;31;47m${msg}\e[0m"
}

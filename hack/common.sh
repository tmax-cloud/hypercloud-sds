#!/bin/bash

# Exit script when any commands failed
set -eo pipefail

#####################
# define common variables
srcDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/../"

templatesDir="${srcDir}/templates"
deployDir="${srcDir}/_deploy"
buildDir="${srcDir}/_build"

pkgDir1="${srcDir}/pkg/test-installation" # TODO it is hardcoded to test temporarily!!!
pkgDir2="${srcDir}/pkg/test-pod-networking" # TODO it is hardcoded to test temporarily!!!

testDir="${buildDir}/test"

build_out1="${testDir}/out1" # TODO it is hardcoded to test temporarily!!!
build_out2="${testDir}/out2"
######################

# include
. ${srcDir}/cluster_config

function kubectl_wait_avail() {
  namespace=$1
  waitFor=$2
  timeoutSecond=$3

  echo "Waiting for available ${waitFor}..."

  for ((seconds = 0; seconds <= timeoutSecond; seconds = seconds + 1)); do
    if kubectl describe -n ${namespace} ${waitFor} &>/dev/null; then break; fi
    sleep 1
  done

  kubectl wait -n ${namespace} --for=condition=available ${waitFor} --timeout=${timeoutSecond}s
}

function kubectl_wait_delete() {
  namespace=$1
  waitFor=$2
  timeoutSecond=$3

  echo "Waiting for deleting ${waitFor}..."

  for ((seconds = 0; seconds <= timeoutSecond; seconds = seconds + 1)); do
    if ! kubectl describe -n ${namespace} ${waitFor} &>/dev/null; then break; fi
    sleep 1
  done
}

function print_red() {
  msg=$1

  echo -e "\e[0;31;47m${msg}\e[0m"
}

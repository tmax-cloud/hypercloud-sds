#!/bin/bash

# Exit script when any commands failed
set -eo pipefail

#####################
# define common variables
srcDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/../"

# shellcheck disable=SC2034
templatesDir="${srcDir}/templates"
# shellcheck disable=SC2034
deployDir="${srcDir}/_deploy"
# shellcheck disable=SC2034
buildDir="${srcDir}/_build"
# shellcheck disable=SC2034
testDir="${srcDir}/tests"
######################

# include
# shellcheck disable=SC1090
. "${srcDir}"/cluster_config

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

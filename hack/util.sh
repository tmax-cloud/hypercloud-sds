#!/bin/bash

# Exit script when any commands failed
set -eo pipefail

# TODO 추후 helm install check 방식 변경되면 해당 파일 삭제
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